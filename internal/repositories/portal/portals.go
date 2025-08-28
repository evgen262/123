package portal

import (
	"cmp"
	"context"
	"fmt"
	"slices"

	portalsv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/portals/portals/v1"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"google.golang.org/grpc/codes"

	entityPortal "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portal"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/repositories"
)

type portalsRepository struct {
	client portalsv1.PortalsAPIClient
	mapper PortalsMapper
}

func NewPortalsRepository(client portalsv1.PortalsAPIClient, mapper PortalsMapper) *portalsRepository {
	return &portalsRepository{
		client: client,
		mapper: mapper,
	}
}

func (pr *portalsRepository) Get(ctx context.Context, id int) (*entityPortal.Portal, error) {
	portals, err := pr.Filter(ctx, entityPortal.PortalsFilterOptions{PortalIDs: entityPortal.PortalIDs{entityPortal.PortalID(id)}})
	if err != nil {
		return nil, fmt.Errorf("portalsRepository.Get: %w", err)
	}
	if len(portals) == 0 {
		return nil, diterrors.ErrNotFound
	}
	return portals[0], nil
}

func (pr *portalsRepository) Filter(ctx context.Context, options entityPortal.PortalsFilterOptions) ([]*entityPortal.Portal, error) {

	req := &portalsv1.FilterRequest{
		Filters: &portalsv1.FilterRequest_Filters{
			Inns:      options.INNs.ToStringSlice(),
			OrgIds:    options.OrgIDs.ToStringSlice(),
			PortalIds: options.PortalIDs.ToToInt32Slice(),
		},
		Options: &portalsv1.FilterRequest_Options{
			WithDeleted: options.WithDeleted,
			OnlyLinked:  options.OnlyLinked,
		},
	}

	resp, err := pr.client.Filter(ctx, req)
	if err != nil {
		msg := diterrors.NewLocalizedError(diterrors.LocalizeLocale, err)
		switch msg.Code() {
		case codes.InvalidArgument:
			return nil, diterrors.NewValidationError(msg)
		case codes.NotFound:
			return nil, repositories.ErrNotFound
		default:
			return nil, fmt.Errorf("can't filter portal: %w", msg)
		}
	}

	sortedPortals := resp.GetPortals()
	slices.SortFunc(sortedPortals, func(a, b *portalsv1.Portal) int {
		return cmp.Compare(a.GetId(), b.GetId())
	})

	return pr.mapper.PortalsToEntity(sortedPortals), nil
}

func (pr *portalsRepository) Add(ctx context.Context, portals []*entityPortal.Portal) ([]*entityPortal.Portal, error) {
	r, err := pr.client.Add(ctx, &portalsv1.AddRequest{
		Portals: pr.mapper.NewPortalsToPb(portals),
	})
	if err != nil {
		msg := diterrors.NewLocalizedError(diterrors.LocalizeLocale, err)
		switch msg.Code() {
		case codes.InvalidArgument:
			return nil, diterrors.NewValidationError(msg)
		case codes.NotFound:
			return nil, repositories.ErrNotFound
		default:
			return nil, fmt.Errorf("can't add portal: %w", msg)
		}
	}

	return pr.mapper.PortalsToEntity(r.GetPortals()), nil
}

func (pr *portalsRepository) Update(ctx context.Context, portal *entityPortal.Portal) (*entityPortal.Portal, error) {
	r, err := pr.client.Update(ctx, &portalsv1.UpdateRequest{
		Portal: pr.mapper.PortalToPb(portal),
	})
	if err != nil {
		msg := diterrors.NewLocalizedError(diterrors.LocalizeLocale, err)
		switch msg.Code() {
		case codes.InvalidArgument:
			return nil, diterrors.NewValidationError(msg)
		case codes.NotFound:
			return nil, repositories.ErrNotFound
		default:
			return nil, fmt.Errorf("can't update portal: %w", msg)
		}
	}

	return pr.mapper.PortalToEntity(r.GetPortal()), nil
}

func (pr *portalsRepository) Delete(ctx context.Context, id entityPortal.PortalID) error {
	_, err := pr.client.Delete(ctx, &portalsv1.DeleteRequest{
		PortalId: int32(id),
	})
	if err != nil {
		msg := diterrors.NewLocalizedError(diterrors.LocalizeLocale, err)
		switch msg.Code() {
		case codes.InvalidArgument:
			return diterrors.NewValidationError(msg)
		case codes.NotFound:
			return repositories.ErrNotFound
		default:
			return fmt.Errorf("can't delete portal with id[%d]: %w", id, err)
		}
	}

	return nil
}
