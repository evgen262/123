package portals

import (
	"time"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view"
)

//go:generate ditgen -source=organization.go

type OrganizationId string // @name OrganizationId

type OrganizationIds []OrganizationId // @name OrganizationIds

type Organization struct {
	Id              string              `json:"id,omitempty"`
	FullName        string              `json:"full_name"`
	ShortName       string              `json:"short_name"`
	RegCode         string              `json:"reg_code"`
	OrgCode         string              `json:"org_code"`
	UNKCode         string              `json:"unk_code"`
	OGRN            string              `json:"ogrn"`
	INN             string              `json:"inn"`
	KPP             string              `json:"kpp"`
	OrgTypeName     string              `json:"org_type_name"`
	OrgTypeCode     string              `json:"org_type_code"`
	UchrezhTypeCode string              `json:"uchrezh_type_code"`
	Grbs            []*OrganizationGrbs `json:"grbs"`
	Email           string              `json:"email"`
	Phone           string              `json:"phone"`
	AdditionalPhone string              `json:"additional_phone"`
	Site            string              `json:"site"`
	IsLiquidated    bool                `json:"is_liquidated"`
	LiquidatedAt    *string             `json:"liquidated_at,omitempty"`
	CreatedAt       *time.Time          `json:"created_at,omitempty"`
	UpdatedAt       *time.Time          `json:"updated_at,omitempty"`
} // @name Organization

// OrganizationGrbs ГРБС организации.
type OrganizationGrbs struct {
	GrbsId    string  `json:"grbs_id"`
	Name      string  `json:"name"`
	Inn       string  `json:"inn"`
	Ogrn      string  `json:"ogrn"`
	StartDate *string `json:"start_date"`
} // @name OrganizationGrbs

type OrganizationLink struct {
	// PortalId ИД портала для привязки
	PortalId int `json:"portal_id"`
	// Organizations список ИД организаций для привязки к порталу
	OrganizationIds OrganizationIds `json:"organization_ids"`
} // @name OrganizationLink

// OrganizationsWithPagination организации с пагинацией
type OrganizationsWithPagination struct {
	Pagination    *view.StringPagination `json:"pagination,omitempty"`
	Organizations []*Organization        `json:"organizations,omitempty"`
} // @name OrganizationsWithPagination

// FilterOrganizations структура объединяющая опции фильтрации и пагинацию.
type FilterOrganizations struct {
	Filters    OrganizationsFilters        `json:"filters,omitempty"`
	Pagination *view.StringPagination      `json:"pagination,omitempty"`
	Options    *OrganizationsFilterOptions `json:"options,omitempty"`
} // @name FilterOrganizations

// OrganizationsFilters опции фильтрации
type OrganizationsFilters struct {
	Ids   []string `json:"ids,omitempty"`
	Names []string `json:"names,omitempty"`
	Ogrns []string `json:"ogrns,omitempty"`
	Inns  []string `json:"inns,omitempty"`
} // @name OrganizationsFilters

// OrganizationsFilterOptions опции фильтрации
type OrganizationsFilterOptions struct {
	WithLiquidated bool `json:"with_liquidated"`
} // @name OrganizationsFilterOptions
