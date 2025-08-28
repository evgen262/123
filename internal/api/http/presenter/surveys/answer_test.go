package surveys

import (
	"testing"

	viewSurveys "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/surveys"
	entitySurvey "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/survey"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSurveyAnswersPresenter_ToNewEntities(t *testing.T) {
	type fields struct {
	}

	type args struct {
		answers *viewSurveys.NewSurveyAnswers
	}

	testUUID := uuid.New()
	rID := entitySurvey.RespondentID(testUUID)

	tests := []struct {
		name string
		args args
		want func(a args, f fields) []*entitySurvey.RespondentAnswer
	}{
		{
			name: "correct",
			args: args{
				answers: &viewSurveys.NewSurveyAnswers{
					RespondentId: &testUUID,
					Answers: []*viewSurveys.NewSurveyAnswer{
						{ChosenVariant: testUUID},
					},
				},
			},
			want: func(a args, f fields) []*entitySurvey.RespondentAnswer {
				return []*entitySurvey.RespondentAnswer{{
					ChosenVariant: entitySurvey.AnswerID(testUUID),
					RespondentId:  &rID,
				}}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fields{}
			want := tt.want(tt.args, f)
			s := NewAnswersPresenter()
			got := s.ToNewEntities(tt.args.answers)
			assert.Equal(t, want, got)
		})
	}
}

func TestSurveyAnswersPresenter_ToViews(t *testing.T) {
	type fields struct {
	}

	type args struct {
		answers []*entitySurvey.RespondentAnswer
	}

	testUUID := uuid.New()
	aID := entitySurvey.AnswerID(testUUID)
	rID := entitySurvey.RespondentID(testUUID)
	qID := entitySurvey.QuestionID(testUUID)
	testContent := "content"

	tests := []struct {
		name string
		args args
		want func(a args, f fields) []*viewSurveys.SurveyAnswer
	}{
		{
			name: "correct",
			args: args{
				answers: []*entitySurvey.RespondentAnswer{{
					ID:            &testUUID,
					QuestionID:    &qID,
					ChosenVariant: aID,
					RespondentId:  &rID,
					Content:       testContent,
				}},
			},
			want: func(a args, f fields) []*viewSurveys.SurveyAnswer {
				return []*viewSurveys.SurveyAnswer{{
					ID:            testUUID,
					RespondentId:  &testUUID,
					QuestionId:    testUUID,
					ChosenVariant: testUUID,
					Content:       testContent,
				}}
			},
		},
		{
			name: "correct nil input",
			args: args{
				answers: []*entitySurvey.RespondentAnswer{{}},
			},
			want: func(a args, f fields) []*viewSurveys.SurveyAnswer {
				return []*viewSurveys.SurveyAnswer{{}}
			},
		},
		{
			name: "correct no respondent id",
			args: args{
				answers: []*entitySurvey.RespondentAnswer{{
					ID:            &testUUID,
					QuestionID:    &qID,
					ChosenVariant: aID,
					Content:       testContent,
				}},
			},
			want: func(a args, f fields) []*viewSurveys.SurveyAnswer {
				return []*viewSurveys.SurveyAnswer{{
					ID:            testUUID,
					RespondentId:  nil,
					QuestionId:    testUUID,
					ChosenVariant: testUUID,
					Content:       testContent,
				}}
			},
		},
		{
			name: "correct no respondent, question id",
			args: args{
				answers: []*entitySurvey.RespondentAnswer{{
					ID:            &testUUID,
					ChosenVariant: aID,
					Content:       testContent,
				}},
			},
			want: func(a args, f fields) []*viewSurveys.SurveyAnswer {
				return []*viewSurveys.SurveyAnswer{{
					ID:            testUUID,
					RespondentId:  nil,
					QuestionId:    uuid.Nil,
					ChosenVariant: testUUID,
					Content:       testContent,
				}}
			},
		},
		{
			name: "correct no respondent, question, main id",
			args: args{
				answers: []*entitySurvey.RespondentAnswer{{
					ChosenVariant: aID,
					Content:       testContent,
				}},
			},
			want: func(a args, f fields) []*viewSurveys.SurveyAnswer {
				return []*viewSurveys.SurveyAnswer{{
					ID:            uuid.Nil,
					RespondentId:  nil,
					QuestionId:    uuid.Nil,
					ChosenVariant: testUUID,
					Content:       testContent,
				}}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fields{}
			want := tt.want(tt.args, f)
			s := NewAnswersPresenter()
			got := s.ToViews(tt.args.answers)
			assert.Equal(t, want, got)
		})
	}
}

func TestSurveyAnswersPresenter_ToShortViews(t *testing.T) {
	type fields struct {
	}

	type args struct {
		ids []uuid.UUID
	}

	testUUID := uuid.New()

	tests := []struct {
		name string
		args args
		want func(a args, f fields) []*viewSurveys.SurveyAnswerInfo
	}{
		{
			name: "correct",
			args: args{
				ids: []uuid.UUID{testUUID},
			},
			want: func(a args, f fields) []*viewSurveys.SurveyAnswerInfo {
				return []*viewSurveys.SurveyAnswerInfo{{
					ID: testUUID,
				}}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fields{}
			want := tt.want(tt.args, f)
			s := NewAnswersPresenter()
			got := s.ToShortViews(tt.args.ids)
			assert.Equal(t, want, got)
		})
	}
}
