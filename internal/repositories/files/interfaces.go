package files

import (
	filev1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/fileservice/file/v1"
	sharedv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/fileservice/shared/v1"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/auth"
	entityFile "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/file"
)

//go:generate mockgen -source=interfaces.go -destination=./repositories_mock.go -package=files

type FilesMapper interface {
	FilePbToEntity(f *filev1.File) *entityFile.File
}

type VisitorMapper interface {
	SessionToVisitorPb(session *auth.Session) *sharedv1.Visitor
}
