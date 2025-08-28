package portals

import (
	"time"
)

type WebPortal struct {
	Name  string
	URL   string
	Image string
}

type Portal struct {
	Id            int                 `json:"id,omitempty"`
	FullName      string              `json:"fullName"`
	ShortName     string              `json:"shortName"`
	Url           string              `json:"url"`
	LogoUrl       string              `json:"logoUrl"`
	ChatUrl       string              `json:"chatUrl"`
	Organizations []*OrganizationInfo `json:"organizations"`
	Sort          int                 `json:"sort,omitempty"`
	CreatedAt     *time.Time          `json:"createdAt,omitempty"`
	UpdatedAt     *time.Time          `json:"updatedAt,omitempty"`
	DeletedAt     *time.Time          `json:"deletedAt,omitempty"`
	Active        bool                `json:"active"`
} // @name Portal

type NewPortal struct {
	FullName  string `json:"fullName"`
	ShortName string `json:"shortName"`
	Url       string `json:"url"`
	LogoUrl   string `json:"logoUrl"`
	ChatUrl   string `json:"chatUrl"`
	Sort      int    `json:"sort,omitempty"`
	Active    bool   `json:"active"`
} // @name NewPortal

type UpdatePortal struct {
	Id        int    `json:"id,omitempty" swaggerignore:"true"`
	FullName  string `json:"fullName"`
	ShortName string `json:"shortName"`
	Url       string `json:"url"`
	LogoUrl   string `json:"logoUrl"`
	ChatUrl   string `json:"chatUrl"`
	Sort      int    `json:"sort,omitempty"`
	Active    bool   `json:"active"`
} // @name UpdatePortal

type PortalInfo struct {
	FullName      string              `json:"fullName"`
	ShortName     string              `json:"shortName"`
	Url           string              `json:"url"`
	LogoUrl       string              `json:"logoUrl"`
	ChatUrl       string              `json:"chatUrl"`
	Organizations []*OrganizationInfo `json:"organizations"`
} // @name PortalInfo

type OrganizationInfo struct {
	ID  string `json:"id"`
	INN string `json:"inn"`
} // @name OrganizationInfo

type PortalsFilterOptions struct {
	PortalIDs   []int    `json:"portal_ids,omitempty"`
	OrgIDs      []string `json:"organization_ids,omitempty"`
	INNs        []string `json:"inns,omitempty"`
	WithDeleted bool     `json:"with_deleted"`
	OnlyLinked  bool     `json:"only_linked"`
} // @name PortalsFilterOptions

func (p *Portal) GetCreatedAt() *time.Time {
	if p == nil {
		return nil
	}
	return p.CreatedAt
}

func (p *Portal) GetUpdatedAt() *time.Time {
	if p == nil {
		return nil
	}
	return p.UpdatedAt
}

func (p *Portal) GetDeletedAt() *time.Time {
	if p == nil {
		return nil
	}
	return p.DeletedAt
}
