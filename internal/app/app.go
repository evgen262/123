package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	metricsv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/analytics/metrics/v1"
	authv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/authfacade/auth/v1"
	redirectsessionv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/authfacade/redirectsession/v1"
	bannersv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/banners/banners/v1"
	employeev1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/employees/employee/v1"
	searchv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/employeessearch/search/v1"
	filev1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/fileservice/file/v1"
	categoryv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/newsfacade/category/v1"
	commentsv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/newsfacade/comment/v1"
	newsv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/newsfacade/news/v1"
	portalsv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/portals/portals/v1"
	organizationsv2 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/portalsv2/organizations/v1"
	portalsv2 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/portalsv2/portals/v1"
	bannerv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/proxyfacade/banner/v1"
	eventv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/proxyfacade/event/v1"
	answerv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/surveys/answer/v1"
	surveysImagesv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/surveys/image/v1"
	surveyv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/surveys/survey/v1"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	timeUtils "git.mos.ru/buch-cloud/moscow-team-2.0/build/time-utils.git"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/tree-alive.git"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	redisPkg "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	complexesfacadev1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/portalsfacade/complexes/v1"
	portalsfacadev1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/portalsfacade/portals/v1"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/cmd/web-api/config"
	httpApi "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/service"
	grpcClient "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/client/grpc"
	mapperAnalytics "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/client/grpc/analytics"
	mapperAuth "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/client/grpc/auth"
	mapperBanners "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/client/grpc/banners"
	mapperEmployees "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/client/grpc/employees"
	mapperEmployeesSearch "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/client/grpc/employees-search"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/client/grpc/files"
	mapperNews "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/client/grpc/news"
	mapperPortals "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/client/grpc/portals"
	mapperPortalsV2 "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/client/grpc/portalsv2"
	mapperProxy "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/client/grpc/proxy"
	mapperSurveys "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/client/grpc/surveys"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/auth"
	repositoryAnalytics "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/repositories/analytics"
	repositoriesAuth "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/repositories/auth"
	repositoryBanners "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/repositories/banners"
	repositoryEmployees "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/repositories/employees"
	repositoriesEmployeesSearch "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/repositories/employees-search"
	repositoryFiles "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/repositories/files"
	repositoryNews "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/repositories/news"
	repositoriesPortal "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/repositories/portal"
	repositoryPortalsv2 "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/repositories/portalsv2"
	repositoryProxy "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/repositories/proxy"
	repositoriesSurvey "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/repositories/survey"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/usecase/analytics"
	usecaseAuth "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/usecase/auth"
	usecaseBanners "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/usecase/banners"
	usecaseEmployees "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/usecase/employees"
	usecaseEmployeesSearch "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/usecase/employees-search"
	usesaceFIles "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/usecase/files"
	usecaseNews "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/usecase/news"
	usecasePortalsv2 "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/usecase/portalsv2"
	usecaseProxy "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/usecase/proxy"
	usecaseSurveys "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/usecase/surveys"
	usecaseUsers "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/usecase/users"
)

type Environment string

const (
	EnvironmentDevelop Environment = "develop"
	EnvironmentProd    Environment = "prod"
)
const (
	DialTimeout = 10 * time.Second
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
	config            *config.Config
	serviceHTTPServer Server
	httpServer        Server
	redisClient       redisPkg.UniversalClient
	tree              tree.Tree
	stop              context.CancelFunc
	logger            ditzap.Logger
}

func NewApp(cfg *config.Config, appInfo *AppInfo, logger ditzap.Logger) *app {
	ApplicationInfo = appInfo
	return &app{
		config: cfg,
		logger: logger,
	}
}

//nolint:funlen
func (a *app) Run(ctx context.Context) {
	appCtx, cancelApp := context.WithCancel(ctx)
	a.stop = cancelApp
	// AfterFunc-функция при завершении контекста для graceful shutdown
	context.AfterFunc(appCtx, func() {
		defer func() {
			if e := recover(); e != nil {
				a.logger.Error("context after-func", zap.Error(fmt.Errorf("%s", e)))
			}

			if gErr := a.GracefulShutdown(appCtx); gErr != nil {
				a.logger.Error("can't graceful shutdown", zap.Error(gErr))
			}
		}()
	})

	environment := EnvironmentDevelop
	if a.config.Environment == "prod" {
		environment = EnvironmentProd
	}
	ApplicationInfo.Instance = environment

	a.tree = tree.NewTree()
	a.tree.Alive()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	serviceHttpBranch := a.tree.GrowBranch("service-http-server")
	go a.startServiceServer(appCtx, serviceHttpBranch, wg)

	tu := timeUtils.NewTimeUtils()

	var accessList auth.AccessList
	if strings.TrimSpace(a.config.AccessListFile) != "" {
		var alErr error
		accessList, alErr = a.readAccessList()
		if alErr != nil {
			return
		}
		httpApi.SharedFields.AccessList = accessList
	}

	// Подключение к зависимым gRPC-сервисам
	portalsConn,
		portalsv2Conn,
		surveysConn,
		authConn,
		proxyConn,
		employeesSearchConn,
		filesConn,
		employeesConn,
		analyticsConn,
		portalsFacadeConn,
		newsConn,
		bannersConn,
		err := a.connectToProviders(appCtx, a.config.Endpoints)
	if err != nil {
		a.logger.Error("can't connection to providers", zap.Error(err))
		return
	}

	// Клиенты
	portalsAPIClient := portalsv1.NewPortalsAPIClient(portalsConn)
	portalsv2APIClient := portalsv2.NewPortalsAPIClient(portalsv2Conn)
	organizationsv2APIClient := organizationsv2.NewOrganizationsAPIClient(portalsv2Conn)
	portalsfacadeAPIClient := portalsfacadev1.NewPortalsAPIClient(portalsFacadeConn)
	complexesfacadeAPIClient := complexesfacadev1.NewComplexesAPIClient(portalsFacadeConn)

	filesAPIClient := filev1.NewFileAPIClient(filesConn)

	surveysAPIClient := surveyv1.NewSurveyAPIClient(surveysConn)
	surveysAnswersAPIClient := answerv1.NewAnswerAPIClient(surveysConn)
	surveysImagesAPIClient := surveysImagesv1.NewImageAPIClient(surveysConn)

	authAPIClient := authv1.NewAuthAPIClient(authConn)

	bannerAPIClient := bannerv1.NewBannerAPIClient(proxyConn)
	eventsAPIClient := eventv1.NewEventAPIClient(proxyConn)

	redirectSessionAPIClient := redirectsessionv1.NewRedirectSessionAPIClient(authConn)

	employeesSearchSearchAPIClient := searchv1.NewSearchAPIClient(employeesSearchConn)

	employeesAPIClient := employeev1.NewEmployeesAPIClient(employeesConn)

	metricsAPIClient := metricsv1.NewMetricsAPIClient(analyticsConn)

	newsCategoryAPIClient := categoryv1.NewCategoryAPIClient(newsConn)
	newsAPIClient := newsv1.NewNewsAPIClient(newsConn)
	commentsAPIClient := commentsv1.NewCommentAPIClient(newsConn)

	bannersAPIClient := bannersv1.NewBannersAPIClient(bannersConn)

	// Мапперы
	sharedMapper := grpcClient.NewSharedMapper(tu)
	portalsMapper := mapperPortals.NewPortalsMapper(tu)

	surveysMapper := mapperSurveys.NewSurveyMapper(tu)
	surveysAnswersMapper := mapperSurveys.NewAnswerMapper()
	surveysImagesMapper := mapperSurveys.NewImageMapper()

	authMapper := mapperAuth.NewAuthMapper(tu)
	redirectSessionMapper := mapperAuth.NewRedirectSessionMapper()

	proxyMapper := mapperProxy.NewProxyMapper()

	employeesMapper := mapperEmployees.NewMapperEmployees(tu)
	employeesSearchMapper := mapperEmployeesSearch.NewEmployeesSearchMapper()

	newsMapper := mapperNews.NewNewsMapper(sharedMapper)
	commentsMapper := mapperNews.NewCommentsMapper(sharedMapper)

	filesMapper := files.NewFileMapper()
	visitorMapper := files.NewVisitorMapper()

	bannersMapper := mapperBanners.NewBannersMapper(tu, sharedMapper)

	// Репозитории
	portalsPortalRepository := repositoriesPortal.NewPortalsRepository(portalsAPIClient, portalsMapper)

	surveysRepository := repositoriesSurvey.NewSurveyRepository(surveysAPIClient, surveysMapper)
	surveysAnswersRepository := repositoriesSurvey.NewAnswerRepository(surveysAnswersAPIClient, surveysAnswersMapper)
	surveysImagesRepository := repositoriesSurvey.NewImageRepository(surveysImagesAPIClient, surveysImagesMapper)

	employeesRepository := repositoryEmployees.NewEmployeesRepository(employeesAPIClient, employeesMapper, tu, a.logger)
	employeesSearchRepository := repositoriesEmployeesSearch.NewEmployeesSearchRepository(employeesSearchSearchAPIClient, employeesSearchMapper, tu, a.logger)

	callbackURL, err := url.Parse(a.config.WebAuthURL)
	if err != nil {
		a.logger.Error("can't parse WEB_AUTH_URL parameter", zap.String("url", a.config.WebAuthURL), zap.Error(err))
		return
	}
	authRepository := repositoriesAuth.NewAuthRepository(authAPIClient, authMapper, *callbackURL, a.config.AppName, a.config.TTL.AccessToken, a.config.TTL.RefreshToken, tu, a.logger)
	redirectSessionRepository := repositoriesAuth.NewRedirectSessionRepository(redirectSessionAPIClient, redirectSessionMapper, a.logger)

	proxyRepository := repositoryProxy.NewProxyRepository(bannerAPIClient, eventsAPIClient, proxyMapper, a.logger)

	filesRepository := repositoryFiles.NewFileRepository(filesAPIClient, filesMapper, visitorMapper)

	analyticsRepository := repositoryAnalytics.NewAnalyticsRepository(metricsAPIClient, mapperAnalytics.NewMetricsMapper())

	// TODO: Переименовать, когшда будет убираться portals v1
	portalsFacadePortalRepository := repositoryPortalsv2.NewPortalsRepository(portalsfacadeAPIClient, mapperPortalsV2.NewPortalsMapper(tu))
	portalsFacadeComplexesRepository := repositoryPortalsv2.NewComplexesRepository(complexesfacadeAPIClient, mapperPortalsV2.NewComplexesMapper(tu))

	newsCategoryRepository := repositoryNews.NewCategoryRepository(newsCategoryAPIClient, sharedMapper)
	newsRepository := repositoryNews.NewNewsRepository(newsAPIClient, portalsv2APIClient, organizationsv2APIClient, newsMapper, sharedMapper, a.logger)
	bannersRepository := repositoryBanners.NewBannersRepository(bannersAPIClient, bannersMapper, a.logger)
	commentsRepository := repositoryNews.NewCommentsRepository(commentsAPIClient, commentsMapper, newsMapper)

	// Интеракторы
	portalsV2Interactor := usecasePortalsv2.NewPortalsUseCase(portalsFacadePortalRepository)
	complexesV2Interactor := usecasePortalsv2.NewComplexesUseCase(portalsFacadeComplexesRepository)

	surveysInteractor := usecaseSurveys.NewSurveysUseCase(surveysRepository)
	surveysAnswersInteractor := usecaseSurveys.NewAnswersUseCase(surveysAnswersRepository)
	surveysImagesInteractor := usecaseSurveys.NewImagesUseCase(surveysImagesRepository)

	authInteractor := usecaseAuth.NewAuthUseCase(authRepository, a.logger, accessList)
	redirectSessionInteractor := usecaseAuth.NewRedirectSessionInteractor(redirectSessionRepository, a.config.WebAuthRedirectURI)

	proxyInteractor := usecaseProxy.NewProxyInteractor(proxyRepository)

	employeesInteractor := usecaseEmployees.NewEmployeesInteractor(employeesRepository, portalsPortalRepository, a.logger)
	employeesSearchInteractor := usecaseEmployeesSearch.NewEmployeesSearchInteractor(employeesSearchRepository, tu)

	usersInteractor := usecaseUsers.NewUsersInteractor(employeesRepository, a.logger)

	filesInteractor := usesaceFIles.NewFileUsecase(filesRepository, a.logger)

	analyticsInteractor := analytics.NewAnalyticsInteractor(analyticsRepository, a.logger)

	newsCategoryInteractor := usecaseNews.NewCategoryInteractor(newsCategoryRepository, employeesRepository, a.logger)
	newsAdminInteractor := usecaseNews.NewNewsAdminInteractor(newsRepository, employeesRepository, a.logger)

	newsInteractor := usecaseNews.NewNewsInteractor(newsRepository, employeesRepository, a.logger)
	newsCommentsInteractor := usecaseNews.NewCommentInteractor(commentsRepository, employeesRepository, a.logger)

	bannersInteractor := usecaseBanners.NewBannersInteractor(bannersRepository, a.logger)

	// Старт HTTP-сервера
	wg.Add(1)
	httpBranch := a.tree.GrowBranch("http-server")
	go a.startHTTPServer(
		appCtx,
		httpBranch,
		wg,
		portalsV2Interactor,
		complexesV2Interactor,
		surveysInteractor,
		surveysAnswersInteractor,
		surveysImagesInteractor,
		authInteractor,
		proxyInteractor,
		redirectSessionInteractor,
		employeesSearchInteractor,
		usersInteractor,
		filesInteractor,
		employeesInteractor,
		analyticsInteractor,
		newsCategoryInteractor,
		newsAdminInteractor,
		newsCommentsInteractor,
		newsInteractor,
		bannersInteractor,
	)

	a.tree.Ready()
	wg.Wait()
}

// GracefulShutdown graceful shutdown приложения
func (a *app) GracefulShutdown(c context.Context) (err error) {
	ctx := context.WithoutCancel(c)
	if a.tree != nil {
		a.tree.Die()
	}
	if a.httpServer != nil {
		httpErr := a.httpServer.Shutdown(ctx)
		if httpErr != nil {
			err = errors.Join(err, fmt.Errorf("can't shutdown http-server: %w", httpErr))
		}
	}
	if a.serviceHTTPServer != nil {
		serviceErr := a.serviceHTTPServer.Shutdown(ctx)
		if serviceErr != nil {
			err = errors.Join(err, fmt.Errorf("can't shutdown service http-server: %w", serviceErr))
		}
	}
	if a.redisClient != nil {
		status := a.redisClient.Shutdown(ctx)
		if status.Err() != nil {
			err = errors.Join(err, fmt.Errorf("can't shutdown service redis: %w", status.Err()))
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

func (a *app) connectToProviders(
	ctx context.Context,
	endpoints *config.Endpoints,
) (
	portalsConn *grpc.ClientConn,
	portalsV2Conn *grpc.ClientConn,
	surveysConn *grpc.ClientConn,
	authConn *grpc.ClientConn,
	proxyConn *grpc.ClientConn,
	employeesSearchConn *grpc.ClientConn,
	filesConn *grpc.ClientConn,
	employeesConn *grpc.ClientConn,
	analyticsConn *grpc.ClientConn,
	portalsFacadeConn *grpc.ClientConn,
	newsConn *grpc.ClientConn,
	bannersConn *grpc.ClientConn,
	err error,
) {
	var (
		dialErr error
	)
	grpc_prometheus.EnableClientHandlingTimeHistogram()
	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(grpc_prometheus.UnaryClientInterceptor)),
		grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(grpc_prometheus.StreamClientInterceptor)),
	}

	// Инициализация подключения к сервису portal
	portalsConn, dialErr = a.dial(ctx, endpoints.PortalsEndpoint, false, dialOpts...)
	if dialErr != nil {
		err = fmt.Errorf("can't connect to posrtals service: %w", dialErr)
		return
	}

	// Инициализация подключения к сервису portal v2
	portalsV2Conn, dialErr = a.dial(ctx, endpoints.PortalsV2Endpoint, false, dialOpts...)
	if dialErr != nil {
		err = fmt.Errorf("can't connect to posrtals v2 service: %w", dialErr)
		return
	}

	// Инициализация подключения к сервису survey
	surveysConn, dialErr = a.dial(ctx, endpoints.SurveysEndpoint, false, dialOpts...)
	if dialErr != nil {
		err = fmt.Errorf("can't connect to surveys service: %w", dialErr)
		return
	}

	// Инициализация подключения к сервису auth-facade
	authConn, dialErr = a.dial(ctx, endpoints.AuthFacadeEndpoint, false, dialOpts...)
	if dialErr != nil {
		err = fmt.Errorf("can't connect to auth-facade service: %w", dialErr)
		return
	}

	proxyConn, dialErr = a.dial(ctx, endpoints.ProxyFacadeEndpoint, false, dialOpts...)
	if dialErr != nil {
		err = fmt.Errorf("can't connect to banners service: %w", dialErr)
		return
	}

	employeesSearchConn, dialErr = a.dial(ctx, endpoints.EmployeesSearchEndpoint, false, dialOpts...)
	if dialErr != nil {
		err = fmt.Errorf("can't connect to employees-search service: %w", dialErr)
		return
	}

	filesConn, dialErr = a.dial(ctx, endpoints.FilesEndpoint, false, dialOpts...)
	if dialErr != nil {
		err = fmt.Errorf("can't connect to files service: %w", dialErr)
		return
	}

	employeesConn, dialErr = a.dial(ctx, endpoints.EmployeesEndpoint, false, dialOpts...)
	if dialErr != nil {
		err = fmt.Errorf("can't connect to employees service: %w", dialErr)
		return
	}

	analyticsConn, dialErr = a.dial(ctx, endpoints.AnalyticsEndpoint, false, dialOpts...)
	if dialErr != nil {
		err = fmt.Errorf("can't connect to analytics service: %w", dialErr)
		return
	}

	portalsFacadeConn, dialErr = a.dial(ctx, endpoints.PortalsFacadeEndpoint, false, dialOpts...)
	if dialErr != nil {
		err = fmt.Errorf("can't connect to portals-facade service: %w", dialErr)
		return
	}

	newsConn, dialErr = a.dial(ctx, endpoints.NewsFacadeEndpoint, false, dialOpts...)
	if dialErr != nil {
		err = fmt.Errorf("can't connect to news service: %w", dialErr)
		return
	}

	bannersConn, dialErr = a.dial(ctx, endpoints.BannersEndpoint, false, dialOpts...)
	if dialErr != nil {
		err = fmt.Errorf("can't connect to banners service: %w", dialErr)
		return
	}

	return
}

func (a *app) initMetrics() {
	// Добавление префикса к метрикам
	appNamePrefix := strings.ReplaceAll(ApplicationInfo.Name, "-", "_")
	metricRegisterer := prometheus.WrapRegistererWithPrefix(appNamePrefix+"_", prometheus.DefaultRegisterer)
	prometheus.DefaultRegisterer = metricRegisterer
	if gather, ok := metricRegisterer.(prometheus.Gatherer); ok {
		prometheus.DefaultGatherer = gather
	}
}

func (a *app) startServiceServer(c context.Context, serviceHttpBranch tree.Branch, wg *sync.WaitGroup) {
	ctx, serverCancel := context.WithCancel(c)
	defer func() {
		if e := recover(); e != nil {
			a.logger.Panic("service http start panic", zap.Error(fmt.Errorf("%s", e)))
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
	err := a.serviceHTTPServer.Run(
		context.WithValue(ctx, &service.ContextKeyServiceAddr, entity.MakeContextStringValue(a.config.ServiceHTTPHost)), //nolint:staticcheck
	)
	// Отменяем контекст, если HTTP-сервер завершил работу
	if err != nil {
		a.logger.Error("service http server is shutdown", zap.Error(err))
	}
}

//nolint:funlen
func (a *app) startHTTPServer(
	c context.Context,
	httpBranch tree.Branch,
	wg *sync.WaitGroup,
	portalsV2Interactor httpApi.PortalsV2Interactor,
	complexesV2Interactor httpApi.ComplexesV2Interactor,
	surveysInteractor httpApi.SurveysSurveysInteractor,
	surveysAnswersInteractor httpApi.SurveysAnswersInteractor,
	surveysImagesInteractor httpApi.SurveysImagesInteractor,
	authInteractor httpApi.AuthInteractor,
	proxyInteractor httpApi.ProxyInteractor,
	redirectSessionInteractor httpApi.RedirectSessionInteractor,
	employeesSearchesInteractor httpApi.EmployeesSearchUseCases,
	usersInteractor httpApi.UsersInteractor,
	filesInteractor httpApi.FilesInteractor,
	employeesInteractor httpApi.EmployeesUseCases,
	analyticsInteractor httpApi.AnalyticsInteractor,
	newsCategoryInteractor httpApi.NewsCategoryInteractor,
	newsAdminInteractor httpApi.NewsAdminInteractor,
	newsCommentsInteractor httpApi.NewsCommentsInteractor,
	newsInteractor httpApi.NewsInteractor,
	bannersInteractor httpApi.BannersInteractor,
) {
	ctx, serverCancel := context.WithCancel(c)
	defer func() {
		if e := recover(); e != nil {
			a.logger.Error("http start panic", zap.Error(fmt.Errorf("%s", e)))
		}
		httpBranch.Die()
		// Отменяем контекст, если HTTP-сервер завершил работу
		serverCancel()
		wg.Done()
	}()

	// Заполняем переменные для http-роутов
	httpApi.SharedFields.AllowOrigins = strings.Split(a.config.HttpServer.AllowOrigins, ",")
	httpApi.SharedFields.ExternalHost = a.config.HttpServer.ExternalHost

	addr := fmt.Sprintf("%s:%d", a.config.HttpServer.Host, a.config.HttpServer.Port)
	var servEnvironment httpApi.Environment
	switch ApplicationInfo.Instance {
	case EnvironmentDevelop:
		servEnvironment = httpApi.EnvironmentDevelop
	case EnvironmentProd:
		servEnvironment = httpApi.EnvironmentProd
	}
	if a.config.DevMode {
		servEnvironment = httpApi.EnvironmentDebug
	}
	middlewareOptions := []*httpApi.MiddlewareOption{
		{
			Name:  "apiKey",
			Value: a.config.ApiKey,
		},
	}
	a.httpServer = httpApi.
		NewServer(addr, a.logger, servEnvironment, httpBranch).
		InitRouter(
			a.config.Path.UploadPath,
			portalsV2Interactor,
			complexesV2Interactor,
			surveysInteractor,
			surveysAnswersInteractor,
			surveysImagesInteractor,
			authInteractor,
			proxyInteractor,
			redirectSessionInteractor,
			employeesSearchesInteractor,
			usersInteractor,
			filesInteractor,
			employeesInteractor,
			analyticsInteractor,
			newsCategoryInteractor,
			newsAdminInteractor,
			newsInteractor,
			newsCommentsInteractor,
			bannersInteractor,
			middlewareOptions...,
		)
	if a.httpServer == nil {
		a.logger.Error("can't create http server")
		return
	}

	if err := a.httpServer.Run(ctx); err != nil {
		a.logger.Error("http server is shutdown", zap.Error(err))
		return
	}
}

func (a *app) readAccessList() (auth.AccessList, error) {
	accessList := make(auth.AccessList, 0)
	file, err := os.Open(a.config.AccessListFile)
	if err != nil {
		a.logger.Error("can't open access list file", zap.Error(err))
		return nil, fmt.Errorf("can't open access list file")
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		a.logger.Error("can't read access list file", zap.Error(err))
		return nil, fmt.Errorf("can't read access list file")
	}
	err = json.Unmarshal(data, &accessList)
	if err != nil && !errors.Is(err, io.EOF) {
		a.logger.Error("can't unmarshal access list file", zap.Error(err))
		return nil, fmt.Errorf("can't unmarshal access list file")
	}
	return accessList, nil
}
