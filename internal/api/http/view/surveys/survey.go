package view

import (
	"time"

	"github.com/google/uuid"
)

type Survey struct {
	ID                uuid.UUID `json:"id"`
	Title             string    `json:"title"`
	Description       string    `json:"description"`
	ActivePeriodStart time.Time `json:"active_period_start"`
	ActivePeriodEnd   time.Time `json:"active_period_end"`
	Author            string    `json:"author"`
	IsPublished       bool      `json:"is_published"`
	// Тип респондента:
	// * 0 - Неизвестный респондент.
	// * 1 - Анонимный респондент.
	// * 2 - Авторизованный респондент.
	RespondentType int                 `json:"respondent_type" enums:"0,1,2"`
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      *time.Time          `json:"updated_at,omitempty"`
	DeletedAt      *time.Time          `json:"deleted_at,omitempty"`
	Respondents    []*SurveyRespondent `json:"respondents,omitempty"`
	Questions      []*SurveyQuestion   `json:"questions,omitempty"`
}

type SurveyRespondent struct {
	ID uuid.UUID `json:"id"`
}

type SurveyQuestion struct {
	ID       uuid.UUID `json:"id"`
	SurveyID uuid.UUID `json:"survey_id"`
	Text     string    `json:"text"`
	// Тип вопроса:
	// * invalid - Некорректный тип.
	// * text - Произвольный текст.
	// * inputNumber - Целочисленный ответ.
	// * checkbox - Множественный выбор.
	// * radio - Одиночный выбор.
	// * checkboxImg - Множественный выбор с изображениями.
	// * radioImg - Одиночный выбор с изображениями.
	Type           string                 `json:"type" enums:"invalid,text,inputNumber,checkbox,radio,checkboxImg,radioImg"`
	IsRequired     bool                   `json:"is_required"`
	IsActive       bool                   `json:"is_active"`
	Weight         int                    `json:"weight"`
	Rules          *SurveyQuestionRules   `json:"rules,omitempty"`
	DeletedAt      *time.Time             `json:"deleted_at,omitempty"`
	AnswerVariants []*SurveyAnswerVariant `json:"answers,omitempty"`
}

type SurveyAnswerVariant struct {
	ID          uuid.UUID `json:"id"`
	QuestionID  uuid.UUID `json:"question_id"`
	Text        string    `json:"text"`
	WithContent bool      `json:"with_content"`
	// Тип варианта ответа:
	// * 0 - Некорректный тип.
	// * 1 - Контент в текстовом формате.
	// * 2 - Контент в целочисленном формате.
	ContentType *int                      `json:"content_type,omitempty" enums:"0,1,2"`
	DeletedAt   *time.Time                `json:"deleted_at,omitempty"`
	Rules       *SurveyAnswerVariantRules `json:"rules,omitempty"`
	Image       *SurveyImage              `json:"image,omitempty"`
	Weight      int                       `json:"weight"`
}

type NewSurvey struct {
	Title             string    `json:"title"`
	Description       string    `json:"description"`
	ActivePeriodStart time.Time `json:"active_period_start"`
	ActivePeriodEnd   time.Time `json:"active_period_end"`
	Author            string    `json:"author"`
	IsPublished       bool      `json:"is_published"`
	// Тип респондента:
	// * 0 - Неизвестный респондент.
	// * 1 - Анонимный респондент.
	// * 2 - Авторизованный респондент.
	RespondentType int                  `json:"respondent_type" enums:"0,1,2"`
	Respondents    []*SurveyRespondent  `json:"respondents,omitempty"`
	Questions      []*NewSurveyQuestion `json:"questions,omitempty"`
}

type NewSurveyQuestion struct {
	Text string `json:"text"`
	// Тип вопроса:
	// * invalid - Некорректный тип.
	// * text - Произвольный текст.
	// * inputNumber - Целочисленный ответ.
	// * checkbox - Множественный выбор.
	// * radio - Одиночный выбор.
	// * checkboxImg - Множественный выбор с изображениями.
	// * radioImg - Одиночный выбор с изображениями.
	Type           string                    `json:"type" enums:"invalid,text,inputNumber,checkbox,radio,checkboxImg,radioImg"`
	IsRequired     bool                      `json:"is_required"`
	IsActive       bool                      `json:"is_active"`
	Weight         int                       `json:"weight"`
	Rules          *SurveyQuestionRules      `json:"rules,omitempty"`
	AnswerVariants []*NewSurveyAnswerVariant `json:"answers,omitempty"`
}

type NewSurveyAnswerVariant struct {
	Text        string `json:"text"`
	WithContent bool   `json:"with_content"`
	// Тип варианта ответа:
	// * 0 - Некорректный тип.
	// * 1 - Контент в текстовом формате.
	// * 2 - Контент в целочисленном формате.
	ContentType *int                      `json:"content_type,omitempty" enums:"0,1,2"`
	Rules       *SurveyAnswerVariantRules `json:"rules,omitempty"`
	Image       *SurveyImage              `json:"image,omitempty"`
	Weight      int                       `json:"weight"`
}

type SurveyQuestionRules struct {
	PickMinCount *int `json:"pick_min_count,omitempty"`
	PickMaxCount *int `json:"pick_max_count,omitempty"`
}

type SurveyAnswerVariantRules struct {
	ContentMinLength *int `json:"content_min_length,omitempty"`
	ContentMaxLength *int `json:"content_max_length,omitempty"`
	ContentMinDigit  *int `json:"content_min_digit,omitempty"`
	ContentMaxDigit  *int `json:"content_max_digit,omitempty"`
}

type SurveyImage struct {
	ID                    string    `json:"id"`
	ImageExternalID       uuid.UUID `json:"external_id"`
	ImageExternalFileName string    `json:"external_file_name"`
	ImageExternalURL      string    `json:"external_url"`
	ImageExternalSize     int       `json:"external_size"`
}

type SurveyInfo struct {
	Title             string                `json:"title"`
	Description       string                `json:"description"`
	ActivePeriodStart time.Time             `json:"active_period_start"`
	ActivePeriodEnd   time.Time             `json:"active_period_end"`
	Questions         []*SurveyQuestionInfo `json:"questions,omitempty"`
}

type SurveyQuestionInfo struct {
	ID   uuid.UUID `json:"id"`
	Text string    `json:"text"`
	// Тип вопроса:
	// * invalid - Некорректный тип.
	// * text - Произвольный текст.
	// * inputNumber - Целочисленный ответ.
	// * checkbox - Множественный выбор.
	// * radio - Одиночный выбор.
	// * checkboxImg - Множественный выбор с изображениями.
	// * radioImg - Одиночный выбор с изображениями.
	Type           string                     `json:"type" enums:"invalid,text,inputNumber,checkbox,radio,checkboxImg,radioImg"`
	IsRequired     bool                       `json:"is_required"`
	Weight         int                        `json:"weight"`
	Rules          *SurveyQuestionRules       `json:"rules,omitempty"`
	AnswerVariants []*SurveyAnswerVariantInfo `json:"answers,omitempty"`
}

type SurveyAnswerVariantInfo struct {
	ID          uuid.UUID `json:"id"`
	Text        string    `json:"text"`
	WithContent bool      `json:"with_content"`
	// Тип варианта ответа:
	// * 0 - Некорректный тип.
	// * 1 - Контент в текстовом формате.
	// * 2 - Контент в целочисленном формате.
	ContentType *int                      `json:"content_type,omitempty" enums:"0,1,2"`
	Rules       *SurveyAnswerVariantRules `json:"rules,omitempty"`
	Image       *SurveyImageInfo          `json:"image,omitempty"`
	Weight      int                       `json:"weight"`
}

type SurveyImageInfo struct {
	ID string `json:"id"`
}

type GetAllSurveys struct {
	SurveyIDs   []uuid.UUID              `json:"survey_ids,omitempty"`
	Respondents *GetAllSurveysRespondent `json:"respondents,omitempty"`
	Options     SurveysOptions           `json:"options"`
	Pagination  GetAllSurveysPagination  `json:"pagination"`
}

type GetAllSurveysRespondent struct {
	// Тип респондента:
	// * 0 - Неизвестный респондент.
	// * 1 - Анонимный респондент.
	// * 2 - Авторизованный респондент.
	RespondentType int         `json:"respondent_type" enums:"0,1,2"`
	RespondentIDs  []uuid.UUID `json:"respondent_ids,omitempty"`
}

type SurveysOptions struct {
	WithQuestions         bool `json:"with_questions"`
	WithAnswers           bool `json:"with_answers"`
	WithDeleted           bool `json:"with_deleted"`
	WithInactiveQuestions bool `json:"with_inactive_questions"`
}

type GetAllSurveysPagination struct {
	Limit    int        `json:"limit"`
	LastId   *uuid.UUID `json:"last_id,omitempty"`
	LastDate *time.Time `json:"last_date,omitempty"`
	Total    *int       `json:"total,omitempty"`
}

type SurveysWithPagination struct {
	Surveys    []*Survey        `json:"survey"`
	Pagination SurveyPagination `json:"pagination"`
}

type SurveyPagination struct {
	Limit    int        `json:"limit"`
	LastId   *uuid.UUID `json:"last_id,omitempty"`
	LastDate *time.Time `json:"last_date,omitempty"`
	Total    int        `json:"total"`
}

type IDResponse struct {
	ID string `json:"id"`
}
