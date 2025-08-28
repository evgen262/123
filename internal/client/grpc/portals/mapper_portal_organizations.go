package portals

import (
	"fmt"

	organizationsv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/portals/organizations/v1"
	sharedv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/portals/shared/v1"
	timeUtils "git.mos.ru/buch-cloud/moscow-team-2.0/build/time-utils.git"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portal"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity"
)

type organizationsMapper struct {
	timeUtils timeUtils.TimeUtils
}

func (pfm organizationsMapper) PaginationToEntity(paginationPb *sharedv1.PaginationResponse) *entity.StringPagination {
	if paginationPb == nil {
		return nil
	}
	var (
		limit *int
		total *int
	)
	if paginationPb.GetLimit() != 0 {
		l := int(paginationPb.GetLimit())
		limit = &l
	}
	if paginationPb.GetTotal() != 0 {
		t := int(paginationPb.GetTotal())
		total = &t
	}
	return &entity.StringPagination{
		Limit:           limit,
		LastId:          portal.OrganizationId(paginationPb.GetLastId()),
		LastCreatedTime: pfm.timeUtils.TimestampToTime(paginationPb.GetLastCreatedTime()),
		Total:           total,
	}
}

func (pfm organizationsMapper) PaginationToPb(pagination *entity.StringPagination) *sharedv1.PaginationRequest {
	if pagination == nil {
		return nil
	}
	var (
		limit uint32
	)
	if pagination.Limit != nil {
		limit = uint32(*pagination.Limit)
	}
	return &sharedv1.PaginationRequest{
		Limit:           limit,
		LastId:          pagination.LastId.String(),
		LastCreatedTime: pfm.timeUtils.TimeToTimestamp(pagination.GetLastCreatedTime()),
	}
}

func NewOrganizationsMapper(timeUtils timeUtils.TimeUtils) *organizationsMapper {
	return &organizationsMapper{
		timeUtils: timeUtils,
	}
}

func (pfm organizationsMapper) OnceGrbsToEntity(onceGrbsPb *organizationsv1.Organization_Grbs) *portal.OrganizationGrbs {
	orgGrbs := &portal.OrganizationGrbs{
		GrbsId: fmt.Sprintf("%s_%s", onceGrbsPb.GetOgrn(), onceGrbsPb.GetInn()),
		Name:   onceGrbsPb.GetName(),
		Inn:    onceGrbsPb.GetInn(),
		Ogrn:   onceGrbsPb.GetOgrn(),
	}
	if onceGrbsPb.GetStartDate() != nil {
		value := onceGrbsPb.GetStartDate().GetValue()
		orgGrbs.StartDate = &value
	}

	return orgGrbs
}

func (pfm organizationsMapper) GrbsToEntity(grbsPb []*organizationsv1.Organization_Grbs) []*portal.OrganizationGrbs {
	grbs := make([]*portal.OrganizationGrbs, 0, len(grbsPb))

	for _, onceGrbs := range grbsPb {
		grbs = append(grbs, pfm.OnceGrbsToEntity(onceGrbs))
	}

	return grbs
}

func (pfm organizationsMapper) OrganizationToEntity(organizationPb *organizationsv1.Organization) *portal.Organization {
	org := &portal.Organization{
		Id:              portal.OrganizationId(organizationPb.GetId()),
		FullName:        organizationPb.GetFullName(),
		ShortName:       organizationPb.GetShortName(),
		RegCode:         organizationPb.GetRegCode(),
		OrgCode:         organizationPb.GetOrgCode(),
		UNKCode:         organizationPb.GetUnkCode(),
		OGRN:            organizationPb.GetOgrn(),
		INN:             portal.OrganizationINN(organizationPb.GetInn()),
		KPP:             organizationPb.GetKpp(),
		OrgTypeName:     organizationPb.GetOrgTypeName(),
		OrgTypeCode:     organizationPb.GetOrgTypeCode(),
		UchrezhTypeCode: organizationPb.GetUchrezhTypeCode(),
		Grbs:            pfm.GrbsToEntity(organizationPb.GetGrbs()),
		Email:           organizationPb.GetEmail(),
		Phone:           organizationPb.GetPhone(),
		AdditionalPhone: organizationPb.GetAdditionalPhone(),
		Site:            organizationPb.GetSite(),
		CreatedAt:       pfm.timeUtils.TimestampToTime(organizationPb.GetCreatedTime()),
		UpdatedAt:       pfm.timeUtils.TimestampToTime(organizationPb.GetUpdatedTime()),
	}
	if organizationPb.GetLiquidatedAt() != nil {
		value := organizationPb.GetLiquidatedAt().GetValue()
		org.LiquidatedAt = &value
		org.IsLiquidated = organizationPb.GetIsLiquidated()
	}

	return org
}

func (pfm organizationsMapper) OrganizationsToEntity(organizationsPb []*organizationsv1.Organization) []*portal.Organization {
	organizations := make([]*portal.Organization, 0, len(organizationsPb))

	for _, organizationPb := range organizationsPb {
		organizations = append(organizations, pfm.OrganizationToEntity(organizationPb))
	}

	return organizations
}

func (pfm organizationsMapper) OptionsToPb(options portal.OrganizationsFilterOptions) *organizationsv1.OrganizationFilterOptions {
	return &organizationsv1.OrganizationFilterOptions{WithLiquidated: options.WithLiquidated}
}
