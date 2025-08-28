package portals

import (
	"context"
	"errors"
	"fmt"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portal"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/repositories"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"go.uber.org/zap"
)

type imagesUseCase struct {
	repo       ImagesRepository
	uploadPath string
	logger     ditzap.Logger
}

func NewImageUseCase(
	repository ImagesRepository,
	uploadPath string,
	logger ditzap.Logger,
) *imagesUseCase {
	return &imagesUseCase{
		repo:       repository,
		uploadPath: uploadPath,
		logger:     logger,
	}
}

func (iuc imagesUseCase) All(ctx context.Context) ([]*portal.Image, error) {
	images, err := iuc.repo.All(ctx)
	if err != nil {
		iuc.logger.Error("can't get all images from repo", zap.Error(err))
		return nil, fmt.Errorf("can't get all images: %w", err)
	}
	return images, nil
}

func (iuc imagesUseCase) Get(ctx context.Context, imageId int) (*portal.Image, error) {
	image, err := iuc.repo.Get(ctx, portal.ImageId(imageId))
	if err != nil {
		switch true {
		case errors.As(err, new(diterrors.ValidationError)):
			fallthrough
		case errors.Is(err, repositories.ErrNotFound):
			return nil, repositories.ErrNotFound
		default:
			iuc.logger.Error("can't get image from repo", zap.Error(err))
			return nil, fmt.Errorf("can't get image: %w", err)
		}
	}
	return image, nil
}

func (iuc imagesUseCase) GetRawImage(ctx context.Context, path string) (portal.ImageData, error) {
	data, err := iuc.repo.GetImageData(ctx, path)
	if err != nil {
		switch true {
		case errors.As(err, new(diterrors.ValidationError)):
			fallthrough
		case errors.Is(err, repositories.ErrNotFound):
			return nil, repositories.ErrNotFound
		default:
			iuc.logger.Error("can't get raw image from repo", zap.Error(err))
			return nil, fmt.Errorf("can't get raw image: %w", err)
		}
	}
	return data, nil
}

func (iuc imagesUseCase) Add(ctx context.Context, image *portal.Image) (*portal.Image, error) {
	image, err := iuc.repo.Add(ctx, image)
	if err != nil {
		switch true {
		case errors.As(err, new(diterrors.ValidationError)):
			fallthrough
		case errors.Is(err, repositories.ErrNotFound):
			return nil, repositories.ErrNotFound
		default:
			iuc.logger.Error("can't add image into repo", zap.Error(err))
			return nil, fmt.Errorf("can't add image: %w", err)
		}
	}
	return image, nil
}

func (iuc imagesUseCase) Delete(ctx context.Context, imageId int) error {
	err := iuc.repo.Delete(ctx, portal.ImageId(imageId))
	if err != nil {
		switch true {
		case errors.As(err, new(diterrors.ValidationError)):
			fallthrough
		case errors.Is(err, repositories.ErrNotFound):
			return repositories.ErrNotFound
		default:
			iuc.logger.Error("can't delete image from repo", zap.Error(err), zap.Int("id", imageId))
			return fmt.Errorf("can't delete image: %w", err)
		}
	}
	return nil
}
