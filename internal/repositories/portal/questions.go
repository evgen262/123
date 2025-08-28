package portal

import (
	"context"
	"fmt"

	questionsv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/portals/questions/v1"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portal"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/repositories"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"google.golang.org/grpc/codes"
)

type questionsRepository struct {
	questionsAPIClient questionsv1.QuestionsAPIClient
	questionsMapper    QuestionsMapper
}

func NewQuestionsRepository(client questionsv1.QuestionsAPIClient, mapper QuestionsMapper) *questionsRepository {
	return &questionsRepository{
		questionsAPIClient: client,
		questionsMapper:    mapper,
	}
}

func (pqr questionsRepository) All(ctx context.Context, withDeleted bool) (*portal.Questions, error) {
	resp, err := pqr.questionsAPIClient.All(ctx, &questionsv1.AllRequest{WithDeleted: withDeleted})
	if err != nil {
		msg := diterrors.NewLocalizedError(diterrors.LocalizeLocale, err)
		switch msg.Code() {
		case codes.InvalidArgument:
			return nil, diterrors.NewValidationError(msg)
		case codes.NotFound:
			return nil, repositories.ErrNotFound
		default:
			return nil, fmt.Errorf("can't get all questions: %w", msg)
		}
	}
	return &portal.Questions{
		SupportEmail: resp.GetSupportEmail(),
		Questions:    pqr.questionsMapper.QuestionsToEntity(resp.GetQuestions()),
	}, nil
}

func (pqr questionsRepository) Get(ctx context.Context, questionId portal.QuestionId, withDeleted bool) (*portal.Question, error) {
	resp, err := pqr.questionsAPIClient.Get(ctx, &questionsv1.GetRequest{Id: int32(questionId), WithDeleted: withDeleted})
	if err != nil {
		msg := diterrors.NewLocalizedError(diterrors.LocalizeLocale, err)
		switch msg.Code() {
		case codes.InvalidArgument:
			return nil, diterrors.NewValidationError(msg)
		case codes.NotFound:
			return nil, repositories.ErrNotFound
		default:
			return nil, fmt.Errorf("can't get question: %w", msg)
		}
	}
	return pqr.questionsMapper.QuestionToEntity(resp.GetQuestion()), nil
}

func (pqr questionsRepository) Add(ctx context.Context, questions []*portal.Question) ([]*portal.Question, error) {
	resp, err := pqr.questionsAPIClient.Add(ctx, &questionsv1.AddRequest{Questions: pqr.questionsMapper.NewQuestionsToPb(questions)})
	if err != nil {
		msg := diterrors.NewLocalizedError(diterrors.LocalizeLocale, err)
		switch msg.Code() {
		case codes.InvalidArgument:
			return nil, diterrors.NewValidationError(msg)
		case codes.NotFound:
			return nil, repositories.ErrNotFound
		default:
			return nil, fmt.Errorf("can't add questions: %w", msg)
		}
	}
	return pqr.questionsMapper.QuestionsToEntity(resp.GetQuestions()), nil
}

func (pqr questionsRepository) Update(ctx context.Context, question *portal.Question) (*portal.Question, error) {
	resp, err := pqr.questionsAPIClient.Update(ctx, &questionsv1.UpdateRequest{Question: pqr.questionsMapper.QuestionToPb(question)})
	if err != nil {
		msg := diterrors.NewLocalizedError(diterrors.LocalizeLocale, err)
		switch msg.Code() {
		case codes.InvalidArgument:
			return nil, diterrors.NewValidationError(msg)
		case codes.NotFound:
			return nil, repositories.ErrNotFound
		default:
			return nil, fmt.Errorf("can't update question: %w", msg)
		}
	}
	return pqr.questionsMapper.QuestionToEntity(resp.GetQuestion()), nil
}

func (pqr questionsRepository) Delete(ctx context.Context, questionId portal.QuestionId) error {
	_, err := pqr.questionsAPIClient.Delete(ctx, &questionsv1.DeleteRequest{Id: int32(questionId)})
	if err != nil {
		msg := diterrors.NewLocalizedError(diterrors.LocalizeLocale, err)
		switch msg.Code() {
		case codes.InvalidArgument:
			return diterrors.NewValidationError(msg)
		case codes.NotFound:
			return repositories.ErrNotFound
		default:
			return fmt.Errorf("can't delete question: %w", msg)
		}
	}
	return nil
}
