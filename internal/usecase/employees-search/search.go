package employees_search

import (
	"context"
	"fmt"
	"time"

	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	timeUtils "git.mos.ru/buch-cloud/moscow-team-2.0/build/time-utils.git"

	entityEmployeesSearch "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/employees-search"
)

// Период выборки статусов
const statusPeriod = 2 * (time.Hour * 24 * 7) // 2 недели

type employeesSearchInteractor struct {
	employeesSearchRepo EmployeesSearchRepository
	tu                  timeUtils.TimeUtils
}

func NewEmployeesSearchInteractor(employeesSearchRepo EmployeesSearchRepository, tu timeUtils.TimeUtils) *employeesSearchInteractor {
	return &employeesSearchInteractor{
		employeesSearchRepo: employeesSearchRepo,
		tu:                  tu,
	}
}

func (e employeesSearchInteractor) Search(ctx context.Context, params *entityEmployeesSearch.SearchParams) (*entityEmployeesSearch.SearchResponse, error) {
	if params == nil {
		return nil, diterrors.NewValidationError(ErrIsEmpty)
	}

	request := &entityEmployeesSearch.SearchRequest{
		Query:   params.Query,
		Options: params.GetOptions().ToOptions(),
		Limit:   params.Limit,
		AfterID: params.AfterID,
	}

	// если указан статус день рождения, то берём текущую дату на сервере (по ФТ)
	if params.GetFilters() != nil {
		now := e.tu.New()
		nowDate := (*timeUtils.Date)(now)
		request.Filters = params.GetFilters().ToFilters()

		if params.GetFilters().Statuses.IsBirthDay {
			birthday := now.Format("01-02")

			request.Filters.BirthDay = &birthday
		}

		periodDate := (timeUtils.Date)(now.Add(+statusPeriod))
		if params.GetFilters().Statuses.IsVacation {
			request.Filters.Absences = append(request.Filters.Absences, &entityEmployeesSearch.SearchAbsence{
				Name: entityEmployeesSearch.AbsenceTypeVacation.String(),
				From: nowDate,
				To:   &periodDate,
			})
		}

		if params.GetFilters().Statuses.IsMaternityLeave {
			request.Filters.Absences = append(request.Filters.Absences, &entityEmployeesSearch.SearchAbsence{
				Name: entityEmployeesSearch.AbsenceTypeDecree.String(),
				From: nowDate,
				To:   &periodDate,
			})
		}
	}

	res, err := e.employeesSearchRepo.Search(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("can't search employees in repository: %w", err)
	}

	// TODO: после реализации сервиса настроек необходимо добавить скрытие полей у физлиц согласно контракту
	return res, nil
}

func (e employeesSearchInteractor) Filters(ctx context.Context, params *entityEmployeesSearch.SearchParams) (*entityEmployeesSearch.FiltersResponse, error) {
	if params == nil {
		return nil, diterrors.NewValidationError(ErrIsEmpty)
	}

	now := e.tu.New()
	nowDate := (*timeUtils.Date)(now)
	periodDate := (timeUtils.Date)(now.Add(+statusPeriod))

	filters := params.GetFilters().ToFilters()
	if params.GetFilters() != nil {
		if params.GetFilters().Statuses.IsBirthDay {
			bd := now.Format("01-02")
			filters.BirthDay = &bd
		}

		if params.GetFilters().Statuses.IsVacation {
			filters.Absences = append(filters.Absences, &entityEmployeesSearch.SearchAbsence{
				Name: entityEmployeesSearch.AbsenceTypeVacation.String(),
				From: nowDate,
				To:   &periodDate,
			})
		}

		if params.GetFilters().Statuses.IsMaternityLeave {
			filters.Absences = append(filters.Absences, &entityEmployeesSearch.SearchAbsence{
				Name: entityEmployeesSearch.AbsenceTypeDecree.String(),
				From: nowDate,
				To:   &periodDate,
			})
		}
	}

	request := &entityEmployeesSearch.FiltersRequest{
		Filters: filters,
		Options: params.GetOptions().ToOptions(),
		Params: entityEmployeesSearch.Params{
			Names: []string{
				entityEmployeesSearch.AbsenceTypeDecree.String(),
				entityEmployeesSearch.AbsenceTypeVacation.String(),
				entityEmployeesSearch.AbsenceTypeMedical.String(),
			},
			From: nowDate,
			To:   &periodDate,
		},
	}

	res, err := e.employeesSearchRepo.Filters(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("can't search filters in repository: %w", err)
	}
	haveAbsence := make([]*entityEmployeesSearch.FilterAbsence, 0, len(res.Absences))
	for _, a := range res.Absences {
		if a.Count > 0 {
			haveAbsence = append(haveAbsence, a)
		}
	}
	res.Absences = haveAbsence
	return res, nil
}
