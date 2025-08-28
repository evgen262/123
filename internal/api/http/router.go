package http

import (
	"fmt"
	"net/http"
	"net/url"

	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	timeUtils "git.mos.ru/buch-cloud/moscow-team-2.0/build/time-utils.git"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	presenter "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/presenter/auth"
	presenterBanners "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/presenter/banners"
	presenterEmployees "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/presenter/employees"
	presenterEmployeesSearch "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/presenter/employees-search"
	presenterNews "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/presenter/news"
	presenterPortalsV2 "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/presenter/portalsv2"
	presenterProxy "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/presenter/proxy"
	presenterSurveys "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/presenter/surveys"
	presenterUsers "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/presenter/users"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/auth"
)

type routerHandlers struct {
	surveysHandlers           SurveysHandlers
	authHandlers              AuthHandlers
	proxyHandlers             ProxyHandlers
	employeesSearchesHandlers EmployeesSearchHandlers
	usersHandlers             UsersHandlers
	redirectSessionHandlers   RedirectSessionHandlers
	filesHandlers             FilesHandlers
	employeesHandlers         EmployeesHandlers
	analyticsHandlers         AnalyticsHandlers
	newsAdminHandlers         NewsAdminHandlers
	newsHandlers              NewsHandlers
	portalsHandlers           PortalsV2Handlers
	bannersHandlers           BannersHandlers
}

type handlersFields struct {
	AllowOrigins    []string
	UploadPath      string
	AuthCallbackUrl *url.URL
	TemplatesPath   string
	ExternalHost    string
	AccessList      auth.AccessList
}

// SharedFields Общие поля для передачи их в пакет http.
//
//	Например, поля конфигурации приложения используемые в хэндлерах.
var SharedFields = new(handlersFields)

type router struct {
	engine            *gin.Engine
	middlewareOptions MiddlewareOptions
	handlers          *routerHandlers
	tu                timeUtils.TimeUtils
	environment       Environment
	logger            ditzap.Logger

	portalsInteractor   PortalsV2Interactor
	complexesInteractor ComplexesV2Interactor

	surveysInteractor        SurveysSurveysInteractor
	surveysAnswersInteractor SurveysAnswersInteractor
	surveysImagesInteractor  SurveysImagesInteractor

	authInteractor AuthInteractor

	proxyInteractor           ProxyInteractor
	redirectSessionInteractor RedirectSessionInteractor

	employeesSearchInteractor EmployeesSearchUseCases

	employeesInteractor EmployeesUseCases

	usersInteractor UsersInteractor

	filesInteractor FilesInteractor

	analyticsInteractor AnalyticsInteractor

	newsCategoryInteractor NewsCategoryInteractor
	commentsInteractor     NewsCommentsInteractor
	newsAdminInteractor    NewsAdminInteractor
	newsInteractor         NewsInteractor

	bannersInteractor BannersInteractor
}

func NewRouter(
	environment Environment,
	logger ditzap.Logger,

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
) *router {
	opts := MiddlewareOptions{opts: middlewareOptions}
	return &router{
		engine:            gin.New(),
		middlewareOptions: opts,
		environment:       environment,
		tu:                timeUtils.NewTimeUtils(),

		portalsInteractor:   portalsInteractor,
		complexesInteractor: complexesInteractor,

		surveysInteractor:        surveysInteractor,
		surveysAnswersInteractor: surveysAnswersInteractor,
		surveysImagesInteractor:  surveysImagesInteractor,

		authInteractor: authInteractor,

		proxyInteractor:           proxyInteractor,
		redirectSessionInteractor: redirectSessionInteractor,

		employeesSearchInteractor: employeesSearchInteractor,

		usersInteractor: usersInteractor,

		filesInteractor: filesInteractor,

		employeesInteractor: employeesInteractor,
		analyticsInteractor: analyticsInteractor,

		newsCategoryInteractor: newsCategoryInteractor,
		newsAdminInteractor:    newsAdminInteractor,
		newsInteractor:         newsInteractor,
		commentsInteractor:     newsCommentsInteractor,
		bannersInteractor:      bannersInteractor,
		logger:                 logger,
	}
}

func (r *router) Init() {
	InitHTTPMetrics()
	r.engine.Use(
		gin.Logger(),
		gin.CustomRecovery(r.recovery),
		cors.New(cors.Config{
			AllowAllOrigins: false,
			AllowOrigins:    SharedFields.AllowOrigins,
			AllowMethods: []string{
				http.MethodHead,
				http.MethodGet,
				http.MethodPost,
				http.MethodOptions,
				http.MethodPut,
				http.MethodPatch,
				http.MethodDelete,
			},
			AllowHeaders:     []string{"Content-Type", "Accept-Encoding", "Authorization", "Origin", "Session-Id"},
			AllowCredentials: true,
			ExposeHeaders:    []string{"Content-Length", "Session-Id", "X-Request-Id", "Authorization"},
		}),
	)
	r.initHandlers().
		registerRoutes()
}

func (r *router) recovery(c *gin.Context, recovered any) {
	r.logger.Error("http server recovery", zap.Error(fmt.Errorf("%s", recovered)))
	c.AbortWithStatus(http.StatusInternalServerError)
}

func (r *router) initHandlers() *router {
	// Portals presenters
	portalsPresenter := presenterPortalsV2.NewPortalsPresenter()
	complexesPresenter := presenterPortalsV2.NewComplexesPresenter()

	ph := NewPortalsV2Handlers(
		r.portalsInteractor,
		portalsPresenter,
		r.complexesInteractor,
		complexesPresenter,
		r.logger,
	)

	// Surveys presenters
	surveysPresenter := presenterSurveys.NewSurveysPresenter(r.logger)
	surveysAnswersPresenter := presenterSurveys.NewAnswersPresenter()
	surveysImagesPresenter := presenterSurveys.NewSurveyImagesPresenter()

	sh := NewSurveysHandlers(
		r.surveysInteractor,
		surveysPresenter,
		r.surveysAnswersInteractor,
		surveysAnswersPresenter,
		r.surveysImagesInteractor,
		surveysImagesPresenter,
		r.logger,
	)

	// Auth presenters
	authPresenter := presenter.NewAuthPresenter()
	ah := NewAuthHandlers(r.authInteractor, authPresenter, r.logger)

	proxyPresenter := presenterProxy.NewProxyPresenter()
	pxh := NewProxyHandlers(r.proxyInteractor, proxyPresenter)

	// EmployeesSearch presenters
	employeesSearchPresenter := presenterEmployeesSearch.NewEmployeesSearchPresenter()
	esh := NewEmployeesSearchHandlers(r.employeesSearchInteractor, employeesSearchPresenter, r.logger)

	// Users presenters
	usersPresenter := presenterUsers.NewUsersPresenter()
	uh := NewUsersHandlers(r.authInteractor, r.usersInteractor, authPresenter, usersPresenter)

	rhs := NewRedirectSessionHandlers(r.redirectSessionInteractor)

	fh := NewFilesHandlers(r.filesInteractor)

	employeesPresenter := presenterEmployees.NewEmployeesPresenter()

	eh := NewEmployeesHandlers(
		r.employeesInteractor,
		employeesPresenter,
		r.logger,
	)

	anh := NewAnalyticsHandlers(r.analyticsInteractor)

	// news presenters
	newsAdminPresenter := presenterNews.NewNewsAdminPresenter()
	newsCommentPresenter := presenterNews.NewCommentsPresenter()
	nah := NewNewsAdminHandlers(r.newsCategoryInteractor, r.newsAdminInteractor, newsAdminPresenter, r.logger)
	nh := NewNewsHandlers(r.newsCategoryInteractor, r.commentsInteractor, newsCommentPresenter, r.newsInteractor, newsAdminPresenter, r.logger)

	bannersPresenter := presenterBanners.NewBannersPresenter()
	bh := NewBannersHandlers(r.bannersInteractor, bannersPresenter, r.logger)

	r.handlers = &routerHandlers{
		surveysHandlers:           sh,
		authHandlers:              ah,
		proxyHandlers:             pxh,
		employeesSearchesHandlers: esh,
		usersHandlers:             uh,
		redirectSessionHandlers:   rhs,
		filesHandlers:             fh,
		employeesHandlers:         eh,
		analyticsHandlers:         anh,
		portalsHandlers:           ph,
		newsAdminHandlers:         nah,
		newsHandlers:              nh,
		bannersHandlers:           bh,
	}

	return r
}

//nolint:funlen
func (r *router) registerRoutes() {
	if r.handlers == nil {
		r.logger.Error("handlers not init")
		return
	}
	switch r.environment {
	case EnvironmentTest:
		fallthrough
	case EnvironmentDebug:
		fallthrough
	case EnvironmentDevelop:
		AppInstance = AppInstanceDevelop
	case EnvironmentProd:
	default:
		AppInstance = AppInstanceProd
	}
	// Выводим ошибку при недопустимых(отсутствующих) методах
	r.engine.NoMethod(NotImplementedHandler)
	r.engine.NoRoute(NotImplementedHandler)

	r.engine.Routes()
	r.engine.Use(
		NewRequestIDMiddleware(r.middlewareOptions),
		NewHeadersMiddleware(r.middlewareOptions),
	)

	// TODO: реализовать логику при взаимодействии с session service
	onlyAuthOpts := r.middlewareOptions
	onlyAuthOpts.opts = append(onlyAuthOpts.opts, &MiddlewareOption{
		Name:  onlyAuthOptKey,
		Value: true,
	})
	sessionMiddleware := NewAuthSessionMiddleware(r.authInteractor, r.tu, onlyAuthOpts)

	api := r.engine.Group("/", NewRequestIDMiddleware(r.middlewareOptions))
	{
		/**
		TODO Убрать данные роуты в рамках задачи https://oblako.mos.ru/jira/browse/TECH-515
		  после реализации фронтом задачи https://oblako.mos.ru/jira/browse/TECH-511
		*/
		api.GET("/oivs", r.handlers.portalsHandlers.getPortals)
		api.GET("/complexes", r.handlers.portalsHandlers.getComplexes)

		portalsModuleGroup := api.Group("/portals")
		{
			portalsV1Group := portalsModuleGroup.Group("/v1")
			{
				portalsV1Group.GET("/oivs", r.handlers.portalsHandlers.getPortals)
				portalsV1Group.GET("/complexes", r.handlers.portalsHandlers.getComplexes)
			}
		}

		/**
		TODO Убрать данные роуты в рамках задачи https://oblako.mos.ru/jira/browse/TECH-515
		  после реализации фронтом задачи https://oblako.mos.ru/jira/browse/TECH-511
		*/
		surveysGroupDeprecated := api.Group("/survey")
		{
			surveysAnswersGroup := surveysGroupDeprecated.Group("/answers")
			{
				surveysAnswersGroup.POST("/", r.handlers.surveysHandlers.addAnswers)
			}
			surveysImagesGroup := surveysGroupDeprecated.Group("/images")
			{
				surveysImagesGroup.GET("/:id", r.handlers.surveysHandlers.getImage)
			}
			surveysGroupDeprecated.GET("/:id", r.handlers.surveysHandlers.getSurvey)
		}

		surveysModuleGroup := api.Group("/surveys")
		{
			surveysV1Group := surveysModuleGroup.Group("/v1")
			{
				surveysAnswersGroup := surveysV1Group.Group("/answers")
				{
					surveysAnswersGroup.POST("/", r.handlers.surveysHandlers.addAnswers)
				}
				surveysImagesGroup := surveysV1Group.Group("/images")
				{
					surveysImagesGroup.GET("/:id", r.handlers.surveysHandlers.getImage)
				}
				surveysGroup := surveysV1Group.Group("/surveys")
				{
					surveysGroup.GET("/:id", r.handlers.surveysHandlers.getSurvey)
				}
			}
		}

		authModuleGroup := api.Group("/auth")
		{
			authV1Group := authModuleGroup.Group("/v1")
			{
				authV1Group.POST("/redirect", sessionMiddleware, r.handlers.redirectSessionHandlers.createSession)
				authV1Group.GET("/auth", r.handlers.authHandlers.auth)
				authV1Group.GET("/logout", sessionMiddleware, r.handlers.authHandlers.logout)
				authV1Group.GET("/refresh", r.handlers.authHandlers.refresh)
			}
		}

		proxyGroup := api.Group("/proxy")
		proxyGroup.Use(sessionMiddleware)
		{
			proxyV1Group := proxyGroup.Group("/v1")
			{
				proxyV1Group.GET("/banners/home-slider", r.handlers.proxyHandlers.listHomeBanners)
				proxyV1Group.GET("/events/list", r.handlers.proxyHandlers.listCalendarEvents)
				proxyV1Group.POST("/events/links", r.handlers.proxyHandlers.listCalendarEventsLinks)
			}
		}

		/**
		TODO Убрать данные роуты в рамках задачи https://oblako.mos.ru/jira/browse/TECH-515
		  после реализации фронтом задачи https://oblako.mos.ru/jira/browse/TECH-511
		*/
		employeesSearchGroupDeprecated := api.Group("/search")
		employeesSearchGroupDeprecated.Use(sessionMiddleware)
		{
			employeesV1Group := employeesSearchGroupDeprecated.Group("/v1")
			{
				employeesV1Group.POST("/employees/search", r.handlers.employeesSearchesHandlers.search)
				employeesV1Group.POST("/employees/filters", r.handlers.employeesSearchesHandlers.filters)
			}
		}

		usersModuleGroup := api.Group("/users")
		usersModuleGroup.Use(sessionMiddleware)
		{
			/**
			TODO Убрать данные роуты в рамках задачи https://oblako.mos.ru/jira/browse/TECH-515
			  после реализации фронтом задачи https://oblako.mos.ru/jira/browse/TECH-511
			*/
			usersModuleGroup.GET("/me", r.handlers.usersHandlers.getMe)
			usersModuleGroup.GET("/changeportal/:id", r.handlers.usersHandlers.changePortal)

			usersV1Group := usersModuleGroup.Group("/v1")
			{
				usersV1Group.GET("/profile", r.handlers.employeesHandlers.getProfile)
				usersV1Group.GET("/me", r.handlers.usersHandlers.getMe)
				usersV1Group.GET("/changeportal/:id", r.handlers.usersHandlers.changePortal)
			}
		}

		employeesModuleGroup := api.Group("/employees")
		employeesModuleGroup.Use(sessionMiddleware)
		{
			employeesV1Group := employeesModuleGroup.Group("/v1")
			{
				employeesGroup := employeesV1Group.Group("/employees")
				{
					employeesGroup.GET("/:id", r.handlers.employeesHandlers.getEmployee)
				}
				// Группа роутов для поиска сотрудников
				employeesSearchGroup := employeesV1Group.Group("/search")
				{
					employeesSearchGroup.POST("", r.handlers.employeesSearchesHandlers.search)
					employeesSearchGroup.POST("/filters", r.handlers.employeesSearchesHandlers.filters)
				}
			}
		}

		/**
		TODO Убрать данные роуты в рамках задачи https://oblako.mos.ru/jira/browse/TECH-515
		  после реализации фронтом задачи https://oblako.mos.ru/jira/browse/TECH-511
		*/
		filesGroupDeprecated := api.Group("/files")
		filesGroupDeprecated.Use(sessionMiddleware)
		{
			filesGroupDeprecated.GET("/:file_id", r.handlers.filesHandlers.get)
		}

		filesModuleGroup := api.Group("/files")
		filesModuleGroup.Use(sessionMiddleware)
		{
			filesV1Group := filesModuleGroup.Group("/v1")
			{
				filesGroup := filesV1Group.Group("/files")
				{
					filesGroup.GET("/:file_id", r.handlers.filesHandlers.get)
				}
			}
		}

		analyticsModuleGroup := api.Group("/analytics")
		analyticsModuleGroup.Use(sessionMiddleware)
		{
			analyticsV1Group := analyticsModuleGroup.Group("/v1")
			{
				analyticsV1Group.POST("/metrics", r.handlers.analyticsHandlers.addMetrics)
			}
		}

		newsModuleGroup := api.Group("/news")
		newsModuleGroup.Use(sessionMiddleware)
		{
			newsV1Group := newsModuleGroup.Group("/v1")
			{
				newsGroup := newsV1Group.Group("/news")
				{
					newsGroup.POST("/search", r.handlers.newsHandlers.searchNews)
					// Используем id!!! для wildcard, но в методе getNews подразумеваем что это slug новости!!!
					newsGroup.GET("/:id", r.handlers.newsHandlers.getNews)
					newsGroup.POST("/:id/comments", r.handlers.newsHandlers.createComment)
					newsGroup.GET("/:id/comments", r.handlers.newsHandlers.listComments)
				}
			}
		}

		// Main роуты main/*
		mainModuleGroup := api.Group("/main")
		mainModuleGroup.Use(sessionMiddleware)
		{
			mainV1Group := mainModuleGroup.Group("/v1")
			{
				mainV1Group.GET("/banners", r.handlers.bannersHandlers.list)
			}
		}
	}

	// Административные роуты /admin/*
	r.registerAdminRoutes(api.Group("/admin"))
}

// Административные роуты
func (r *router) registerAdminRoutes(adminGroup *gin.RouterGroup) {
	authMiddleware := NewAuthSessionMiddleware(r.authInteractor, r.tu, r.middlewareOptions)
	adminGroup.Use(authMiddleware)

	newsModuleGroup := adminGroup.Group("/news")
	{
		newsV1Group := newsModuleGroup.Group("/v1")
		{
			categoryGroup := newsV1Group.Group("/category")
			{
				categoryGroup.POST("", r.handlers.newsAdminHandlers.createCategory)
				categoryGroup.PUT("/:id", r.handlers.newsAdminHandlers.updateCategory)
				categoryGroup.DELETE("/:id", r.handlers.newsAdminHandlers.deleteCategory)
				categoryGroup.POST("/search", r.handlers.newsAdminHandlers.searchCategory)
				categoryGroup.GET("/:id", r.handlers.newsAdminHandlers.getCategory)
			}
			newsGroup := newsV1Group.Group("/news")
			{
				newsGroup.POST("", r.handlers.newsAdminHandlers.createNews)
				newsGroup.GET("/:id", r.handlers.newsAdminHandlers.getNews)
				newsGroup.PUT("/:id", r.handlers.newsAdminHandlers.updateNews)
				newsGroup.DELETE("/:id", r.handlers.newsAdminHandlers.deleteNews)
				newsGroup.POST("/:id/status", r.handlers.newsAdminHandlers.setStatusNews)
				newsGroup.PATCH("/:id/flags", r.handlers.newsAdminHandlers.setFlagsNews)
				searchGroup := newsGroup.Group("/search")
				{
					searchGroup.POST("", r.handlers.newsAdminHandlers.searchNews)
				}
			}
		}
	}

	bannersModuleGroup := adminGroup.Group("/banners")
	{
		bannersV1Group := bannersModuleGroup.Group("/v1")
		{
			bannersGroup := bannersV1Group.Group("/banners")
			{
				bannersGroup.PUT("", r.handlers.bannersHandlers.set)
			}
		}
	}
}
