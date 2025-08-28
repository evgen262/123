package portalsv2

import (
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/portalsv2"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portalv2"
)

type portalsPresenter struct {
}

func NewPortalsPresenter() *portalsPresenter {
	return &portalsPresenter{}
}

func (p portalsPresenter) PortalsWithCountToView(portalsWithCounts []*portalv2.PortalWithCounts) []*portalsv2.Portal {
	result := make([]*portalsv2.Portal, 0, len(portalsWithCounts))
	for _, portalWithCounts := range portalsWithCounts {
		if portalWithCounts == nil || portalWithCounts.Portal == nil {
			continue
		}
		result = append(result, p.portalWithCountToPortalView(portalWithCounts))
	}

	return result
}

func (p portalsPresenter) PortalsFilterToEntity(filter *portalsv2.PortalsFilterRequest) *portalv2.FilterPortalsFilters {
	portalIDs := make([]int, 0, len(filter.PortalIDs))
	for _, id := range filter.PortalIDs {
		portalIDs = append(portalIDs, int(id))
	}
	return &portalv2.FilterPortalsFilters{
		IDs: portalIDs,
	}
}

func (p portalsPresenter) portalWithCountToPortalView(portalWithCounts *portalv2.PortalWithCounts) *portalsv2.Portal {
	if portalWithCounts == nil || portalWithCounts.Portal == nil {
		return nil
	}

	icon := ""
	if portalWithCounts.Portal.ImageID != nil {
		icon = *portalWithCounts.Portal.ImageID
	}

	return &portalsv2.Portal{
		ID:          portalWithCounts.Portal.ID,
		Name:        portalWithCounts.Portal.ShortName,
		IconID:      icon,
		Description: portalWithCounts.Portal.Name,
		Count: portalsv2.Count{
			Employees: portalWithCounts.EmployeesCount,
			Podved:    portalWithCounts.OrgsCount,
		},
		// TODO: реаизовать после описания БЛ (откуда брать данные)
		SctructureType: []string{"staffpositions", "management"},
		Head:           p.portalManagerToHeadView(portalWithCounts.Portal.Manager),
	}
}

func (p portalsPresenter) portalManagerToHeadView(portalManager *portalv2.PortalManager) *portalsv2.Head {
	if portalManager == nil {
		return nil
	}

	middleName := ""
	if portalManager.MiddleName != nil {
		middleName = *portalManager.MiddleName
	}

	image := ""
	if portalManager.ImageID != nil {
		image = *portalManager.ImageID
	}

	return &portalsv2.Head{
		FirstName:  portalManager.FirstName,
		LastName:   portalManager.LastName,
		MiddleName: middleName,
		ImageID:    image,
		// TODO: пока будет пустым так как из portals-v2 не приходит
		Description: portalManager.Prosition,
	}
}
