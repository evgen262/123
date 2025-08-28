package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view"
	viewSession "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/auth"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity"
	entitySession "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/auth"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/usecase/auth"
)

type redirectSessionHandler struct {
	redirectSessionInteractor RedirectSessionInteractor
}

func NewRedirectSessionHandlers(redirectSessionInteractor RedirectSessionInteractor) *redirectSessionHandler {
	return &redirectSessionHandler{
		redirectSessionInteractor: redirectSessionInteractor,
	}
}

func (rsh redirectSessionHandler) createSession(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c.Request.Context(), RequestTimeOut)
	defer cancelCtx()

	s, err := entity.SessionFromContext(ctx)
	if err != nil {
		c.Header(StatusCodeHeader, "RSH_01")
		c.JSON(http.StatusForbidden, view.NewErrorResponse(view.ErrMessageUnauthorized))
		return
	}

	type request struct {
		PortalURL string `json:"portalUrl" binding:"required"`
		TargetURL string `json:"targetUrl" binding:"required"`
	}

	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Header(StatusCodeHeader, "RSH_02")
		c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		return
	}

	redirectURL, err := rsh.redirectSessionInteractor.CreateSession(ctx, &entitySession.RedirectSessionUserInfo{
		Email:     s.GetUser().Email,
		SNILS:     s.GetUser().SNILS,
		IP:        s.UserIP.String(),
		UserAgent: s.Device.UserAgent,
		PortalURL: req.PortalURL,
		TargetURL: req.TargetURL,
	})
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrUserInfoRequired):
			c.Header(StatusCodeHeader, "RSH_03")
			c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
			return
		default:
			c.JSON(http.StatusInternalServerError, view.NewErrorResponse(view.ErrMessageInternalError))
			return
		}
	}

	c.JSON(http.StatusOK, view.NewSuccessResponse(&viewSession.CreateSessionSessionResponse{
		RedirectURL: redirectURL,
	}))
}
