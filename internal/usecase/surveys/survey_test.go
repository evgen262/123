package surveys

import (
	"context"
	"errors"
	"fmt"
	"testing"

	surveys "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/survey"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_surveysUseCase_Get(t *testing.T) {
	type fields struct {
		repo *MockSurveyRepository
	}
	type args struct {
		ctx     context.Context
		id      surveys.SurveyID
		options surveys.SurveyFilterOptions
	}

	ctx := context.TODO()
	testErr := errors.New("some error")
	testUUID := uuid.New()
	answerID := surveys.AnswerID(testUUID)
	questionID := surveys.QuestionID(testUUID)
	surveyID := surveys.SurveyID(testUUID)

	tests := []struct {
		name string
		args args
		want func(a args, f fields) (*surveys.Survey, error)
	}{
		{
			name: "correct",
			args: args{
				ctx:     ctx,
				id:      surveyID,
				options: surveys.SurveyFilterOptions{true, true, false, false},
			},
			want: func(a args, f fields) (*surveys.Survey, error) {
				surveyEntity := &surveys.Survey{
					ID: &surveyID,
					Questions: []*surveys.Question{
						{
							ID:       &questionID,
							SurveyID: &surveyID,
							Answers: []*surveys.AnswerVariant{
								{
									ID:         &answerID,
									QuestionID: &questionID,
									Image: &surveys.Image{
										ID: "test id",
										ExternalImageInfo: &surveys.ExternalProperties{
											ID:   testUUID,
											Size: 1,
										},
									},
									WithContent: false,
								},
							},
						},
					},
					Respondent: &surveys.SurveyRespondent{
						Type: surveys.RespondentTypeAnonymous,
					},
				}
				f.repo.EXPECT().Get(a.ctx, a.id, a.options).Return(surveyEntity, nil)
				return surveyEntity, nil
			},
		},
		{
			name: "err",
			args: args{
				ctx:     ctx,
				id:      surveyID,
				options: surveys.SurveyFilterOptions{true, true, false, false},
			},
			want: func(a args, f fields) (*surveys.Survey, error) {
				f.repo.EXPECT().Get(a.ctx, a.id, a.options).Return(nil, testErr)
				return nil, fmt.Errorf("can't get survey from repository: %w", testErr)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				repo: NewMockSurveyRepository(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			sr := NewSurveysUseCase(f.repo)
			got, err := sr.Get(tt.args.ctx, tt.args.id, tt.args.options)
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
