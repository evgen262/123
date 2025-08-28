package employees_search

import (
	timeUtils "git.mos.ru/buch-cloud/moscow-team-2.0/build/time-utils.git"
	"github.com/google/uuid"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity"
)

//go:generate ditgen -source=search.go

type SearchResponse struct {
	Employees []*Employee
	Total     int
	AfterID   string
}

type SearchRequest struct {
	Query   string
	Filters *Filters
	Options *Options
	Limit   int
	AfterID *string
}

type Filters struct {
	OIVs          []int
	Organizations []uuid.UUID
	Products      []uuid.UUID
	Subdivisions  []uuid.UUID
	Positions     []string
	Genders       []entity.Gender
	BirthDay      *string
	Absences      []*SearchAbsence
}

type SearchAbsence struct {
	Name string
	From *timeUtils.Date
	To   *timeUtils.Date
}

type Options struct {
	WithFired *FiredDateRange
}

type FiredDateRange struct {
	From *timeUtils.Date `json:"from"`
	To   *timeUtils.Date `json:"to"`
}

type Params struct {
	Names []string
	From  *timeUtils.Date
	To    *timeUtils.Date
}
