package files

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/auth"
	entityFile "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/file"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/usecase"
)

func Test_fileUsecase_Get(t *testing.T) {
	type fields struct {
		fileRepository *MockFileRepository
		logger         *ditzap.MockLogger
	}
	type args struct {
		ctx    context.Context
		fileId uuid.UUID
	}

	ctx := context.TODO()

	tests := []struct {
		name string
		args args
		want func(f fields, a args) (*entityFile.File, error)
	}{
		{
			name: "success",
			args: args{
				ctx:    entity.WithSession(ctx, &auth.Session{}),
				fileId: uuid.New(),
			},
			want: func(f fields, a args) (*entityFile.File, error) {
				testSession, err := entity.SessionFromContext(a.ctx)
				if err != nil {
					return nil, err
				}
				file := &entityFile.File{}
				f.fileRepository.EXPECT().Get(a.ctx, a.fileId, testSession).Return(file, nil)
				return file, nil
			},
		},
		{
			name: "get session error",
			args: args{
				ctx:    ctx,
				fileId: uuid.New(),
			},
			want: func(f fields, a args) (*entityFile.File, error) {
				testErr := errors.New("session context not found")
				f.logger.EXPECT().Error(usecase.ErrGetSessionFromContext.Error(),
					ditzap.UUID("file_id", a.fileId),
					zap.Error(testErr),
				)
				return nil, diterrors.ErrUnauthenticated
			},
		},
		{
			name: "error unauthenticated",
			args: args{
				ctx:    entity.WithSession(ctx, &auth.Session{}),
				fileId: uuid.New(),
			},
			want: func(f fields, a args) (*entityFile.File, error) {
				testSession, err := entity.SessionFromContext(a.ctx)
				if err != nil {
					return nil, err
				}
				f.fileRepository.EXPECT().Get(a.ctx, a.fileId, testSession).Return(nil, diterrors.ErrUnauthenticated)
				f.logger.EXPECT().Debug("fileUsecase.Get: can't get public file from repository", zap.Error(diterrors.ErrUnauthenticated))
				return nil, fmt.Errorf("can't get public file from repository: %w", diterrors.ErrUnauthenticated)
			},
		},
		{
			name: "error permission denied",
			args: args{
				ctx:    entity.WithSession(ctx, &auth.Session{}),
				fileId: uuid.New(),
			},
			want: func(f fields, a args) (*entityFile.File, error) {
				testSession, err := entity.SessionFromContext(a.ctx)
				if err != nil {
					return nil, err
				}
				f.fileRepository.EXPECT().Get(a.ctx, a.fileId, testSession).Return(nil, diterrors.ErrPermissionDenied)
				f.logger.EXPECT().Debug("fileUsecase.Get: can't get public file from repository", zap.Error(diterrors.ErrPermissionDenied))
				return nil, fmt.Errorf("can't get public file from repository: %w", diterrors.ErrPermissionDenied)
			},
		},
		{
			name: "error not found",
			args: args{
				ctx:    entity.WithSession(ctx, &auth.Session{}),
				fileId: uuid.New(),
			},
			want: func(f fields, a args) (*entityFile.File, error) {
				testSession, err := entity.SessionFromContext(a.ctx)
				if err != nil {
					return nil, err
				}
				f.fileRepository.EXPECT().Get(a.ctx, a.fileId, testSession).Return(nil, diterrors.ErrNotFound)
				f.logger.EXPECT().Debug("fileUsecase.Get: can't get public file from repository", zap.Error(diterrors.ErrNotFound))
				return nil, fmt.Errorf("can't get public file from repository: %w", diterrors.ErrNotFound)
			},
		},
		{
			name: "error unimplemented",
			args: args{
				ctx:    entity.WithSession(ctx, &auth.Session{}),
				fileId: uuid.New(),
			},
			want: func(f fields, a args) (*entityFile.File, error) {
				testSession, err := entity.SessionFromContext(a.ctx)
				if err != nil {
					return nil, err
				}
				f.fileRepository.EXPECT().Get(a.ctx, a.fileId, testSession).Return(nil, diterrors.ErrUnimplemented)
				f.logger.EXPECT().Debug("fileUsecase.Get: can't get public file from repository", zap.Error(diterrors.ErrUnimplemented))
				return nil, fmt.Errorf("can't get public file from repository: %w", diterrors.ErrUnimplemented)
			},
		},
		{
			name: "error validation",
			args: args{
				ctx:    entity.WithSession(ctx, &auth.Session{}),
				fileId: uuid.New(),
			},
			want: func(f fields, a args) (*entityFile.File, error) {
				validationErr := diterrors.NewValidationError(assert.AnError)
				testSession, err := entity.SessionFromContext(a.ctx)
				if err != nil {
					return nil, err
				}
				f.fileRepository.EXPECT().Get(a.ctx, a.fileId, testSession).Return(nil, validationErr)
				f.logger.EXPECT().Debug("fileUsecase.Get: can't get public file from repository", zap.Error(validationErr))
				return nil, fmt.Errorf("can't get public file from repository: %w", validationErr)
			},
		},
		{
			name: "error unknown",
			args: args{
				ctx:    entity.WithSession(ctx, &auth.Session{}),
				fileId: uuid.New(),
			},
			want: func(f fields, a args) (*entityFile.File, error) {
				testSession, err := entity.SessionFromContext(a.ctx)
				if err != nil {
					return nil, err
				}
				f.fileRepository.EXPECT().Get(a.ctx, a.fileId, testSession).Return(nil, assert.AnError)
				f.logger.EXPECT().Error("fileUsecase.Get: can't get public file from repository", zap.Error(assert.AnError))
				return nil, fmt.Errorf("can't get public file from repository: %w", assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			f := fields{
				fileRepository: NewMockFileRepository(ctrl),
				logger:         ditzap.NewMockLogger(ctrl),
			}

			expectedFile, expectedErr := tt.want(f, tt.args)

			uc := NewFileUsecase(
				f.fileRepository,
				f.logger,
			)

			actualFile, actualErr := uc.Get(tt.args.ctx, tt.args.fileId)
			if expectedErr != nil {
				assert.Error(t, actualErr)
				assert.EqualError(t, actualErr, expectedErr.Error())
			} else {
				assert.NoError(t, actualErr)
				assert.Equal(t, expectedFile, actualFile)
			}
		})
	}
}
