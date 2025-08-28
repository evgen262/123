package portalsv2

import (
	portalsv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/portalsfacade/portals/v1"
	entityPortalsV2 "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portalv2"
	timeUtils "git.mos.ru/buch-cloud/moscow-team-2.0/build/time-utils.git"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type portalsMapper struct {
	tu timeUtils.TimeUtils
}

func NewPortalsMapper(tu timeUtils.TimeUtils) *portalsMapper {
	return &portalsMapper{
		tu: tu,
	}
}

func (m portalsMapper) PortalsWithCountsToEntity(portalsWithCountsPb []*portalsv1.PortalWithCounts) []*entityPortalsV2.PortalWithCounts {
	if portalsWithCountsPb == nil {
		return nil
	}
	portals := make([]*entityPortalsV2.PortalWithCounts, 0, len(portalsWithCountsPb))
	for _, portalWithCounts := range portalsWithCountsPb {
		if portalWithCounts == nil {
			continue
		}

		portal := portalWithCounts.GetPortal()
		if portal == nil {
			continue
		}

		portals = append(portals, &entityPortalsV2.PortalWithCounts{
			Portal:         m.portalToEntity(portal),
			OrgsCount:      int(portalWithCounts.GetOrgsCount()),
			EmployeesCount: int(portalWithCounts.GetEmployeesCount()),
		})
	}
	return portals
}

func (m portalsMapper) PortalsFiltersToPb(filters *entityPortalsV2.FilterPortalsFilters) *portalsv1.FilterRequest_Filters {
	if filters == nil {
		return nil
	}
	portalIDs := make([]int32, 0, len(filters.IDs))
	for _, id := range filters.IDs {
		portalIDs = append(portalIDs, int32(id))
	}
	return &portalsv1.FilterRequest_Filters{
		Ids: portalIDs,
	}
}
func (m portalsMapper) PortalsOptionsToPb(options *entityPortalsV2.FilterPortalsOptions) *portalsv1.FilterRequest_Options {
	if options == nil {
		return nil
	}
	return &portalsv1.FilterRequest_Options{
		WithEmployeesCount: options.WithEmployeesCount,
	}
}

func (m portalsMapper) portalToEntity(p *portalsv1.Portal) *entityPortalsV2.Portal {
	if p == nil {
		return nil
	}

	return &entityPortalsV2.Portal{
		ID:                int(p.GetId()),
		Name:              p.GetName(),
		ShortName:         p.GetShortName(),
		Url:               p.GetUrl(),
		ImageID:           m.stringValueToStringPtr(p.GetImageId()),
		Status:            m.portalStatusToEntity(p.GetStatus()),
		StatusUpdatedTime: m.tu.TimestampToTime(p.GetStatusUpdatedTime()),
		IsDisabled:        p.GetIsDisabled(),
		Sort:              int(p.GetSort()),
		Manager:           m.portalManagerToEntity(p.GetManager()),
		CreatedAt:         p.GetCreatedTime().AsTime(),
		UpdatedAt:         m.tu.TimestampToTime(p.GetUpdatedTime()),
		DeletedAt:         m.tu.TimestampToTime(p.GetDeletedTime()),
	}
}

func (m portalsMapper) portalManagerToEntity(p *portalsv1.PortalManager) *entityPortalsV2.PortalManager {
	if p == nil {
		return nil
	}

	// TODO: добавить Position когда будет реализовано в portals-v2
	return &entityPortalsV2.PortalManager{
		FirstName:  p.GetFirstName(),
		LastName:   p.GetLastName(),
		MiddleName: m.stringValueToStringPtr(p.GetMiddleName()),
		ImageID:    m.stringValueToStringPtr(p.GetImageId()),
	}
}

func (m portalsMapper) portalStatusToEntity(s portalsv1.PortalStatus) entityPortalsV2.PortalStatus {
	switch s {
	case portalsv1.PortalStatus_PORTAL_STATUS_ACTIVE:
		return entityPortalsV2.PortalStatusActive
	case portalsv1.PortalStatus_PORTAL_STATUS_INACTIVE:
		return entityPortalsV2.PortalStatusInactive
	case portalsv1.PortalStatus_PORTAL_STATUS_MAINTENANCE:
		return entityPortalsV2.PortalStatusMaintenance
	default:
		return entityPortalsV2.PortalStatusInvalid
	}
}

func (m portalsMapper) stringValueToStringPtr(s *wrapperspb.StringValue) *string {
	if s == nil {
		return nil
	}
	val := s.GetValue()
	return &val
}
