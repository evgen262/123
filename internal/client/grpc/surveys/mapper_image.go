package surveys

import (
	"fmt"

	imagev1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/surveys/image/v1"
	surveys "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/survey"
	"github.com/google/uuid"
)

type imageMapper struct {
}

func NewImageMapper() *imageMapper {
	return &imageMapper{}
}

func (im imageMapper) NewImageToPb(image *surveys.Image) *imagev1.AddRequest_Image {
	newImage := &imagev1.AddRequest_Image{
		Payload: image.Data,
	}
	if image.ExternalImageInfo != nil {
		newImage.ExternalImageInfo = &imagev1.AddRequest_Image_ExternalInfo{
			Id:       image.ExternalImageInfo.ID.String(),
			Filename: image.ExternalImageInfo.FileName,
			Url:      image.ExternalImageInfo.URL,
			Size:     image.ExternalImageInfo.Size,
		}
	}

	return newImage
}

func (im imageMapper) ImageToEntity(image *imagev1.Image) (*surveys.Image, error) {
	externalInfoId, err := uuid.Parse(image.GetExternalImageInfo().GetId())
	if err != nil {
		return nil, fmt.Errorf("can't parse uuid: %w", err)
	}
	newImage := &surveys.Image{
		ID: surveys.ImageID(image.Id),
		ExternalImageInfo: &surveys.ExternalProperties{
			ID:       externalInfoId,
			FileName: image.GetExternalImageInfo().GetFilename(),
			URL:      image.GetExternalImageInfo().GetUrl(),
			Size:     image.GetExternalImageInfo().GetSize(),
		},
	}

	return newImage, nil
}
