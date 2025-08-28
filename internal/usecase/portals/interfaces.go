package portals

import (
	"context"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity"
	entityPortal "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portal"
)

//go:generate mockgen -source=interfaces.go -destination=./portals_mock.go -package=portals

type PortalRepository interface {
	Filter(ctx context.Context, options entityPortal.PortalsFilterOptions) ([]*entityPortal.Portal, error)
	Add(ctx context.Context, newPortal []*entityPortal.Portal) ([]*entityPortal.Portal, error)
	Update(ctx context.Context, newPortal *entityPortal.Portal) (*entityPortal.Portal, error)
	Delete(ctx context.Context, id entityPortal.PortalID) error
}

type QuestionsRepository interface {
	All(ctx context.Context, withDeleted bool) (*entityPortal.Questions, error)
	Get(ctx context.Context, questionId entityPortal.QuestionId, withDeleted bool) (*entityPortal.Question, error)
	Add(ctx context.Context, questions []*entityPortal.Question) ([]*entityPortal.Question, error)
	Update(ctx context.Context, question *entityPortal.Question) (*entityPortal.Question, error)
	Delete(ctx context.Context, questionId entityPortal.QuestionId) error
}

type ImagesRepository interface {
	All(ctx context.Context) ([]*entityPortal.Image, error)
	Get(ctx context.Context, imageId entityPortal.ImageId) (*entityPortal.Image, error)
	GetImageData(ctx context.Context, path string) (entityPortal.ImageData, error)
	Add(ctx context.Context, image *entityPortal.Image) (*entityPortal.Image, error)
	Delete(ctx context.Context, imageId entityPortal.ImageId) error
}

type FeaturesRepository interface {
	All(ctx context.Context, withDisabled bool) ([]*entityPortal.Feature, error)
	Get(ctx context.Context, featureId entityPortal.FeatureId, withDisabled bool) (*entityPortal.Feature, error)
	Add(ctx context.Context, features []*entityPortal.Feature) ([]*entityPortal.Feature, error)
	Update(ctx context.Context, feature *entityPortal.Feature) (*entityPortal.Feature, error)
	Delete(ctx context.Context, featureId entityPortal.FeatureId) error
}

type OrganizationsRepository interface {
	Filter(
		ctx context.Context,
		filters entityPortal.OrganizationsFilters,
		pagination *entity.StringPagination,
		options entityPortal.OrganizationsFilterOptions,
	) (*entityPortal.OrganizationsWithPagination, error)
	LinkOrganizationsToPortal(ctx context.Context, portalId entityPortal.PortalID, ids entityPortal.OrganizationIDs) error
	UnlinkOrganizations(ctx context.Context, ids entityPortal.OrganizationIDs) error
}
