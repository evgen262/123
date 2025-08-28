package http

import (
	"context"
	"errors"
	"net/http"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/portalsv2"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portalv2"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type portalsV2Handlers struct {
	portalsV2Interactor PortalsV2Interactor
	portalsV2Presenter  PortalsV2Presenter

	complexesV2Interactor ComplexesV2Interactor
	complexesV2Presenter  ComplexesV2Presenter

	logger ditzap.Logger
}

func NewPortalsV2Handlers(
	portalsV2Interactor PortalsV2Interactor,
	portalsV2Presenter PortalsV2Presenter,
	complexesV2Interactor ComplexesV2Interactor,
	complexesV2Presenter ComplexesV2Presenter,
	logger ditzap.Logger) *portalsV2Handlers {
	return &portalsV2Handlers{
		portalsV2Interactor:   portalsV2Interactor,
		complexesV2Interactor: complexesV2Interactor,
		complexesV2Presenter:  complexesV2Presenter,
		portalsV2Presenter:    portalsV2Presenter,
		logger:                logger,
	}
}

func (ph *portalsV2Handlers) getPortals(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c, RequestTimeOut)
	defer cancelCtx()

	var filter portalsv2.PortalsFilterRequest
	if err := c.ShouldBindJSON(&filter); err != nil {
		ph.logger.Debug("can't unbind PortalFilterRequest json", zap.Error(err))
		c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		return
	}

	portalsWithCounts, err := ph.portalsV2Interactor.Filter(ctx,
		ph.portalsV2Presenter.PortalsFilterToEntity(&filter),
		&portalv2.FilterPortalsOptions{
			WithEmployeesCount: true,
		})
	if err != nil {
		switch {
		case errors.As(err, new(diterrors.ValidationError)):
			c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		default:
			c.JSON(http.StatusInternalServerError, view.NewErrorResponse(err))
		}
		return
	}

	c.JSON(http.StatusOK, ph.portalsV2Presenter.PortalsWithCountToView(portalsWithCounts))
}

func (h *portalsV2Handlers) getComplexes(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c, RequestTimeOut)
	defer cancelCtx()

	complexes, err := h.complexesV2Interactor.Filter(ctx, nil, nil)
	if err != nil {
		switch {
		case errors.As(err, new(diterrors.ValidationError)):
			c.JSON(http.StatusBadRequest, view.NewErrorResponse(err.(diterrors.ValidationError).Unwrap()))
		default:
			c.JSON(http.StatusInternalServerError, view.NewErrorResponse(err))
		}
		return
	}
	c.JSON(http.StatusOK, h.complexesV2Presenter.ComplexesToView(complexes))
}
