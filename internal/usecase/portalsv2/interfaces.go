package portalsv2

import (
	"context"

	entityPortalsV2 "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portalv2"
)

//go:generate mockgen -source=interfaces.go -destination=./usecases_mock.go -package=portalsv2

type PortalsRepository interface {
	Filter(
		ctx context.Context,
		filters *entityPortalsV2.FilterPortalsFilters,
		options *entityPortalsV2.FilterPortalsOptions,
	) ([]*entityPortalsV2.PortalWithCounts, error)
}

type ComplexesRepository interface {
	Filter(
		ctx context.Context,
		filters *entityPortalsV2.FilterComplexesFilters,
		options *entityPortalsV2.FilterComplexesOptions,
	) ([]*entityPortalsV2.Complex, error)
}
