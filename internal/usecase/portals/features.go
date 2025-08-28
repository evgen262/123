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

type featuresUseCase struct {
	repo   FeaturesRepository
	logger ditzap.Logger
}

func NewFeatureUseCase(repository FeaturesRepository, logger ditzap.Logger) *featuresUseCase {
	return &featuresUseCase{
		repo:   repository,
		logger: logger,
	}
}

func (fuc featuresUseCase) All(ctx context.Context, withDisabled bool) ([]*portal.Feature, error) {
	features, err := fuc.repo.All(ctx, withDisabled)
	if err != nil {
		switch {
		case errors.Is(err, repositories.ErrNotFound):
			return nil, nil
		case errors.As(err, new(diterrors.ValidationError)):
			fallthrough
		default:
			fuc.logger.Error("can't get all features from repo", zap.Error(err))
		}
		return nil, fmt.Errorf("can't get all features: %w", err)
	}
	return features, nil
}

func (fuc featuresUseCase) Get(ctx context.Context, featureId int, withDisabled bool) (*portal.Feature, error) {
	feature, err := fuc.repo.Get(ctx, portal.FeatureId(featureId), withDisabled)
	if err != nil {
		if !errors.Is(err, repositories.ErrNotFound) {
			fuc.logger.Error("can't get feature from repo", zap.Error(err))
		}
		return nil, fmt.Errorf("can't get feature: %w", err)
	}
	return feature, nil
}

func (fuc featuresUseCase) MultipleAdd(ctx context.Context, features []*portal.Feature) ([]*portal.Feature, error) {
	newFeatures, err := fuc.repo.Add(ctx, features)
	if err != nil {
		fuc.logger.Error("can't add features into repo", zap.Error(err))
		return nil, fmt.Errorf("can't add features: %w", err)
	}
	return newFeatures, nil
}

func (fuc featuresUseCase) Add(ctx context.Context, feature *portal.Feature) (*portal.Feature, error) {
	newFeature, err := fuc.repo.Add(ctx, []*portal.Feature{feature})
	if err != nil {
		fuc.logger.Error("can't add feature into repo", zap.Error(err))
		return nil, fmt.Errorf("can't add feature: %w", err)
	}
	if len(newFeature) == 0 {
		return nil, nil
	}
	return newFeature[0], nil
}

func (fuc featuresUseCase) Update(ctx context.Context, feature *portal.Feature) (*portal.Feature, error) {
	updatedFeature, err := fuc.repo.Update(ctx, feature)
	if err != nil {
		fuc.logger.Error("can't update feature into repo", zap.String("feature_name", feature.Name), zap.Error(err))
		return nil, fmt.Errorf("can't update feature: %w", err)
	}
	return updatedFeature, nil
}

func (fuc featuresUseCase) Delete(ctx context.Context, featureId int) error {
	if err := fuc.repo.Delete(ctx, portal.FeatureId(featureId)); err != nil {
		fuc.logger.Error("can't delete feature from repo", zap.Int("feature_id", featureId), zap.Error(err))
		return fmt.Errorf("can't delete feature: %w", err)
	}

	return nil
}
