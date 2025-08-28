package portals

import (
	portalsv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/portals/portals/v1"
	timeUtils "git.mos.ru/buch-cloud/moscow-team-2.0/build/time-utils.git"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portal"
)

type portalsMapper struct {
	timeUtils timeUtils.TimeUtils
}

func NewPortalsMapper(timeUtils timeUtils.TimeUtils) *portalsMapper {
	return &portalsMapper{
		timeUtils: timeUtils,
	}
}

func (pm portalsMapper) NewPortalsToPb(portals []*portal.Portal) []*portalsv1.AddRequest_Portal {
	portalsPb := make([]*portalsv1.AddRequest_Portal, 0, len(portals))
	for _, portal := range portals {
		p := pm.NewPortalToPb(portal)
		portalsPb = append(portalsPb, p)
	}
	return portalsPb
}

func (pm portalsMapper) NewPortalToPb(portal *portal.Portal) *portalsv1.AddRequest_Portal {
	return &portalsv1.AddRequest_Portal{
		FullName:  portal.FullName,
		ShortName: portal.ShortName,
		Url:       portal.Url,
		LogoUrl:   portal.LogoUrl,
		ChatUrl:   portal.ChatUrl,
		Sort:      int32(portal.Sort),
	}
}

func (pm portalsMapper) PortalsToPb(portals []*portal.Portal) []*portalsv1.Portal {
	portalsPb := make([]*portalsv1.Portal, 0, len(portals))
	for _, portal := range portals {
		p := pm.PortalToPb(portal)
		portalsPb = append(portalsPb, p)
	}
	return portalsPb
}

func (pm portalsMapper) PortalToPb(portal *portal.Portal) *portalsv1.Portal {
	return &portalsv1.Portal{
		Id:            int32(portal.Id),
		FullName:      portal.FullName,
		ShortName:     portal.ShortName,
		Url:           portal.Url,
		LogoUrl:       portal.LogoUrl,
		ChatUrl:       portal.ChatUrl,
		Organizations: pm.PortalOrganizationsToPb(portal.Organizations),
		Sort:          int32(portal.Sort),
		CreatedTime:   pm.timeUtils.TimeToTimestamp(portal.CreatedAt),
		UpdatedTime:   pm.timeUtils.TimeToTimestamp(portal.UpdatedAt),
		DeletedTime:   pm.timeUtils.TimeToTimestamp(portal.DeletedAt),
		IsDeleted:     portal.IsDeleted,
	}
}

func (pm portalsMapper) PortalOrganizationsToPb(orgs []*portal.PortalOrganization) []*portalsv1.Portal_Organization {
	orgsPb := make([]*portalsv1.Portal_Organization, 0, len(orgs))
	for _, org := range orgs {
		orgsPb = append(orgsPb, pm.PortalOrganizationToPb(org))
	}
	return orgsPb
}

func (pm portalsMapper) PortalOrganizationToPb(org *portal.PortalOrganization) *portalsv1.Portal_Organization {
	return &portalsv1.Portal_Organization{
		OrgId: string(org.Id),
		Inn:   org.INN,
	}
}

func (pm portalsMapper) PortalsToEntity(portalsPb []*portalsv1.Portal) []*portal.Portal {
	portals := make([]*portal.Portal, 0, len(portalsPb))

	for _, portalPb := range portalsPb {
		portals = append(portals, pm.PortalToEntity(portalPb))
	}

	return portals
}

func (pm portalsMapper) PortalToEntity(portalPb *portalsv1.Portal) *portal.Portal {
	if portalPb == nil {
		return nil
	}

	return &portal.Portal{
		Id:            portal.PortalID(portalPb.GetId()),
		FullName:      portalPb.GetFullName(),
		ShortName:     portalPb.GetShortName(),
		Url:           portalPb.GetUrl(),
		LogoUrl:       portalPb.GetLogoUrl(),
		ChatUrl:       portalPb.GetChatUrl(),
		Organizations: pm.PortalOrganizationsToEntity(portalPb.GetOrganizations()),
		IsDeleted:     portalPb.GetIsDeleted(),
	}
}

func (pm portalsMapper) PortalOrganizationsToEntity(orgsPb []*portalsv1.Portal_Organization) []*portal.PortalOrganization {
	orgs := make([]*portal.PortalOrganization, 0, len(orgsPb))
	for _, orgPb := range orgsPb {
		orgs = append(orgs, pm.PortalOrganizationToEntity(orgPb))
	}
	return orgs
}

func (pm portalsMapper) PortalOrganizationToEntity(orgPb *portalsv1.Portal_Organization) *portal.PortalOrganization {
	return &portal.PortalOrganization{
		Id:  portal.OrganizationId(orgPb.GetOrgId()),
		INN: orgPb.GetInn(),
	}
}
