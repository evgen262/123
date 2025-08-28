package survey

import (
	"context"
	"fmt"

	answerv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/surveys/answer/v1"
	surveys "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/survey"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
)

type answerRepository struct {
	client answerv1.AnswerAPIClient
	mapper AnswerMapper
}

func NewAnswerRepository(client answerv1.AnswerAPIClient, mapper AnswerMapper) *answerRepository {
	return &answerRepository{
		client: client,
		mapper: mapper,
	}
}

func (ar answerRepository) Add(ctx context.Context, answers []*surveys.RespondentAnswer) ([]uuid.UUID, error) {
	var respondent *surveys.RespondentID
	if len(answers) != 0 && answers[0].RespondentId != nil {
		respondent = answers[0].RespondentId
	}
	req := &answerv1.AddRequest{
		Answers: ar.mapper.NewAnswersToPb(answers),
	}
	if respondent != nil {
		req.Respondent = &answerv1.Respondent{
			Id: respondent.String(),
		}
	}

	resp, err := ar.client.Add(ctx, req)
	if err != nil {
		msg := diterrors.NewLocalizedError(diterrors.LocalizeLocale, err)
		switch msg.Code() {
		case codes.InvalidArgument:
			return nil, diterrors.NewValidationError(msg)
		default:
			return nil, fmt.Errorf("can't add answers: %w", msg)
		}
	}

	result, err := ar.mapper.IDsToUUIDs(resp.GetAnswerIds())
	if err != nil {
		return nil, fmt.Errorf("can't convert answer ids to uuids: %w", err)
	}

	return result, nil
}
