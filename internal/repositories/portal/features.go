package portal

import (
	"context"
	"fmt"

	featuresv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/portals/features/v1"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portal"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/repositories"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"google.golang.org/grpc/codes"
)

type featuresRepository struct {
	client featuresv1.FeaturesAPIClient
	mapper FeaturesMapper
}

func NewFeaturesRepository(client featuresv1.FeaturesAPIClient, mapper FeaturesMapper) *featuresRepository {
	return &featuresRepository{
		client: client,
		mapper: mapper,
	}
}

func (pqr featuresRepository) All(ctx context.Context, withDeleted bool) ([]*portal.Feature, error) {
	resp, err := pqr.client.All(ctx, &featuresv1.AllRequest{WithDeleted: withDeleted})
	if err != nil {
		msg := diterrors.NewLocalizedError(diterrors.LocalizeLocale, err)
		if msg.Code() == codes.NotFound {
			return nil, repositories.ErrNotFound
		}
		return nil, fmt.Errorf("can't get all features: %w", err)
	}
	return pqr.mapper.FeaturesToEntity(resp.GetFeatures()), nil
}

func (pqr featuresRepository) Get(ctx context.Context, questionId portal.FeatureId, withDeleted bool) (*portal.Feature, error) {
	resp, err := pqr.client.Get(ctx, &featuresv1.GetRequest{Id: int32(questionId), WithDeleted: withDeleted})
	if err != nil {
		return nil, fmt.Errorf("can't get feature: %w", err)
	}
	return pqr.mapper.FeatureToEntity(resp.GetFeature()), nil
}

func (pqr featuresRepository) Add(ctx context.Context, features []*portal.Feature) ([]*portal.Feature, error) {
	resp, err := pqr.client.Add(ctx, &featuresv1.AddRequest{Features: pqr.mapper.NewFeaturesToPb(features)})
	if err != nil {
		return nil, fmt.Errorf("can't add features: %w", err)
	}
	return pqr.mapper.FeaturesToEntity(resp.GetFeatures()), nil
}

func (pqr featuresRepository) Update(ctx context.Context, feature *portal.Feature) (*portal.Feature, error) {
	resp, err := pqr.client.Update(ctx, &featuresv1.UpdateRequest{Feature: pqr.mapper.FeatureToPb(feature)})
	if err != nil {
		return nil, fmt.Errorf("can't update feature: %w", err)
	}
	return pqr.mapper.FeatureToEntity(resp.GetFeature()), nil
}

func (pqr featuresRepository) Delete(ctx context.Context, featureId portal.FeatureId) error {
	_, err := pqr.client.Delete(ctx, &featuresv1.DeleteRequest{Id: int32(featureId)})
	if err != nil {
		return fmt.Errorf("can't delete feature: %w", err)
	}
	return nil
}
