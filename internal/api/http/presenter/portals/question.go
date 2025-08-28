package portals

import (
	viewPortals "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/portals"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portal"
)

type questionPresenter struct {
}

func NewQuestionPresenter() *questionPresenter {
	return &questionPresenter{}
}

func (qp questionPresenter) ToNewEntities(questions []*viewPortals.NewQuestion) []*portal.Question {
	resQuestions := make([]*portal.Question, 0, len(questions))
	for _, question := range questions {
		resQuestions = append(resQuestions, qp.ToNewEntity(question))
	}

	return resQuestions
}

func (qp questionPresenter) ToNewEntity(question *viewPortals.NewQuestion) *portal.Question {
	return &portal.Question{
		Name:        question.Name,
		Description: question.Description,
		Sort:        question.Sort,
		IsDeleted:   question.IsDeleted,
	}
}

func (qp questionPresenter) ToEntities(questions []*viewPortals.UpdateQuestion) []*portal.Question {
	resQuestions := make([]*portal.Question, 0, len(questions))
	for _, portal := range questions {
		resQuestions = append(resQuestions, qp.ToEntity(portal))
	}

	return resQuestions
}

func (qp questionPresenter) ToEntity(question *viewPortals.UpdateQuestion) *portal.Question {
	return &portal.Question{
		Id:          portal.QuestionId(question.Id),
		Name:        question.Name,
		Description: question.Description,
		Sort:        question.Sort,
		IsDeleted:   question.IsDeleted,
	}
}

func (qp questionPresenter) ToViews(questions []*portal.Question) []*viewPortals.Question {
	resQuestions := make([]*viewPortals.Question, 0, len(questions))
	for _, question := range questions {
		resQuestions = append(resQuestions, qp.ToView(question))
	}

	return resQuestions
}

func (qp questionPresenter) ToView(question *portal.Question) *viewPortals.Question {
	return &viewPortals.Question{
		Id:          int(question.Id),
		Name:        question.Name,
		Description: question.Description,
		Sort:        question.Sort,
		CreatedAt:   question.CreatedAt,
		UpdatedAt:   question.UpdatedAt,
		DeletedAt:   question.DeletedAt,
		IsDeleted:   question.IsDeleted,
	}
}

func (qp questionPresenter) ToShortViews(questions []*portal.Question) []*viewPortals.QuestionInfo {
	resQuestions := make([]*viewPortals.QuestionInfo, 0, len(questions))
	for _, question := range questions {
		resQuestions = append(resQuestions, qp.ToShortView(question))
	}

	return resQuestions
}

func (qp questionPresenter) ToShortView(question *portal.Question) *viewPortals.QuestionInfo {
	return &viewPortals.QuestionInfo{
		Name:        question.Name,
		Description: question.Description,
	}
}
