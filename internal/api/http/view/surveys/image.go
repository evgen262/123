package view

import (
	"github.com/google/uuid"
)

type SurveyImageObject struct {
	ID                string                    `json:"id"`
	ExternalImageInfo SurveyExternalImageObject `json:"external_image_info"`
}

type SurveyExternalImageObject struct {
	ID       uuid.UUID `json:"id"`
	FileName string    `json:"file_name"`
	URL      string    `json:"url"`
	Size     int       `json:"size"`
}

type NewSurveyImageObject struct {
	Payload           string                    `json:"payload"`
	ExternalImageInfo SurveyExternalImageObject `json:"external_image_info"`
}
