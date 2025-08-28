package portals

import (
	viewPortals "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/portals"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portal"
)

type imagePresenter struct {
}

func NewImagePresenter() *imagePresenter {
	return &imagePresenter{}
}

func (ip imagePresenter) ToNewEntity(image *viewPortals.NewImage) *portal.Image {
	return &portal.Image{
		Name: image.Name,
		Data: portal.ImageData(image.Data),
	}
}

func (ip imagePresenter) ToEntities(images []*viewPortals.Image) []*portal.Image {
	imagesView := make([]*portal.Image, 0, len(images))
	for _, image := range images {
		imagesView = append(imagesView, ip.ToEntity(image))
	}
	return imagesView
}

func (ip imagePresenter) ToEntity(image *viewPortals.Image) *portal.Image {
	return &portal.Image{
		Id:        portal.ImageId(image.Id),
		Name:      image.Name,
		Path:      image.Path,
		Data:      portal.ImageData(image.Data),
		CreatedAt: image.CreatedAt,
		UpdatedAt: image.UpdatedAt,
	}
}

func (ip imagePresenter) ToViews(images []*portal.Image) []*viewPortals.Image {
	imagesView := make([]*viewPortals.Image, 0, len(images))
	for _, image := range images {
		imagesView = append(imagesView, ip.ToView(image))
	}
	return imagesView
}

func (ip imagePresenter) ToView(image *portal.Image) *viewPortals.Image {
	return &viewPortals.Image{
		Id:        int(image.Id),
		Name:      image.Name,
		Path:      image.Path,
		Data:      string(image.Data),
		CreatedAt: image.CreatedAt,
		UpdatedAt: image.UpdatedAt,
	}
}

func (ip imagePresenter) ToShortViews(images []*portal.Image) []*viewPortals.ImageInfo {
	imagesView := make([]*viewPortals.ImageInfo, 0, len(images))
	for _, image := range images {
		imagesView = append(imagesView, ip.ToShortView(image))
	}
	return imagesView
}

func (ip imagePresenter) ToShortView(image *portal.Image) *viewPortals.ImageInfo {
	return &viewPortals.ImageInfo{
		Id:   int(image.Id),
		Name: image.Name,
		Path: image.Path,
	}
}
