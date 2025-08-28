package employees

import (
	timeUtils "git.mos.ru/buch-cloud/moscow-team-2.0/build/time-utils.git"
)

const (
	GenderInvalidView = "invalid"
	GenderMaleView    = "male"
	GenderFemaleView  = "female"
)

const (
	OperationTypeHiringView            = "Прием на работу"
	OperationTypeTransferView          = "Перемещение"
	OperationTypeDismissalView         = "Увольнение"
	OperationTypeWorkingConditionsView = "Изменение условий труда"
	OperationTypeNameChangeView        = "Изменение фио" // ФИО -> фио для нижнего регистра
	OperationTypeMaternityLeaveView    = "Отпуск по беременности и родам"
	OperationTypeChildcareLeaveView    = "Отпуск по уходу за ребенком"
	OperationTypeInvalidView           = "Неизвестный тип операции"
)

type Employee struct {
	FullName         string          `json:"fullName"`
	Gender           string          `json:"gender"`
	ImageID          string          `json:"photoUrl"`
	Birthday         string          `json:"birthday"`
	Email            string          `json:"email"`
	DateOfEmployment *timeUtils.Date `json:"dateOfEmployment,omitempty"`
	Position         string          `json:"position"`
	Statuses         Statuses        `json:"statuses"`
	OIV              OIV             `json:"oiv"`
	Product          *Product        `json:"product,omitempty"`
	LegalEntity      LegalEntity     `json:"legalEntity"`
	Structure        Structure       `json:"structure"`
	WorkPhone        *string         `json:"workPhone,omitempty"`
	AddPhone         *string         `json:"addPhone,omitempty"`
	MobilePhone      string          `json:"mobilePhone"`
	Workplace        *Workplace      `json:"workplace,omitempty"`
	HeadOfOrg        *OrgHead        `json:"headOfOrgStructure,omitempty"`
	HeadOfManage     *ManageHead     `json:"headOfManageStructure,omitempty"`
}

type Statuses struct {
	IsFired bool     `json:"isFired"`
	Absence *Absence `json:"absence,omitempty"`
}

type Absence struct {
	From *timeUtils.Date `json:"from"`
	To   *timeUtils.Date `json:"to"`
	Type string          `json:"type"`
}

type OIV struct {
	Name   string `json:"name"`
	IconID string `json:"icon"`
}

type Product struct {
	Name   string `json:"name"`
	IconID string `json:"icon"`
}

type LegalEntity struct {
	Name   string `json:"name"`
	IconID string `json:"icon"`
}

type Structure struct {
	Position    string  `json:"position"`
	SubUnit     *string `json:"subUnit,omitempty"`
	LegalEntity string  `json:"legalEntity"`
	OIV         string  `json:"oiv"`
	Product     *string `json:"product,omitempty"`
}

type Workplace struct {
	Address *string `json:"address,omitempty"`
	Floor   *string `json:"floor,omitempty"`
	Cabinet *string `json:"cabinet,omitempty"`
}

type OrgHead struct {
	ID         string  `json:"id"`
	ImageID    *string `json:"photo,omitempty"`
	FirstName  string  `json:"firstName"`
	MiddleName string  `json:"middleName"`
	LastName   string  `json:"lastName"`
	Gender     string  `json:"gender"`
}

type ManageHead struct {
	ID         string  `json:"id"`
	ImageID    *string `json:"photo,omitempty"`
	FirstName  string  `json:"firstName"`
	MiddleName string  `json:"middleName"`
	LastName   string  `json:"lastName"`
	Gender     string  `json:"gender"`
}
