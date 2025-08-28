package surveys

import (
	"fmt"

	answerv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/surveys/answer/v1"
	surveys "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/survey"
	"github.com/google/uuid"
)

type answerMapper struct {
}

func NewAnswerMapper() *answerMapper {
	return &answerMapper{}
}

func (am answerMapper) AnswersToEntities(answers []*answerv1.Answer) ([]*surveys.RespondentAnswer, error) {
	answersArray := make([]*surveys.RespondentAnswer, 0, len(answers))
	for _, answer := range answers {
		answerEntity, err := am.answerToEntity(answer)
		if err != nil {
			return nil, fmt.Errorf("can't convert answer to entity: %w", err)
		}
		answersArray = append(answersArray, answerEntity)
	}

	return answersArray, nil
}

func (am answerMapper) answerToEntity(answer *answerv1.Answer) (*surveys.RespondentAnswer, error) {
	ID, err := uuid.Parse(answer.GetId())
	if err != nil {
		return nil, fmt.Errorf("can't parse answer id: %w", err)
	}
	variantID, err := uuid.Parse(answer.GetChosenVariant())
	if err != nil {
		return nil, fmt.Errorf("can't parse chosen variant id: %w", err)
	}
	questionID, err := uuid.Parse(answer.GetQuestionId())
	if err != nil {
		return nil, fmt.Errorf("can't parse question id: %w", err)
	}
	qID := surveys.QuestionID(questionID)

	newAnswer := &surveys.RespondentAnswer{
		ID:            &ID,
		QuestionID:    &qID,
		ChosenVariant: surveys.AnswerID(variantID),
		Content:       answer.Content,
	}

	if answer.GetRespondent() != nil {
		id, err := uuid.Parse(answer.GetRespondent().GetId())
		if err != nil {
			return nil, fmt.Errorf("can't parse respondent id: %w", err)
		}
		respondentId := surveys.RespondentID(id)
		newAnswer.RespondentId = &respondentId
	}

	return newAnswer, nil
}

func (am answerMapper) NewAnswersToPb(answers []*surveys.RespondentAnswer) []*answerv1.AddRequest_Answer {
	result := make([]*answerv1.AddRequest_Answer, 0, len(answers))
	for _, answer := range answers {
		entityAnswer := &answerv1.AddRequest_Answer{
			ChosenVariant: answer.ChosenVariant.String(),
			Content:       answer.Content,
		}
		result = append(result, entityAnswer)
	}

	return result
}

func (am answerMapper) IDsToUUIDs(answerIDs []string) ([]uuid.UUID, error) {
	arr := make([]uuid.UUID, 0, len(answerIDs))
	for _, answerID := range answerIDs {
		id, err := uuid.Parse(answerID)
		if err != nil {
			return nil, fmt.Errorf("can't parse answer id: %w", err)
		}
		arr = append(arr, id)
	}

	return arr, nil
}
