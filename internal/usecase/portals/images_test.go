package portals

import (
	"context"
	"fmt"
	"testing"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portal"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func Test_imagesUseCase_AddImage(t *testing.T) {
	type fields struct {
		repo       *MockImagesRepository
		uploadPath string
		logger     *ditzap.MockLogger
	}
	type args struct {
		ctx   context.Context
		image *portal.Image
	}
	ctx := context.TODO()
	testErr := fmt.Errorf("testErr")
	tests := []struct {
		name string
		args args
		want func(a args, f fields) (*portal.Image, error)
	}{
		{
			name: "err",
			args: args{
				ctx:   ctx,
				image: &portal.Image{Name: "test"},
			},
			want: func(a args, f fields) (*portal.Image, error) {
				f.repo.EXPECT().Add(a.ctx, a.image).Return(nil, testErr)
				f.logger.EXPECT().Error("can't add image into repo", zap.Error(testErr))
				return nil, fmt.Errorf("can't add image: %w", testErr)
			},
		},
		{
			name: "correct",
			args: args{
				ctx:   ctx,
				image: &portal.Image{Name: "test"},
			},
			want: func(a args, f fields) (*portal.Image, error) {
				f.repo.EXPECT().Add(a.ctx, a.image).Return(&portal.Image{Id: 1, Name: a.image.Name}, nil)
				return &portal.Image{Id: 1, Name: a.image.Name}, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				repo:       NewMockImagesRepository(ctrl),
				uploadPath: "uploads",
				logger:     ditzap.NewMockLogger(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			fuc := NewImageUseCase(f.repo, f.uploadPath, f.logger)
			got, err := fuc.Add(tt.args.ctx, tt.args.image)
			if wantErr != nil {
				assert.EqualError(t, err, wantErr.Error())
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, want, got)
			}
		})
	}
}

func Test_imagesUseCase_DeleteImage(t *testing.T) {
	type fields struct {
		repo       *MockImagesRepository
		uploadPath string
		logger     *ditzap.MockLogger
	}
	type args struct {
		ctx     context.Context
		imageId int
	}
	ctx := context.TODO()
	testErr := fmt.Errorf("testErr")
	tests := []struct {
		name string
		args args
		want func(a args, f fields) error
	}{
		{
			name: "err",
			args: args{
				ctx:     ctx,
				imageId: 1,
			},
			want: func(a args, f fields) error {
				f.repo.EXPECT().Delete(a.ctx, portal.ImageId(a.imageId)).Return(testErr)
				f.logger.EXPECT().Error("can't delete image from repo", zap.Error(testErr))
				return fmt.Errorf("can't delete image: %w", testErr)
			},
		},
		{
			name: "correct",
			args: args{
				ctx:     ctx,
				imageId: 1,
			},
			want: func(a args, f fields) error {
				f.repo.EXPECT().Delete(a.ctx, portal.ImageId(a.imageId)).Return(nil)
				return nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				repo:       NewMockImagesRepository(ctrl),
				uploadPath: "uploads",
				logger:     ditzap.NewMockLogger(ctrl),
			}
			wantErr := tt.want(tt.args, f)
			fuc := NewImageUseCase(f.repo, f.uploadPath, f.logger)
			err := fuc.Delete(tt.args.ctx, tt.args.imageId)
			if wantErr != nil {
				assert.EqualError(t, err, wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_imagesUseCase_GetAllImages(t *testing.T) {
	type fields struct {
		repo       *MockImagesRepository
		uploadPath string
		logger     *ditzap.MockLogger
	}
	type args struct {
		ctx context.Context
	}
	ctx := context.TODO()
	testErr := fmt.Errorf("testErr")
	tests := []struct {
		name string
		args args
		want func(a args, f fields) ([]*portal.Image, error)
	}{
		{
			name: "err",
			args: args{
				ctx: ctx,
			},
			want: func(a args, f fields) ([]*portal.Image, error) {
				f.repo.EXPECT().All(a.ctx).Return(nil, testErr)
				f.logger.EXPECT().Error("can't get all images from repo", zap.Error(testErr))
				return nil, fmt.Errorf("can't get all images: %w", testErr)
			},
		},
		{
			name: "correct",
			args: args{
				ctx: ctx,
			},
			want: func(a args, f fields) ([]*portal.Image, error) {
				images := []*portal.Image{
					{Name: "test"},
				}
				f.repo.EXPECT().All(a.ctx).Return(images, nil)
				return images, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				repo:       NewMockImagesRepository(ctrl),
				uploadPath: "uploads",
				logger:     ditzap.NewMockLogger(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			fuc := NewImageUseCase(f.repo, f.uploadPath, f.logger)
			got, err := fuc.All(tt.args.ctx)
			if wantErr != nil {
				assert.EqualError(t, err, wantErr.Error())
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, want, got)
			}
		})
	}
}

func Test_imagesUseCase_GetImage(t *testing.T) {
	type fields struct {
		repo       *MockImagesRepository
		uploadPath string
		logger     *ditzap.MockLogger
	}
	type args struct {
		ctx     context.Context
		imageId int
	}
	ctx := context.TODO()
	testErr := fmt.Errorf("testErr")
	tests := []struct {
		name string
		args args
		want func(a args, f fields) (*portal.Image, error)
	}{
		{
			name: "err",
			args: args{
				ctx:     ctx,
				imageId: 1,
			},
			want: func(a args, f fields) (*portal.Image, error) {
				f.repo.EXPECT().Get(a.ctx, portal.ImageId(a.imageId)).Return(nil, testErr)
				f.logger.EXPECT().Error("can't get image from repo", zap.Error(testErr))
				return nil, fmt.Errorf("can't get image: %w", testErr)
			},
		},
		{
			name: "correct",
			args: args{
				ctx:     ctx,
				imageId: 1,
			},
			want: func(a args, f fields) (*portal.Image, error) {
				image := &portal.Image{
					Name: "test",
				}
				f.repo.EXPECT().Get(a.ctx, portal.ImageId(a.imageId)).Return(image, nil)
				return image, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				repo:       NewMockImagesRepository(ctrl),
				uploadPath: "uploads",
				logger:     ditzap.NewMockLogger(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			fuc := NewImageUseCase(f.repo, f.uploadPath, f.logger)
			got, err := fuc.Get(tt.args.ctx, tt.args.imageId)
			if wantErr != nil {
				assert.EqualError(t, err, wantErr.Error())
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, want, got)
			}
		})
	}
}

func Test_imagesUseCase_GetImageData(t *testing.T) {
	type fields struct {
		repo       *MockImagesRepository
		uploadPath string
		logger     *ditzap.MockLogger
	}
	type args struct {
		ctx  context.Context
		path string
	}
	ctx := context.TODO()
	testErr := fmt.Errorf("testErr")
	tests := []struct {
		name string
		args args
		want func(a args, f fields) (portal.ImageData, error)
	}{
		{
			name: "err",
			args: args{
				ctx:  ctx,
				path: "testPath",
			},
			want: func(a args, f fields) (portal.ImageData, error) {
				f.repo.EXPECT().GetImageData(a.ctx, a.path).Return(nil, testErr)
				f.logger.EXPECT().Error("can't get raw image from repo", zap.Error(testErr))
				return nil, fmt.Errorf("can't get raw image: %w", testErr)
			},
		},
		{
			name: "correct",
			args: args{
				ctx:  ctx,
				path: "testPath",
			},
			want: func(a args, f fields) (portal.ImageData, error) {
				imageData := portal.ImageData{123, 125}
				f.repo.EXPECT().GetImageData(a.ctx, a.path).Return(imageData, nil)
				return imageData, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				repo:       NewMockImagesRepository(ctrl),
				uploadPath: "uploads",
				logger:     ditzap.NewMockLogger(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			fuc := NewImageUseCase(f.repo, f.uploadPath, f.logger)
			got, err := fuc.GetRawImage(tt.args.ctx, tt.args.path)
			if wantErr != nil {
				assert.EqualError(t, err, wantErr.Error())
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, want, got)
			}
		})
	}
}
