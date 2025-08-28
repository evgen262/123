package http

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"github.com/gin-gonic/gin"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view"
	authView "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/auth"
	entityAuth "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/auth"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/usecase"
	authUseCase "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/usecase/auth"
	usersUseCase "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/usecase/users"
)

type usersHandlers struct {
	authInteractor  AuthInteractor
	usersInteractor UsersInteractor
	authPresenter   AuthPresenter
	usersPresenter  UsersPresenter
}

func NewUsersHandlers(
	authInteractor AuthInteractor,
	usersInteractor UsersInteractor,
	authPresenter AuthPresenter,
	usersPresenter UsersPresenter,
) *usersHandlers {
	return &usersHandlers{
		authInteractor:  authInteractor,
		usersInteractor: usersInteractor,
		authPresenter:   authPresenter,
		usersPresenter:  usersPresenter,
	}
}

func (uh usersHandlers) getMe(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c.Request.Context(), RequestTimeOut)
	defer cancelCtx()

	userInfo, err := uh.usersInteractor.GetMe(ctx)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrGetSessionFromContext):
			c.JSON(http.StatusUnauthorized, view.NewErrorResponse(authView.ErrMessageFrontUnauthenticated))
		case errors.Is(err, usersUseCase.ErrEmptySessionPersonID):
			c.JSON(http.StatusUnauthorized, view.NewErrorResponse(authView.ErrMessageFrontUnauthenticated))
		default:
			c.JSON(http.StatusInternalServerError, view.NewErrorResponse(authView.ErrMessageFrontInternal))
		}
		return
	}
	c.JSON(http.StatusOK, view.NewSuccessResponse(uh.usersPresenter.ShortUserToView(&userInfo.User)))
}

func (uh usersHandlers) changePortal(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c.Request.Context(), RequestTimeOut)
	defer cancelCtx()

	selectedPortalID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		return
	}
	if selectedPortalID < 1 {
		c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidId))
		return
	}

	portals, selectedPortalSID, err := uh.authInteractor.ChangePortal(ctx, selectedPortalID)
	if err != nil {
		switch {
		case errors.As(err, new(diterrors.ValidationError)):
			c.JSON(http.StatusBadRequest, view.NewErrorResponse(authView.ErrMessageFrontSamePortals))
		case errors.Is(err, usecase.ErrGetSessionFromContext):
			c.JSON(http.StatusUnauthorized, view.NewErrorResponse(authView.ErrMessageFrontUnauthenticated))
		case errors.Is(err, authUseCase.ErrPortalsNotFound):
			c.JSON(http.StatusBadRequest, view.NewErrorResponse(authView.ErrMessageFrontPortalForUserNotFound))
		case errors.Is(err, authUseCase.ErrUnavailablePortal):
			c.JSON(http.StatusBadRequest, view.NewErrorResponse(authView.ErrMessageFrontUnavailablePortal))
		case errors.Is(err, authUseCase.ErrEmptyPortalURL):
			c.JSON(http.StatusInternalServerError, view.NewErrorResponse(authView.ErrMessageFrontEmptyPortalURL))
		case errors.Is(err, diterrors.ErrUnauthenticated):
			c.JSON(http.StatusUnauthorized, view.NewErrorResponse(authView.ErrMessageFrontUnauthenticated))
		case errors.Is(err, diterrors.ErrNotFound):
			c.JSON(http.StatusNotFound, view.NewErrorResponse(authView.ErrMessageFrontUserNotIntoPortal))
		default:
			c.JSON(http.StatusInternalServerError, view.NewErrorResponse(authView.ErrMessageFrontInternal))
		}
		return
	}
	c.JSON(http.StatusOK, view.NewSuccessResponse(uh.authPresenter.AuthToView(&entityAuth.Auth{PortalSession: selectedPortalSID, Portals: portals})))
}
