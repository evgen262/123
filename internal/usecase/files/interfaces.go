package files

import (
	"context"

	"github.com/google/uuid"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/auth"
	entityFile "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/file"
)

//go:generate mockgen -source=interfaces.go -destination=./usecases_mock.go -package=files
type FileRepository interface {
	Get(ctx context.Context, fileId uuid.UUID, session *auth.Session) (*entityFile.File, error)
}
