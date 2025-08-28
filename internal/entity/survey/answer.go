package entitySurveys

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

//go:generate ditgen -source=answer.go

type AnswerID uuid.UUID

func (a AnswerID) String() string {
	return uuid.UUID(a).String()
}

type ContentType int

const (
	// ContentTypeInvalid некорректный тип
	ContentTypeInvalid ContentType = 0
	// ContentTypeText контент в текстовом формате.
	ContentTypeText ContentType = 1
	// ContentTypeDigit контент в целочисленном формате.
	ContentTypeDigit ContentType = 2
)

// AnswerRules ограничения контента ответа.
type AnswerRules struct {
	// Минимальная длина ответа в символах. Используется в случае, если создается вариант ответа "Другое".
	ContentMinLength *int
	// Максимальная длина ответа в символах. Используется в случае, если создается вариант ответа "Другое".
	ContentMaxLength *int
	// Минимальное число. Используется в случае типа ответа inputNumber.
	ContentMinDigit *int
	// Максимальное число. Используется в случае варианта ответа inputNumber.
	ContentMaxDigit *int
}

// AnswerVariant вариант ответа.
type AnswerVariant struct {
	// Id ответа.
	// Id = ключ объекта в минио (название объекта).
	ID *AnswerID
	// Id вопроса.
	QuestionID *QuestionID
	// Текст ответа.
	Text string
	// Информация о картинке, соответствующей ответу.
	Image *Image
	// Ограничения контента ответа.
	// Используется, если WithContent = true.
	Rules *AnswerRules
	// Использовать ли поле Content в ответе пользователя.
	WithContent bool
	// Тип контента.
	//  Используется, если WithContent = true.
	ContentType ContentType
	// Время удаления варианта ответа.
	DeletedAt *time.Time
	// Вес варианта ответа с целью сортировки. Отображение происходит по возрастанию веса.
	Weight int
}

// RespondentAnswer ответ пользователя.
type RespondentAnswer struct {
	// uuid.
	ID *uuid.UUID
	// uuid вопроса, для которого выбран ответ.
	QuestionID *QuestionID
	// uuid выбранного ответа.
	ChosenVariant AnswerID
	// uuid респондента.
	RespondentId *RespondentID
	// Контент ответа. Используется в случае, если выбран вариант ответа "Другое".
	Content string
}

type AnswerIDs []AnswerID

func (ai AnswerIDs) ToStringSlice() []string {
	stringIDs := make([]string, 0, len(ai))
	for _, id := range ai {
		stringIDs = append(stringIDs, uuid.UUID(id).String())
	}
	return stringIDs
}

func (ai AnswerIDs) ToString() string {
	return strings.Join(ai.ToStringSlice(), ",")
}
