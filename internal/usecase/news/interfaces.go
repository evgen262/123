package news

import (
	"context"

	"github.com/google/uuid"

	dtoNews "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/dto/news"
	entityEmployee "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/employee"
	entityNews "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/news"
)

//go:generate mockgen -source=interfaces.go -destination=./news_mock.go -package=news

type CategoryRepository interface {
	Create(ctx context.Context, nc *dtoNews.NewCategory) (*entityNews.Category, error)
	Update(ctx context.Context, c *dtoNews.UpdateCategory) (*entityNews.Category, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Filter(ctx context.Context, filter *dtoNews.FilterCategory) (*entityNews.CategoriesWithPagination, error)
}

type EmployeesRepository interface {
	Get(ctx context.Context, id uuid.UUID) (*entityEmployee.Employee, error)
	GetByExtIDAndPortalID(ctx context.Context, extID string, portalID int) (*entityEmployee.Employee, error)
}

type NewsRepository interface {
	Create(ctx context.Context, news *dtoNews.NewNews) (uuid.UUID, error)
	Update(ctx context.Context, id uuid.UUID, news *dtoNews.UpdateNews) (*entityNews.News, error)
	Search(ctx context.Context, search *dtoNews.SearchNews) (*dtoNews.SearchNewsResult, error)
	Get(ctx context.Context, id uuid.UUID) (*entityNews.NewsFull, error)
	GetBySlug(ctx context.Context, slug string) (*entityNews.NewsFull, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type CommentsRepository interface {
	Create(ctx context.Context, in dtoNews.NewComment) (uuid.UUID, int, error)
	List(ctx context.Context, params *dtoNews.FilterComments) ([]*entityNews.NewsComment, int, error)
}
