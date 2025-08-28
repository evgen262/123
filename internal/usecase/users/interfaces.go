package users

import (
	"context"

	entityEmployee "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/employee"
)

//go:generate mockgen -source=interfaces.go -destination=./users_mock.go -package=users

type EmployeesRepository interface {
	GetByExtIDAndPortalID(ctx context.Context, extID string, portalID int) (*entityEmployee.Employee, error)
}
