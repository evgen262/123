package portals

import (
	imagesv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/portals/images/v1"
	timeUtils "git.mos.ru/buch-cloud/moscow-team-2.0/build/time-utils.git"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portal"
)

type imagesMapper struct {
	timeUtils timeUtils.TimeUtils
}

func NewImagesMapper(timeUtils timeUtils.TimeUtils) *imagesMapper {
	return &imagesMapper{
		timeUtils: timeUtils,
	}
}

func (pfm imagesMapper) mapImageTypeToPb(imageType portal.ImageType) imagesv1.ImageType {
	switch imageType {
	case portal.ImageTypeJpeg:
		return imagesv1.ImageType_IMAGE_TYPE_JPG
	case portal.ImageTypePng:
		return imagesv1.ImageType_IMAGE_TYPE_PNG
	case portal.ImageTypeSvg:
		return imagesv1.ImageType_IMAGE_TYPE_SVG
	case portal.ImageTypeGif:
		return imagesv1.ImageType_IMAGE_TYPE_GIF
	}
	return imagesv1.ImageType_IMAGE_TYPE_INVALID
}

func (pfm imagesMapper) mapImageTypeToEntity(imageType imagesv1.ImageType) portal.ImageType {
	switch imageType {
	case imagesv1.ImageType_IMAGE_TYPE_JPG:
		return portal.ImageTypeJpeg
	case imagesv1.ImageType_IMAGE_TYPE_PNG:
		return portal.ImageTypePng
	case imagesv1.ImageType_IMAGE_TYPE_SVG:
		return portal.ImageTypeSvg
	case imagesv1.ImageType_IMAGE_TYPE_GIF:
		return portal.ImageTypeGif
	}
	return portal.ImageTypeUnknown
}
func (pfm imagesMapper) NewImageToPb(image *portal.Image) *imagesv1.AddRequest {
	return &imagesv1.AddRequest{
		Name:  image.Name,
		Type:  pfm.mapImageTypeToPb(image.Type),
		Image: image.Data,
	}
}

func (pfm imagesMapper) ImageToPb(image *portal.Image) *imagesv1.Image {
	return &imagesv1.Image{
		Id:          int32(image.Id),
		Name:        image.Name,
		Path:        image.Path,
		Type:        pfm.mapImageTypeToPb(image.Type),
		Image:       image.Data,
		CreatedTime: pfm.timeUtils.TimeToTimestamp(image.CreatedAt),
		UpdatedTime: pfm.timeUtils.TimeToTimestamp(image.UpdatedAt),
	}
}

func (pfm imagesMapper) ImagesToPb(images []*portal.Image) []*imagesv1.Image {
	imagesPb := make([]*imagesv1.Image, 0, len(images))

	for _, image := range images {
		imagesPb = append(imagesPb, pfm.ImageToPb(image))
	}

	return imagesPb
}

func (pfm imagesMapper) ImageToEntity(imagePb *imagesv1.Image) *portal.Image {
	return &portal.Image{
		Id:        portal.ImageId(imagePb.GetId()),
		Name:      imagePb.GetName(),
		Data:      imagePb.GetImage(),
		Path:      imagePb.GetPath(),
		Type:      pfm.mapImageTypeToEntity(imagePb.GetType()),
		CreatedAt: pfm.timeUtils.TimestampToTime(imagePb.GetCreatedTime()),
		UpdatedAt: pfm.timeUtils.TimestampToTime(imagePb.GetUpdatedTime()),
	}
}

func (pfm imagesMapper) ImagesToEntity(imagesPb []*imagesv1.Image) []*portal.Image {
	images := make([]*portal.Image, 0, len(imagesPb))

	for _, imagePb := range imagesPb {
		images = append(images, pfm.ImageToEntity(imagePb))
	}

	return images
}
