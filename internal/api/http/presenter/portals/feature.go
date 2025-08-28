package portals

import (
	viewPortals "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/portals"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portal"
)

type featurePresenter struct {
}

func NewFeaturePresenter() *featurePresenter {
	return &featurePresenter{}
}

func (fp featurePresenter) ToNewEntities(features []*viewPortals.NewFeature) []*portal.Feature {
	var featuresView []*portal.Feature
	for _, feature := range features {
		featuresView = append(featuresView, fp.ToNewEntity(feature))
	}

	return featuresView
}

func (fp featurePresenter) ToNewEntity(feature *viewPortals.NewFeature) *portal.Feature {
	return &portal.Feature{
		Name:    feature.Name,
		Version: feature.Version,
		Enabled: feature.Enabled,
	}
}

func (fp featurePresenter) ToEntities(features []*viewPortals.UpdateFeature) []*portal.Feature {
	var featuresView []*portal.Feature
	for _, feature := range features {
		featuresView = append(featuresView, fp.ToEntity(feature))
	}

	return featuresView
}

func (fp featurePresenter) ToEntity(feature *viewPortals.UpdateFeature) *portal.Feature {
	return &portal.Feature{
		Id:      portal.FeatureId(feature.Id),
		Name:    feature.Name,
		Version: feature.Version,
		Enabled: feature.Enabled,
	}
}

func (fp featurePresenter) ToViews(features []*portal.Feature) []*viewPortals.Feature {
	featuresView := make([]*viewPortals.Feature, 0, len(features))
	for _, feature := range features {
		featuresView = append(featuresView, fp.ToView(feature))
	}

	return featuresView
}

func (fp featurePresenter) ToView(feature *portal.Feature) *viewPortals.Feature {
	return &viewPortals.Feature{
		Id:        int(feature.Id),
		Name:      feature.Name,
		Version:   feature.Version,
		CreatedAt: feature.CreatedAt,
		UpdatedAt: feature.UpdatedAt,
		Enabled:   feature.Enabled,
	}
}

func (fp featurePresenter) ToShortViews(features []*portal.Feature) viewPortals.Features {
	featuresView := make(viewPortals.Features, len(features))
	for _, feature := range features {
		featuresView[feature.Name] = fp.ToShortView(feature)
	}

	return featuresView
}

func (fp featurePresenter) ToShortView(feature *portal.Feature) *viewPortals.FeatureInfo {
	return &viewPortals.FeatureInfo{
		Id:      int(feature.Id),
		Version: feature.Version,
		Enabled: feature.Enabled,
	}
}
