package surveys

import (
	viewSurveys "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/surveys"
	entitySurvey "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/survey"

	"github.com/google/uuid"
)

type surveyAnswersPresenter struct {
}

func NewAnswersPresenter() *surveyAnswersPresenter {
	return &surveyAnswersPresenter{}
}

func (sap surveyAnswersPresenter) ToNewEntities(answers *viewSurveys.NewSurveyAnswers) []*entitySurvey.RespondentAnswer {
	answersResult := make([]*entitySurvey.RespondentAnswer, 0, len(answers.Answers))
	for _, answer := range answers.Answers {
		newAnswer := &entitySurvey.RespondentAnswer{
			ChosenVariant: entitySurvey.AnswerID(answer.ChosenVariant),
			Content:       answer.Content,
		}
		if answers.RespondentId != nil {
			id := entitySurvey.RespondentID(*answers.RespondentId)
			newAnswer.RespondentId = &id
		}

		answersResult = append(answersResult, newAnswer)
	}

	return answersResult
}

func (sap surveyAnswersPresenter) ToViews(answers []*entitySurvey.RespondentAnswer) []*viewSurveys.SurveyAnswer {
	answersResult := make([]*viewSurveys.SurveyAnswer, 0, len(answers))
	for _, answer := range answers {
		newAnswer := &viewSurveys.SurveyAnswer{
			ChosenVariant: uuid.UUID(answer.ChosenVariant),
			Content:       answer.Content,
		}
		if answer.GetID() != nil {
			newAnswer.ID = *answer.GetID()
		}
		if answer.GetQuestionID() != nil {
			newAnswer.QuestionId = uuid.UUID(*answer.GetQuestionID())
		}
		if answer.GetRespondentId() != nil {
			rID := uuid.UUID(*answer.RespondentId)
			newAnswer.RespondentId = &rID
		}

		answersResult = append(answersResult, newAnswer)
	}

	return answersResult
}

func (sap surveyAnswersPresenter) ToShortViews(ids []uuid.UUID) []*viewSurveys.SurveyAnswerInfo {
	answersResult := make([]*viewSurveys.SurveyAnswerInfo, 0, len(ids))
	for _, id := range ids {
		newAnswer := &viewSurveys.SurveyAnswerInfo{
			ID: id,
		}

		answersResult = append(answersResult, newAnswer)
	}

	return answersResult
}
