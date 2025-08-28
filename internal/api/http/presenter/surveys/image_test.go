package surveys

import (
	"testing"

	viewSurveys "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/surveys"
	entitySurvey "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/survey"

	"github.com/stretchr/testify/assert"
)

func TestSurveyImagesPresenter_ToNewEntity(t *testing.T) {
	type fields struct {
	}

	type args struct {
		image *viewSurveys.NewSurveyImageObject
	}

	tests := []struct {
		name string
		args args
		want func(a args, f fields) *entitySurvey.Image
	}{
		{
			name: "correct",
			args: args{
				image: &viewSurveys.NewSurveyImageObject{},
			},
			want: func(a args, f fields) *entitySurvey.Image {
				return &entitySurvey.Image{
					ExternalImageInfo: &entitySurvey.ExternalProperties{},
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fields{}
			want := tt.want(tt.args, f)
			s := NewSurveyImagesPresenter()
			got := s.ToNewEntity(tt.args.image)
			assert.Equal(t, want, got)
		})
	}
}

func TestSurveyImagesPresenter_ToView(t *testing.T) {
	type fields struct {
	}

	type args struct {
		image *entitySurvey.Image
	}

	tests := []struct {
		name string
		args args
		want func(a args, f fields) *viewSurveys.SurveyImageObject
	}{
		{
			name: "correct",
			args: args{
				image: &entitySurvey.Image{
					ExternalImageInfo: &entitySurvey.ExternalProperties{},
				},
			},
			want: func(a args, f fields) *viewSurveys.SurveyImageObject {
				return &viewSurveys.SurveyImageObject{
					ExternalImageInfo: viewSurveys.SurveyExternalImageObject{},
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fields{}
			want := tt.want(tt.args, f)
			s := NewSurveyImagesPresenter()
			got := s.ToView(tt.args.image)
			assert.Equal(t, want, got)
		})
	}
}
