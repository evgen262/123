package portal

import (
	"context"
	"fmt"

	imagesv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/portals/images/v1"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portal"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/repositories"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"google.golang.org/grpc/codes"
)

type imagesRepository struct {
	client imagesv1.ImagesAPIClient
	mapper ImagesMapper
}

func NewImagesRepository(client imagesv1.ImagesAPIClient, mapper ImagesMapper) *imagesRepository {
	return &imagesRepository{
		client: client,
		mapper: mapper,
	}
}

func (pir imagesRepository) All(ctx context.Context) ([]*portal.Image, error) {
	resp, err := pir.client.All(ctx, &imagesv1.AllRequest{})
	if err != nil {
		return nil, fmt.Errorf("can't get all images: %w", err)
	}

	return pir.mapper.ImagesToEntity(resp.GetImages()), nil
}

func (pir imagesRepository) Get(ctx context.Context, imageId portal.ImageId) (*portal.Image, error) {
	resp, err := pir.client.Get(ctx, &imagesv1.GetRequest{
		Id: int32(imageId),
	})
	if err != nil {
		msg := diterrors.NewLocalizedError(diterrors.LocalizeLocale, err)
		switch msg.Code() {
		case codes.NotFound:
			return nil, repositories.ErrNotFound
		case codes.InvalidArgument:
			return nil, diterrors.NewValidationError(err)
		default:
			return nil, fmt.Errorf("can't get image: %w", err)
		}
	}

	return pir.mapper.ImageToEntity(resp.GetImage()), nil
}

func (pir imagesRepository) GetImageData(ctx context.Context, path string) (portal.ImageData, error) {
	resp, err := pir.client.GetImage(ctx, &imagesv1.GetImageRequest{
		Path: path,
	})
	if err != nil {
		msg := diterrors.NewLocalizedError(diterrors.LocalizeLocale, err)
		switch msg.Code() {
		case codes.NotFound:
			return nil, repositories.ErrNotFound
		case codes.InvalidArgument:
			return nil, diterrors.NewValidationError(err)
		default:
			return nil, fmt.Errorf("can't get image data: %w", err)
		}
	}

	return resp.GetImage(), nil
}

func (pir imagesRepository) Add(ctx context.Context, image *portal.Image) (*portal.Image, error) {
	resp, err := pir.client.Add(ctx, pir.mapper.NewImageToPb(image))
	if err != nil {
		msg := diterrors.NewLocalizedError(diterrors.LocalizeLocale, err)
		switch msg.Code() {
		case codes.NotFound:
			return nil, repositories.ErrNotFound
		case codes.InvalidArgument:
			return nil, diterrors.NewValidationError(err)
		default:
			return nil, fmt.Errorf("can't add image to db: %w", err)
		}
	}

	return pir.mapper.ImageToEntity(resp.GetImage()), nil
}

func (pir imagesRepository) Delete(ctx context.Context, imageId portal.ImageId) error {
	_, err := pir.client.Delete(ctx, &imagesv1.DeleteRequest{Id: int32(imageId)})
	if err != nil {
		msg := diterrors.NewLocalizedError(diterrors.LocalizeLocale, err)
		switch msg.Code() {
		case codes.NotFound:
			return repositories.ErrNotFound
		case codes.InvalidArgument:
			return diterrors.NewValidationError(err)
		default:
			return fmt.Errorf("can't delete image: %w", err)
		}
	}

	return nil
}
