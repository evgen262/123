package employees_search

import (
	"github.com/google/uuid"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity"
)

//go:generate ditgen -source=params.go

type SearchParams struct {
	Query   string
	Filters *FiltersParams
	Options *OptionsParams
	Limit   int
	AfterID *string
}

type OptionsParams struct {
	WithFired *FiredDateRange
}

func (op *OptionsParams) ToOptions() *Options {
	if op == nil {
		return nil
	}

	o := &Options{}
	if op.WithFired != nil {
		o.WithFired = &FiredDateRange{
			From: op.GetWithFired().GetFrom(),
			To:   op.GetWithFired().GetTo(),
		}
	}

	return o
}

type FiltersParams struct {
	OIVs          []int
	Organizations []uuid.UUID
	Products      []uuid.UUID
	Subdivisions  []uuid.UUID
	Positions     []string
	Genders       []entity.Gender
	Statuses      FilterStatuses
}

func (fp *FiltersParams) ToFilters() *Filters {
	if fp == nil {
		return nil
	}

	return &Filters{
		OIVs:          fp.OIVs,
		Organizations: fp.Organizations,
		Products:      fp.Products,
		Subdivisions:  fp.Subdivisions,
		Positions:     fp.Positions,
		Genders:       fp.Genders,
	}
}

type FilterStatuses struct {
	IsBirthDay       bool
	IsMaternityLeave bool // Декретный отпуск
	IsVacation       bool // Отпуск
}
