package surveys

import (
	viewSurveys "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/surveys"
	entitySurvey "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/survey"
)

type surveyImagesPresenter struct {
}

func NewSurveyImagesPresenter() *surveyImagesPresenter {
	return &surveyImagesPresenter{}
}

func (sip surveyImagesPresenter) ToNewEntity(image *viewSurveys.NewSurveyImageObject) *entitySurvey.Image {
	return &entitySurvey.Image{
		Data: image.Payload,
		ExternalImageInfo: &entitySurvey.ExternalProperties{
			ID:       image.ExternalImageInfo.ID,
			FileName: image.ExternalImageInfo.FileName,
			URL:      image.ExternalImageInfo.URL,
			Size:     int64(image.ExternalImageInfo.Size),
		},
	}
}

func (sip surveyImagesPresenter) ToView(image *entitySurvey.Image) *viewSurveys.SurveyImageObject {
	return &viewSurveys.SurveyImageObject{
		ID: string(image.ID),
		ExternalImageInfo: viewSurveys.SurveyExternalImageObject{
			ID:       image.ExternalImageInfo.ID,
			FileName: image.ExternalImageInfo.FileName,
			URL:      image.ExternalImageInfo.URL,
			Size:     int(image.ExternalImageInfo.Size),
		},
	}
}
