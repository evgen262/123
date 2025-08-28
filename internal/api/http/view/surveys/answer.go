package view

import "github.com/google/uuid"

type SurveyAnswer struct {
	ID            uuid.UUID  `json:"id"`
	QuestionId    uuid.UUID  `json:"question_id"`
	ChosenVariant uuid.UUID  `json:"chosen_variant_id"`
	RespondentId  *uuid.UUID `json:"respondent_id,omitempty"`
	Content       string     `json:"content,omitempty"`
}

type NewSurveyAnswer struct {
	ChosenVariant uuid.UUID `json:"chosen_variant_id"`
	Content       string    `json:"content,omitempty"`
}

type NewSurveyAnswers struct {
	RespondentId *uuid.UUID         `json:"respondent_id,omitempty"`
	Answers      []*NewSurveyAnswer `json:"answers"`
}

type SurveyAnswerInfo struct {
	ID uuid.UUID `json:"id"`
}
