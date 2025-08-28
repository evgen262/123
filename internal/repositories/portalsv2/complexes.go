package portalsv2

import (
	"context"
	"fmt"

	complexesv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/portalsfacade/complexes/v1"
	entityPortalsV2 "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portalv2"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
)

type complexesRepository struct {
	complexesClient complexesv1.ComplexesAPIClient
	complexesMapper ComplexesMapper
}

func NewComplexesRepository(client complexesv1.ComplexesAPIClient, mapper ComplexesMapper) *complexesRepository {
	return &complexesRepository{
		complexesClient: client,
		complexesMapper: mapper,
	}
}

func (r *complexesRepository) Filter(ctx context.Context, filters *entityPortalsV2.FilterComplexesFilters, options *entityPortalsV2.FilterComplexesOptions) ([]*entityPortalsV2.Complex, error) {
	resp, err := r.complexesClient.Filter(ctx, &complexesv1.FilterRequest{
		Filters: r.complexesMapper.ComplexesFiltersToPb(filters),
		Options: r.complexesMapper.ComplexesOptionsToPb(options),
	})
	if err != nil {
		return nil, fmt.Errorf("complexesClient.Filter: can't get complexes: %w", diterrors.GrpcErrorToError(err))
	}

	return r.complexesMapper.ComplexesToEntity(resp.GetComplexes()), nil
}
