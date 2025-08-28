package surveys

import (
	"context"
	"fmt"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/survey"
)

type surveysUseCase struct {
	repo SurveyRepository
}

func NewSurveysUseCase(repository SurveyRepository) *surveysUseCase {
	return &surveysUseCase{repo: repository}
}

func (suc surveysUseCase) Get(
	ctx context.Context,
	id entitySurveys.SurveyID,
	options entitySurveys.SurveyFilterOptions,
) (*entitySurveys.Survey, error) {
	result, err := suc.repo.Get(ctx, id, options)
	if err != nil {
		return nil, fmt.Errorf("can't get survey from repository: %w", err)
	}

	return result, nil
}
