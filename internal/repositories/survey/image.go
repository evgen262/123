package survey

import (
	"context"
	"fmt"

	imagev1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/surveys/image/v1"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
)

type imageRepository struct {
	client imagev1.ImageAPIClient
	mapper ImageMapper
}

func NewImageRepository(client imagev1.ImageAPIClient, mapper ImageMapper) *imageRepository {
	return &imageRepository{
		client: client,
		mapper: mapper,
	}
}

func (ir imageRepository) Get(ctx context.Context, imageName string) ([]byte, error) {
	resp, err := ir.client.Get(ctx, &imagev1.GetRequest{Id: imageName})
	if err != nil {
		msg := diterrors.NewLocalizedError(diterrors.LocalizeLocale, err)
		return nil, fmt.Errorf("can't get image: %w", msg)
	}

	return resp.GetImage(), nil
}
