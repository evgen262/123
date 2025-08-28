package portals

import (
	"testing"
	"time"

	imagesv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/portals/images/v1"
	timeUtils "git.mos.ru/buch-cloud/moscow-team-2.0/build/time-utils.git"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/timestamppb"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portal"
)

func Test_portalMapper_MapImageTypeToPb(t *testing.T) {
	type fields struct {
		timeUtils *timeUtils.MockTimeUtils
	}
	type args struct {
		imageType portal.ImageType
	}

	tests := []struct {
		name string
		args args
		want imagesv1.ImageType
	}{
		{
			name: "invalid",
			args: args{
				imageType: 0,
			},
			want: imagesv1.ImageType_IMAGE_TYPE_INVALID,
		},
		{
			name: "jpg",
			args: args{
				imageType: 1,
			},
			want: imagesv1.ImageType_IMAGE_TYPE_JPG,
		},
		{
			name: "png",
			args: args{
				imageType: 2,
			},
			want: imagesv1.ImageType_IMAGE_TYPE_PNG,
		},
		{
			name: "svg",
			args: args{
				imageType: 3,
			},
			want: imagesv1.ImageType_IMAGE_TYPE_SVG,
		},
		{
			name: "gif",
			args: args{
				imageType: 4,
			},
			want: imagesv1.ImageType_IMAGE_TYPE_GIF,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{timeUtils.NewMockTimeUtils(ctrl)}
			pfm := NewImagesMapper(f.timeUtils)
			assert.Equalf(t, tt.want, pfm.mapImageTypeToPb(tt.args.imageType), "ToPb(%v)", tt.args.imageType)
		})
	}
}

func Test_portalMapper_MapImageTypeToEntity(t *testing.T) {
	type fields struct {
		timeUtils *timeUtils.MockTimeUtils
	}
	type args struct {
		imageType imagesv1.ImageType
	}

	tests := []struct {
		name string
		args args
		want portal.ImageType
	}{
		{
			name: "invalid",
			args: args{
				imageType: 0,
			},
			want: portal.ImageTypeUnknown,
		},
		{
			name: "jpg",
			args: args{
				imageType: 1,
			},
			want: portal.ImageTypeJpeg,
		},
		{
			name: "png",
			args: args{
				imageType: 2,
			},
			want: portal.ImageTypePng,
		},
		{
			name: "svg",
			args: args{
				imageType: 3,
			},
			want: portal.ImageTypeSvg,
		},
		{
			name: "gif",
			args: args{
				imageType: 4,
			},
			want: portal.ImageTypeGif,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{timeUtils.NewMockTimeUtils(ctrl)}
			pfm := NewImagesMapper(f.timeUtils)
			assert.Equalf(t, tt.want, pfm.mapImageTypeToEntity(tt.args.imageType), "ToPb(%v)", tt.args.imageType)
		})
	}
}

func TestPortalsMapper_NewImageToPb(t *testing.T) {
	type fields struct {
		timeUtils *timeUtils.MockTimeUtils
	}
	type args struct {
		image *portal.Image
	}

	testT := time.Date(2023, 10, 31, 13, 00, 00, 00, time.UTC)

	tests := []struct {
		name string
		args args
		want func(a args, f fields) *imagesv1.AddRequest
	}{
		{
			name: "correct",
			args: args{
				image: &portal.Image{
					Id:        1,
					Name:      "Test ImageID Name",
					Path:      "Test ImageID Path",
					Data:      []byte{},
					Type:      1,
					CreatedAt: &testT,
					UpdatedAt: &testT,
				},
			},
			want: func(a args, f fields) *imagesv1.AddRequest {
				return &imagesv1.AddRequest{
					Name:  "Test ImageID Name",
					Type:  imagesv1.ImageType_IMAGE_TYPE_JPG,
					Image: a.image.Data,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{timeUtils.NewMockTimeUtils(ctrl)}
			pm := NewImagesMapper(f.timeUtils)
			want := tt.want(tt.args, f)
			got := pm.NewImageToPb(tt.args.image)

			assert.Equal(t, want, got)
		})
	}
}

func TestPortalsMapper_ImageToPb(t *testing.T) {
	type fields struct {
		timeUtils *timeUtils.MockTimeUtils
	}
	type args struct {
		image *portal.Image
	}

	testT := time.Date(2023, 10, 31, 13, 00, 00, 00, time.UTC)

	tests := []struct {
		name string
		args args
		want func(a args, f fields) *imagesv1.Image
	}{
		{
			name: "correct",
			args: args{
				image: &portal.Image{
					Id:        1,
					Name:      "Test ImageID Name",
					Path:      "Test ImageID Path",
					Data:      []byte{},
					Type:      1,
					CreatedAt: &testT,
					UpdatedAt: &testT,
				},
			},
			want: func(a args, f fields) *imagesv1.Image {
				t := timestamppb.New(testT)
				f.timeUtils.EXPECT().TimeToTimestamp(a.image.CreatedAt).Return(t)
				f.timeUtils.EXPECT().TimeToTimestamp(a.image.UpdatedAt).Return(t)
				return &imagesv1.Image{
					Id:          1,
					Name:        "Test ImageID Name",
					Path:        "Test ImageID Path",
					Type:        imagesv1.ImageType_IMAGE_TYPE_JPG,
					Image:       a.image.Data,
					CreatedTime: t,
					UpdatedTime: t,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{timeUtils.NewMockTimeUtils(ctrl)}
			pm := NewImagesMapper(f.timeUtils)
			want := tt.want(tt.args, f)
			got := pm.ImageToPb(tt.args.image)

			assert.Equal(t, want, got)
		})
	}
}

func TestPortalsMapper_ImagesToPb(t *testing.T) {
	type fields struct {
		timeUtils *timeUtils.MockTimeUtils
	}
	type args struct {
		images []*portal.Image
	}

	testT := time.Date(2023, 10, 31, 13, 00, 00, 00, time.UTC)

	tests := []struct {
		name string
		args args
		want func(a args, f fields) []*imagesv1.Image
	}{
		{
			name: "correct",
			args: args{
				images: []*portal.Image{{
					Id:        1,
					Name:      "Test ImageID Name",
					Path:      "Test ImageID Path",
					Data:      []byte{},
					Type:      1,
					CreatedAt: &testT,
					UpdatedAt: &testT,
				}},
			},
			want: func(a args, f fields) []*imagesv1.Image {
				t := timestamppb.New(testT)
				f.timeUtils.EXPECT().TimeToTimestamp(a.images[0].CreatedAt).Return(t)
				f.timeUtils.EXPECT().TimeToTimestamp(a.images[0].UpdatedAt).Return(t)
				return []*imagesv1.Image{{
					Id:          1,
					Name:        "Test ImageID Name",
					Path:        "Test ImageID Path",
					Type:        1,
					Image:       a.images[0].Data,
					CreatedTime: t,
					UpdatedTime: t,
				}}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{timeUtils.NewMockTimeUtils(ctrl)}
			pm := NewImagesMapper(f.timeUtils)
			want := tt.want(tt.args, f)
			got := pm.ImagesToPb(tt.args.images)

			assert.Equal(t, want, got)
		})
	}
}

func TestPortalsMapper_ImageToEntity(t *testing.T) {
	type fields struct {
		timeUtils *timeUtils.MockTimeUtils
	}
	type args struct {
		imagePb *imagesv1.Image
	}

	testT := time.Date(2023, 10, 31, 13, 00, 00, 00, time.UTC)
	testTime := timestamppb.New(testT)

	tests := []struct {
		name string
		args args
		want func(a args, f fields) *portal.Image
	}{
		{
			name: "correct",
			args: args{
				imagePb: &imagesv1.Image{
					Id:          1,
					Name:        "Test ImageID Name",
					Path:        "Test ImageID Path",
					Type:        1,
					Image:       []byte{},
					CreatedTime: testTime,
					UpdatedTime: testTime,
				},
			},
			want: func(a args, f fields) *portal.Image {
				f.timeUtils.EXPECT().TimestampToTime(a.imagePb.CreatedTime).Return(&testT)
				f.timeUtils.EXPECT().TimestampToTime(a.imagePb.UpdatedTime).Return(&testT)
				return &portal.Image{
					Id:        1,
					Name:      "Test ImageID Name",
					Path:      "Test ImageID Path",
					Data:      a.imagePb.Image,
					Type:      1,
					CreatedAt: &testT,
					UpdatedAt: &testT,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{timeUtils.NewMockTimeUtils(ctrl)}
			pm := NewImagesMapper(f.timeUtils)
			want := tt.want(tt.args, f)
			got := pm.ImageToEntity(tt.args.imagePb)

			assert.Equal(t, want, got)
		})
	}
}

func TestPortalsMapper_ImagesToEntity(t *testing.T) {
	type fields struct {
		timeUtils *timeUtils.MockTimeUtils
	}
	type args struct {
		imagesPb []*imagesv1.Image
	}

	testT := time.Date(2023, 10, 31, 13, 00, 00, 00, time.UTC)
	testTime := timestamppb.New(testT)

	tests := []struct {
		name string
		args args
		want func(a args, f fields) []*portal.Image
	}{
		{
			name: "correct",
			args: args{
				imagesPb: []*imagesv1.Image{{
					Id:          1,
					Name:        "Test ImageID Name",
					Path:        "Test ImageID Path",
					Type:        1,
					Image:       []byte{},
					CreatedTime: testTime,
					UpdatedTime: testTime,
				}},
			},
			want: func(a args, f fields) []*portal.Image {
				f.timeUtils.EXPECT().TimestampToTime(a.imagesPb[0].CreatedTime).Return(&testT)
				f.timeUtils.EXPECT().TimestampToTime(a.imagesPb[0].UpdatedTime).Return(&testT)
				return []*portal.Image{{
					Id:        1,
					Name:      "Test ImageID Name",
					Path:      "Test ImageID Path",
					Data:      a.imagesPb[0].Image,
					Type:      1,
					CreatedAt: &testT,
					UpdatedAt: &testT,
				}}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{timeUtils.NewMockTimeUtils(ctrl)}
			pm := NewImagesMapper(f.timeUtils)
			want := tt.want(tt.args, f)
			got := pm.ImagesToEntity(tt.args.imagesPb)

			assert.Equal(t, want, got)
		})
	}
}
