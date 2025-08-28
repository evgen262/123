package surveys

import (
	"fmt"
	"testing"

	imagev1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/surveys/image/v1"
	surveys "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/survey"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestImageMapper_NewImageToPb(t *testing.T) {
	type args struct {
		image *surveys.Image
	}

	testUUID := uuid.New()
	tests := []struct {
		name string
		args args
		want func(a args) *imagev1.AddRequest_Image
	}{
		{
			name: "correct",
			args: args{image: &surveys.Image{
				Data: "test payload",
				ExternalImageInfo: &surveys.ExternalProperties{
					ID:       testUUID,
					FileName: "test filename",
					URL:      "test url",
					Size:     1,
				},
			}},
			want: func(a args) *imagev1.AddRequest_Image {
				return &imagev1.AddRequest_Image{
					Payload: "test payload",
					ExternalImageInfo: &imagev1.AddRequest_Image_ExternalInfo{
						Id:       testUUID.String(),
						Filename: "test filename",
						Url:      "test url",
						Size:     1,
					},
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewImageMapper()
			want := tt.want(tt.args)
			got := p.NewImageToPb(tt.args.image)
			assert.Equal(t, want, got)
		})
	}
}

func TestImageMapper_ImageToEntity(t *testing.T) {
	type args struct {
		image *imagev1.Image
	}
	testUUID := uuid.New()
	tests := []struct {
		name string
		args args
		want func(a args) (*surveys.Image, error)
	}{
		{
			name: "correct",
			args: args{image: &imagev1.Image{
				Id: "testID",
				ExternalImageInfo: &imagev1.Image_ExternalInfo{
					Id:       testUUID.String(),
					Filename: "test filename",
					Url:      "test url",
					Size:     1,
				},
			}},
			want: func(a args) (*surveys.Image, error) {
				return &surveys.Image{
					ID: "testID",
					ExternalImageInfo: &surveys.ExternalProperties{
						ID:       testUUID,
						FileName: "test filename",
						URL:      "test url",
						Size:     1,
					},
				}, nil
			},
		},
		{
			name: "parse err",
			args: args{image: &imagev1.Image{
				Id: "testID",
				ExternalImageInfo: &imagev1.Image_ExternalInfo{
					Id:       "test",
					Filename: "test filename",
					Url:      "test url",
					Size:     1,
				},
			}},
			want: func(a args) (*surveys.Image, error) {
				return nil, fmt.Errorf("can't parse uuid: invalid UUID length: 4")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewImageMapper()
			want, wantErr := tt.want(tt.args)
			got, err := p.ImageToEntity(tt.args.image)
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
