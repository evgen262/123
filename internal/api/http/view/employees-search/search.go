package employees_search

import "github.com/google/uuid"

type SearchRequest struct {
	Query   string   `json:"query,omitempty"`
	Filters *Filters `json:"filters,omitempty"`
	Options *Options `json:"options,omitempty"`
	Limit   int      `json:"limit"`
	AfterID *string  `json:"afterId,omitempty"`
}

type SearchResponse struct {
	Employees []*Employee `json:"data"`
	Total     int         `json:"total"`
	AfterID   string      `json:"afterId"`
}

type FiltersRequest struct {
	Filters *Filters `json:"filters"`
	Options *Options `json:"options"`
}

type FiltersResponse struct {
	Filters FiltersResult `json:"filters"`
	Options Options       `json:"options"`
}

type FiltersResult struct {
	OIVs          []*FilterOIV          `json:"oiv,omitempty"`
	Organizations []*FilterOrganization `json:"legalEntities,omitempty"`
	Products      []*FilterProduct      `json:"products,omitempty"`
	Subdivisions  []*FilterSubdivision  `json:"subUnits,omitempty"`
	Positions     []*FilterPosition     `json:"positions,omitempty"`
	Genders       []*FilterGender       `json:"genders,omitempty"`
	Statuses      []string              `json:"statuses,omitempty"`
	// TODO: добавить адреса офисов и этажи после того как они будут добавлены в employees-search
	// addresses {id, name} - адрес офиса
	// locations {id, name} - этажи офисов
}

type FilterOIV struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type FilterOrganization struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type FilterProduct struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type FilterSubdivision struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type FilterPosition struct {
	// ID   uuid.UUID `json:"id"`
	Name string `json:"name"`
}

type FilterGender struct {
	Name string `json:"name"`
	ID   string `json:"id"`
	// IsDisabled bool   `json:"isDisabled"`
}
