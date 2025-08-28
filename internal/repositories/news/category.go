package news

import (
	"context"
	"fmt"

	categoryv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/newsfacade/category/v1"
	sharedv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/newsfacade/shared/v1"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/wrapperspb"

	dtoNews "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/dto/news"
	entityNews "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/news"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/repositories"
)

type categoryRepository struct {
	categoryAPIClient categoryv1.CategoryAPIClient
	sharedMapper      repositories.SharedMapper
}

func NewCategoryRepository(
	categoryAPIClient categoryv1.CategoryAPIClient,
	sharedMapper repositories.SharedMapper,
) *categoryRepository {
	return &categoryRepository{
		categoryAPIClient: categoryAPIClient,
		sharedMapper:      sharedMapper,
	}
}

func (c *categoryRepository) Create(ctx context.Context, nc *dtoNews.NewCategory) (*entityNews.Category, error) {
	if nc == nil {
		return nil, diterrors.ErrInputEmpty
	}

	res, err := c.categoryAPIClient.Create(ctx, &categoryv1.CreateRequest{
		Name:       nc.Name,
		Visibility: &categoryv1.Visibility{PortalIds: c.sharedMapper.IntSliceToInt32(nc.GetVisibility().OIVs)},
	})
	if err != nil {
		gErr := diterrors.NewLocalizedError(diterrors.LocalizeLocale, err)
		switch gErr.Code() {
		case codes.InvalidArgument:
			return nil, fmt.Errorf("categoryRepository.Create: invalid argument : %w", diterrors.NewValidationError(err))
		case codes.AlreadyExists:
			return nil, ErrCategoryAlreadyExists
		default:
			return nil, fmt.Errorf("categoryRepository.Create: %w", diterrors.GrpcErrorToError(err))
		}
	}

	id, err := uuid.Parse(res.GetId())
	if err != nil || id == uuid.Nil {
		return nil, fmt.Errorf("categoryRepository.Create: can't parse result id=%s: %w", res.GetId(), err)
	}

	return &entityNews.Category{
		ID: id,
		Visibility: entityNews.CategoryVisibility{
			Condition:        nc.GetVisibility().Condition,
			ComplexIDs:       nc.GetVisibility().ComplexIDs,
			OIVs:             nc.GetVisibility().OIVs,
			OrgIDs:           nc.GetVisibility().OrgIDs,
			ProductIDs:       nc.GetVisibility().ProductIDs,
			SubdivisionNames: nc.GetVisibility().SubdivisionNames,
			PositionNames:    nc.GetVisibility().PositionNames,
			EmployeeIDs:      nc.GetVisibility().EmployeeIDs,
			RoleNames:        nc.GetVisibility().RoleNames,
		},
		Name: nc.Name,
	}, nil
}

func (c *categoryRepository) Update(ctx context.Context, category *dtoNews.UpdateCategory) (*entityNews.Category, error) {
	if category == nil {
		return nil, diterrors.ErrInputEmpty
	}

	res, err := c.categoryAPIClient.Update(ctx, &categoryv1.UpdateRequest{
		Id:         category.ID.String(),
		Name:       wrapperspb.String(category.Name),
		Visibility: &categoryv1.Visibility{PortalIds: c.sharedMapper.IntSliceToInt32(category.GetVisibility().OIVs)},
	})
	if err != nil {
		gErr := diterrors.NewLocalizedError(diterrors.LocalizeLocale, err)
		switch gErr.Code() {
		case codes.InvalidArgument:
			return nil, fmt.Errorf("categoryRepository.Update: invalid argument : %w", diterrors.NewValidationError(err))
		case codes.AlreadyExists:
			return nil, ErrCategoryAlreadyExists
		case codes.NotFound:
			return nil, ErrCategoryNotFound
		default:
			return nil, fmt.Errorf("categoryRepository.Update: %w", diterrors.GrpcErrorToError(err))
		}
	}

	id, err := uuid.Parse(res.GetCategory().GetId())
	if err != nil || id == uuid.Nil {
		return nil, fmt.Errorf("categoryRepository.Update: can't parse result id=%s: %w", res.GetCategory().GetId(), err)
	}

	return &entityNews.Category{
		ID: id,
		Visibility: entityNews.CategoryVisibility{
			OIVs:       c.sharedMapper.Int32SliceToInt(res.GetCategory().GetVisibility().GetPortalIds()),
			ComplexIDs: c.sharedMapper.Int32SliceToInt(res.GetCategory().GetVisibility().GetComplexIds()),
		},
		Name: res.GetCategory().GetName(),
	}, nil
}

func (c *categoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := c.categoryAPIClient.Delete(ctx, &categoryv1.DeleteRequest{
		Id: id.String(),
	})
	if err != nil {
		gErr := diterrors.NewLocalizedError(diterrors.LocalizeLocale, err)
		switch gErr.Code() {
		case codes.InvalidArgument:
			return fmt.Errorf("categoryRepository.Delete: invalid argument: %w", diterrors.NewValidationError(err))
		case codes.NotFound:
			return ErrCategoryNotFound
		default:
			return fmt.Errorf("categoryRepository.Delete: %w", diterrors.GrpcErrorToError(err))
		}
	}
	return nil
}

func (c *categoryRepository) Filter(ctx context.Context, filter *dtoNews.FilterCategory) (*entityNews.CategoriesWithPagination, error) {
	if filter == nil {
		return nil, diterrors.ErrInputEmpty
	}

	visitor := filter.GetVisitor()
	reqPb := &categoryv1.FilterRequest{
		Visitor: &sharedv1.Visitor{
			PortalId: int32(visitor.GetPortalID()),
		},
	}

	switch filter.By {
	case dtoNews.FilterCategoryByName:
		reqPb.Filters = &categoryv1.FilterRequest_Name{
			Name: *filter.Name,
		}
	case dtoNews.FilterCategoryByIDs:
		if filter.IDs == nil {
			return nil, fmt.Errorf("categoryRepository.Filter: empty ids")
		}
		ids := make([]string, 0, len(filter.IDs))
		for _, id := range filter.IDs {
			ids = append(ids, id.String())
		}
		reqPb.Filters = &categoryv1.FilterRequest_Ids{
			Ids: &sharedv1.IdsUUID{Ids: ids},
		}
	default:
		return nil, fmt.Errorf("categoryRepository.Filter: unknown filter by %d", filter.By)
	}

	resPb, err := c.categoryAPIClient.Filter(ctx, reqPb)
	if err != nil {
		gErr := diterrors.NewLocalizedError(diterrors.LocalizeLocale, err)
		switch gErr.Code() {
		case codes.InvalidArgument:
			return nil, fmt.Errorf("categoryRepository.Filter: invalid argument: %w", diterrors.NewValidationError(err))
		default:
			return nil, fmt.Errorf("categoryRepository.Filter: %w", diterrors.GrpcErrorToError(err))
		}
	}

	res := &entityNews.CategoriesWithPagination{}

	eRes := make([]*entityNews.Category, 0, len(resPb.GetCategories()))

	for _, category := range resPb.GetCategories() {
		if category == nil {
			continue
		}
		if id, err := uuid.Parse(category.GetId()); err == nil && id != uuid.Nil {
			eRes = append(eRes, &entityNews.Category{
				ID:   id,
				Name: category.GetName(),
				Visibility: entityNews.CategoryVisibility{
					ComplexIDs: c.sharedMapper.Int32SliceToInt(category.GetVisibility().GetComplexIds()),
					OIVs:       c.sharedMapper.Int32SliceToInt(category.GetVisibility().GetPortalIds()),
				},
			})
		}
	}

	res.Categories = eRes

	return res, nil
}
