package portals

import (
	"context"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portal"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"go.uber.org/zap"
)

type questionsUseCase struct {
	repo   QuestionsRepository
	logger ditzap.Logger
}

func NewQuestionUseCase(repository QuestionsRepository, logger ditzap.Logger) *questionsUseCase {
	return &questionsUseCase{
		repo:   repository,
		logger: logger,
	}
}

func (quc questionsUseCase) GetAllQuestions(ctx context.Context, withDeleted bool) (*portal.Questions, error) {
	questions, err := quc.repo.All(ctx, withDeleted)
	if err != nil {
		quc.logger.Error("can't get questions", zap.Error(err))
		return nil, err
	}
	return questions, nil
}

func (quc questionsUseCase) GetQuestion(ctx context.Context, questionId int, withDeleted bool) (*portal.Question, error) {
	question, err := quc.repo.Get(ctx, portal.QuestionId(questionId), withDeleted)
	if err != nil {
		quc.logger.Error("can't get question", zap.Error(err))
		return nil, err
	}
	return question, nil
}

func (quc questionsUseCase) AddQuestions(ctx context.Context, questions []*portal.Question) ([]*portal.Question, error) {
	newQuestions, err := quc.repo.Add(ctx, questions)
	if err != nil {
		quc.logger.Error("can't add questions", zap.Error(err))
		return nil, err
	}
	return newQuestions, nil
}

func (quc questionsUseCase) AddQuestion(ctx context.Context, question *portal.Question) (*portal.Question, error) {
	newQuestion, err := quc.repo.Add(ctx, []*portal.Question{question})
	if err != nil {
		quc.logger.Error("can't add question", zap.Error(err))
		return nil, err
	}
	if len(newQuestion) == 0 {
		return nil, nil
	}
	return newQuestion[0], nil
}

func (quc questionsUseCase) UpdateQuestion(ctx context.Context, question *portal.Question) (*portal.Question, error) {
	q, err := quc.repo.Update(ctx, question)
	if err != nil {
		quc.logger.Error("can't update question", zap.Int("question_id", int(question.Id)), zap.Error(err))
		return nil, err
	}
	return q, nil
}

func (quc questionsUseCase) DeleteQuestion(ctx context.Context, questionId int) error {
	if err := quc.repo.Delete(ctx, portal.QuestionId(questionId)); err != nil {
		quc.logger.Error("can't delete question", zap.Int("question_id", questionId), zap.Error(err))
		return err
	}
	return nil
}
