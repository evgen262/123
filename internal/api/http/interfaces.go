package http

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	viewBanners "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/banners"
	viewNews "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/news"
	dtoBanners "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/dto/banners"
	dtoNews "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/dto/news"
	entityBanners "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/banners"
	entityNews "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/news"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view"
	viewAuth "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/auth"
	viewBanner "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/banner"
	viewEmployees "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/employees"
	viewEmployeesSearch "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/employees-search"
	viewEvents "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/events"
	viewPortals "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/portals"
	viewPortalsV2 "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/portalsv2"
	viewSurveys "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/surveys"
	viewUsers "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/users"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/analytics"
	entityAuth "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/auth"
	entitySession "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/auth"
	entityBanner "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/banner"
	entityEmployee "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/employee"
	entityEmployeesSearch "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/employees-search"
	entityEvent "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/event"
	entityFile "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/file"
	entityPortal "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portal"
	entityPortalsV2 "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portalv2"
	entitySurvey "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/survey"
	entityUser "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/user"
)

//go:generate mockgen -source=interfaces.go -destination=./http_mock.go -package=http

type HTTPSrv interface {
	Shutdown(ctx context.Context) error
	ListenAndServe() error
}

/**
Хендлеры
*/

/**
Модуль порталов
*/

// PortalsHandlers хендлеры модуля portal
type PortalsHandlers interface {
	PortalsFeatureHandlers
	PortalsImageHandlers
	PortalsPortalHandlers
	PortalsOrganizationHandlers
	PortalsQuestionHandlers
}

// PortalsFeatureHandlers ручки по функционалу
type PortalsFeatureHandlers interface {
	getFeature(c *gin.Context)
	getFeatures(c *gin.Context)
	getFeatureByAdmin(c *gin.Context)
	getFeaturesByAdmin(c *gin.Context)
	addFeature(c *gin.Context)
	addFeatures(c *gin.Context)
	updateFeature(c *gin.Context)
	deleteFeature(c *gin.Context)
}

// PortalsImageHandlers ручки по изображениям
type PortalsImageHandlers interface {
	getImages(c *gin.Context)
	getImageByAdmin(c *gin.Context)
	getImagesByAdmin(c *gin.Context)
	addImage(c *gin.Context)
	getRawImage(c *gin.Context)
	deleteImage(c *gin.Context)
}

// PortalsPortalHandlers ручки по порталам
type PortalsPortalHandlers interface {
	filterPortals(c *gin.Context)
	getPortals(c *gin.Context)
	getPortal(c *gin.Context)
	filterPortalsByAdmin(c *gin.Context)
	getPortalByAdmin(c *gin.Context)
	getPortalsByAdmin(c *gin.Context)
	addPortal(c *gin.Context)
	addPortals(c *gin.Context)
	updatePortal(c *gin.Context)
	deletePortal(c *gin.Context)
}

// PortalsQuestionHandlers ручки по вопросам
type PortalsQuestionHandlers interface {
	getQuestions(c *gin.Context)
	getQuestion(c *gin.Context)
	getQuestionByAdmin(c *gin.Context)
	getQuestionsByAdmin(c *gin.Context)
	addQuestion(c *gin.Context)
	addQuestions(c *gin.Context)
	updateQuestion(c *gin.Context)
	deleteQuestion(c *gin.Context)
}

type PortalsOrganizationHandlers interface {
	filterOrganizations(c *gin.Context)
	linkOrganizations(c *gin.Context)
	unlinkOrganizations(c *gin.Context)
}

/**
Аутентификация
*/

// AuthHandlers ручки для аутентификации
type AuthHandlers interface {
	auth(c *gin.Context)
	logout(c *gin.Context)
	refresh(c *gin.Context)
}

/**
Модуль proxy-facade
*/

// ProxyHandlers ручки для proxy-facade
type ProxyHandlers interface {
	listHomeBanners(c *gin.Context)
	listCalendarEvents(c *gin.Context)
	listCalendarEventsLinks(c *gin.Context)
}

type ProxyInteractor interface {
	ListHomeBanners(ctx context.Context, sessionID string, portalURL string) (*entityBanner.BannersList, error)
	ListCalendarEvents(ctx context.Context, req entityEvent.CalendarEventRequest) (*entityEvent.CalendarEventsList, error)
	ListCalendarEventsLinks(ctx context.Context, req entityEvent.CalendarEventLinksRequest) ([]*entityEvent.CalendarEventLink, error)
}

type ProxyPresenter interface {
	BannersListToView(bannersList *entityBanner.BannersList) *viewBanner.BannersList
	EventsListToView(eventsList *entityEvent.CalendarEventsList) *viewEvents.CalendarEventsList
	EventsLinksToView(eventsLinks []*entityEvent.CalendarEventLink) []*viewEvents.CalendarEventLink
}

/**
Модуль редиректов
*/

type RedirectSessionHandlers interface {
	createSession(c *gin.Context)
}

type RedirectSessionInteractor interface {
	CreateSession(ctx context.Context, userInfo *entitySession.RedirectSessionUserInfo) (string, error)
}

/**
Модуль опросов
*/

// SurveysHandlers хендлеры модуля survey
type SurveysHandlers interface {
	SurveysAnswerHandlers
	SurveysImageHandlers
	SurveysSurveyHandlers
}

// SurveysSurveyHandlers ручки по опросам
type SurveysSurveyHandlers interface {
	getSurvey(c *gin.Context)
}

// SurveysAnswerHandlers ручки по ответам на опрос
type SurveysAnswerHandlers interface {
	addAnswers(c *gin.Context)
}

// SurveysImageHandlers ручки по изображениям опроса
type SurveysImageHandlers interface {
	getImage(c *gin.Context)
}

/**
Модуль поиска сотрудников
*/

type EmployeesSearchHandlers interface {
	search(c *gin.Context)
	filters(c *gin.Context)
}

/**
Модуль пользователей
*/

// UsersHandlers хендлеры модуля пользователей
type UsersHandlers interface {
	getMe(c *gin.Context)
	changePortal(c *gin.Context)
}

/**
Модуль сотрудников
*/

// EmployeesHandlers хендлеры модуля сотрудников
type EmployeesHandlers interface {
	getEmployee(c *gin.Context)
	getProfile(c *gin.Context)
}

/*
Use-case'ы
*/

/**
Модуль порталов
*/

// PortalsPortalsInteractor use-кейсы методов порталов
type PortalsPortalsInteractor interface {
	GetByEmployees(ctx context.Context, employees []entityPortal.EmployeeInfo) ([]*entityPortal.Portal, error)
	Filter(ctx context.Context, opts entityPortal.PortalsFilterOptions) ([]*entityPortal.Portal, error)
	GetAll(ctx context.Context, opts entityPortal.GetAllOptions) ([]*entityPortal.Portal, error)
	Get(ctx context.Context, id int, withDeleted bool) (*entityPortal.Portal, error)
	MultiplyAdd(ctx context.Context, entityPortal []*entityPortal.Portal) ([]*entityPortal.Portal, error)
	Add(ctx context.Context, newPortal *entityPortal.Portal) (*entityPortal.Portal, error)
	Update(ctx context.Context, newPortal *entityPortal.Portal) (*entityPortal.Portal, error)
	Delete(ctx context.Context, id int) error
}

// PortalsQuestionsInteractor use-кейсы методов вопросов
type PortalsQuestionsInteractor interface {
	GetAllQuestions(ctx context.Context, withDeleted bool) (*entityPortal.Questions, error)
	GetQuestion(ctx context.Context, questionId int, withDeleted bool) (*entityPortal.Question, error)
	AddQuestions(ctx context.Context, questions []*entityPortal.Question) ([]*entityPortal.Question, error)
	AddQuestion(ctx context.Context, question *entityPortal.Question) (*entityPortal.Question, error)
	UpdateQuestion(ctx context.Context, question *entityPortal.Question) (*entityPortal.Question, error)
	DeleteQuestion(ctx context.Context, questionId int) error
}

// PortalsImagesInteractor use-кейсы методов images
type PortalsImagesInteractor interface {
	All(ctx context.Context) ([]*entityPortal.Image, error)
	Get(ctx context.Context, imageId int) (*entityPortal.Image, error)
	GetRawImage(ctx context.Context, path string) (entityPortal.ImageData, error)
	Add(ctx context.Context, image *entityPortal.Image) (*entityPortal.Image, error)
	Delete(ctx context.Context, imageId int) error
}

// PortalsFeaturesInteractor use-кейсы методов features
type PortalsFeaturesInteractor interface {
	All(ctx context.Context, withDisabled bool) ([]*entityPortal.Feature, error)
	Get(ctx context.Context, featureId int, withDisabled bool) (*entityPortal.Feature, error)
	MultipleAdd(ctx context.Context, features []*entityPortal.Feature) ([]*entityPortal.Feature, error)
	Add(ctx context.Context, feature *entityPortal.Feature) (*entityPortal.Feature, error)
	Update(ctx context.Context, feature *entityPortal.Feature) (*entityPortal.Feature, error)
	Delete(ctx context.Context, featureId int) error
}

// PortalsOrganizationsInteractor use-кейсы методов организаций
type PortalsOrganizationsInteractor interface {
	Filter(
		ctx context.Context,
		filters entityPortal.OrganizationsFilters,
		pagination *entity.StringPagination,
		options entityPortal.OrganizationsFilterOptions,
	) (
		*entityPortal.OrganizationsWithPagination, error)
	Link(
		ctx context.Context,
		portalId entityPortal.PortalID,
		orgIds entityPortal.OrganizationIDs,
	) error
	Unlink(ctx context.Context, orgIds entityPortal.OrganizationIDs) error
}

/**
Модуль опросов
*/

// SurveysSurveysInteractor use-кейсы методов опросов
type SurveysSurveysInteractor interface {
	Get(
		ctx context.Context,
		id entitySurvey.SurveyID,
		options entitySurvey.SurveyFilterOptions,
	) (*entitySurvey.Survey, error)
}

// SurveysAnswersInteractor use-кейсы методов ответов на опрос
type SurveysAnswersInteractor interface {
	Add(ctx context.Context, answers []*entitySurvey.RespondentAnswer) ([]uuid.UUID, error)
}

// SurveysImagesInteractor use-кейсы методов изобржений для опроса
type SurveysImagesInteractor interface {
	Get(ctx context.Context, imageName string) ([]byte, error)
}

type AuthInteractor interface {
	// GetAuthURL
	//  URL для перенаправления пользователя для авторизации в СУДИР
	GetAuthURL(ctx context.Context, callbackURI string) (string, error)
	// Auth авторизация web пользователя
	//  метод возвращает информацию о пользователе в СУДИР
	//  и oauth2 токены
	Auth(ctx context.Context, code, state, callbackURI string) (*entityAuth.Auth, error)
	GetSession(ctx context.Context, accessToken string) (*entityAuth.Session, error)
	Logout(ctx context.Context, accessToken, refreshToken string) error
	// ChangePortal метод для смены активного портала. На вход принимает идентификатор выбранного портала 1С и сессию. Возвращает порталы и сессию портала 1С
	ChangePortal(ctx context.Context, selectedPortalID int) ([]*entityAuth.Portal, string, error)
	RefreshTokensPair(ctx context.Context, accessToken, refreshToken string) (*entityAuth.TokensPair, error)
}

/**
Модуль поиска сотрудников
*/

type EmployeesSearchUseCases interface {
	Search(ctx context.Context, request *entityEmployeesSearch.SearchParams) (*entityEmployeesSearch.SearchResponse, error)
	Filters(ctx context.Context, request *entityEmployeesSearch.SearchParams) (*entityEmployeesSearch.FiltersResponse, error)
}

/**
Модуль сотрудников
*/

type UsersInteractor interface {
	GetMe(ctx context.Context) (*entityUser.UserInfo, error)
}

/**
Презентеры
*/

type StringPaginationPresenter interface {
	PaginationToView(pagination *entity.StringPagination) *view.StringPagination
	PaginationToEntity(pagination *view.StringPagination) *entity.StringPagination
}

/**
Модуль порталов
*/

// PortalsPortalsPresenter презентер методов порталов
type PortalsPortalsPresenter interface {
	ToNewEntities(entityPortal []*viewPortals.NewPortal) []*entityPortal.Portal
	ToNewEntity(newPortal *viewPortals.NewPortal) *entityPortal.Portal
	ToEntities(entityPortal []*viewPortals.UpdatePortal) []*entityPortal.Portal
	ToEntity(viewPortal *viewPortals.UpdatePortal) *entityPortal.Portal
	ToViews(entityPortal []*entityPortal.Portal) []*viewPortals.Portal
	ToView(portal *entityPortal.Portal) *viewPortals.Portal
	ToShortViews(entityPortal []*entityPortal.Portal) []*viewPortals.PortalInfo
	ToShortView(portal *entityPortal.Portal) *viewPortals.PortalInfo
	ToWebViews(entityPortal []*entityPortal.Portal) []*viewPortals.WebPortal
	ToWebView(portal *entityPortal.Portal) *viewPortals.WebPortal
	FilterOptionsToEntity(options viewPortals.PortalsFilterOptions) entityPortal.PortalsFilterOptions
}

// PortalsQuestionsPresenter презентер методов вопросов
type PortalsQuestionsPresenter interface {
	ToNewEntities(questions []*viewPortals.NewQuestion) []*entityPortal.Question
	ToNewEntity(question *viewPortals.NewQuestion) *entityPortal.Question
	ToEntities(questions []*viewPortals.UpdateQuestion) []*entityPortal.Question
	ToEntity(question *viewPortals.UpdateQuestion) *entityPortal.Question
	ToViews(questions []*entityPortal.Question) []*viewPortals.Question
	ToView(question *entityPortal.Question) *viewPortals.Question
	ToShortViews(questions []*entityPortal.Question) []*viewPortals.QuestionInfo
	ToShortView(question *entityPortal.Question) *viewPortals.QuestionInfo
}

// PortalsImagesPresenter презентер методов изображений
type PortalsImagesPresenter interface {
	ToNewEntity(image *viewPortals.NewImage) *entityPortal.Image
	ToEntities(images []*viewPortals.Image) []*entityPortal.Image
	ToEntity(image *viewPortals.Image) *entityPortal.Image
	ToViews(images []*entityPortal.Image) []*viewPortals.Image
	ToView(image *entityPortal.Image) *viewPortals.Image
	ToShortViews(images []*entityPortal.Image) []*viewPortals.ImageInfo
	ToShortView(image *entityPortal.Image) *viewPortals.ImageInfo
}

// PortalsFeaturesPresenter презентер методов feature
type PortalsFeaturesPresenter interface {
	ToNewEntities(features []*viewPortals.NewFeature) []*entityPortal.Feature
	ToNewEntity(feature *viewPortals.NewFeature) *entityPortal.Feature
	ToEntities(features []*viewPortals.UpdateFeature) []*entityPortal.Feature
	ToEntity(feature *viewPortals.UpdateFeature) *entityPortal.Feature
	ToViews(features []*entityPortal.Feature) []*viewPortals.Feature
	ToView(feature *entityPortal.Feature) *viewPortals.Feature
	ToShortViews(features []*entityPortal.Feature) viewPortals.Features
	ToShortView(feature *entityPortal.Feature) *viewPortals.FeatureInfo
}

// PortalsOrganizationsPresenter презентер методов organization
type PortalsOrganizationsPresenter interface {
	StringPaginationPresenter
	OrganizationIdsToEntity(ids viewPortals.OrganizationIds) entityPortal.OrganizationIDs
	OrganizationsToView(orgs []*entityPortal.Organization) []*viewPortals.Organization
	FiltersToEntity(filters viewPortals.OrganizationsFilters) entityPortal.OrganizationsFilters
	OptionsToEntity(options *viewPortals.OrganizationsFilterOptions) entityPortal.OrganizationsFilterOptions
}

/**
Модуль опросов
*/

// SurveysPresenter презентер методов опросов
type SurveysPresenter interface {
	ToNewEntity(s *viewSurveys.NewSurvey) *entitySurvey.Survey
	ToEntity(s *viewSurveys.Survey) *entitySurvey.Survey
	ToView(s *entitySurvey.Survey) *viewSurveys.Survey
	ToShortView(s *entitySurvey.Survey) *viewSurveys.SurveyInfo
	IDsToEntities(ids []uuid.UUID) entitySurvey.SurveyIDs
	RespondentToEntity(respondent *viewSurveys.GetAllSurveysRespondent) *entitySurvey.SurveyRespondent
	OptionsToEntity(options *viewSurveys.SurveysOptions) entitySurvey.SurveyFilterOptions
	PaginationToEntity(pagination *viewSurveys.GetAllSurveysPagination) entitySurvey.Pagination
	SurveysWithPaginationToView(surveysWithPagination *entitySurvey.SurveysWithPagination) *viewSurveys.SurveysWithPagination
	IDToView(ID *entitySurvey.SurveyID) *viewSurveys.IDResponse
}

// SurveysAnswersPresenter презентер методов ответов на опрос
type SurveysAnswersPresenter interface {
	ToNewEntities(answers *viewSurveys.NewSurveyAnswers) []*entitySurvey.RespondentAnswer
	ToViews(answers []*entitySurvey.RespondentAnswer) []*viewSurveys.SurveyAnswer
	ToShortViews(ids []uuid.UUID) []*viewSurveys.SurveyAnswerInfo
}

// SurveysImagesPresenter презентер методов изображений для опроса
type SurveysImagesPresenter interface {
	ToNewEntity(image *viewSurveys.NewSurveyImageObject) *entitySurvey.Image
	ToView(image *entitySurvey.Image) *viewSurveys.SurveyImageObject
}

type AuthPresenter interface {
	AuthToView(authInfo *entityAuth.Auth) *viewAuth.AuthResponse
}

/**
Модуль поиска сотрудников
*/

type EmployeesSearchPresenter interface {
	SearchRequestToEntity(search *viewEmployeesSearch.SearchRequest) *entityEmployeesSearch.SearchParams
	SearchResponseToView(search *entityEmployeesSearch.SearchResponse) *viewEmployeesSearch.SearchResponse
	FiltersRequestToEntity(filter *viewEmployeesSearch.FiltersRequest) *entityEmployeesSearch.SearchParams
	FiltersResponseToView(filter *entityEmployeesSearch.FiltersResponse) *viewEmployeesSearch.FiltersResponse
	GenderToEntity(gender string) entity.Gender
	GenderToView(gender entity.Gender) string
}

/**
Модуль сотрудников
*/

type EmployeesUseCases interface {
	Get(ctx context.Context, id uuid.UUID) (*entityEmployee.Employee, error)
	GetByExtIDAndPortalID(ctx context.Context, extID string, portalID int) (*entityEmployee.Employee, error)
}

type UsersPresenter interface {
	ShortUserToView(employee *entityUser.ShortUser) *viewUsers.ShortUser
}

type FilesHandlers interface {
	get(c *gin.Context)
}

// FilesInteractor use-кейсы методов изобржений для опроса
type FilesInteractor interface {
	Get(ctx context.Context, fileId uuid.UUID) (*entityFile.File, error)
}

type EmployeesPresenter interface {
	EmployeeToView(composite *entityEmployee.Employee) *viewEmployees.Employee
}

/**
Модуль аналитики
*/

type AnalyticsHandlers interface {
	addMetrics(c *gin.Context)
}

type AnalyticsInteractor interface {
	AddMetrics(ctx context.Context, headers analytics.XCFCUserAgentHeader, body []byte) (string, error)
}

/**
Модуль порталов v2
*/

type PortalsV2Presenter interface {
	PortalsWithCountToView(portalsWithCounts []*entityPortalsV2.PortalWithCounts) []*viewPortalsV2.Portal
	PortalsFilterToEntity(filter *viewPortalsV2.PortalsFilterRequest) *entityPortalsV2.FilterPortalsFilters
}

type ComplexesV2Presenter interface {
	ComplexesToView(complexes []*entityPortalsV2.Complex) []*viewPortalsV2.Complex
}

type PortalsV2Interactor interface {
	Filter(
		ctx context.Context,
		filters *entityPortalsV2.FilterPortalsFilters,
		options *entityPortalsV2.FilterPortalsOptions,
	) ([]*entityPortalsV2.PortalWithCounts, error)
}

type ComplexesV2Interactor interface {
	Filter(
		ctx context.Context,
		filters *entityPortalsV2.FilterComplexesFilters,
		options *entityPortalsV2.FilterComplexesOptions,
	) ([]*entityPortalsV2.Complex, error)
}

type PortalsV2Handlers interface {
	getPortals(c *gin.Context)
	getComplexes(c *gin.Context)
}

/*
Модуль news
*/
type NewsCategoryInteractor interface {
	Create(ctx context.Context, nc *dtoNews.NewCategory) (*entityNews.Category, error)
	Update(ctx context.Context, c *dtoNews.UpdateCategory) (*entityNews.Category, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Search(crx context.Context, c *dtoNews.SearchCategory) (*entityNews.CategoriesWithPagination, error)
	Get(ctx context.Context, id uuid.UUID) (*entityNews.Category, error)
}
type NewsCommentsInteractor interface {
	Create(ctx context.Context, in dtoNews.NewComment) (uuid.UUID, int, error)
	List(ctx context.Context, params *dtoNews.FilterComments) ([]*entityNews.NewsComment, int, error)
}
type NewsAdminInteractor interface {
	Create(ctx context.Context, news *dtoNews.NewNews) (uuid.UUID, error)
	Update(ctx context.Context, id uuid.UUID, updateNews *dtoNews.UpdateNews) (*entityNews.News, error)
	ChangeStatus(ctx context.Context, id uuid.UUID, status entityNews.NewsStatus) (*entityNews.News, error)
	Get(ctx context.Context, id uuid.UUID) (*entityNews.NewsFull, error)
	Search(ctx context.Context, search *dtoNews.SearchNews) (*dtoNews.SearchNewsResult, error)
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateFlags(ctx context.Context, id uuid.UUID, updateNews *dtoNews.UpdateFlags) (*entityNews.News, error)
}

type NewsInteractor interface {
	Get(ctx context.Context, slug string) (*entityNews.NewsFull, error)
	Search(ctx context.Context, search *dtoNews.SearchNews) (*dtoNews.SearchNewsResult, error)
}

type NewsAdminPresenter interface {
	NewsCategoryPresenter
	NewsCommentsPresenter
	NewNewsToDTO(news *viewNews.NewNews) *dtoNews.NewNews
	UpdateNewsToDTO(updateNews *viewNews.UpdateNews) *dtoNews.UpdateNews
	SearchNewsToDTO(search *viewNews.SearchNewsRequest) *dtoNews.SearchNews
	StatusToEntity(status viewNews.NewsStatus) entityNews.NewsStatus
	FullNewsToView(n *entityNews.NewsFull) *viewNews.News
	NewsCategoryToView(category *entityNews.Category) *viewNews.NewsCategory
	NewsOrganizationToView(organization *entityNews.NewsOrganization) *viewNews.NewsOrganization
	NewsProductToView(product *entityNews.NewsProduct) *viewNews.NewsProduct
	AuthorToView(author entityNews.Author) viewNews.Author
	StatusToView(status entityNews.NewsStatus) viewNews.NewsStatus
	ParticipantsToView(participants []*entityNews.Participant) []*viewNews.NewsParticipants
	ParticipantToView(participant *entityNews.Participant) *viewNews.NewsParticipants
	FullNewsToSearchItems(n []*entityNews.NewsFull) []*viewNews.SearchNewsResponseItem
	FullNewsToSearchItem(n *entityNews.NewsFull) *viewNews.SearchNewsResponseItem
	UpdateFlagsToDTO(n *viewNews.UpdateNewsFlags) *dtoNews.UpdateFlags
}

type NewsCategoryPresenter interface {
	CategoryToView(c *entityNews.Category) *viewNews.Category
	NewCategoryToDTO(nc *viewNews.NewCategory) *dtoNews.NewCategory
	UpdateCategoryToDTO(c *viewNews.UpdateCategory) *dtoNews.UpdateCategory
	CategoryToResult(c *entityNews.Category) *viewNews.CategoryResult
}

type NewsCommentsPresenter interface {
	NewCommentToDTO(newsID uuid.UUID, v *viewNews.NewNewsComment) dtoNews.NewComment
	CommentsToView(list []*entityNews.NewsComment) []*viewNews.NewsComment
}

type NewsAdminHandlers interface {
	createCategory(c *gin.Context)
	updateCategory(c *gin.Context)
	searchCategory(c *gin.Context)
	getCategory(c *gin.Context)
	deleteCategory(c *gin.Context)
	createNews(c *gin.Context)
	updateNews(c *gin.Context)
	searchNews(c *gin.Context)
	getNews(c *gin.Context)
	deleteNews(c *gin.Context)
	setStatusNews(c *gin.Context)
	setFlagsNews(c *gin.Context)
}

type NewsHandlers interface {
	searchNews(c *gin.Context)
	getNews(c *gin.Context)
	createComment(c *gin.Context)
	listComments(c *gin.Context)
}

/*
Модуль banners
*/
type BannersHandlers interface {
	list(c *gin.Context)
	set(c *gin.Context)
}

type BannersPresenter interface {
	SetBannerToDTO(view *viewBanners.SetBanner) *dtoBanners.SetBanner
	SetBannersToDTOs(view *viewBanners.SetBanners) []*dtoBanners.SetBanner

	BannerTypeToEntity(t viewBanners.BannerType) entityBanners.BannerType
	BannerTypeToView(t entityBanners.BannerType) viewBanners.BannerType

	BannerInfoToView(banner *entityBanners.BannerInfo) *viewBanners.BannerInfo
	BannerInfosToViews(banners []*entityBanners.BannerInfo) []*viewBanners.BannerInfo

	BannerToView(banner *entityBanners.Banner) *viewBanners.Banner
	ContentToView(content entityBanners.Content) viewBanners.Content
	BannersToViews(banner []*entityBanners.Banner) []*viewBanners.Banner
}

type BannersInteractor interface {
	List(ctx context.Context) ([]*entityBanners.Banner, []*entityBanners.Banner, []*entityBanners.Banner, error)
	Set(ctx context.Context, banners []*dtoBanners.SetBanner) ([]*entityBanners.BannerInfo, error)
}
