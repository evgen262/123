package portals

import (
	viewPortals "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/portals"
	entityPortal "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portal"
)

type portalsPresenter struct {
}

func NewPortalsPresenter() *portalsPresenter {
	return &portalsPresenter{}
}

func (p portalsPresenter) ToNewEntities(portalsArr []*viewPortals.NewPortal) []*entityPortal.Portal {
	resPortals := make([]*entityPortal.Portal, 0, len(portalsArr))
	for _, portal := range portalsArr {
		resPortals = append(resPortals, p.ToNewEntity(portal))
	}

	return resPortals
}

func (p portalsPresenter) ToNewEntity(portal *viewPortals.NewPortal) *entityPortal.Portal {
	return &entityPortal.Portal{
		FullName:  portal.FullName,
		ShortName: portal.ShortName,
		Url:       portal.Url,
		LogoUrl:   portal.LogoUrl,
		ChatUrl:   portal.ChatUrl,
		Sort:      portal.Sort,
		IsDeleted: !portal.Active,
	}
}

func (p portalsPresenter) ToEntities(portalsArr []*viewPortals.UpdatePortal) []*entityPortal.Portal {
	resPortals := make([]*entityPortal.Portal, 0, len(portalsArr))
	for _, portal := range portalsArr {
		resPortals = append(resPortals, p.ToEntity(portal))
	}

	return resPortals
}

func (p portalsPresenter) ToEntity(portal *viewPortals.UpdatePortal) *entityPortal.Portal {
	return &entityPortal.Portal{
		Id:        entityPortal.PortalID(portal.Id),
		FullName:  portal.FullName,
		ShortName: portal.ShortName,
		Url:       portal.Url,
		LogoUrl:   portal.LogoUrl,
		ChatUrl:   portal.ChatUrl,
		Sort:      portal.Sort,
		IsDeleted: !portal.Active,
	}
}

func (p portalsPresenter) ToViews(portals []*entityPortal.Portal) []*viewPortals.Portal {
	resPortals := make([]*viewPortals.Portal, 0, len(portals))
	for _, portal := range portals {
		resPortals = append(resPortals, p.ToView(portal))
	}

	return resPortals
}

func (p portalsPresenter) ToView(portal *entityPortal.Portal) *viewPortals.Portal {
	return &viewPortals.Portal{
		Id:            int(portal.Id),
		FullName:      portal.FullName,
		ShortName:     portal.ShortName,
		Url:           portal.Url,
		LogoUrl:       portal.LogoUrl,
		ChatUrl:       portal.ChatUrl,
		Sort:          portal.Sort,
		CreatedAt:     portal.CreatedAt,
		UpdatedAt:     portal.UpdatedAt,
		DeletedAt:     portal.DeletedAt,
		Active:        !portal.IsDeleted,
		Organizations: p.organizationsToView(portal.Organizations),
	}
}

func (p portalsPresenter) organizationsToView(organizations []*entityPortal.PortalOrganization) []*viewPortals.OrganizationInfo {
	result := make([]*viewPortals.OrganizationInfo, 0, len(organizations))
	for _, organization := range organizations {
		result = append(result, &viewPortals.OrganizationInfo{
			ID:  string(organization.Id),
			INN: organization.INN,
		})
	}

	return result
}

func (p portalsPresenter) ToShortViews(portals []*entityPortal.Portal) []*viewPortals.PortalInfo {
	resPortals := make([]*viewPortals.PortalInfo, 0, len(portals))
	for _, portal := range portals {
		resPortals = append(resPortals, p.ToShortView(portal))
	}

	return resPortals
}

func (p portalsPresenter) ToShortView(portal *entityPortal.Portal) *viewPortals.PortalInfo {
	return &viewPortals.PortalInfo{
		FullName:      portal.FullName,
		ShortName:     portal.ShortName,
		Url:           portal.Url,
		LogoUrl:       portal.LogoUrl,
		ChatUrl:       portal.ChatUrl,
		Organizations: p.organizationsToView(portal.Organizations),
	}
}

func (p portalsPresenter) ToWebView(portal *entityPortal.Portal) *viewPortals.WebPortal {
	return &viewPortals.WebPortal{
		Name:  portal.FullName,
		URL:   portal.Url,
		Image: portal.LogoUrl,
	}
}

func (p portalsPresenter) ToWebViews(portals []*entityPortal.Portal) []*viewPortals.WebPortal {
	ps := make([]*viewPortals.WebPortal, 0, len(portals))

	for _, portal := range portals {
		pl := p.ToWebView(portal)
		ps = append(ps, pl)
	}

	return ps
}

func (p portalsPresenter) FilterOptionsToEntity(options viewPortals.PortalsFilterOptions) entityPortal.PortalsFilterOptions {
	portalIDs := make(entityPortal.PortalIDs, 0, len(options.PortalIDs))
	for _, portalID := range options.PortalIDs {
		portalIDs = append(portalIDs, entityPortal.PortalID(portalID))
	}
	orgIDs := make(entityPortal.OrganizationIDs, 0, len(options.OrgIDs))
	for _, orgID := range options.OrgIDs {
		orgIDs = append(orgIDs, entityPortal.OrganizationId(orgID))
	}
	inns := make(entityPortal.OrganizationINNs, 0, len(options.INNs))
	for _, inn := range options.INNs {
		inns = append(inns, entityPortal.OrganizationINN(inn))
	}

	return entityPortal.PortalsFilterOptions{
		PortalIDs:   portalIDs,
		OrgIDs:      orgIDs,
		INNs:        inns,
		WithDeleted: options.WithDeleted,
		OnlyLinked:  options.OnlyLinked,
	}
}
