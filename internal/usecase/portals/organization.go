package portals

import (
	"context"
	"fmt"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portal"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"go.uber.org/zap"
)

type organizationsUseCase struct {
	repo   OrganizationsRepository
	logger ditzap.Logger
}

func NewOrganizationsUseCase(repository OrganizationsRepository, logger ditzap.Logger) *organizationsUseCase {
	return &organizationsUseCase{
		repo:   repository,
		logger: logger,
	}
}

func (ouc organizationsUseCase) Filter(
	ctx context.Context,
	filters portal.OrganizationsFilters,
	pagination *entity.StringPagination,
	options portal.OrganizationsFilterOptions,
) (*portal.OrganizationsWithPagination, error) {
	orgs, err := ouc.repo.Filter(ctx, filters, pagination, options)
	if err != nil {
		ouc.logger.Error("can't filter organizations", zap.Error(err))
		return nil, fmt.Errorf("can't filter organizations: %w", err)
	}
	return orgs, nil
}

func (ouc organizationsUseCase) Link(
	ctx context.Context,
	portalId portal.PortalID,
	orgIds portal.OrganizationIDs,
) error {
	if err := ouc.repo.LinkOrganizationsToPortal(ctx, portalId, orgIds); err != nil {
		ouc.logger.Error(
			"can't link organizations",
			zap.Int("portal_id", int(portalId)),
			zap.String("org_ids", orgIds.ToString()),
			zap.Error(err),
		)
		return fmt.Errorf("can't link organizations: %w", err)
	}
	return nil
}

func (ouc organizationsUseCase) Unlink(ctx context.Context, ids portal.OrganizationIDs) error {
	if err := ouc.repo.UnlinkOrganizations(ctx, ids); err != nil {
		ouc.logger.Error(
			"can't unlink organizations",
			zap.String("org_ids", ids.ToString()),
			zap.Error(err),
		)
		return fmt.Errorf("can't unlink organizations: %w", err)
	}
	return nil
}
