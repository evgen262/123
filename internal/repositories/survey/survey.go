package survey

import (
	"context"
	"fmt"

	surveyv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/surveys/survey/v1"
	surveys "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/survey"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"google.golang.org/grpc/codes"
)

type surveyRepository struct {
	client surveyv1.SurveyAPIClient
	mapper SurveyMapper
}

func NewSurveyRepository(client surveyv1.SurveyAPIClient, mapper SurveyMapper) *surveyRepository {
	return &surveyRepository{
		client: client,
		mapper: mapper,
	}
}

func (sr surveyRepository) Get(
	ctx context.Context,
	id surveys.SurveyID,
	options surveys.SurveyFilterOptions,
) (*surveys.Survey, error) {
	resp, err := sr.client.Get(ctx, &surveyv1.GetRequest{
		Id:      id.String(),
		Options: sr.mapper.OptionsToPb(&options),
	})
	if err != nil {
		msg := diterrors.NewLocalizedError(diterrors.LocalizeLocale, err)
		switch msg.Code() {
		case codes.InvalidArgument:
			return nil, diterrors.NewValidationError(msg)
		case codes.NotFound:
			return nil, diterrors.ErrNotFound
		default:
			return nil, fmt.Errorf("can't get survey: %w", msg)
		}
	}

	result, err := sr.mapper.SurveyToEntity(resp.GetSurvey())
	if err != nil {
		return nil, fmt.Errorf("can't convert survey to entity: %w", err)
	}

	return result, nil
}
