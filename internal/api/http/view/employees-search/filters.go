package employees_search

import (
	"github.com/google/uuid"
)

const (
	StatusTypeVacation       = "vacation"       // Отпуск
	StatusTypeMaternityLeave = "maternityLeave" // Декретный отпуск
	StatusTypeBirthday       = "birthday"       // День рождения
)

type Filters struct {
	OIVs          []int       `json:"oiv,omitempty"`
	Organizations []uuid.UUID `json:"legalEntities,omitempty"`
	Products      []uuid.UUID `json:"products,omitempty"`
	Subdivisions  []uuid.UUID `json:"subUnits,omitempty"`
	Positions     []string    `json:"positions,omitempty"`
	Gender        string      `json:"genders,omitempty"`
	Status        string      `json:"statuses,omitempty"`

	// TODO: добавить адреса офисов и этажи после того как они будут добавлены в employees-search
	// Addresses string - Адреса (уточнить по контракту) почему указан 1 адрес тип - string
	// Locations string - id Этажа (уточнить по контракту) почему указан 1 этаж тип - string
	// TODO: добавить фильтр по дате приёма на работу после того как они будут добавлены в employees-search
	// EmployedDateRange { from, to } -диапазон дат приёма на работу
}
