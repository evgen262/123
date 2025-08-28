package portals

import (
	"context"
	"errors"
	"fmt"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portal"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/repositories"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"go.uber.org/zap"
)

type portalsUseCase struct {
	logger     ditzap.Logger
	repository PortalRepository
}

func NewPortalsUseCase(
	repository PortalRepository,
	logger ditzap.Logger,
) *portalsUseCase {
	return &portalsUseCase{
		logger:     logger,
		repository: repository,
	}
}

func (p *portalsUseCase) GetByEmployees(ctx context.Context, employees []portal.EmployeeInfo) ([]*portal.Portal, error) {
	inns := make(portal.OrganizationINNs, 0, len(employees))

	for _, employee := range employees {
		inns = append(inns, portal.OrganizationINN(employee.Inn))
	}

	opts := portal.PortalsFilterOptions{
		INNs:       inns,
		OnlyLinked: true,
	}

	portals, err := p.repository.Filter(ctx, opts)
	if err != nil {
		switch true {
		case errors.As(err, new(diterrors.LocalizedError)):
			p.logger.Error("can't filter portal by employees", zap.Error(err.(diterrors.LocalizedError).Unwrap()))
		case errors.As(err, new(diterrors.ValidationError)):
			return nil, err
		case errors.Is(err, repositories.ErrNotFound):
			return nil, nil
		default:
			p.logger.Error("can't filter portal by employees", zap.Error(err))
		}
		return nil, fmt.Errorf("can't filter portal by employees: %w", err)
	}

	return portals, nil
}

func (p *portalsUseCase) Filter(ctx context.Context, opts portal.PortalsFilterOptions) ([]*portal.Portal, error) {
	portals, err := p.repository.Filter(ctx, opts)
	if err != nil {
		switch true {
		case errors.As(err, new(diterrors.ValidationError)):
			return nil, err
		case errors.Is(err, repositories.ErrNotFound):
			return nil, nil
		case errors.As(err, new(diterrors.LocalizedError)):
			p.logger.Error("can't filter portal", zap.Error(err.(diterrors.LocalizedError).Unwrap()))
		default:
			p.logger.Error("can't filter portal", zap.Error(err))
		}
		return nil, fmt.Errorf("can't filter portal: %w", err)
	}

	return portals, nil
}

func (p *portalsUseCase) GetAll(ctx context.Context, opts portal.GetAllOptions) ([]*portal.Portal, error) {
	portals, err := p.repository.Filter(ctx, portal.PortalsFilterOptions{
		OnlyLinked:  opts.OnlyLinked,
		WithDeleted: opts.WithDeleted,
	})
	if err != nil {
		var (
			localizedError diterrors.LocalizedError
			validationErr  diterrors.ValidationError
		)
		switch {
		case errors.As(err, &validationErr):
			return nil, validationErr
		case errors.Is(err, repositories.ErrNotFound):
			return nil, nil
		case errors.As(err, &localizedError):
			p.logger.Error("can't get all portal", zap.Error(localizedError.Unwrap()))
		default:
			p.logger.Error("can't get all portal", zap.Error(err))
		}
		return nil, fmt.Errorf("can't get all portal: %w", err)
	}

	return portals, nil
}

func (p *portalsUseCase) Get(ctx context.Context, id int, withDeleted bool) (*portal.Portal, error) {
	portals, err := p.repository.Filter(ctx, portal.PortalsFilterOptions{
		PortalIDs:   portal.PortalIDs{portal.PortalID(id)},
		WithDeleted: withDeleted,
	})
	if err != nil {
		var (
			localizedError diterrors.LocalizedError
			validationErr  diterrors.ValidationError
		)
		switch {
		case errors.As(err, &validationErr):
			return nil, validationErr
		case errors.Is(err, repositories.ErrNotFound):
			return nil, repositories.ErrNotFound
		case errors.As(err, &localizedError):
			p.logger.Error("can't get portal", zap.Error(localizedError.Unwrap()))
		default:
			p.logger.Error("can't get portal", zap.Error(err))
		}
		return nil, fmt.Errorf("can't get portal: %w", err)
	}
	if len(portals) == 0 {
		return nil, repositories.ErrNotFound
	}

	return portals[0], nil
}

func (p *portalsUseCase) MultiplyAdd(ctx context.Context, portals []*portal.Portal) ([]*portal.Portal, error) {
	newPortals, err := p.repository.Add(ctx, portals)
	if err != nil {
		switch true {
		case errors.As(err, new(diterrors.ValidationError)):
			return nil, err
		case errors.Is(err, repositories.ErrNotFound):
			return nil, repositories.ErrNotFound
		case errors.As(err, new(diterrors.LocalizedError)):
			p.logger.Error("can't add portal", zap.Error(err.(diterrors.LocalizedError).Unwrap()))
		default:
			p.logger.Error("can't add portal", zap.Error(err))
		}
		return nil, fmt.Errorf("can't add portal: %w", err)
	}

	return newPortals, nil
}

func (p *portalsUseCase) Add(ctx context.Context, newPortal *portal.Portal) (*portal.Portal, error) {
	newPortals, err := p.repository.Add(ctx, []*portal.Portal{newPortal})
	if err != nil {
		switch true {
		case errors.As(err, new(diterrors.ValidationError)):
			return nil, err
		case errors.Is(err, repositories.ErrNotFound):
			return nil, repositories.ErrNotFound
		case errors.As(err, new(diterrors.LocalizedError)):
			p.logger.Error("can't add portal", zap.Error(err.(diterrors.LocalizedError).Unwrap()))
		default:
			p.logger.Error("can't add portal", zap.Error(err))
		}
		return nil, fmt.Errorf("can't add portal: %w", err)
	}

	if len(newPortals) == 0 {
		return nil, nil
	}

	return newPortals[0], nil
}

func (p *portalsUseCase) Update(ctx context.Context, portal *portal.Portal) (*portal.Portal, error) {
	updatedPortal, err := p.repository.Update(ctx, portal)
	if err != nil {
		switch true {
		case errors.As(err, new(diterrors.ValidationError)):
			return nil, err
		case errors.Is(err, repositories.ErrNotFound):
			return nil, repositories.ErrNotFound
		case errors.As(err, new(diterrors.LocalizedError)):
			p.logger.Error("can't update portal", zap.Error(err.(diterrors.LocalizedError).Unwrap()))
		default:
			p.logger.Error("can't update portal", zap.Error(err))
		}
		return nil, fmt.Errorf("can't update portal: %w", err)
	}

	return updatedPortal, nil
}

func (p *portalsUseCase) Delete(ctx context.Context, id int) error {
	err := p.repository.Delete(ctx, portal.PortalID(id))
	if err != nil {
		switch true {
		case errors.As(err, new(diterrors.ValidationError)):
			return err
		case errors.Is(err, repositories.ErrNotFound):
			return repositories.ErrNotFound
		case errors.As(err, new(diterrors.LocalizedError)):
			p.logger.Error("can't delete portal", zap.Error(err.(diterrors.LocalizedError).Unwrap()))
		default:
			p.logger.Error("can't delete portal", zap.Error(err))
		}
		return fmt.Errorf("can't delete portal: %w", err)
	}

	return nil
}
