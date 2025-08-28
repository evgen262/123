package employees_search

import (
	"strings"

	"github.com/google/uuid"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity"
)

//go:generate ditgen -source=filter.go

type FiltersRequest struct {
	Filters *Filters
	Options *Options
	Params  Params
}

type FiltersResponse struct {
	OIVs          []*FilterOIV
	Organizations []*FilterOrganization
	Products      []*FilterProduct
	Subdivisions  []*FilterSubdivision
	Positions     []*FilterPosition
	Genders       []*FilterGender
	BirthDayCount int
	Absences      []*FilterAbsence
}

type FilterOIV struct {
	ID   int
	Name string
}

type FilterOrganization struct {
	ID   uuid.UUID
	Name string
}

type FilterProduct struct {
	ID   uuid.UUID
	Name string
}

type FilterSubdivision struct {
	ID   uuid.UUID
	Name string
}

type FilterPosition struct {
	ID   uuid.UUID
	Name string
}

type FilterGender struct {
	Gender     entity.Gender
	IsDisabled bool
}

type FilterAbsence struct {
	// Наименование причины отсутствия
	Name string
	// Тип причины отсутствия
	_type AbsenceType
	// Количество сотрудников
	Count int
}

func (f FilterAbsence) Type() AbsenceType {
	return f._type
}

func MakeFilterAbsence(name string, count int) *FilterAbsence {
	a := &FilterAbsence{
		Name:  name,
		Count: count,
	}
	switch strings.ToLower(strings.TrimSpace(name)) {
	case AbsenceTypeDecree.String():
		a._type = AbsenceTypeDecree
	case AbsenceTypeDecreeWork.String():
		a._type = AbsenceTypeDecreeWork
	case AbsenceTypeVacation.String():
		a._type = AbsenceTypeVacation
	case AbsenceTypeBusinessTrip.String():
		a._type = AbsenceTypeBusinessTrip
	case AbsenceTypeMedical.String():
		a._type = AbsenceTypeMedical
	}
	return a
}
