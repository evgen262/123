package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/tree-alive.git"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//go:generate mockgen -source=server.go -destination=./server_mock.go -package=http

// RequestTimeOut Тайм-аут запросов
const RequestTimeOut = 30 * time.Second

type Environment uint8

const (
	EnvironmentTest Environment = iota
	EnvironmentDebug
	EnvironmentDevelop
	EnvironmentProd
)

type Server interface {
	Run(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

func NewServer(
	addr string,
	logger ditzap.Logger,
	environment Environment,
	treeBranch tree.Branch,
) *server {
	switch environment {
	case EnvironmentDebug:
		fallthrough
	case EnvironmentDevelop:
		gin.SetMode(gin.DebugMode)
	case EnvironmentProd:
		gin.SetMode(gin.ReleaseMode)
	case EnvironmentTest:
		fallthrough
	default:
		gin.SetMode(gin.TestMode)
	}
	s := &server{
		treeBranch: treeBranch,
		enviroment: environment,
		logger:     logger,
	}

	httpServer := &http.Server{
		Addr:              addr,
		ReadHeaderTimeout: RequestTimeOut,
	}
	s.server = httpServer

	return s
}

type server struct {
	server     HTTPSrv
	treeBranch tree.Branch
	enviroment Environment
	logger     ditzap.Logger
}

func (s *server) Server() *http.Server {
	if srv, ok := s.server.(*http.Server); ok {
		return srv
	}
	return nil
}

func (s *server) InitRouter(
	uploadPath string,
	portalsInteractor PortalsV2Interactor,
	complexesInteractor ComplexesV2Interactor,
	surveysInteractor SurveysSurveysInteractor,
	surveysAnswersInteractor SurveysAnswersInteractor,
	surveysImagesInteractor SurveysImagesInteractor,
	authInteractor AuthInteractor,
	proxyInteractor ProxyInteractor,
	redirectSessionInteractor RedirectSessionInteractor,
	employeesSearchInteractor EmployeesSearchUseCases,
	usersInteractor UsersInteractor,
	filesInteractor FilesInteractor,
	employeesInteractor EmployeesUseCases,
	analyticsInteractor AnalyticsInteractor,
	newsCategoryInteractor NewsCategoryInteractor,
	newsAdminInteractor NewsAdminInteractor,
	newsInteractor NewsInteractor,
	newsCommentsInteractor NewsCommentsInteractor,
	bannersInteractor BannersInteractor,
	middlewareOptions ...*MiddlewareOption,
) *server {
	if s == nil {
		return nil
	}

	SharedFields.UploadPath = uploadPath

	r := NewRouter(
		s.enviroment,
		s.logger,
		portalsInteractor,
		complexesInteractor,
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
		newsInteractor,
		newsCommentsInteractor,
		bannersInteractor,
		middlewareOptions...,
	)
	r.Init()
	s.Server().Handler = r.engine
	return s
}

func (s *server) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		err := s.Shutdown(ctx)
		if err != nil {
			s.logger.Error("can't shutdown http-server", zap.Error(err))
			return
		}
	}()

	s.treeBranch.Ready()
	s.logger.Info("http-server started", zap.String("server-host", s.Server().Addr))
	return s.server.ListenAndServe() //nolint:wrapcheck
}

func (s *server) Shutdown(ctx context.Context) error {
	err := s.server.Shutdown(ctx)
	s.treeBranch.Die()
	if err != nil {
		return fmt.Errorf("http server shutdown error: %w", err)
	}
	return nil
}
