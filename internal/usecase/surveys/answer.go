package surveys

import (
	"context"
	"fmt"

	surveys "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/survey"
	"github.com/google/uuid"
)

type answersUseCase struct {
	repo AnswersRepository
}

func NewAnswersUseCase(repository AnswersRepository) *answersUseCase {
	return &answersUseCase{repo: repository}
}

func (auc answersUseCase) Add(ctx context.Context, answers []*surveys.RespondentAnswer) ([]uuid.UUID, error) {
	result, err := auc.repo.Add(ctx, answers)
	if err != nil {
		return nil, fmt.Errorf("can't add answers to repository: %w", err)
	}

	return result, nil
}
