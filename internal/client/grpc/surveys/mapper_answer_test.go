package surveys

import (
	"fmt"
	"testing"

	answerv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/surveys/answer/v1"
	surveys "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/survey"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAnswerMapper_AnswerToEntity(t *testing.T) {
	type args struct {
		answer *answerv1.Answer
	}

	testContent := "some content"

	testID := uuid.New()
	aID := surveys.AnswerID(testID)
	qID := surveys.QuestionID(testID)
	rID := surveys.RespondentID(testID)

	tests := []struct {
		name string
		args args
		want func(a args) (*surveys.RespondentAnswer, error)
	}{
		{
			name: "correct respondent authorized with empty content",
			args: args{answer: &answerv1.Answer{
				Id:            testID.String(),
				ChosenVariant: testID.String(),
				Respondent:    &answerv1.Respondent{Id: testID.String()},
				QuestionId:    testID.String(),
			}},
			want: func(a args) (*surveys.RespondentAnswer, error) {
				return &surveys.RespondentAnswer{
					ID:            &testID,
					ChosenVariant: aID,
					RespondentId:  &rID,
					QuestionID:    &qID,
				}, nil
			},
		},
		{
			name: "correct respondent authorized with content",
			args: args{answer: &answerv1.Answer{
				Id:            testID.String(),
				ChosenVariant: testID.String(),
				Respondent:    &answerv1.Respondent{Id: testID.String()},
				Content:       testContent,
				QuestionId:    testID.String(),
			}},
			want: func(a args) (*surveys.RespondentAnswer, error) {
				return &surveys.RespondentAnswer{
					ID:            &testID,
					ChosenVariant: aID,
					RespondentId:  &rID,
					Content:       testContent,
					QuestionID:    &qID,
				}, nil
			},
		},
		{
			name: "correct respondent anonymous with empty content ",
			args: args{answer: &answerv1.Answer{
				Id:            testID.String(),
				ChosenVariant: testID.String(),
				QuestionId:    testID.String(),
			}},
			want: func(a args) (*surveys.RespondentAnswer, error) {
				return &surveys.RespondentAnswer{
					ID:            &testID,
					ChosenVariant: aID,
					QuestionID:    &qID,
				}, nil
			},
		},
		{
			name: "incorrect answer id",
			args: args{answer: &answerv1.Answer{
				Id:            "not corrected",
				ChosenVariant: testID.String(),
				Content:       testContent,
				QuestionId:    testID.String(),
			}},
			want: func(a args) (*surveys.RespondentAnswer, error) {
				return nil, fmt.Errorf("can't parse answer id: invalid UUID length: 13")
			},
		},
		{
			name: "incorrect chosen variant id",
			args: args{answer: &answerv1.Answer{
				Id:            testID.String(),
				ChosenVariant: "test",
				QuestionId:    testID.String(),
			}},
			want: func(a args) (*surveys.RespondentAnswer, error) {
				return nil, fmt.Errorf("can't parse chosen variant id: invalid UUID length: 4")
			},
		},
		{
			name: "incorrect question id",
			args: args{answer: &answerv1.Answer{
				Id:            testID.String(),
				ChosenVariant: testID.String(),
				QuestionId:    "test",
			}},
			want: func(a args) (*surveys.RespondentAnswer, error) {
				return nil, fmt.Errorf("can't parse question id: invalid UUID length: 4")
			},
		},
		{
			name: "incorrect user id",
			args: args{answer: &answerv1.Answer{
				Id:            testID.String(),
				ChosenVariant: testID.String(),
				Respondent:    &answerv1.Respondent{Id: "test"},
				QuestionId:    testID.String(),
			}},
			want: func(a args) (*surveys.RespondentAnswer, error) {
				return nil, fmt.Errorf("can't parse respondent id: invalid UUID length: 4")

			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want, wantErr := tt.want(tt.args)
			s := NewAnswerMapper()
			got, err := s.answerToEntity(tt.args.answer)
			if wantErr != nil {
				assert.EqualError(t, err, wantErr.Error())
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, want, got)
			}
		})
	}
}

func TestAnswerMapper_AnswersToEntities(t *testing.T) {
	type args struct {
		answers []*answerv1.Answer
	}

	testID := uuid.New()
	aID := surveys.AnswerID(testID)
	qID := surveys.QuestionID(testID)
	rID := surveys.RespondentID(testID)

	tests := []struct {
		name string
		args args
		want func(a args) ([]*surveys.RespondentAnswer, error)
	}{
		{
			name: "correct respondent authorized with empty content",
			args: args{answers: []*answerv1.Answer{{
				Id:            testID.String(),
				ChosenVariant: testID.String(),
				Respondent:    &answerv1.Respondent{Id: testID.String()},
				QuestionId:    testID.String(),
			}}},
			want: func(a args) ([]*surveys.RespondentAnswer, error) {
				return []*surveys.RespondentAnswer{{
					ID:            &testID,
					ChosenVariant: aID,
					RespondentId:  &rID,
					QuestionID:    &qID,
				}}, nil
			},
		},
		{
			name: "err",
			args: args{answers: []*answerv1.Answer{{
				Id: "test",
			}}},
			want: func(a args) ([]*surveys.RespondentAnswer, error) {
				return nil, fmt.Errorf("can't convert answer to entity: can't parse answer id: invalid UUID length: 4")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want, wantErr := tt.want(tt.args)
			s := NewAnswerMapper()
			got, err := s.AnswersToEntities(tt.args.answers)
			if wantErr != nil {
				assert.EqualError(t, err, wantErr.Error())
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, want, got)
			}
		})
	}
}

func TestAnswerMapper_NewAnswersToPb(t *testing.T) {
	type args struct {
		answer []*surveys.RespondentAnswer
	}

	testContent := "some content"
	testChoseId := uuid.New()

	tests := []struct {
		name string
		args args
		want []*answerv1.AddRequest_Answer
	}{
		{
			name: "correct",
			args: args{[]*surveys.RespondentAnswer{
				{
					ChosenVariant: surveys.AnswerID(testChoseId),
					Content:       testContent,
				},
			}},
			want: []*answerv1.AddRequest_Answer{
				{
					ChosenVariant: testChoseId.String(),
					Content:       testContent,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ap := NewAnswerMapper()
			assert.Equal(t, tt.want, ap.NewAnswersToPb(tt.args.answer))
		})
	}
}

func TestAnswerMapper_IDsToUUIDs(t *testing.T) {
	type args struct {
		answerIDs []string
	}

	id := "6ba7b814-9dad-11d1-80b4-00c04fd430c9"
	testUuid, err := uuid.Parse(id)
	assert.NoError(t, err)

	tests := []struct {
		name string
		args args
		want func(a args) ([]uuid.UUID, error)
	}{
		{
			name: "correct",
			args: args{[]string{id}},
			want: func(a args) ([]uuid.UUID, error) {
				return []uuid.UUID{testUuid}, nil
			},
		},
		{
			name: "err",
			args: args{[]string{"test"}},
			want: func(a args) ([]uuid.UUID, error) {
				return nil, fmt.Errorf("can't parse answer id: invalid UUID length: 4")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qp := NewAnswerMapper()
			want, wantErr := tt.want(tt.args)
			got, err := qp.IDsToUUIDs(tt.args.answerIDs)
			if wantErr != nil {
				assert.EqualError(t, err, wantErr.Error())
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, want, got)
			}
		})
	}
}
