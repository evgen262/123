package news

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"go.uber.org/zap"

	dtoNews "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/dto/news"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity"
	entityNews "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/news"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/repositories/employees"
	repositoryNews "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/repositories/news"
)

type categoryInteractor struct {
	categoryRepository  CategoryRepository
	employeesRepository EmployeesRepository
	logger              ditzap.Logger
}

func NewCategoryInteractor(
	categoryRepository CategoryRepository,
	employeesRepository EmployeesRepository,
	logger ditzap.Logger,
) *categoryInteractor {
	return &categoryInteractor{
		categoryRepository:  categoryRepository,
		employeesRepository: employeesRepository,
		logger:              logger,
	}
}

func (c *categoryInteractor) Create(ctx context.Context, nc *dtoNews.NewCategory) (*entityNews.Category, error) {
	if nc == nil {
		c.logger.Debug("categoryInteractor.Create: nil request")
		return nil, diterrors.ErrInputEmpty
	}

	session, err := entity.SessionFromContext(ctx)
	if err != nil {
		c.logger.Error("categoryInteractor.Create: can't get session from context", zap.String("author_id", nc.AuthorID.String()), zap.String("category_name", nc.Name))
		return nil, fmt.Errorf("categoryInteractor.Create: can't get session from context: %w", err)
	}

	// TODO: использовать GetByID когда добавим в сессию ID
	emp, err := c.employeesRepository.GetByExtIDAndPortalID(ctx, session.GetUser().GetEmployee().GetExtID(), session.GetActivePortal().GetPortalID())
	if err != nil {
		switch {
		case errors.Is(err, employees.ErrEmployeeNotFound):
			c.logger.Error("categoryInteractor.Create: author not found", zap.String("author_id", nc.AuthorID.String()), zap.String("category_name", nc.Name))
			return nil, ErrAuthorNotFound
		default:
			c.logger.Error("categoryInteractor.Create: can't get author", zap.String("author_id", nc.AuthorID.String()), zap.String("category_name", nc.Name), zap.Error(err))
			return nil, fmt.Errorf("categoryInteractor.Create: %w", err)
		}
	}

	nc.AuthorID = emp.ID
	nc.PortalID = emp.Portal.ID

	res, err := c.categoryRepository.Create(ctx, nc)
	if err != nil {
		var valErr diterrors.ValidationError
		switch {
		case errors.As(err, &valErr):
			c.logger.Debug("categoryInteractor: invalid argument", zap.String("author_id", nc.AuthorID.String()), zap.String("category_name", nc.Name))
			return nil, valErr
		case errors.Is(err, repositoryNews.ErrCategoryAlreadyExists):
			return nil, ErrCategoryAlreadyExists
		default:
			c.logger.Error("categoryInteractor.Create: can't create category", zap.String("category_name", nc.Name), zap.Error(err))
			return nil, fmt.Errorf("categoryInteractor.Create: %w", err)
		}

	}
	return res, nil
}

func (c *categoryInteractor) Update(ctx context.Context, category *dtoNews.UpdateCategory) (*entityNews.Category, error) {
	if category == nil {
		c.logger.Debug("categoryInteractor.Update: nil request")
		return nil, diterrors.ErrInputEmpty
	}

	session, err := entity.SessionFromContext(ctx)
	if err != nil {
		c.logger.Error("categoryInteractor.Update: can't get session from context", zap.String("category_id", category.ID.String()), zap.String("author_id", category.AuthorID.String()), zap.String("category_name", category.Name))
		return nil, fmt.Errorf("categoryInteractor.Create: can't get session from context: %w", err)
	}

	if err = category.Validate(); err != nil {
		return nil, diterrors.NewValidationError(fmt.Errorf("categoryInteractor.Update: %w", err))
	}

	// TODO: использовать GetByID когда добавим в сессию ID
	emp, err := c.employeesRepository.GetByExtIDAndPortalID(ctx, session.GetUser().GetEmployee().GetExtID(), session.GetActivePortal().GetPortalID())
	if err != nil {
		switch {
		case errors.Is(err, employees.ErrEmployeeNotFound):
			c.logger.Error("categoryInteractor.Update: author not found", zap.String("category_id", category.ID.String()), zap.String("author_id", category.AuthorID.String()), zap.String("category_name", category.Name))
			return nil, ErrAuthorNotFound
		default:
			c.logger.Error("categoryInteractor.Update: can't get author", zap.String("category_id", category.ID.String()), zap.String("author_id", category.AuthorID.String()), zap.String("category_name", category.Name), zap.Error(err))
			return nil, fmt.Errorf("categoryInteractor.Update: %w", err)
		}
	}

	category.AuthorID = emp.ID
	category.PortalID = emp.Portal.ID

	res, err := c.categoryRepository.Update(ctx, category)
	if err != nil {
		switch {
		case errors.Is(err, repositoryNews.ErrCategoryNotFound):
			return nil, ErrCategoryNotFound
		case errors.Is(err, repositoryNews.ErrCategoryAlreadyExists):
			return nil, ErrCategoryAlreadyExists
		default:
			c.logger.Error("categoryInteractor.Update: can't update category", zap.String("category_id", category.ID.String()), zap.String("category_name", category.Name), zap.Error(err))
			return nil, fmt.Errorf("categoryInteractor.Update: %w", err)
		}
	}
	return res, nil
}

func (c *categoryInteractor) Delete(ctx context.Context, id uuid.UUID) error {
	err := c.categoryRepository.Delete(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, repositoryNews.ErrCategoryNotFound):
			return ErrCategoryNotFound
		default:
			c.logger.Error("categoryInteractor.Delete: can't delete category", zap.String("category_id", id.String()))
			return fmt.Errorf("categoryInteractor.Delete: %w", err)
		}
	}
	return nil
}

func (c *categoryInteractor) Search(ctx context.Context, s *dtoNews.SearchCategory) (*entityNews.CategoriesWithPagination, error) {
	if s == nil {
		c.logger.Debug("categoryInteractor.Search: nil request")
		return nil, diterrors.ErrInputEmpty
	}

	session, err := entity.SessionFromContext(ctx)
	if err != nil {
		c.logger.Error("categoryInteractor.Search", zap.Error(err))
		return nil, fmt.Errorf("categoryInteractor.Search: %w", err)
	}

	filter := &dtoNews.FilterCategory{
		By:         dtoNews.FilterCategoryByName,
		Name:       &s.Query,
		Pagination: &s.Pagination,
		Visitor: &entityNews.Visitor{
			PortalID: session.GetActivePortal().GetPortalID(),
		},
	}

	res, err := c.categoryRepository.Filter(ctx, filter)
	if err != nil {
		c.logger.Error("categoryInteractor.Search", zap.String("ext_employee_id", session.GetUser().GetEmployee().GetExtID()), zap.Error(err))
		return nil, fmt.Errorf("categoryInteractor.Search: %w", err)
	}
	return res, nil
}

func (c *categoryInteractor) Get(ctx context.Context, id uuid.UUID) (*entityNews.Category, error) {
	filter := &dtoNews.FilterCategory{
		By:  dtoNews.FilterCategoryByIDs,
		IDs: []uuid.UUID{id},
	}

	res, err := c.categoryRepository.Filter(ctx, filter)
	if err != nil {
		c.logger.Error("categoryInteractor.Get", zap.String("category_id", id.String()), zap.Error(err))
		return nil, fmt.Errorf("categoryInteractor.Get: %w", err)
	}

	if len(res.Categories) == 0 {
		return nil, ErrCategoryNotFound
	}
	return res.Categories[0], nil
}
