package portals

import (
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view"
	viewPortals "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/portals"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portal"
)

type organizationPresenter struct {
}

func NewOrganizationPresenter() *organizationPresenter {
	return &organizationPresenter{}
}

func (op organizationPresenter) OrganizationIdsToEntity(ids viewPortals.OrganizationIds) portal.OrganizationIDs {
	entIds := make(portal.OrganizationIDs, 0, len(ids))
	for _, id := range ids {
		entIds = append(entIds, portal.OrganizationId(id))
	}
	return entIds
}

func (op organizationPresenter) FiltersToEntity(filters viewPortals.OrganizationsFilters) portal.OrganizationsFilters {
	return portal.OrganizationsFilters{
		Ids:   filters.Ids,
		Names: filters.Names,
		Ogrns: filters.Ogrns,
		Inns:  filters.Inns,
	}
}

func (op organizationPresenter) OrganizationsToView(orgs []*portal.Organization) []*viewPortals.Organization {
	organizations := make([]*viewPortals.Organization, 0, len(orgs))
	for _, org := range orgs {
		grbs := make([]*viewPortals.OrganizationGrbs, 0)
		for _, onceGrbs := range org.Grbs {
			grbs = append(grbs, &viewPortals.OrganizationGrbs{
				GrbsId:    onceGrbs.GrbsId,
				Name:      onceGrbs.Name,
				Inn:       onceGrbs.Inn,
				Ogrn:      onceGrbs.Ogrn,
				StartDate: onceGrbs.GetStartDate(),
			})
		}
		organizations = append(organizations, &viewPortals.Organization{
			Id:              string(org.Id),
			FullName:        org.FullName,
			ShortName:       org.ShortName,
			RegCode:         org.RegCode,
			OrgCode:         org.OrgCode,
			UNKCode:         org.UNKCode,
			OGRN:            org.OGRN,
			INN:             string(org.INN),
			KPP:             org.KPP,
			OrgTypeName:     org.OrgTypeName,
			OrgTypeCode:     org.OrgTypeCode,
			UchrezhTypeCode: org.UchrezhTypeCode,
			Grbs:            grbs,
			Email:           org.Email,
			Phone:           org.Phone,
			AdditionalPhone: org.AdditionalPhone,
			Site:            org.Site,
			CreatedAt:       org.CreatedAt,
			UpdatedAt:       org.UpdatedAt,
			LiquidatedAt:    org.GetLiquidatedAt(),
			IsLiquidated:    org.IsLiquidated,
		})
	}
	return organizations
}

func (op organizationPresenter) PaginationToView(pagination *entity.StringPagination) *view.StringPagination {
	if pagination == nil || pagination.LastId == nil {
		return nil
	}
	return &view.StringPagination{
		LastId:          pagination.LastId.String(),
		LastCreatedTime: pagination.GetLastCreatedTime(),
		Limit:           pagination.GetLimit(),
		Total:           pagination.GetTotal(),
	}
}

func (op organizationPresenter) PaginationToEntity(pagination *view.StringPagination) *entity.StringPagination {
	if pagination == nil {
		return nil
	}
	return &entity.StringPagination{
		LastId:          portal.OrganizationId(pagination.LastId),
		LastCreatedTime: pagination.GetLastCreatedTime(),
		Limit:           pagination.GetLimit(),
		Total:           pagination.GetTotal(),
	}
}

func (op organizationPresenter) OptionsToEntity(options *viewPortals.OrganizationsFilterOptions) portal.OrganizationsFilterOptions {
	opts := portal.OrganizationsFilterOptions{}
	if options == nil {
		return opts
	}
	opts.WithLiquidated = options.WithLiquidated

	return opts
}
