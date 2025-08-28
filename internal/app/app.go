package app

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"git.mos.ru/buch-cloud/moscow-team-2.0/pud/auth.git/cmd/auth/config"
	employeev1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/employees/employee/v1"
	personv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/employees/person/v1"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/tree-alive.git"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	grpcApi "git.mos.ru/buch-cloud/moscow-team-2.0/pud/auth.git/internal/api/grpc"
	"git.mos.ru/buch-cloud/moscow-team-2.0/pud/auth.git/internal/api/http/service"
	"git.mos.ru/buch-cloud/moscow-team-2.0/pud/auth.git/internal/client/http/kadry"
	"git.mos.ru/buch-cloud/moscow-team-2.0/pud/auth.git/internal/client/http/sudir"
	"git.mos.ru/buch-cloud/moscow-team-2.0/pud/auth.git/internal/repositories"
	"git.mos.ru/buch-cloud/moscow-team-2.0/pud/auth.git/internal/source/cache"
	"git.mos.ru/buch-cloud/moscow-team-2.0/pud/auth.git/internal/usecase"
)

type Environment string

const (
	EnvironmentDevelop Environment = "develop"
	EnvironmentProd    Environment = "prod"
)

type AppInfo struct {
	Name      string
	Instance  Environment
	BuildTime string
	Commit    string
	Release   string
}

func (ai AppInfo) GetReleaseVersion() string {
	return fmt.Sprintf("%s@%s", ai.Name, ai.Release)
}

var ApplicationInfo *AppInfo

type app struct {
	tree              tree.Tree
	config            *config.Config
	serviceHTTPServer Server
	grpcServer        Server
	redisClient       cache.Redis
	logger            ditzap.Logger
	stop              context.CancelFunc
}

func NewApp(cfg *config.Config, appInfo *AppInfo, logger ditzap.Logger) *app {
	ApplicationInfo = appInfo
	return &app{
		config: cfg,
		logger: logger,
	}
}

func (a *app) Run(ctx context.Context) {
	appCtx, cancelApp := context.WithCancel(ctx)
	a.stop = cancelApp
	defer func() {
		if e := recover(); e != nil {
			a.logger.Error("application shutdown", zap.Error(fmt.Errorf("%s", e)))
			cancelApp()
		}
	}()
	// AfterFunc-функция при завершении контекста для graceful shutdown
	context.AfterFunc(appCtx, func() {
		defer func() {
			if e := recover(); e != nil {
				a.logger.Error("context after-func", zap.Error(fmt.Errorf("%s", e)))
				cancelApp()
			}
		}()
		a.contextCancelFunc(appCtx)
	})

	environment := EnvironmentDevelop
	if a.config.Environment == "prod" {
		environment = EnvironmentProd
	}
	ApplicationInfo.Instance = environment

	employeesConn, err := a.connectToProviders(ctx, a.config.Endpoints)
	if err != nil {
		a.logger.Error("connect to providers error", zap.Error(err))
	}

	employeesClient := employeev1.NewEmployeesAPIClient(employeesConn)
	personClient := personv1.NewPersonAPIClient(employeesConn)

	sudirClient := sudir.NewClient(
		a.config.OAuth.URL,
		a.config.OAuth.ClientID,
		a.config.OAuth.ClientSecret,
		a.logger,
	)

	kadryClient := kadry.NewClient(
		a.config.SKS.URL,
		a.config.SKS.SubscriberID,
		a.config.SKS.Secret,
		a.config.SKS.UserID,
		a.logger,
	)

	a.redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", a.config.Redis.Host, a.config.Redis.Port),
		Username: a.config.Redis.Username,
		Password: a.config.Redis.Password,
		DB:       a.config.Redis.DB,
	})

	cacheSource := cache.NewCacheSource(a.redisClient)

	stateRepository := repositories.NewStateRepository(a.config.Redis.Prefix, cacheSource, a.logger)
	tokenRepository := repositories.NewTokenRepository(a.config.Redis.Prefix, cacheSource, a.logger)
	employeeRepository := repositories.NewEmployeeRepository(
		a.config.Redis.Prefix,
		cacheSource,
		employeesClient,
		personClient,
		a.logger,
	)

	authUseCase := usecase.NewAuthUseCase(
		sudirClient,
		kadryClient,
		stateRepository,
		tokenRepository,
		employeeRepository,
		a.logger,
	)

	a.tree = tree.NewTree()
	a.tree.Alive()

	a.initMetrics()

	wg := &sync.WaitGroup{}

	// старт сервисного сервера
	wg.Add(1)
	serviceHttpBranch := a.tree.GrowBranch("service-http-server")
	go a.startServiceServer(appCtx, serviceHttpBranch, wg)

	// Старт GRPC-сервера
	wg.Add(1)
	grpcBranch := a.tree.GrowBranch("grpc-server")
	go a.startGRPCServer(
		appCtx,
		grpcBranch,
		wg,
		authUseCase,
	)

	a.tree.Ready()
	wg.Wait()
}

// GracefulShutdown graceful shutdown приложения.
func (a *app) GracefulShutdown(c context.Context) (err error) {
	ctx := context.WithoutCancel(c)
	if a.tree != nil {
		a.tree.Die()
	}

	if a.serviceHTTPServer != nil {
		if srvErr := a.serviceHTTPServer.Shutdown(ctx); srvErr != nil {
			err = errors.Join(err, fmt.Errorf("can't shutdown service server: %w", srvErr))
		}
	}

	if a.grpcServer != nil {
		grpcErr := a.grpcServer.Shutdown(ctx)
		if grpcErr != nil {
			err = errors.Join(err, fmt.Errorf("can't shutdown grpc-server: %w", grpcErr))
		}
	}

	if a.redisClient != nil {
		if redisErr := a.redisClient.Close(); redisErr != nil {
			err = errors.Join(err, fmt.Errorf("can't close redis connection: %w", redisErr))
		}
	}
	return
}

func (a *app) contextCancelFunc(ctx context.Context) {
	err := a.GracefulShutdown(ctx)
	if err != nil {
		a.logger.Error("graceful shutdown error", zap.Error(err))
	}
}

func (a *app) connectToProviders(ctx context.Context, endpoints config.Endpoints) (
	employeesConn *grpc.ClientConn,
	err error,
) {
	var (
		dialErr error
	)
	grpcPrometheus.EnableClientHandlingTimeHistogram()

	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	// Инициализация подключения к сервису employees
	employeesConn, dialErr = a.dial(
		ctx,
		endpoints.EmployeesEndpoint,
		false,
		dialOpts...,
	)
	if dialErr != nil {
		err = fmt.Errorf("can't connect to employees service: %w", dialErr)
		return
	}

	return
}

func (a *app) initMetrics() {
	// Добавление префикса к метрикам
	appName := strings.Replace(ApplicationInfo.Name, "-", "_", -1)
	metricRegisterer := prometheus.WrapRegistererWithPrefix(appName+"_", prometheus.DefaultRegisterer)
	prometheus.DefaultRegisterer = metricRegisterer
	if gather, ok := metricRegisterer.(prometheus.Gatherer); ok {
		prometheus.DefaultGatherer = gather
	}
}

func (a *app) startServiceServer(c context.Context, serviceHttpBranch tree.Branch, wg *sync.WaitGroup) {
	ctx, serverCancel := context.WithCancel(c)
	defer func() {
		if e := recover(); e != nil {
			a.logger.Error("service http start panic", zap.Error(fmt.Errorf("%s", e)))
		}
		serviceHttpBranch.Die()
		serverCancel()
		wg.Done()
	}()
	a.serviceHTTPServer = service.NewServer(serviceHttpBranch, &service.AppInfo{
		Name:      ApplicationInfo.Name,
		Instance:  string(ApplicationInfo.Instance),
		BuildTime: ApplicationInfo.BuildTime,
		Commit:    ApplicationInfo.Commit,
		Release:   ApplicationInfo.Release,
	})
	err := a.serviceHTTPServer.Run(context.WithValue(ctx, service.ContextKeyServiceAddr, a.config.ServiceHTTPHost)) //nolint:staticcheck
	// Отменяем контекст, если HTTP-сервер завершил работу
	if err != nil {
		a.logger.Error("service http server is shutdown", zap.Error(err))
	}
}

func (a *app) startGRPCServer(
	c context.Context,
	grpcBranch tree.Branch,
	wg *sync.WaitGroup,
	authInteractor grpcApi.AuthInteractor,
) {
	ctx, serverCancel := context.WithCancel(c)
	defer func() {
		if e := recover(); e != nil {
			a.logger.Error("grpc start panic", zap.Error(fmt.Errorf("%s", e)))
		}
		grpcBranch.Die()
		// Отменяем контекст, если GRPC-сервер завершил работу
		serverCancel()
		wg.Done()
	}()

	addr := fmt.Sprintf("%s:%d", a.config.GrpcServer.Host, a.config.GrpcServer.Port)
	s := grpcApi.NewServer(addr, grpcBranch, a.logger)

	// Регистрация серверов gRPC-сервисов
	s.RegisterServers(authInteractor)
	a.grpcServer = s
	err := s.Run(ctx)
	if err != nil {
		a.logger.Error("grpc server is shutdown", zap.Error(err))
		return
	}
}
