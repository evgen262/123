package employees_search

import (
	"strings"

	"github.com/google/uuid"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity"
)

//go:generate ditgen -source=employee.go

type Employee struct {
	ID           uuid.UUID
	FullName     string
	FirstName    string
	MiddleName   *string
	LastName     string
	Gender       entity.Gender
	ImageID      string
	Position     *Position
	OIV          *OIV
	Product      *Product
	Organization *Organization
	Structure    *Structure
	Statuses     *Statuses
}

type Position struct {
	Name string
}

func (p *Position) GetName() string {
	if p == nil {
		return ""
	}
	return p.Name
}

type OIV struct {
	Name string
}

type Statuses struct {
	IsFired    bool
	IsBirthday bool
	Absences   []*Absence
}

func (o *OIV) GetName() string {
	if o == nil {
		return ""
	}
	return o.Name
}

type Product struct {
	Name string
}

type Organization struct {
	Name string
}

func (o *Organization) GetName() string {
	if o == nil {
		return ""
	}
	return o.Name
}

type Subdivision struct {
	Name string
}

func (s *Subdivision) GetName() string {
	if s == nil {
		return ""
	}
	return s.Name
}

type Structure struct {
	Position     *Position
	Subdivision  *Subdivision
	Organization *Organization
	OIV          *OIV
}

type AbsenceType string // Тип причины отсутствия сотрудника

func (r AbsenceType) String() string {
	return strings.ToLower(string(r))
}

const (
	AbsenceTypeDecree       AbsenceType = "декрет"
	AbsenceTypeDecreeWork   AbsenceType = "работа в декрете"
	AbsenceTypeVacation     AbsenceType = "отпуск"
	AbsenceTypeBusinessTrip AbsenceType = "командировка"
	AbsenceTypeMedical      AbsenceType = "больничный"
)
