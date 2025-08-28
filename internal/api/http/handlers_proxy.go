package http

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"github.com/gin-gonic/gin"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/event"
	reposiotryProxy "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/repositories/proxy"
)

type proxyHandlers struct {
	proxyInteractor ProxyInteractor
	proxyPresenter  ProxyPresenter
}

func NewProxyHandlers(proxyInteractor ProxyInteractor, proxyPresenter ProxyPresenter) *proxyHandlers {
	return &proxyHandlers{
		proxyInteractor: proxyInteractor,
		proxyPresenter:  proxyPresenter,
	}
}

func (h proxyHandlers) listHomeBanners(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c.Request.Context(), RequestTimeOut)
	defer cancelCtx()

	session, err := entity.SessionFromContext(ctx)
	if err != nil {
		c.Header(StatusCodeHeader, "PBH_01")
		c.JSON(http.StatusForbidden, view.NewErrorResponse(view.ErrMessageUnauthorized))
		return
	}

	if session == nil || session.GetActivePortal() == nil {
		c.Header(StatusCodeHeader, "PBH_02")
		c.JSON(http.StatusForbidden, view.NewErrorResponse(view.ErrMessageUnauthorized))
		return
	}

	sessionID := session.GetActivePortal().SID
	if sessionID == "" {
		c.Header(StatusCodeHeader, "PBH_03")
		c.JSON(http.StatusForbidden, view.NewErrorResponse(view.ErrMessageUnauthorized))
		return
	}

	portalURL := session.GetActivePortal().GetPortal().URL
	if portalURL == "" {
		c.Header(StatusCodeHeader, "PBH_04")
		c.JSON(http.StatusForbidden, view.NewErrorResponse(view.ErrMessageUnauthorized))
		return
	}

	bannersList, err := h.proxyInteractor.ListHomeBanners(ctx, sessionID, portalURL)
	if err != nil {
		switch {
		case errors.Is(err, reposiotryProxy.ErrNotFound):
			// Баннеры не найдены
			c.Header(StatusCodeHeader, "PBH_05")
			c.JSON(http.StatusNotFound, view.NewErrorResponse(view.ErrMessageNotFound))
		case errors.As(err, new(diterrors.ValidationError)):
			fallthrough
		case errors.Is(err, reposiotryProxy.ErrBadRequest):
			// Не верный запрос
			c.Header(StatusCodeHeader, "PBH_06")
			c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		case errors.Is(err, reposiotryProxy.ErrUnauthorized):
			// Пользователь не авторизован
			c.Header(StatusCodeHeader, "PBH_07")
			c.JSON(http.StatusUnauthorized, view.NewErrorResponse(view.ErrMessageUnauthorized))
		case errors.Is(err, reposiotryProxy.ErrPermissionDenied):
			// Недостаточно прав
			c.Header(StatusCodeHeader, "PBH_08")
			c.JSON(http.StatusForbidden, view.NewErrorResponse(view.ErrPermissionDenied))
		default:
			// Внутренняя ошибка
			c.Header(StatusCodeHeader, "PBH_09")
			c.JSON(http.StatusInternalServerError, view.NewErrorResponse(view.ErrMessageInternalError))
		}
		return
	}

	c.JSON(http.StatusOK, view.NewSuccessResponse(h.proxyPresenter.BannersListToView(bannersList)))
}

func (h proxyHandlers) listCalendarEvents(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c.Request.Context(), RequestTimeOut)
	defer cancelCtx()

	session, err := entity.SessionFromContext(ctx)
	if err != nil {
		c.Header(StatusCodeHeader, "PE_01")
		c.JSON(http.StatusForbidden, view.NewErrorResponse(view.ErrMessageUnauthorized))
		return
	}

	if session == nil || session.GetActivePortal() == nil {
		c.Header(StatusCodeHeader, "PE_02")
		c.JSON(http.StatusForbidden, view.NewErrorResponse(view.ErrMessageUnauthorized))
		return
	}

	sessionID := session.GetActivePortal().SID
	if sessionID == "" {
		c.Header(StatusCodeHeader, "PE_03")
		c.JSON(http.StatusForbidden, view.NewErrorResponse(view.ErrMessageUnauthorized))
		return
	}

	portalURL := session.GetActivePortal().GetPortal().URL
	if portalURL == "" {
		c.Header(StatusCodeHeader, "PE_04")
		c.JSON(http.StatusForbidden, view.NewErrorResponse(view.ErrMessageUnauthorized))
		return
	}

	year, err := strconv.Atoi(c.Query("year"))
	if year <= 0 || err != nil {
		c.Header(StatusCodeHeader, "PE_05")
		c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		return
	}

	month, err := strconv.Atoi(c.Query("month"))
	if (month <= 0 && month > 12) || err != nil {
		c.Header(StatusCodeHeader, "PE_06")
		c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		return
	}

	eventsList, err := h.proxyInteractor.ListCalendarEvents(ctx, event.CalendarEventRequest{
		SessionID: sessionID,
		PortalURL: portalURL,
		Year:      year,
		Month:     month,
	})
	if err != nil {
		switch {
		case errors.Is(err, reposiotryProxy.ErrNotFound):
			// События не найдены
			c.Header(StatusCodeHeader, "PE_07")
			c.JSON(http.StatusNotFound, view.NewErrorResponse(view.ErrMessageNotFound))
		case errors.As(err, new(diterrors.ValidationError)):
			fallthrough
		case errors.Is(err, reposiotryProxy.ErrBadRequest):
			// Не верный запрос
			c.Header(StatusCodeHeader, "PE_08")
			c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		case errors.Is(err, reposiotryProxy.ErrUnauthorized):
			// Пользователь не авторизован
			c.Header(StatusCodeHeader, "PE_09")
			c.JSON(http.StatusUnauthorized, view.NewErrorResponse(view.ErrMessageUnauthorized))
		case errors.Is(err, reposiotryProxy.ErrPermissionDenied):
			// Недостаточно прав
			c.Header(StatusCodeHeader, "PE_10")
			c.JSON(http.StatusForbidden, view.NewErrorResponse(view.ErrPermissionDenied))
		default:
			// Внутренняя ошибка
			c.Header(StatusCodeHeader, "PE_11")
			c.JSON(http.StatusInternalServerError, view.NewErrorResponse(view.ErrMessageInternalError))
		}
		return
	}

	c.JSON(http.StatusOK, view.NewSuccessResponse(h.proxyPresenter.EventsListToView(eventsList)))
}

func (h proxyHandlers) listCalendarEventsLinks(c *gin.Context) {
	type request struct {
		EventIDs []string `json:"ids"`
	}

	ctx, cancelCtx := context.WithTimeout(c.Request.Context(), RequestTimeOut)
	defer cancelCtx()

	session, err := entity.SessionFromContext(ctx)
	if err != nil {
		c.Header(StatusCodeHeader, "PEL_01")
		c.JSON(http.StatusForbidden, view.NewErrorResponse(view.ErrMessageUnauthorized))
		return
	}

	if session == nil || session.GetActivePortal() == nil {
		c.Header(StatusCodeHeader, "PEL_02")
		c.JSON(http.StatusForbidden, view.NewErrorResponse(view.ErrMessageUnauthorized))
		return
	}

	sessionID := session.GetActivePortal().SID
	if sessionID == "" {
		c.Header(StatusCodeHeader, "PEL_03")
		c.JSON(http.StatusForbidden, view.NewErrorResponse(view.ErrMessageUnauthorized))
		return
	}

	portalURL := session.GetActivePortal().GetPortal().URL
	if portalURL == "" {
		c.Header(StatusCodeHeader, "PEL_04")
		c.JSON(http.StatusForbidden, view.NewErrorResponse(view.ErrMessageUnauthorized))
		return
	}

	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Header(StatusCodeHeader, "PEL_05")
		c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		return
	}

	if len(req.EventIDs) == 0 {
		c.Header(StatusCodeHeader, "PEL_06")
		c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		return
	}

	eventsLinks, err := h.proxyInteractor.ListCalendarEventsLinks(ctx, event.CalendarEventLinksRequest{
		SessionID: sessionID,
		PortalURL: portalURL,
		EventIDs:  req.EventIDs,
	})
	if err != nil {
		switch {
		case errors.Is(err, reposiotryProxy.ErrNotFound):
			// Ссылки на события не нашлись
			c.Header(StatusCodeHeader, "PEL_07")
			c.JSON(http.StatusNotFound, view.NewErrorResponse(view.ErrMessageNotFound))
		case errors.As(err, new(diterrors.ValidationError)):
			fallthrough
		case errors.Is(err, reposiotryProxy.ErrBadRequest):
			// Неправильные параметры запроса
			c.Header(StatusCodeHeader, "PEL_08")
			c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		case errors.Is(err, reposiotryProxy.ErrUnauthorized):
			// Пользователь не авторизован
			c.Header(StatusCodeHeader, "PEL_09")
			c.JSON(http.StatusUnauthorized, view.NewErrorResponse(view.ErrMessageUnauthorized))
		case errors.Is(err, reposiotryProxy.ErrPermissionDenied):
			// Нет прав доступа
			c.Header(StatusCodeHeader, "PEL_10")
			c.JSON(http.StatusForbidden, view.NewErrorResponse(view.ErrPermissionDenied))
		default:
			//  Внутренняя ошибка
			c.Header(StatusCodeHeader, "PEL_11")
			c.JSON(http.StatusInternalServerError, view.NewErrorResponse(view.ErrMessageInternalError))
		}
		return
	}

	c.JSON(http.StatusOK, view.NewSuccessResponse(h.proxyPresenter.EventsLinksToView(eventsLinks)))
}
