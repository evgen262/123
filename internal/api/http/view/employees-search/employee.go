package employees_search

import (
	employees_search "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/employees-search"
)

const (
	GenderInvalidView = "invalid"
	GenderMaleView    = "male"
	GenderFemaleView  = "female"
)

const (
	GenderInvalidStringName = "Пол не указан"
	GenderMaleStringName    = "Мужчины"
	GenderFemaleStringName  = "Женщины"
)

type Employee struct {
	ID           string        `json:"id"`
	FullName     string        `json:"fullName"`
	FirstName    string        `json:"firstName"`
	MiddleName   *string       `json:"middleName"`
	LastName     string        `json:"lastName"`
	Gender       string        `json:"gender"`
	ImageID      string        `json:"photoUrl,omitempty"`
	Position     string        `json:"position"`
	OIV          *OIV          `json:"oiv"`
	Product      *Product      `json:"product,omitempty"`
	Organization *Organization `json:"legalEntity"`
	Structure    *Structure    `json:"structure"`
	Statuses     *Statuses     `json:"statuses"`
}

type OIV struct {
	Name   string `json:"name"`
	IconID string `json:"icon"`
}

type Product struct {
	Name   string `json:"name"`
	IconID string `json:"icon"`
}

type Organization struct {
	Name   string `json:"name"`
	IconID string `json:"icon"`
}

type Structure struct {
	Position     string `json:"position"`
	Subdivision  string `json:"subUnit"`
	Organization string `json:"legalEntity"`
	OIV          string `json:"oiv"`
}

type Statuses struct {
	IsCOVID19Vaccinated bool                        `json:"isCOVID19Vaccinated"`
	IsFired             bool                        `json:"isFired"`
	IsBirthday          bool                        `json:"isBirthday"`
	Absences            []*employees_search.Absence `json:"absences,omitempty"`
}
