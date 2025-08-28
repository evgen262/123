package survey

import (
	answerv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/surveys/answer/v1"
	imagev1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/surveys/image/v1"
	respondentv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/surveys/respondent/v1"
	sharedv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/surveys/shared/v1"
	surveyv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/surveys/survey/v1"
	surveys "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/survey"
	"github.com/google/uuid"
)

//go:generate mockgen -source=interfaces.go -destination=./repositories_mock.go -package=survey

type SurveyMapper interface {
	PaginationToEntity(response *sharedv1.PaginationResponse) (*surveys.Pagination, error)
	PaginationToPb(pagination *surveys.Pagination) *sharedv1.PaginationRequest
	RespondentIDsToEntity(respondentIDs []string) (surveys.RespondentIDs, error)
	RespondentToPb(respondent *surveys.SurveyRespondent) *respondentv1.Respondent
	RespondentToEntity(respondent *respondentv1.Respondent) (*surveys.SurveyRespondent, error)
	OptionsToPb(options *surveys.SurveyFilterOptions) *sharedv1.Options
	SurveyToPb(survey *surveys.Survey) *surveyv1.Survey
	NewSurveyToPb(survey *surveys.Survey) *surveyv1.AddRequest_Survey
	SurveyToEntity(surveyPb *surveyv1.Survey) (*surveys.Survey, error)
	SurveysToEntities(surveysArr []*surveyv1.Survey) ([]*surveys.Survey, error)
}

type AnswerMapper interface {
	AnswersToEntities(answers []*answerv1.Answer) ([]*surveys.RespondentAnswer, error)
	NewAnswersToPb(answers []*surveys.RespondentAnswer) []*answerv1.AddRequest_Answer
	IDsToUUIDs(answerIDs []string) ([]uuid.UUID, error)
}

type ImageMapper interface {
	NewImageToPb(image *surveys.Image) *imagev1.AddRequest_Image
	ImageToEntity(image *imagev1.Image) (*surveys.Image, error)
}
