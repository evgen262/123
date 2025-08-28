package portalsv2

import (
	portalsv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/portalsfacade/portals/v1"
	complexesv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/portalsfacade/complexes/v1"
	entityPortalv2 "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portalv2"
)

//go:generate mockgen -source=interfaces.go -destination=./repositories_mock.go -package=portalsv2

type PortalsMapper interface {
	PortalsWithCountsToEntity(portalsWithCountsPb []*portalsv1.PortalWithCounts) []*entityPortalv2.PortalWithCounts
	PortalsFiltersToPb(filters *entityPortalv2.FilterPortalsFilters) *portalsv1.FilterRequest_Filters
	PortalsOptionsToPb(options *entityPortalv2.FilterPortalsOptions) *portalsv1.FilterRequest_Options
}

type ComplexesMapper interface {
	ComplexesFiltersToPb(filters *entityPortalv2.FilterComplexesFilters) *complexesv1.FilterRequest_Filters
	ComplexesOptionsToPb(options *entityPortalv2.FilterComplexesOptions) *complexesv1.FilterRequest_Options
	ComplexesToEntity(complexesPb []*complexesv1.Complex) []*entityPortalv2.Complex
}
