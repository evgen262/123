package portals

import (
	featuresv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/portals/features/v1"
	timeUtils "git.mos.ru/buch-cloud/moscow-team-2.0/build/time-utils.git"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portal"
)

type featuresMapper struct {
	timeUtils timeUtils.TimeUtils
}

func NewFeaturesMapper(timeUtils timeUtils.TimeUtils) *featuresMapper {
	return &featuresMapper{
		timeUtils: timeUtils,
	}
}

func (pfm featuresMapper) NewFeatureToPb(feature *portal.Feature) *featuresv1.AddRequest_Feature {
	return &featuresv1.AddRequest_Feature{
		Name:    feature.Name,
		Version: feature.Version,
		Enabled: feature.Enabled,
	}
}

func (pfm featuresMapper) NewFeaturesToPb(features []*portal.Feature) []*featuresv1.AddRequest_Feature {
	featuresPb := make([]*featuresv1.AddRequest_Feature, 0, len(features))

	for _, feature := range features {
		featuresPb = append(featuresPb, pfm.NewFeatureToPb(feature))
	}

	return featuresPb
}

func (pfm featuresMapper) FeatureToPb(feature *portal.Feature) *featuresv1.Feature {
	return &featuresv1.Feature{
		Id:          int32(feature.Id),
		Name:        feature.Name,
		Version:     feature.Version,
		CreatedTime: pfm.timeUtils.TimeToTimestamp(feature.CreatedAt),
		UpdatedTime: pfm.timeUtils.TimeToTimestamp(feature.UpdatedAt),
		Enabled:     false,
	}
}

func (pfm featuresMapper) FeaturesToPb(features []*portal.Feature) []*featuresv1.Feature {
	featuresPb := make([]*featuresv1.Feature, 0, len(features))

	for _, feature := range features {
		featuresPb = append(featuresPb, pfm.FeatureToPb(feature))
	}

	return featuresPb
}

func (pfm featuresMapper) FeatureToEntity(featurePb *featuresv1.Feature) *portal.Feature {
	return &portal.Feature{
		Id:        portal.FeatureId(featurePb.GetId()),
		Name:      featurePb.GetName(),
		Version:   featurePb.GetVersion(),
		CreatedAt: pfm.timeUtils.TimestampToTime(featurePb.GetCreatedTime()),
		UpdatedAt: pfm.timeUtils.TimestampToTime(featurePb.GetUpdatedTime()),
		Enabled:   featurePb.GetEnabled(),
	}
}

func (pfm featuresMapper) FeaturesToEntity(featuresPb []*featuresv1.Feature) []*portal.Feature {
	features := make([]*portal.Feature, 0, len(featuresPb))

	for _, featurePb := range featuresPb {
		features = append(features, pfm.FeatureToEntity(featurePb))
	}

	return features
}
