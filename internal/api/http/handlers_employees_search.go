package http

import (
	"context"
	"errors"
	"net/http"

	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view"
	viewEmployeesSearch "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/employees-search"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity"
)

const (
	SearchLimitMin = 0
	SearchLimitMax = 50
)

type employeesSearchHandlers struct {
	interactor EmployeesSearchUseCases
	presenter  EmployeesSearchPresenter
	logger     ditzap.Logger
}

func NewEmployeesSearchHandlers(interactor EmployeesSearchUseCases, presenter EmployeesSearchPresenter, logger ditzap.Logger) *employeesSearchHandlers {
	return &employeesSearchHandlers{interactor: interactor, presenter: presenter, logger: logger}
}

func (s employeesSearchHandlers) search(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c.Request.Context(), RequestTimeOut)
	defer cancelCtx()

	session, err := entity.SessionFromContext(ctx)
	if err != nil {
		c.Header(StatusCodeHeader, "ESS_01")
		c.JSON(http.StatusForbidden, view.NewErrorResponse(view.ErrMessageUnauthorized))
		return
	}
	ditzap.WithFields(s.logger, zap.String("login", session.User.GetLogin()), zap.Int("portal_id", session.ActivePortal.GetPortal().ID))

	req := new(viewEmployeesSearch.SearchRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		s.logger.Debug("can't parse SearchRequest to json", zap.Error(err))
		c.Header(StatusCodeHeader, "ESS_02")
		c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		return
	}

	if req.Limit < SearchLimitMin || req.Limit > SearchLimitMax {
		s.logger.Debug("invalid limit in filter", zap.Int("limit", req.Limit))
		c.Header(StatusCodeHeader, "ESS_03")
		c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		return
	}

	reqEntity := s.presenter.SearchRequestToEntity(req)
	if reqEntity == nil {
		s.logger.Debug("empty request")
		c.Header(StatusCodeHeader, "ESS_04")
		c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		return
	}

	res, err := s.interactor.Search(ctx, reqEntity)
	if err != nil {
		var valErr diterrors.ValidationError
		if errors.As(err, &valErr) {
			c.Header(StatusCodeHeader, "ESS_06")
			s.logger.Debug("validation error in search employees", zap.Error(valErr))
			c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
			return
		}
		s.logger.Error("can't search employees", zap.Error(err))
		c.JSON(http.StatusInternalServerError, view.NewErrorResponse(view.ErrMessageInternalError))
		return
	}

	resView := s.presenter.SearchResponseToView(res)
	c.JSON(http.StatusOK, view.NewSuccessResponse(resView))
}

func (s employeesSearchHandlers) filters(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c.Request.Context(), RequestTimeOut)
	defer cancelCtx()

	session, err := entity.SessionFromContext(ctx)
	if err != nil {
		c.Header(StatusCodeHeader, "ESF_01")
		c.JSON(http.StatusForbidden, view.NewErrorResponse(view.ErrMessageUnauthorized))
		return
	}
	ditzap.WithFields(s.logger, zap.String("login", session.User.GetLogin()), zap.Int("portal_id", session.ActivePortal.GetPortal().ID))

	req := new(viewEmployeesSearch.FiltersRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		s.logger.Debug("can't parse FilterRequest to json", zap.Error(err))
		c.Header(StatusCodeHeader, "ESF_02")
		c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		return
	}

	reqEntity := s.presenter.FiltersRequestToEntity(req)
	if reqEntity == nil {
		s.logger.Debug("empty request")
		c.Header(StatusCodeHeader, "ESF_03")
		c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		return
	}

	res, err := s.interactor.Filters(ctx, reqEntity)
	if err != nil {
		var valErr diterrors.ValidationError
		if errors.As(err, &valErr) {
			c.Header(StatusCodeHeader, "ESF_05")
			s.logger.Debug("validation error in search employees", zap.Error(valErr))
			c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
			return
		}
		s.logger.Error("can't search filters", zap.Error(err))
		c.JSON(http.StatusInternalServerError, view.NewErrorResponse(view.ErrMessageInternalError))
		return
	}

	resView := s.presenter.FiltersResponseToView(res)
	c.JSON(http.StatusOK, view.NewSuccessResponse(resView))
}
