package surveys

import (
	"context"

	entitySurveys "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/survey"
	"github.com/google/uuid"
)

//go:generate mockgen -source=interfaces.go -destination=./surveys_mock.go -package=surveys

type SurveyRepository interface {
	Get(
		ctx context.Context,
		id entitySurveys.SurveyID,
		options entitySurveys.SurveyFilterOptions,
	) (*entitySurveys.Survey, error)
}

type AnswersRepository interface {
	Add(ctx context.Context, answers []*entitySurveys.RespondentAnswer) ([]uuid.UUID, error)
}

type ImagesRepository interface {
	Get(ctx context.Context, imageName string) ([]byte, error)
}
