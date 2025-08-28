package http

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity"
	entityEmployee "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/employee"
)

type employeesHandlers struct {
	employeesInteractor EmployeesUseCases
	employeesPresenter  EmployeesPresenter
	logger              ditzap.Logger
}

func NewEmployeesHandlers(
	employeesInteractor EmployeesUseCases,
	employeesPresenter EmployeesPresenter,
	logger ditzap.Logger,
) *employeesHandlers {
	return &employeesHandlers{
		employeesInteractor: employeesInteractor,
		employeesPresenter:  employeesPresenter,
		logger:              logger,
	}
}

func (e *employeesHandlers) getEmployee(c *gin.Context) {
	sID := strings.TrimSpace(c.Param("id"))
	sPortalID := c.Query("portalId")

	if sID == "" {
		e.logger.Debug("employeesHandlers.getEmployee: ID is empty")
		c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		return
	}

	var (
		employee *entityEmployee.Employee
		err      error
	)

	// Если указан ID портала, то запрашивается сотрудник по extID и portalID
	if sPortalID != "" {
		portalID, pErr := strconv.Atoi(sPortalID)
		if pErr != nil {
			e.logger.Debug("employeesHandlers.getEmployee: portalID is invalid", zap.String("portal_id", sPortalID))
			c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
			return
		}

		employee, err = e.employeesInteractor.GetByExtIDAndPortalID(c, sID, portalID)
		if err != nil {
			switch {
			case errors.Is(err, diterrors.ErrNotFound):
				e.logger.Debug("employeesHandlers.getEmployee: employee not found", zap.String("ext_id", sID), zap.Int("portal_id", portalID))
				c.JSON(http.StatusNotFound, view.NewErrorResponse(view.ErrMessageNotFound))
				return
			}
			e.logger.Error("employeesHandlers.getEmployee: can't get employee", zap.Error(err), zap.String("ext_id", sID), zap.Int("portal_id", portalID))
			c.JSON(http.StatusInternalServerError, view.NewErrorResponse(view.ErrMessageInternalError))
			return
		}
	} else {
		id, pErr := uuid.Parse(sID)
		if pErr != nil {
			e.logger.Debug("employeesHandlers.getEmployee: id is invalid", zap.String("id", sID))
			c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
			return
		}

		employee, err = e.employeesInteractor.Get(c, id)
		if err != nil {
			switch {
			case errors.Is(err, diterrors.ErrNotFound):
				e.logger.Debug("employeesHandlers.getEmployee: employee not found", zap.String("id", sID))
				c.JSON(http.StatusNotFound, view.NewErrorResponse(view.ErrMessageNotFound))
				return
			}
			e.logger.Error("employeesHandlers.getEmployee: can't get employee", zap.Error(err), zap.String("id", sID))
			c.JSON(http.StatusInternalServerError, view.NewErrorResponse(view.ErrMessageInternalError))
			return
		}
	}

	c.JSON(http.StatusOK, e.employeesPresenter.EmployeeToView(employee))
}

func (e *employeesHandlers) getProfile(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c.Request.Context(), RequestTimeOut)
	defer cancelCtx()

	session, err := entity.SessionFromContext(ctx)
	if err != nil {
		e.logger.Error("employeesHandlers.getProfile: can't get session", zap.Error(err))
		c.JSON(http.StatusUnauthorized, view.NewErrorResponse(view.ErrMessageUnauthorized))
		return
	}
	extID := session.User.GetEmployee().GetExtID()
	portalID := session.ActivePortal.GetPortal().ID

	employee, err := e.employeesInteractor.GetByExtIDAndPortalID(c, extID, portalID)
	if err != nil {
		switch {
		case errors.Is(err, diterrors.ErrNotFound):
			e.logger.Debug("employeesHandlers.getProfile: employee not found", zap.String("session_id", session.ID.String()), zap.String("ext_id", extID), zap.Int("portal_id", portalID))
			c.JSON(http.StatusNotFound, view.NewErrorResponse(view.ErrMessageNotFound))
			return
		}
		e.logger.Error("employeesHandlers.getProfile: can't get employee", zap.Error(err), zap.String("session_id", session.ID.String()), zap.String("ext_id", extID), zap.Int("portal_id", portalID))
		c.JSON(http.StatusInternalServerError, view.NewErrorResponse(view.ErrMessageInternalError))
		return
	}
	c.JSON(http.StatusOK, e.employeesPresenter.EmployeeToView(employee))
}
