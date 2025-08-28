package surveys

import (
	"context"
	"fmt"
)

type imagesUseCase struct {
	repo ImagesRepository
}

func NewImagesUseCase(repository ImagesRepository) *imagesUseCase {
	return &imagesUseCase{repo: repository}
}

func (iuc imagesUseCase) Get(ctx context.Context, imageName string) ([]byte, error) {
	result, err := iuc.repo.Get(ctx, imageName)
	if err != nil {
		return nil, fmt.Errorf("can't get image from repository: %w", err)
	}

	return result, nil
}
