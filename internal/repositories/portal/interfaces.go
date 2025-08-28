package portal

import (
	featuresv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/portals/features/v1"
	imagesv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/portals/images/v1"
	organizationsv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/portals/organizations/v1"
	portalsv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/portals/portals/v1"
	questionsv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/portals/questions/v1"
	sharedv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/portals/shared/v1"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portal"
)

//go:generate mockgen -source=interfaces.go -destination=./repositories_mock.go -package=portal

type PortalsMapper interface {
	NewPortalsToPb(portals []*portal.Portal) []*portalsv1.AddRequest_Portal
	NewPortalToPb(portal *portal.Portal) *portalsv1.AddRequest_Portal
	PortalsToPb(portals []*portal.Portal) []*portalsv1.Portal
	PortalToPb(portal *portal.Portal) *portalsv1.Portal
	PortalOrganizationsToPb(orgs []*portal.PortalOrganization) []*portalsv1.Portal_Organization
	PortalOrganizationToPb(org *portal.PortalOrganization) *portalsv1.Portal_Organization
	PortalsToEntity(portalsPb []*portalsv1.Portal) []*portal.Portal
	PortalToEntity(portalPb *portalsv1.Portal) *portal.Portal
	PortalOrganizationsToEntity(orgsPb []*portalsv1.Portal_Organization) []*portal.PortalOrganization
	PortalOrganizationToEntity(orgPb *portalsv1.Portal_Organization) *portal.PortalOrganization
}

type QuestionsMapper interface {
	NewQuestionToPb(question *portal.Question) *questionsv1.AddRequest_Question
	NewQuestionsToPb(questions []*portal.Question) []*questionsv1.AddRequest_Question
	QuestionToPb(question *portal.Question) *questionsv1.Question
	QuestionsToPb(questions []*portal.Question) []*questionsv1.Question
	QuestionToEntity(questionPb *questionsv1.Question) *portal.Question
	QuestionsToEntity(questionsPb []*questionsv1.Question) []*portal.Question
}

type FeaturesMapper interface {
	NewFeatureToPb(feature *portal.Feature) *featuresv1.AddRequest_Feature
	NewFeaturesToPb(features []*portal.Feature) []*featuresv1.AddRequest_Feature
	FeatureToPb(feature *portal.Feature) *featuresv1.Feature
	FeaturesToPb(features []*portal.Feature) []*featuresv1.Feature
	FeatureToEntity(featurePb *featuresv1.Feature) *portal.Feature
	FeaturesToEntity(featuresPb []*featuresv1.Feature) []*portal.Feature
}

type OrganizationsMapper interface {
	OnceGrbsToEntity(onceGrbsPb *organizationsv1.Organization_Grbs) *portal.OrganizationGrbs
	GrbsToEntity(grbsPb []*organizationsv1.Organization_Grbs) []*portal.OrganizationGrbs
	OrganizationToEntity(organizationPb *organizationsv1.Organization) *portal.Organization
	OrganizationsToEntity(organizationsPb []*organizationsv1.Organization) []*portal.Organization
	PaginationToEntity(paginationPb *sharedv1.PaginationResponse) *entity.StringPagination
	PaginationToPb(pagination *entity.StringPagination) *sharedv1.PaginationRequest
	OptionsToPb(options portal.OrganizationsFilterOptions) *organizationsv1.OrganizationFilterOptions
}

type ImagesMapper interface {
	NewImageToPb(image *portal.Image) *imagesv1.AddRequest
	ImageToPb(image *portal.Image) *imagesv1.Image
	ImagesToPb(images []*portal.Image) []*imagesv1.Image
	ImageToEntity(imagePb *imagesv1.Image) *portal.Image
	ImagesToEntity(imagesPb []*imagesv1.Image) []*portal.Image
}
