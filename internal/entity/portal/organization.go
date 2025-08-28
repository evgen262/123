package portal

import (
	"strings"
	"time"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity"
)

//go:generate ditgen -source=organization.go

type OrganizationId string

func (oid OrganizationId) String() string {
	return string(oid)
}

type PortalOrganization struct {
	Id  OrganizationId
	INN string
}

type OrganizationINN string

func (i OrganizationINN) String() string {
	return string(i)
}

type OrganizationINNs []OrganizationINN

func (oinns OrganizationINNs) ToStringSlice() []string {
	if len(oinns) == 0 {
		return nil
	}
	s := make([]string, 0, len(oinns))
	for _, inn := range oinns {
		s = append(s, inn.String())
	}
	return s
}

func (oinns OrganizationINNs) ToString() string {
	return strings.Join(oinns.ToStringSlice(), ",")
}

// Organization Организация
type Organization struct {
	Id              OrganizationId
	FullName        string
	ShortName       string
	RegCode         string
	OrgCode         string
	UNKCode         string
	OGRN            string
	INN             OrganizationINN
	KPP             string
	OrgTypeName     string
	OrgTypeCode     string
	UchrezhTypeCode string
	Grbs            []*OrganizationGrbs
	Email           string
	Phone           string
	AdditionalPhone string
	Site            string
	IsLiquidated    bool
	LiquidatedAt    *string
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
}

// OrganizationGrbs ГРБС организации.
type OrganizationGrbs struct {
	GrbsId    string
	Name      string
	Inn       string
	Ogrn      string
	StartDate *string
}

type OrganizationIDs []OrganizationId

func (ois OrganizationIDs) ToStringSlice() []string {
	if len(ois) == 0 {
		return nil
	}
	s := make([]string, 0, len(ois))
	for _, id := range ois {
		s = append(s, string(id))
	}
	return s
}

func (ois OrganizationIDs) ToString() string {
	return strings.Join(ois.ToStringSlice(), ",")
}

// OrganizationsWithPagination организации с пагинацией
type OrganizationsWithPagination struct {
	Pagination    *entity.StringPagination
	Organizations []*Organization
}

// OrganizationsFilters опции фильтрации организаций.
type OrganizationsFilters struct {
	Ids   []string
	Names []string
	Ogrns []string
	Inns  []string
}

// OrganizationsFilterOptions опции запроса фильтрации.
type OrganizationsFilterOptions struct {
	WithLiquidated bool
}
