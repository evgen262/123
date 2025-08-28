package http

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"github.com/gin-gonic/gin"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view"
	authView "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/auth"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/repositories"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/usecase"
	authUseCase "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/usecase/auth"
)

// GET параметры ручки аутентификации
const (
	// paramState
	//  Get параметр state обязательный уникальный идентификатор сессии аутентификации
	//  передаётся в СУДИР и возвращается после успешной аутентификации
	paramState = "state"
	// paramCode
	//  Get параметр code обязательный уникальный код аутентификации пользователя
	//  возвращается после успешной аутентификации в СУДИР
	paramCode = "code"
	// paramCallbackURI
	//  Get параметр callback_uri страница портала на которую будет
	//  перенаправлен пользователь после успешной аутентификации в СУДИР
	paramCallbackURI = "callback_uri"
)

const defaultCookieTTL = time.Minute

type authHandlers struct {
	interactor AuthInteractor
	presenter  AuthPresenter
	logger     ditzap.Logger
}

func NewAuthHandlers(
	interactor AuthInteractor,
	presenter AuthPresenter,
	logger ditzap.Logger,
) *authHandlers {
	return &authHandlers{
		interactor: interactor,
		presenter:  presenter,
		logger:     logger,
	}
}

// @Summary Получение отфильтрованного списка порталов
// @Description Выдаются порталы, удовлетворяющие переданным фильтрам.
// @Tags     Порталы
// @Produce  json
// @Param    options body PortalsFilterOptions true "фильтры"
// @Router   /auth/v1/auth [get]
// @Success  200 {object} AuthResponse
// @Failure  400,401,403,404,500 {object} ErrorResponse
func (ah *authHandlers) auth(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c.Request.Context(), RequestTimeOut)
	defer cancelCtx()

	stateParam := c.Query(paramState)
	codeParam := c.Query(paramCode)
	callbackURIParam := c.Query(paramCallbackURI)

	if stateParam == "" {
		redirectURL, err := ah.interactor.GetAuthURL(ctx, callbackURIParam)
		if err != nil {
			// Ошибка получения URL для перенаправления в СУДИР
			c.Header(StatusCodeHeader, "WAA_01")
			c.JSON(http.StatusInternalServerError, view.NewErrorResponse(authView.ErrMessageFrontInternal))
			return
		}
		c.Redirect(http.StatusFound, redirectURL)
		return
	}

	if codeParam == "" {
		// Отсутствует get-параметр code
		c.Header(StatusCodeHeader, "WAA_02")
		c.JSON(http.StatusBadRequest, view.NewErrorResponse(authView.ErrMessageFrontInvalidParams))
		return
	}

	authInfo, err := ah.interactor.Auth(ctx, codeParam, stateParam, callbackURIParam)
	if err != nil {
		switch {
		case errors.Is(err, diterrors.ErrUnauthenticated):
			c.JSON(http.StatusUnauthorized, view.NewErrorResponse(authView.ErrMessageFrontUnauthenticated))
		case errors.Is(err, diterrors.ErrPermissionDenied):
			c.JSON(http.StatusForbidden, view.NewErrorResponse(view.ErrMessageUnauthorized))
		case errors.Is(err, authUseCase.ErrUserAccessDenied):
			c.JSON(http.StatusForbidden, view.NewErrorResponse(authView.ErrMessageFrontUserAccessDenied))
		case errors.As(err, new(diterrors.ValidationError)):
			c.JSON(http.StatusBadRequest, view.NewErrorResponse(authView.ErrMessageFrontInvalidSUDIRRedirect))
		case errors.Is(err, authUseCase.ErrPortalsNotFound):
			c.JSON(http.StatusForbidden, view.NewErrorResponse(authView.ErrMessageFrontPortalForUserNotFound))
		case errors.Is(err, authUseCase.ErrEmployeesNotFound):
			c.JSON(http.StatusForbidden, view.NewErrorResponse(authView.ErrMessageFrontSKSEmployeeNotFound))
		case errors.Is(err, authUseCase.ErrSUDIRNoCloudID):
			c.JSON(http.StatusForbidden, view.NewErrorResponse(authView.ErrMessageFrontSUDIRUserNotFound))
		case errors.Is(err, diterrors.ErrNotFound):
			c.JSON(http.StatusNotFound, view.NewErrorResponse(authView.ErrMessageFrontUserNotIntoPortal))
		case errors.Is(err, authUseCase.ErrInvalidDevice):
			c.JSON(http.StatusBadRequest, view.NewErrorResponse(authView.ErrMessageFrontInvalidParams))
		default:
			c.JSON(http.StatusInternalServerError, view.NewErrorResponse(authView.ErrMessageFrontInternal))
		}
		return
	}
	// TODO: Добавить геттеры
	if authInfo.JWTToken.Value == "" || authInfo.RefreshToken.Value == "" || authInfo.PortalSession == "" {
		// Неверный один из токенов или сессия портала
		c.Header(StatusCodeHeader, "WAA_12")
		c.JSON(http.StatusInternalServerError, view.NewErrorResponse(authView.ErrMessageFrontInternal))
		return
	}

	cookieTTL := defaultCookieTTL
	if authInfo.RefreshToken.GetExpiredAt() != nil {
		cookieTTL = authInfo.RefreshToken.GetExpiredAt().Sub(time.Now())
	}

	c.SetSameSite(http.SameSiteNoneMode)

	c.SetCookie(
		RefreshTokenCookie,
		authInfo.RefreshToken.Value,
		int(cookieTTL.Seconds()),
		"/web/auth",
		SharedFields.ExternalHost,
		true,
		true,
	)

	c.Header(JWTTokenHeader, strings.Join([]string{"Bearer", authInfo.JWTToken.Value}, " "))
	c.JSON(http.StatusOK, view.NewSuccessResponse(ah.presenter.AuthToView(authInfo)))
}

func (ah *authHandlers) logout(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c.Request.Context(), RequestTimeOut)
	defer cancelCtx()

	accessToken := strings.TrimPrefix(c.GetHeader(JWTTokenHeader), "Bearer ")
	refreshToken, err := c.Cookie(RefreshTokenCookie)
	if err != nil {
		// Не передан refresh-token
		c.Header(StatusCodeHeader, "WAL_01")
		c.JSON(http.StatusBadRequest, view.NewErrorResponse(authView.ErrMessageFrontInvalidParams))
		return
	}

	if err = ah.interactor.Logout(ctx, accessToken, refreshToken); err != nil {
		switch {
		case errors.As(err, new(diterrors.ValidationError)):
			// Не верный один из параметров
			c.Header(StatusCodeHeader, "WAL_02")
			c.JSON(http.StatusBadRequest, view.NewErrorResponse(authView.ErrMessageFrontInvalidParams))
		case errors.Is(err, diterrors.ErrNotFound) || errors.Is(err, usecase.ErrGetSessionFromContext):
			// Не найдена сессия
			c.Header(StatusCodeHeader, "WAL_03")
			c.JSON(http.StatusUnauthorized, view.NewErrorResponse(authView.ErrMessageFrontUnauthenticated))
		case errors.Is(err, authUseCase.ErrInvalidSession):
			// Сессия не валидна
			c.Header(StatusCodeHeader, "WAL_04")
			c.JSON(http.StatusUnauthorized, view.NewErrorResponse(authView.ErrMessageFrontUnauthenticated))
		default:
			// Внутренняя ошибка сервиса
			c.Header(StatusCodeHeader, "WAL_05")
			c.JSON(http.StatusInternalServerError, view.NewErrorResponse(authView.ErrMessageFrontInternal))
		}
	}

	c.Status(http.StatusOK)
}

func (ah *authHandlers) refresh(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c.Request.Context(), RequestTimeOut)
	defer cancelCtx()

	accessToken := strings.TrimPrefix(c.GetHeader(JWTTokenHeader), "Bearer ")
	refreshToken, err := c.Cookie(RefreshTokenCookie)
	if err != nil {
		// Не передан refresh-token
		c.Header(StatusCodeHeader, "WAR_01")
		c.JSON(http.StatusBadRequest, view.NewErrorResponse(authView.ErrMessageFrontInvalidParams))
		return
	}

	tokensPair, err := ah.interactor.RefreshTokensPair(ctx, accessToken, refreshToken)
	if err != nil {
		var errDetails *repositories.DetailsError
		switch {
		case errors.As(err, &errDetails) && !errDetails.GetReauthRequired():
			// Некорретный токен доступа, токен обновления или идентификатор сессии
			c.Header(StatusCodeHeader, "WAR_02")
			c.JSON(http.StatusBadRequest, view.NewErrorResponse(authView.ErrMessageFrontInvalidParams))
		case errors.As(err, &errDetails) && errDetails.GetReauthRequired():
			// Истекший или уже использованный токен обновления
			c.Header(StatusCodeHeader, "WAR_03")
			c.JSON(http.StatusUnauthorized, view.NewErrorResponse(authView.ErrMessageFrontUnauthenticated))
		case errors.Is(err, diterrors.ErrNotFound):
			// Не найден токен обновления или сессия
			c.Header(StatusCodeHeader, "WAR_04")
			c.JSON(http.StatusUnauthorized, view.NewErrorResponse(authView.ErrMessageFrontUnauthenticated))
		default:
			// Внутренняя ошибка сервиса
			c.Header(StatusCodeHeader, "WAR_05")
			c.JSON(http.StatusInternalServerError, view.NewErrorResponse(authView.ErrMessageFrontInternal))
		}
		return
	}

	cookieTTL := defaultCookieTTL
	if tokensPair.RefreshToken.GetExpiredAt() != nil {
		cookieTTL = tokensPair.RefreshToken.GetExpiredAt().Sub(time.Now())
	}

	c.SetSameSite(http.SameSiteNoneMode)

	c.SetCookie(
		RefreshTokenCookie,
		tokensPair.RefreshToken.Value,
		int(cookieTTL.Seconds()),
		"/web/auth",
		SharedFields.ExternalHost,
		true,
		true,
	)

	c.Header(JWTTokenHeader, strings.Join([]string{"Bearer", tokensPair.AccessToken.Value}, " "))
	c.Status(http.StatusOK)
}
