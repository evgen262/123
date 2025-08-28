package portal

import (
	"context"
	"fmt"

	organizationsv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/portals/organizations/v1"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portal"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/repositories"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"google.golang.org/grpc/codes"
)

type organizationsRepository struct {
	client organizationsv1.OrganizationsAPIClient
	mapper OrganizationsMapper
}

func NewOrganizationsRepository(
	client organizationsv1.OrganizationsAPIClient,
	mapper OrganizationsMapper,
) *organizationsRepository {
	return &organizationsRepository{
		client: client,
		mapper: mapper,
	}
}

func (por organizationsRepository) Filter(
	ctx context.Context,
	filters portal.OrganizationsFilters,
	pagination *entity.StringPagination,
	options portal.OrganizationsFilterOptions,
) (*portal.OrganizationsWithPagination, error) {
	var filterPb *organizationsv1.FilterRequest

	switch true {
	case len(filters.Inns) != 0:
		filterPb = &organizationsv1.FilterRequest{
			Filter: &organizationsv1.FilterRequest_Inns{
				Inns: &organizationsv1.FilterRequest_Inn{
					Inns: filters.Inns,
				},
			},
		}
	case len(filters.Ogrns) != 0:
		filterPb = &organizationsv1.FilterRequest{
			Filter: &organizationsv1.FilterRequest_Ogrns{
				Ogrns: &organizationsv1.FilterRequest_Ogrn{
					Ogrns: filters.Ogrns,
				},
			},
		}
	case len(filters.Names) != 0:
		filterPb = &organizationsv1.FilterRequest{
			Filter: &organizationsv1.FilterRequest_Names_{
				Names: &organizationsv1.FilterRequest_Names{
					Names: filters.Names,
				},
			},
		}
	default:
		filterPb = &organizationsv1.FilterRequest{
			Filter: &organizationsv1.FilterRequest_Ids_{
				Ids: &organizationsv1.FilterRequest_Ids{
					Ids: filters.Ids,
				},
			},
		}
	}
	filterPb.Pagination = por.mapper.PaginationToPb(pagination)
	filterPb.Options = por.mapper.OptionsToPb(options)
	res, err := por.client.Filter(ctx, filterPb)
	if err != nil {
		localError := diterrors.NewLocalizedError("ru-RU", err)
		switch localError.Code() {
		case codes.InvalidArgument:
			return nil, diterrors.NewValidationError(localError)
		case codes.NotFound:
			return nil, repositories.ErrNotFound
		default:
			return nil, fmt.Errorf("can't filter organizations: %w", localError.Unwrap())
		}
	}

	return &portal.OrganizationsWithPagination{
		Organizations: por.mapper.OrganizationsToEntity(res.GetOrganizations()),
		Pagination:    por.mapper.PaginationToEntity(res.GetPagination()),
	}, nil
}

func (por organizationsRepository) LinkOrganizationsToPortal(
	ctx context.Context,
	portalId portal.PortalID,
	ids portal.OrganizationIDs,
) error {
	_, err := por.client.Link(ctx, &organizationsv1.LinkRequest{
		PortalId: int32(portalId),
		OrgIds:   ids.ToStringSlice(),
	})
	if err != nil {
		return fmt.Errorf("can't link to portal #%d: %w", portalId, err)
	}
	return nil
}

func (por organizationsRepository) UnlinkOrganizations(ctx context.Context, ids portal.OrganizationIDs) error {
	_, err := por.client.Unlink(ctx, &organizationsv1.UnlinkRequest{OrgIds: ids.ToStringSlice()})
	if err != nil {
		return fmt.Errorf("can't unlink organizations: %w", err)
	}
	return nil
}
