package employees

import (
	"context"

	"github.com/google/uuid"

	entityEmployee "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/employee"
	entityPortal "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portal"
)

//go:generate mockgen -source=interfaces.go -destination=./usecases_mock.go -package=employees

type EmployeesRepository interface {
	Get(ctx context.Context, id uuid.UUID) (*entityEmployee.Employee, error)
	GetByExtIDAndPortalID(ctx context.Context, extID string, portalID int) (*entityEmployee.Employee, error)
}

// TODO Порталы для получения иконок порталов. Далее должно быть в фасаде эмплоев забираться иконки.
type PortalsRepository interface {
	Get(ctx context.Context, id int) (*entityPortal.Portal, error)
}
