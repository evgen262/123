package employees

import (
	"context"
	"fmt"
	"sort"

	entityEmployee "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/employee"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	timeUtils "git.mos.ru/buch-cloud/moscow-team-2.0/build/time-utils.git"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type employeesInteractor struct {
	employeesRepository EmployeesRepository
	// Deprecated: TODO переключить использование сервиса portals на использование portals-facade
	portalsRepository PortalsRepository
	logger            ditzap.Logger
}

// NewEmployeesInteractor TODO переключить использование сервиса portals на использование portals-facade
//
//	Реализовать в рамках задачи https://oblako.mos.ru/jira/browse/GO-1277
func NewEmployeesInteractor(
	employeesRepository EmployeesRepository,
	// Deprecated: TODO переключить использование сервиса portals на использование portals-facade
	portalsRepository PortalsRepository,
	logger ditzap.Logger,
) *employeesInteractor {
	return &employeesInteractor{
		employeesRepository: employeesRepository,
		portalsRepository:   portalsRepository,
		logger:              logger,
	}
}

func (e *employeesInteractor) Get(ctx context.Context, id uuid.UUID) (*entityEmployee.Employee, error) {
	employee, err := e.employeesRepository.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("can't get employee in repositpry: %w", err)
	}

	e.identityEmploymentDate(employee)

	if employee.Portal.ID != 0 {
		portal, err := e.portalsRepository.Get(ctx, employee.Portal.ID)
		if err != nil {
			e.logger.Warn("employeesInteractor.Get: can't get portal: %w", zap.Error(err))
		} else {
			employee.Portal.IconID = portal.LogoUrl
		}
	}

	return employee, nil
}

func (e *employeesInteractor) GetByExtIDAndPortalID(ctx context.Context, extID string, portalID int) (*entityEmployee.Employee, error) {
	employee, err := e.employeesRepository.GetByExtIDAndPortalID(ctx, extID, portalID)
	if err != nil {
		return nil, fmt.Errorf("can't get employee in repository: %w", err)
	}
	e.identityEmploymentDate(employee)

	return employee, nil
}

func (e *employeesInteractor) identityEmploymentDate(employee *entityEmployee.Employee) {
	histories := employee.History

	sort.Slice(histories, func(i, j int) bool {
		if histories[i].EventTime == nil && histories[j].EventTime == nil {
			return false
		}
		if histories[i].EventTime == nil {
			return false
		}
		if histories[j].EventTime == nil {
			return true
		}
		return (*histories[i].EventTime).After(*histories[j].EventTime)
	})

	for _, history := range histories {
		if history == nil {
			continue
		}

		if history.EventType == entityEmployee.OperationTypeHiring {
			if history.EventTime == nil {
				e.logger.Warn("employeesInteractor.identityEmploymentDate: history.EventTime is nil", ditzap.UUID("employee_id", employee.ID), ditzap.UUID("history_id", history.ID))
				continue
			}
			employee.DateOfEmployment = (*timeUtils.Date)(history.EventTime)
			break
		}
	}
}
