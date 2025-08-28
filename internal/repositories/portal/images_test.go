package portal

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	imagesv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/portals/images/v1"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portal"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/repositories"
)

func Test_portalImagesRepository_All(t *testing.T) {
	type fields struct {
		client *imagesv1.MockImagesAPIClient
		mapper *MockImagesMapper
	}
	type args struct {
		ctx context.Context
	}
	ctx := context.TODO()
	testT := time.Now()
	testTime := timestamppb.New(testT)
	testErr := errors.New("error")
	tests := []struct {
		name string
		args args
		want func(a args, f fields) ([]*portal.Image, error)
	}{
		{
			name: "correct",
			args: args{
				ctx: ctx,
			},
			want: func(a args, f fields) ([]*portal.Image, error) {
				images := []*imagesv1.Image{
					{
						Id:          1,
						Name:        "ImageID Test Name",
						Path:        "ImageID Test Path",
						Type:        imagesv1.ImageType_IMAGE_TYPE_GIF,
						Image:       []byte{1, 1, 1, 1},
						CreatedTime: testTime,
						UpdatedTime: testTime,
					},
				}
				entityImage := []*portal.Image{{
					Id:        1,
					Name:      "ImageID Test Name",
					Path:      "ImageID Test Path",
					Data:      portal.ImageData{1, 1, 1, 1},
					Type:      portal.ImageTypeGif,
					CreatedAt: &testT,
					UpdatedAt: &testT,
				}}
				f.client.EXPECT().All(a.ctx, &imagesv1.AllRequest{}).Return(&imagesv1.AllResponse{
					Images: images,
				}, nil)
				f.mapper.EXPECT().ImagesToEntity(images).Return(entityImage)
				return entityImage, nil
			},
		},
		{
			name: "get all images from portal service error",
			args: args{
				ctx: ctx,
			},
			want: func(a args, f fields) ([]*portal.Image, error) {
				f.client.EXPECT().All(a.ctx, &imagesv1.AllRequest{}).Return(nil, testErr)
				return nil, fmt.Errorf("can't get all images: %w", testErr)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				client: imagesv1.NewMockImagesAPIClient(ctrl),
				mapper: NewMockImagesMapper(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			repo := NewImagesRepository(f.client, f.mapper)
			got, err := repo.All(tt.args.ctx)
			if wantErr != nil {
				assert.Empty(t, got)
				assert.EqualError(t, err, wantErr.Error())
			} else {
				assert.Equal(t, want, got)
				assert.NoError(t, err)
			}
		})
	}
}

func Test_portalImagesRepository_Get(t *testing.T) {
	type fields struct {
		client *imagesv1.MockImagesAPIClient
		mapper *MockImagesMapper
	}
	type args struct {
		ctx     context.Context
		imageId portal.ImageId
	}
	ctx := context.TODO()
	testT := time.Now()
	testTime := timestamppb.New(testT)
	testErr := errors.New("testErr")
	tests := []struct {
		name string
		args args
		want func(a args, f fields) (*portal.Image, error)
	}{
		{
			name: "correct",
			args: args{
				ctx: ctx,
			},
			want: func(a args, f fields) (*portal.Image, error) {
				image := &imagesv1.Image{
					Id:          1,
					Name:        "ImageID Test Name",
					Path:        "ImageID Test Path",
					Type:        imagesv1.ImageType_IMAGE_TYPE_GIF,
					Image:       []byte{1, 1, 1, 1},
					CreatedTime: testTime,
					UpdatedTime: testTime,
				}
				entityImage := &portal.Image{
					Id:        1,
					Name:      "ImageID Test Name",
					Path:      "ImageID Test Path",
					Data:      portal.ImageData{1, 1, 1, 1},
					Type:      portal.ImageTypeGif,
					CreatedAt: &testT,
					UpdatedAt: &testT,
				}
				f.client.EXPECT().Get(a.ctx, &imagesv1.GetRequest{
					Id: int32(a.imageId),
				}).Return(&imagesv1.GetResponse{
					Image: image,
				}, nil)
				f.mapper.EXPECT().ImageToEntity(image).Return(entityImage)
				return entityImage, nil
			},
		},
		{
			name: "get image from portal service error",
			args: args{
				ctx: ctx,
			},
			want: func(a args, f fields) (*portal.Image, error) {
				f.client.EXPECT().Get(a.ctx, &imagesv1.GetRequest{}).Return(nil, testErr)
				return nil, fmt.Errorf("can't get image: %w", testErr)
			},
		},
		{
			name: "get image not found",
			args: args{
				ctx:     ctx,
				imageId: 1,
			},
			want: func(a args, f fields) (*portal.Image, error) {
				f.client.EXPECT().Get(a.ctx, &imagesv1.GetRequest{Id: int32(a.imageId)}).Return(nil, status.Error(codes.NotFound, testErr.Error()))
				return nil, repositories.ErrNotFound
			},
		},
		{
			name: "get image invalid argument",
			args: args{
				ctx:     ctx,
				imageId: 1,
			},
			want: func(a args, f fields) (*portal.Image, error) {
				f.client.EXPECT().Get(a.ctx, &imagesv1.GetRequest{Id: int32(a.imageId)}).
					Return(nil, status.Error(codes.InvalidArgument, testErr.Error()))
				return nil, diterrors.NewValidationError(status.Error(codes.InvalidArgument, testErr.Error()))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				client: imagesv1.NewMockImagesAPIClient(ctrl),
				mapper: NewMockImagesMapper(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			repo := NewImagesRepository(f.client, f.mapper)
			got, err := repo.Get(tt.args.ctx, tt.args.imageId)
			if wantErr != nil {
				assert.Empty(t, got)
				assert.EqualError(t, err, wantErr.Error())
			} else {
				assert.Equal(t, want, got)
				assert.NoError(t, err)
			}
		})
	}
}

func Test_portalImagesRepository_GetImageData(t *testing.T) {
	type fields struct {
		client *imagesv1.MockImagesAPIClient
		mapper *MockImagesMapper
	}
	type args struct {
		ctx  context.Context
		path string
	}
	ctx := context.TODO()
	testErr := errors.New("error")
	tests := []struct {
		name string
		args args
		want func(a args, f fields) (portal.ImageData, error)
	}{
		{
			name: "correct",
			args: args{
				ctx:  ctx,
				path: "Test Path",
			},
			want: func(a args, f fields) (portal.ImageData, error) {
				image := []byte{1, 1, 1, 1}
				f.client.EXPECT().GetImage(a.ctx, &imagesv1.GetImageRequest{
					Path: a.path,
				}).Return(&imagesv1.GetImageResponse{
					Image: image,
				}, nil)
				return image, nil
			},
		},
		{
			name: "get image data from portal service error",
			args: args{
				ctx:  ctx,
				path: "Test Path",
			},
			want: func(a args, f fields) (portal.ImageData, error) {
				f.client.EXPECT().GetImage(a.ctx, &imagesv1.GetImageRequest{
					Path: a.path,
				}).Return(nil, testErr)
				return nil, fmt.Errorf("can't get image data: %w", testErr)
			},
		},
		{
			name: "get image data not found",
			args: args{
				ctx:  ctx,
				path: "Test Path",
			},
			want: func(a args, f fields) (portal.ImageData, error) {
				f.client.EXPECT().GetImage(a.ctx, &imagesv1.GetImageRequest{Path: a.path}).Return(nil, diterrors.NewApiError(codes.NotFound, testErr.Error(), testErr))
				return nil, repositories.ErrNotFound
			},
		},
		{
			name: "get image data invalid argument",
			args: args{
				ctx:  ctx,
				path: "Test Path",
			},
			want: func(a args, f fields) (portal.ImageData, error) {
				f.client.EXPECT().GetImage(a.ctx, &imagesv1.GetImageRequest{Path: a.path}).Return(nil, status.Error(codes.InvalidArgument, testErr.Error()))
				return nil, diterrors.NewValidationError(status.Error(codes.InvalidArgument, testErr.Error()))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				client: imagesv1.NewMockImagesAPIClient(ctrl),
				mapper: NewMockImagesMapper(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			repo := NewImagesRepository(f.client, f.mapper)
			got, err := repo.GetImageData(tt.args.ctx, tt.args.path)
			if wantErr != nil {
				assert.Empty(t, got)
				assert.EqualError(t, err, wantErr.Error())
			} else {
				assert.Equal(t, want, got)
				assert.NoError(t, err)
			}
		})
	}
}

func Test_portalImagesRepository_Add(t *testing.T) {
	type fields struct {
		client *imagesv1.MockImagesAPIClient
		mapper *MockImagesMapper
	}
	type args struct {
		ctx   context.Context
		image *portal.Image
	}
	ctx := context.TODO()
	testT := time.Now()
	testTime := timestamppb.New(testT)
	testErr := errors.New("error")
	tests := []struct {
		name string
		args args
		want func(a args, f fields) (*portal.Image, error)
	}{
		{
			name: "correct",
			args: args{
				ctx: ctx,
				image: &portal.Image{
					Id:        1,
					Name:      "ImageID Test Name",
					Path:      "ImageID Test Path",
					Data:      portal.ImageData{1, 1, 1, 1},
					Type:      portal.ImageTypeGif,
					CreatedAt: &testT,
					UpdatedAt: &testT,
				},
			},
			want: func(a args, f fields) (*portal.Image, error) {
				image := &imagesv1.Image{
					Id:          1,
					Name:        "ImageID Test Name",
					Path:        "ImageID Test Path",
					Type:        imagesv1.ImageType_IMAGE_TYPE_GIF,
					Image:       []byte{1, 1, 1, 1},
					CreatedTime: testTime,
					UpdatedTime: testTime,
				}
				entityImage := &portal.Image{
					Id:        1,
					Name:      "ImageID Test Name",
					Path:      "ImageID Test Path",
					Data:      portal.ImageData{1, 1, 1, 1},
					Type:      portal.ImageTypeGif,
					CreatedAt: &testT,
					UpdatedAt: &testT,
				}
				imageAddRequest := &imagesv1.AddRequest{
					Name:  "Test ImageID Name",
					Type:  imagesv1.ImageType_IMAGE_TYPE_GIF,
					Image: []byte{1, 1, 1, 1},
				}
				f.mapper.EXPECT().NewImageToPb(a.image).Return(imageAddRequest)
				f.client.EXPECT().Add(a.ctx, imageAddRequest).Return(&imagesv1.AddResponse{
					Image: image,
				}, nil)
				f.mapper.EXPECT().ImageToEntity(image).Return(entityImage)
				return entityImage, nil
			},
		},
		{
			name: "add image to db service error",
			args: args{
				ctx: ctx,
				image: &portal.Image{
					Id:        1,
					Name:      "ImageID Test Name",
					Path:      "ImageID Test Path",
					Data:      portal.ImageData{},
					Type:      portal.ImageTypeGif,
					CreatedAt: &testT,
					UpdatedAt: &testT,
				},
			},
			want: func(a args, f fields) (*portal.Image, error) {
				imageAddRequest := &imagesv1.AddRequest{
					Name:  "Test ImageID Name",
					Type:  imagesv1.ImageType_IMAGE_TYPE_GIF,
					Image: []byte{1, 1, 1, 1},
				}
				f.mapper.EXPECT().NewImageToPb(a.image).Return(imageAddRequest)
				f.client.EXPECT().Add(a.ctx, imageAddRequest).Return(nil, testErr)
				return nil, fmt.Errorf("can't add image to db: %w", testErr)
			},
		},
		{
			name: "add image not found",
			args: args{
				ctx: ctx,
				image: &portal.Image{
					Id:        1,
					Name:      "ImageID Test Name",
					Path:      "ImageID Test Path",
					Data:      portal.ImageData{1, 1, 1, 1},
					Type:      portal.ImageTypeGif,
					CreatedAt: &testT,
					UpdatedAt: &testT,
				},
			},
			want: func(a args, f fields) (*portal.Image, error) {
				imageAddRequest := &imagesv1.AddRequest{
					Name:  "Test ImageID Name",
					Type:  imagesv1.ImageType_IMAGE_TYPE_GIF,
					Image: []byte{1, 1, 1, 1},
				}
				f.mapper.EXPECT().NewImageToPb(a.image).Return(imageAddRequest)
				f.client.EXPECT().Add(a.ctx, imageAddRequest).Return(nil, status.Error(codes.NotFound, testErr.Error()))
				return nil, repositories.ErrNotFound
			},
		},
		{
			name: "add image invalid argument",
			args: args{
				ctx: ctx,
				image: &portal.Image{
					Id:        1,
					Name:      "ImageID Test Name",
					Path:      "ImageID Test Path",
					Data:      portal.ImageData{1, 1, 1, 1},
					Type:      portal.ImageTypeGif,
					CreatedAt: &testT,
					UpdatedAt: &testT,
				},
			},
			want: func(a args, f fields) (*portal.Image, error) {
				imageAddRequest := &imagesv1.AddRequest{
					Name:  "Test ImageID Name",
					Type:  imagesv1.ImageType_IMAGE_TYPE_GIF,
					Image: []byte{1, 1, 1, 1},
				}
				f.mapper.EXPECT().NewImageToPb(a.image).Return(imageAddRequest)
				f.client.EXPECT().Add(a.ctx, imageAddRequest).Return(nil, status.Error(codes.InvalidArgument, testErr.Error()))
				return nil, diterrors.NewValidationError(status.Error(codes.InvalidArgument, testErr.Error()))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				client: imagesv1.NewMockImagesAPIClient(ctrl),
				mapper: NewMockImagesMapper(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			repo := NewImagesRepository(f.client, f.mapper)
			got, err := repo.Add(tt.args.ctx, tt.args.image)
			if wantErr != nil {
				assert.Empty(t, got)
				assert.EqualError(t, err, wantErr.Error())
			} else {
				assert.Equal(t, want, got)
				assert.NoError(t, err)
			}
		})
	}
}

func Test_portalImagesRepository_Delete(t *testing.T) {
	type fields struct {
		client *imagesv1.MockImagesAPIClient
		mapper *MockImagesMapper
	}
	type args struct {
		ctx     context.Context
		imageId portal.ImageId
	}
	ctx := context.TODO()
	testErr := errors.New("error")
	tests := []struct {
		name string
		args args
		want func(a args, f fields) error
	}{
		{
			name: "correct",
			args: args{
				ctx:     ctx,
				imageId: 1,
			},
			want: func(a args, f fields) error {
				f.client.EXPECT().Delete(a.ctx, &imagesv1.DeleteRequest{
					Id: int32(a.imageId),
				}).Return(nil, nil)
				return nil
			},
		},
		{
			name: "delete image service error",
			args: args{
				ctx:     ctx,
				imageId: 1,
			},
			want: func(a args, f fields) error {
				f.client.EXPECT().Delete(a.ctx, &imagesv1.DeleteRequest{
					Id: int32(a.imageId),
				}).Return(nil, testErr)
				return fmt.Errorf("can't delete image: %w", testErr)
			},
		},
		{
			name: "delete image not found",
			args: args{
				ctx:     ctx,
				imageId: 1,
			},
			want: func(a args, f fields) error {
				f.client.EXPECT().Delete(a.ctx, &imagesv1.DeleteRequest{Id: int32(a.imageId)}).Return(nil, status.Error(codes.NotFound, testErr.Error()))
				return repositories.ErrNotFound
			},
		},
		{
			name: "delete image invalid argument",
			args: args{
				ctx:     ctx,
				imageId: 1,
			},
			want: func(a args, f fields) error {
				f.client.EXPECT().Delete(a.ctx, &imagesv1.DeleteRequest{Id: int32(a.imageId)}).Return(nil, status.Error(codes.InvalidArgument, testErr.Error()))
				return diterrors.NewValidationError(status.Error(codes.InvalidArgument, testErr.Error()))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				client: imagesv1.NewMockImagesAPIClient(ctrl),
				mapper: NewMockImagesMapper(ctrl),
			}
			wantErr := tt.want(tt.args, f)
			repo := NewImagesRepository(f.client, f.mapper)
			err := repo.Delete(tt.args.ctx, tt.args.imageId)
			assert.Equal(t, wantErr, err)
		})
	}
}
