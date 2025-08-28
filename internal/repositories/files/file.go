package files

import (
	"context"
	"fmt"

	filev1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/fileservice/file/v1"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"github.com/google/uuid"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/auth"
	entityFile "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/file"
)

type fileRepositoty struct {
	client        filev1.FileAPIClient
	filesMapper   FilesMapper
	visitorMapper VisitorMapper
}

func NewFileRepository(
	client filev1.FileAPIClient,
	fileMapper FilesMapper,
	visitorMapper VisitorMapper,
) *fileRepositoty {
	return &fileRepositoty{
		client:        client,
		filesMapper:   fileMapper,
		visitorMapper: visitorMapper,
	}
}

func (r *fileRepositoty) Get(ctx context.Context, fileId uuid.UUID, session *auth.Session) (*entityFile.File, error) {
	if session == nil {
		return nil, fmt.Errorf("fileRepository.Get: %w", diterrors.NewValidationError(ErrSessionIsEmpty, diterrors.ErrValidationFields{
			Field:   "session",
			Message: ErrSessionIsEmpty.Error(),
		}))
	}

	if fileId == uuid.Nil {
		return nil, fmt.Errorf("fileRepository.Get: %w", diterrors.NewValidationError(ErrFileIdIsEmpty, diterrors.ErrValidationFields{
			Field:   "fileId",
			Message: ErrFileIdIsEmpty.Error(),
		}))
	}

	visitor := r.visitorMapper.SessionToVisitorPb(session)
	if visitor == nil {
		return nil, fmt.Errorf("visitorMapper.SessionToVisitorPb: %w", diterrors.NewValidationError(ErrVisitorIsEmpty, diterrors.ErrValidationFields{
			Field:   "visitor",
			Message: "empty visitor",
		}))
	}

	resp, err := r.client.Get(ctx, &filev1.GetRequest{
		Id:      fileId.String(),
		Visitor: visitor,
	})
	if err != nil {
		return nil, fmt.Errorf("client.Get: %w", diterrors.GrpcErrorToError(err))
	}

	result := r.filesMapper.FilePbToEntity(resp.GetFile())
	if result == nil {
		return nil, fmt.Errorf("fileMapper.FilePbToEntity: %w", ErrFileIsEmpty)
	}

	if result.Id == uuid.Nil {
		return nil, fmt.Errorf("fileMapper.FilePbToEntity: %w", ErrFileIdIsEmpty)
	}

	return result, nil
}
