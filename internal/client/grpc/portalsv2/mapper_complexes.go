package portalsv2

import (
	complexesv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/portalsfacade/complexes/v1"
	entityPortalsv2 "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portalv2"
	timeUtils "git.mos.ru/buch-cloud/moscow-team-2.0/build/time-utils.git"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type complexesMapper struct {
	tu timeUtils.TimeUtils
}

func NewComplexesMapper(tu timeUtils.TimeUtils) *complexesMapper {
	return &complexesMapper{tu: tu}
}

func (c complexesMapper) ComplexesFiltersToPb(filters *entityPortalsv2.FilterComplexesFilters) *complexesv1.FilterRequest_Filters {
	if filters == nil {
		return nil
	}

	pb := &complexesv1.FilterRequest_Filters{
		Ids:       make([]int32, 0, len(filters.IDs)),
		PortalIds: make([]int32, 0, len(filters.PortalIDs)),
	}

	for _, id := range filters.IDs {
		pb.Ids = append(pb.Ids, int32(id))
	}

	for _, id := range filters.PortalIDs {
		pb.PortalIds = append(pb.PortalIds, int32(id))
	}

	return pb
}

func (c complexesMapper) ComplexesOptionsToPb(options *entityPortalsv2.FilterComplexesOptions) *complexesv1.FilterRequest_Options {
	if options == nil {
		return nil
	}

	return &complexesv1.FilterRequest_Options{
		WithDeleted:  options.WithDeleted,
		WithDisabled: options.WithDisabled,
	}
}

func (c complexesMapper) ComplexesToEntity(complexesPb []*complexesv1.Complex) []*entityPortalsv2.Complex {
	if complexesPb == nil {
		return nil
	}
	entities := make([]*entityPortalsv2.Complex, 0, len(complexesPb))
	for _, complexPb := range complexesPb {
		if complexPb == nil {
			continue
		}

		entities = append(entities, c.complexToEntity(complexPb))
	}
	return entities
}

func (c complexesMapper) complexToEntity(complexPb *complexesv1.Complex) *entityPortalsv2.Complex {
	if complexPb == nil {
		return nil
	}

	return &entityPortalsv2.Complex{
		ID:           int(complexPb.GetId()),
		Name:         complexPb.GetName(),
		Description:  c.stringValueToStringPtr(complexPb.GetDescription()),
		ImageID:      c.stringValueToStringPtr(complexPb.GetImageId()),
		ComplexGroup: int(complexPb.GetComplexGroup()),
		Sort:         int(complexPb.GetSort()),
		IsDisabled:   complexPb.GetIsDisabled(),
		Responsible:  c.complexResponsibleToEntity(complexPb.GetResponsible()),
		Portals:      c.complexPortalsToEntity(complexPb.GetPortals()),
		CreatedAt:    complexPb.GetCreatedTime().AsTime(),
		UpdatedAt:    c.tu.TimestampToTime(complexPb.GetUpdatedTime()),
		DeletedAt:    c.tu.TimestampToTime(complexPb.GetDeletedTime()),
	}
}

func (c complexesMapper) complexResponsibleToEntity(complexPb *complexesv1.Complex_Responsible) *entityPortalsv2.ComplexResponsible {
	if complexPb == nil {
		return nil
	}

	return &entityPortalsv2.ComplexResponsible{
		FirstName:   complexPb.GetFirstName(),
		LastName:    complexPb.GetLastName(),
		MiddleName:  c.stringValueToStringPtr(complexPb.GetMiddleName()),
		ImageID:     c.stringValueToStringPtr(complexPb.GetImageId()),
		Description: complexPb.GetDescription(),
	}
}

func (c complexesMapper) complexPortalsToEntity(complexPortalsPb []*complexesv1.Complex_Portal) []*entityPortalsv2.ComplexPortal {
	if complexPortalsPb == nil {
		return nil
	}
	portals := make([]*entityPortalsv2.ComplexPortal, 0, len(complexPortalsPb))
	for _, portal := range complexPortalsPb {
		if portal == nil {
			continue
		}

		portals = append(portals, &entityPortalsv2.ComplexPortal{
			ID:   int(portal.GetId()),
			Sort: int(portal.GetSort()),
		})
	}
	return portals
}

func (c complexesMapper) stringValueToStringPtr(s *wrapperspb.StringValue) *string {
	if s == nil {
		return nil
	}
	val := s.GetValue()
	return &val
}
