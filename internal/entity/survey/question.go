package entitySurveys

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

//go:generate ditgen -source=question.go

type QuestionID uuid.UUID
type QuestionType string

func (q QuestionID) String() string {
	return uuid.UUID(q).String()
}

const (
	// QuestionTypeInvalid Некорректный тип вопроса
	QuestionTypeInvalid QuestionType = "invalid"
	// QuestionTypeText ответ в виде произвольного текста.
	QuestionTypeText QuestionType = "text"
	// QuestionTypeNumber целочисленный ответ.
	QuestionTypeNumber QuestionType = "inputNumber"
	// QuestionTypeCheckbox ответ с множественным выбором.
	QuestionTypeCheckbox QuestionType = "checkbox"
	// QuestionTypeRadio ответ с одиночным выбором.
	QuestionTypeRadio QuestionType = "radio"
	// QuestionTypeCheckboxImg ответ с множественным выбором с картинками.
	QuestionTypeCheckboxImg QuestionType = "checkboxImg"
	// QuestionTypeRadioImg ответ с одиночным выбором с картинками.
	QuestionTypeRadioImg QuestionType = "radioImg"
)

// QuestionRules ограничения вопроса.
type QuestionRules struct {
	// Минимальное количество возможных ответов.
	// Используется в случае типа вопроса checkbox и checkboxImg.
	PickMinCount *int
	// Максимальное количество возможных ответов.
	PickMaxCount *int
}

// Question вопрос.
type Question struct {
	// Id вопроса.
	ID *QuestionID
	// Id опроса.
	SurveyID *SurveyID
	// Текст вопроса.
	Text string
	// Тип вопроса.
	Type QuestionType
	// Обязательность получения ответа на вопрос.
	IsRequired bool
	// Активность вопроса в опросе.
	IsActive bool
	// Вес вопроса с целью сортировки. Отображение происходит по возрастанию веса.
	Weight int
	// Ограничения вопроса.
	Rules *QuestionRules
	// Список возможных ответов на вопрос.
	Answers []*AnswerVariant
	// Время удаления вопроса.
	DeletedAt *time.Time
}

type QuestionIDs []QuestionID

func (qi QuestionIDs) ToStringSlice() []string {
	stringIDs := make([]string, 0, len(qi))
	for _, id := range qi {
		stringIDs = append(stringIDs, uuid.UUID(id).String())
	}
	return stringIDs
}

func (qi QuestionIDs) ToString() string {
	return strings.Join(qi.ToStringSlice(), ",")
}
