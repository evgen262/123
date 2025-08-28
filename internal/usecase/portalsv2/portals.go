package portalsv2

import (
	"context"
	"fmt"

	entityPortalsV2 "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portalv2"
)

type portalsInteractor struct {
	portalsRepository PortalsRepository
}

func NewPortalsUseCase(portalsRepository PortalsRepository) *portalsInteractor {
	return &portalsInteractor{
		portalsRepository: portalsRepository,
	}
}

func (p *portalsInteractor) Filter(
	ctx context.Context,
	filters *entityPortalsV2.FilterPortalsFilters,
	options *entityPortalsV2.FilterPortalsOptions,
) ([]*entityPortalsV2.PortalWithCounts, error) {
	portalsWithCounts, err := p.portalsRepository.Filter(ctx, filters, options)
	if err != nil {
		return nil, fmt.Errorf("portalsRepository.Filter: can't filter portals: %w", err)
	}

	return portalsWithCounts, nil
}
