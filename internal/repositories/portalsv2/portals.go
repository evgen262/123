package portalsv2

import (
	"context"
	"fmt"

	portalsv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/portalsfacade/portals/v1"
	entityPortalsV2 "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portalv2"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
)

type portalsRepository struct {
	portalsClient portalsv1.PortalsAPIClient
	portalsMapper PortalsMapper
}

func NewPortalsRepository(client portalsv1.PortalsAPIClient, mapper PortalsMapper) *portalsRepository {
	return &portalsRepository{
		portalsClient: client,
		portalsMapper: mapper,
	}
}

func (pr *portalsRepository) Filter(
	ctx context.Context,
	filters *entityPortalsV2.FilterPortalsFilters,
	options *entityPortalsV2.FilterPortalsOptions,
) ([]*entityPortalsV2.PortalWithCounts, error) {
	resp, err := pr.portalsClient.Filter(ctx,
		&portalsv1.FilterRequest{
			Filters: pr.portalsMapper.PortalsFiltersToPb(filters),
			Options: pr.portalsMapper.PortalsOptionsToPb(options),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("portalsClient.Filter: can't get portals: %w", diterrors.GrpcErrorToError(err))
	}

	return pr.portalsMapper.PortalsWithCountsToEntity(resp.GetPortals()), nil
}
