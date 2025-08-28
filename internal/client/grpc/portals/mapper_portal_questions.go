package portals

import (
	questionsv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/portals/questions/v1"
	timeUtils "git.mos.ru/buch-cloud/moscow-team-2.0/build/time-utils.git"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portal"
)

type questionsMapper struct {
	timeUtils timeUtils.TimeUtils
}

func NewQuestionsMapper(timeUtils timeUtils.TimeUtils) *questionsMapper {
	return &questionsMapper{
		timeUtils: timeUtils,
	}
}

func (qm questionsMapper) NewQuestionToPb(question *portal.Question) *questionsv1.AddRequest_Question {
	return &questionsv1.AddRequest_Question{
		Name:        question.Name,
		Description: question.Description,
		Sort:        int32(question.Sort),
		IsDeleted:   question.IsDeleted,
	}
}

func (qm questionsMapper) NewQuestionsToPb(questions []*portal.Question) []*questionsv1.AddRequest_Question {
	questionsPb := make([]*questionsv1.AddRequest_Question, 0, len(questions))

	for _, question := range questions {
		questionsPb = append(questionsPb, qm.NewQuestionToPb(question))
	}

	return questionsPb
}

func (qm questionsMapper) QuestionToPb(question *portal.Question) *questionsv1.Question {
	return &questionsv1.Question{
		Id:          int32(question.Id),
		Name:        question.Name,
		Description: question.Description,
		Sort:        int32(question.Sort),
		CreatedTime: qm.timeUtils.TimeToTimestamp(question.CreatedAt),
		UpdatedTime: qm.timeUtils.TimeToTimestamp(question.UpdatedAt),
		DeletedTime: qm.timeUtils.TimeToTimestamp(question.DeletedAt),
		IsDeleted:   question.IsDeleted,
	}
}

func (qm questionsMapper) QuestionsToPb(questions []*portal.Question) []*questionsv1.Question {
	questionsPb := make([]*questionsv1.Question, 0, len(questions))

	for _, question := range questions {
		questionsPb = append(questionsPb, qm.QuestionToPb(question))
	}

	return questionsPb
}

func (qm questionsMapper) QuestionToEntity(questionPb *questionsv1.Question) *portal.Question {
	return &portal.Question{
		Id:          portal.QuestionId(questionPb.GetId()),
		Name:        questionPb.GetName(),
		Description: questionPb.GetDescription(),
		Sort:        int(questionPb.GetSort()),
		CreatedAt:   qm.timeUtils.TimestampToTime(questionPb.GetCreatedTime()),
		UpdatedAt:   qm.timeUtils.TimestampToTime(questionPb.GetUpdatedTime()),
		DeletedAt:   qm.timeUtils.TimestampToTime(questionPb.GetDeletedTime()),
		IsDeleted:   questionPb.IsDeleted,
	}
}

func (qm questionsMapper) QuestionsToEntity(questionsPb []*questionsv1.Question) []*portal.Question {
	questions := make([]*portal.Question, 0, len(questionsPb))

	for _, questionPb := range questionsPb {
		questions = append(questions, qm.QuestionToEntity(questionPb))
	}

	return questions
}
