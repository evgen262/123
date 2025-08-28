package entitySurveys

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

//go:generate ditgen -source=survey.go

type SurveyID uuid.UUID

func (s SurveyID) String() string {
	return uuid.UUID(s).String()
}

type SurveyRespondent struct {
	Ids  RespondentIDs
	Type RespondentType
}

type RespondentType int

const (
	// RespondentTypeAll используется при запросе опросов вне зависимости от respondent_type
	RespondentTypeAll RespondentType = -1 + iota
	RespondentTypeUnknown
	RespondentTypeAnonymous
	RespondentTypeUser
)

// Survey опрос.
//
//	Структура опроса и дочерних сущностей реализована с учетом того, что опросы
//
// приходят обогащенными от 1С.
type Survey struct {
	// uuid опроса.
	ID *SurveyID
	// Название опроса.
	Title string
	// Описание опроса.
	Description string
	// Дата начала проведения опроса.
	ActivePeriodStart time.Time
	// Дата окончания проведения опроса.
	ActivePeriodEnd time.Time
	// Автор опроса.
	Author string
	// Статус опубликованности опроса.
	IsPublished bool
	// Дата создания опроса.
	CreatedAt *time.Time
	// Дата изменения опроса.
	UpdatedAt *time.Time
	// Дата удаления опроса.
	DeletedAt *time.Time
	// Список вопросов опроса.
	Questions []*Question
	// Респондент опроса.
	Respondent *SurveyRespondent
}

type SurveyIDs []SurveyID

func (si SurveyIDs) ToStringSlice() []string {
	stringIDs := make([]string, 0, len(si))
	for _, id := range si {
		stringIDs = append(stringIDs, uuid.UUID(id).String())
	}
	return stringIDs
}

func (si SurveyIDs) ToString() string {
	return strings.Join(si.ToStringSlice(), ",")
}

// SurveyFilterOptions набор опций для фильтрации опросов
type SurveyFilterOptions struct {
	WithQuestions         bool
	WithAnswers           bool
	WithDeleted           bool
	WithInactiveQuestions bool
}

// Pagination параметры для пагинации опросов
type Pagination struct {
	Limit    uint32
	LastId   *SurveyID
	LastDate *time.Time
	Total    int64
}

// SurveysWithPagination
type SurveysWithPagination struct {
	Pagination
	Surveys []*Survey
}

func (sr *SurveyRespondent) ToString() string {
	if sr == nil {
		return ""
	}
	return fmt.Sprintf("%d%s",
		sr.Type,
		sr.Ids.ToString(),
	)
}

func (sfo SurveyFilterOptions) ToString() string {
	return fmt.Sprintf("%s%s%s%s",
		strconv.FormatBool(sfo.WithQuestions),
		strconv.FormatBool(sfo.WithAnswers),
		strconv.FormatBool(sfo.WithDeleted),
		strconv.FormatBool(sfo.WithInactiveQuestions),
	)
}
