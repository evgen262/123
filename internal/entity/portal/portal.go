package portal

import (
	"strconv"
	"strings"
	"time"
)

//go:generate ditgen -source=portal.go

type PortalID int

type PortalIDs []PortalID

func (pi PortalIDs) ToStringSlice() []string {
	if len(pi) == 0 {
		return nil
	}
	stringIds := make([]string, 0, len(pi))
	for _, id := range pi {
		stringIds = append(stringIds, strconv.Itoa(int(id)))
	}
	return stringIds
}

func (pi PortalIDs) ToString() string {
	return strings.Join(pi.ToStringSlice(), ",")
}

type Portal struct {
	Id            PortalID
	FullName      string
	ShortName     string
	Url           string
	LogoUrl       string
	ChatUrl       string
	Sort          int
	Organizations []*PortalOrganization
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
	DeletedAt     *time.Time
	IsDeleted     bool
}

type GetAllOptions struct {
	WithDeleted bool
	OnlyLinked  bool
}

type PortalsFilterOptions struct {
	PortalIDs PortalIDs
	OrgIDs    OrganizationIDs
	INNs      OrganizationINNs

	WithDeleted bool
	OnlyLinked  bool
}

func (pi PortalIDs) ToToInt32Slice() []int32 {
	if len(pi) == 0 {
		return nil
	}
	intIds := make([]int32, 0, len(pi))
	for _, id := range pi {
		intIds = append(intIds, int32(id))
	}
	return intIds
}
