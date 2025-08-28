package files

import (
	"context"
	"errors"
	"fmt"

	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity"
	entityFile "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/file"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/usecase"
)

type fileUsecase struct {
	repository FileRepository
	logger     ditzap.Logger
}

func NewFileUsecase(repository FileRepository, logger ditzap.Logger) *fileUsecase {
	return &fileUsecase{repository: repository, logger: logger}
}

func (fuc fileUsecase) Get(ctx context.Context, fileId uuid.UUID) (*entityFile.File, error) {
	session, err := entity.SessionFromContext(ctx)
	if err != nil {
		fuc.logger.Error(usecase.ErrGetSessionFromContext.Error(),
			ditzap.UUID("file_id", fileId),
			zap.Error(err),
		)
		return nil, diterrors.ErrUnauthenticated
	}

	result, err := fuc.repository.Get(ctx, fileId, session)
	if err != nil {
		switch {
		case errors.Is(err, diterrors.ErrUnauthenticated):
			fallthrough
		case errors.Is(err, diterrors.ErrPermissionDenied):
			fallthrough
		case errors.Is(err, diterrors.ErrNotFound):
			fallthrough
		case errors.Is(err, diterrors.ErrUnimplemented):
			fallthrough
		case errors.As(err, new(diterrors.ValidationError)):
			fuc.logger.Debug("fileUsecase.Get: can't get public file from repository", zap.Error(err))
		default:
			fuc.logger.Error("fileUsecase.Get: can't get public file from repository", zap.Error(err))
		}
		return nil, fmt.Errorf("can't get public file from repository: %w", err)
	}

	return result, nil
}
