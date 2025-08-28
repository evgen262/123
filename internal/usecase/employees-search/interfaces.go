package employees_search

import (
	"context"

	entityEmployeesSearch "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/employees-search"
)

//go:generate mockgen -source=interfaces.go -destination=./usecase_mock.go -package=employees_search

type EmployeesSearchRepository interface {
	Search(ctx context.Context, request *entityEmployeesSearch.SearchRequest) (*entityEmployeesSearch.SearchResponse, error)
	Filters(ctx context.Context, request *entityEmployeesSearch.FiltersRequest) (*entityEmployeesSearch.FiltersResponse, error)
}
