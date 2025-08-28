package portalsv2

import (
	viewPortalsV2 "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/portalsv2"
	entityPortalsV2 "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portalv2"
)

type complexesPresenter struct {
}

func NewComplexesPresenter() *complexesPresenter {
	return &complexesPresenter{}
}

func (p complexesPresenter) ComplexesToView(complexes []*entityPortalsV2.Complex) []*viewPortalsV2.Complex {
	if complexes == nil {
		return nil
	}
	viewComplexes := make([]*viewPortalsV2.Complex, 0, len(complexes))
	for _, c := range complexes {
		vc := &viewPortalsV2.Complex{
			ID:    c.ID,
			Sort:  c.Sort,
			Group: c.ComplexGroup,
		}

		if c.Responsible != nil {
			p.fillResponsibleData(c, vc)
		}

		vc.PortalIDs = make([]int, 0, len(c.Portals))
		for _, portal := range c.Portals {
			vc.PortalIDs = append(vc.PortalIDs, portal.ID)
		}

		viewComplexes = append(viewComplexes, vc)
	}
	return viewComplexes
}

func (p complexesPresenter) fillResponsibleData(c *entityPortalsV2.Complex, vc *viewPortalsV2.Complex) {
	vc.FirstName = c.Responsible.FirstName
	vc.LastName = c.Responsible.LastName

	if c.Responsible.MiddleName != nil {
		vc.MiddleName = *c.Responsible.MiddleName
	} else {
		vc.MiddleName = ""
	}

	if c.Responsible.ImageID != nil {
		vc.IconID = *c.Responsible.ImageID
	} else {
		vc.IconID = ""
	}

	vc.HeadDescription = c.Responsible.Description
}
