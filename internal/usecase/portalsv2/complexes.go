package portalsv2

import (
	"context"
	"fmt"

	entityPortalsV2 "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portalv2"
)

type compexesInteractor struct {
	complexesRepository ComplexesRepository
}

func NewComplexesUseCase(complexesRepository ComplexesRepository) *compexesInteractor {
	return &compexesInteractor{
		complexesRepository: complexesRepository,
	}
}

func (c *compexesInteractor) Filter(
	ctx context.Context,
	filters *entityPortalsV2.FilterComplexesFilters,
	options *entityPortalsV2.FilterComplexesOptions,
) ([]*entityPortalsV2.Complex, error) {
	complexes, err := c.complexesRepository.Filter(ctx, filters, options)
	if err != nil {
		return nil, fmt.Errorf("complexesRepository.Filter: can't filter complexes: %w", err)
	}

	return complexes, nil
}
