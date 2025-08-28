package files

import (
	"context"
	"errors"
	"fmt"
	"testing"

	filev1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/fileservice/file/v1"
	sharedv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/fileservice/shared/v1"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/auth"
	entityFile "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/file"
)

func Test_fileRepositoty_Get(t *testing.T) {
	type fields struct {
		client        *filev1.MockFileAPIClient
		filesMapper   *MockFilesMapper
		visitorMapper *MockVisitorMapper
	}
	type args struct {
		ctx     context.Context
		fileId  uuid.UUID
		session *auth.Session
	}
	tests := []struct {
		name string
		args args
		want func(f fields, a args) (*entityFile.File, error)
	}{
		{
			name: "success",
			args: args{
				ctx:     context.Background(),
				fileId:  uuid.New(),
				session: &auth.Session{},
			},
			want: func(f fields, a args) (*entityFile.File, error) {
				visitor := &sharedv1.Visitor{}
				filePb := &filev1.File{}
				fileEntity := &entityFile.File{
					Id: uuid.New(),
				}

				f.visitorMapper.EXPECT().SessionToVisitorPb(a.session).Return(visitor)
				f.client.EXPECT().Get(a.ctx, &filev1.GetRequest{
					Id:      a.fileId.String(),
					Visitor: visitor,
				}).Return(&filev1.GetResponse{File: filePb}, nil)
				f.filesMapper.EXPECT().FilePbToEntity(filePb).Return(fileEntity)

				return fileEntity, nil
			},
		},
		{
			name: "error session is nil",
			args: args{
				ctx:     context.Background(),
				fileId:  uuid.New(),
				session: nil,
			},
			want: func(f fields, a args) (*entityFile.File, error) {
				return nil, fmt.Errorf("fileRepository.Get: %w", diterrors.NewValidationError(ErrSessionIsEmpty, diterrors.ErrValidationFields{
					Field:   "session",
					Message: ErrSessionIsEmpty.Error(),
				}))
			},
		},
		{
			name: "error visitor is nil",
			args: args{
				ctx:     context.Background(),
				fileId:  uuid.New(),
				session: &auth.Session{},
			},
			want: func(f fields, a args) (*entityFile.File, error) {
				f.visitorMapper.EXPECT().SessionToVisitorPb(a.session).Return(nil)
				return nil, fmt.Errorf("visitorMapper.SessionToVisitorPb: %w", diterrors.NewValidationError(ErrVisitorIsEmpty, diterrors.ErrValidationFields{
					Field:   "visitor",
					Message: "empty visitor",
				}))
			},
		},
		{
			name: "error fileId is nil",
			args: args{
				ctx:     context.Background(),
				fileId:  uuid.Nil,
				session: &auth.Session{},
			},
			want: func(f fields, a args) (*entityFile.File, error) {
				return nil, fmt.Errorf("fileRepository.Get: %w", diterrors.NewValidationError(ErrFileIdIsEmpty, diterrors.ErrValidationFields{
					Field:   "fileId",
					Message: ErrFileIdIsEmpty.Error(),
				}))
			},
		},
		{
			name: "error client get fails",
			args: args{
				ctx:     context.Background(),
				fileId:  uuid.New(),
				session: &auth.Session{},
			},
			want: func(f fields, a args) (*entityFile.File, error) {
				visitor := &sharedv1.Visitor{}
				f.visitorMapper.EXPECT().SessionToVisitorPb(a.session).Return(visitor)
				f.client.EXPECT().Get(a.ctx, &filev1.GetRequest{
					Id:      a.fileId.String(),
					Visitor: visitor,
				}).Return(nil, errors.New("grpc error"))

				return nil, fmt.Errorf("client.Get: %w", diterrors.GrpcErrorToError(errors.New("grpc error")))
			},
		},
		{
			name: "error file is empty",
			args: args{
				ctx:     context.Background(),
				fileId:  uuid.New(),
				session: &auth.Session{},
			},
			want: func(f fields, a args) (*entityFile.File, error) {
				visitor := &sharedv1.Visitor{}
				f.visitorMapper.EXPECT().SessionToVisitorPb(a.session).Return(visitor)
				f.client.EXPECT().Get(a.ctx, &filev1.GetRequest{
					Id:      a.fileId.String(),
					Visitor: visitor,
				}).Return(&filev1.GetResponse{File: nil}, nil)
				f.filesMapper.EXPECT().FilePbToEntity(nil).Return(nil)

				return nil, fmt.Errorf("fileMapper.FilePbToEntity: %w", ErrFileIsEmpty)
			},
		},
		{
			name: "error uuid validation fails",
			args: args{
				ctx:     context.Background(),
				fileId:  uuid.New(),
				session: &auth.Session{},
			},
			want: func(f fields, a args) (*entityFile.File, error) {
				visitor := &sharedv1.Visitor{}
				filePb := &filev1.File{}
				fileEntity := &entityFile.File{Id: uuid.Nil}

				f.visitorMapper.EXPECT().SessionToVisitorPb(a.session).Return(visitor)
				f.client.EXPECT().Get(a.ctx, &filev1.GetRequest{
					Id:      a.fileId.String(),
					Visitor: visitor,
				}).Return(&filev1.GetResponse{File: filePb}, nil)
				f.filesMapper.EXPECT().FilePbToEntity(filePb).Return(fileEntity)

				return nil, fmt.Errorf("fileMapper.FilePbToEntity: %w", ErrFileIdIsEmpty)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				client:        filev1.NewMockFileAPIClient(ctrl),
				filesMapper:   NewMockFilesMapper(ctrl),
				visitorMapper: NewMockVisitorMapper(ctrl),
			}

			expectedFile, expectedErr := tt.want(f, tt.args)

			repo := NewFileRepository(
				f.client,
				f.filesMapper,
				f.visitorMapper,
			)

			actualFile, actualErr := repo.Get(tt.args.ctx, tt.args.fileId, tt.args.session)
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
