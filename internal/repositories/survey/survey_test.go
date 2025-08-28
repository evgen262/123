package survey

import (
	"context"
	"errors"
	"fmt"
	"testing"

	answervariantv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/surveys/answervariant/v1"
	questionv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/surveys/question/v1"
	respondentv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/surveys/respondent/v1"
	surveyv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/surveys/survey/v1"
	surveys "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/survey"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	timeUtils "git.mos.ru/buch-cloud/moscow-team-2.0/build/time-utils.git"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	mapper "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/client/grpc/surveys"
)

func Test_surveyRepository_Get(t *testing.T) {
	type fields struct {
		client *surveyv1.MockSurveyAPIClient
		mapper *MockSurveyMapper
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
	tUtils := timeUtils.NewTimeUtils()
	localMapper := mapper.NewSurveyMapper(tUtils)

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
				optionsMapped := localMapper.OptionsToPb(&a.options)
				surveyPb := &surveyv1.Survey{
					Id: testUUID.String(),
					Questions: []*questionv1.Question{
						{
							Id:       testUUID.String(),
							SurveyId: testUUID.String(),
							AnswerVariants: []*answervariantv1.AnswerVariant{
								{
									Id:          testUUID.String(),
									QuestionId:  testUUID.String(),
									WithContent: false,
									Image: &answervariantv1.AnswerVariant_Image{
										Id: "test id",
										ExternalImageInfo: &answervariantv1.AnswerVariant_Image_ExternalInfo{
											Id:   testUUID.String(),
											Size: 1,
										},
									},
								},
							},
						},
					},
					Respondent: &respondentv1.Respondent{
						Respondent: &respondentv1.Respondent_Anonymous_{},
					},
				}
				f.client.EXPECT().Get(a.ctx, &surveyv1.GetRequest{
					Id:      a.id.String(),
					Options: optionsMapped,
				}).Return(&surveyv1.GetResponse{
					Survey: surveyPb,
				}, nil)
				f.mapper.EXPECT().OptionsToPb(&a.options).Return(optionsMapped)
				surveyEntity := &surveys.Survey{
					ID:                &surveyID,
					ActivePeriodStart: surveyPb.GetActivePeriodStartTime().AsTime(),
					ActivePeriodEnd:   surveyPb.GetActivePeriodEndTime().AsTime(),
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
				f.mapper.EXPECT().SurveyToEntity(surveyPb).Return(surveyEntity, nil)
				return surveyEntity, nil
			},
		},
		{
			name: "convert survey error",
			args: args{
				ctx:     ctx,
				id:      surveyID,
				options: surveys.SurveyFilterOptions{true, true, false, false},
			},
			want: func(a args, f fields) (*surveys.Survey, error) {
				optionsMapped := localMapper.OptionsToPb(&a.options)
				surveyPb := &surveyv1.Survey{
					Id: testUUID.String(),
					Questions: []*questionv1.Question{
						{
							Id:       testUUID.String(),
							SurveyId: testUUID.String(),
							AnswerVariants: []*answervariantv1.AnswerVariant{
								{
									Id:          testUUID.String(),
									QuestionId:  testUUID.String(),
									WithContent: false,
									Image: &answervariantv1.AnswerVariant_Image{
										Id: "test id",
										ExternalImageInfo: &answervariantv1.AnswerVariant_Image_ExternalInfo{
											Id:   testUUID.String(),
											Size: 1,
										},
									},
								},
							},
						},
					},
					Respondent: &respondentv1.Respondent{
						Respondent: &respondentv1.Respondent_Anonymous_{},
					},
				}
				f.client.EXPECT().Get(a.ctx, &surveyv1.GetRequest{
					Id:      a.id.String(),
					Options: optionsMapped,
				}).Return(&surveyv1.GetResponse{
					Survey: surveyPb,
				}, nil)
				f.mapper.EXPECT().OptionsToPb(&a.options).Return(optionsMapped)
				f.mapper.EXPECT().SurveyToEntity(surveyPb).Return(nil, testErr)
				return nil, fmt.Errorf("can't convert survey to entity: %w", testErr)
			},
		},
		{
			name: "get survey default error",
			args: args{
				ctx:     ctx,
				id:      surveyID,
				options: surveys.SurveyFilterOptions{true, true, false, false},
			},
			want: func(a args, f fields) (*surveys.Survey, error) {
				optionsMapped := localMapper.OptionsToPb(&a.options)
				f.client.EXPECT().Get(a.ctx, &surveyv1.GetRequest{
					Id:      a.id.String(),
					Options: optionsMapped,
				}).Return(nil, status.Error(codes.Internal, testErr.Error()))
				f.mapper.EXPECT().OptionsToPb(&a.options).Return(optionsMapped)
				return nil, fmt.Errorf("can't get survey: %s", testErr.Error())
			},
		},
		{
			name: "get survey invalid argument error",
			args: args{
				ctx:     ctx,
				id:      surveyID,
				options: surveys.SurveyFilterOptions{true, true, false, false},
			},
			want: func(a args, f fields) (*surveys.Survey, error) {
				optionsMapped := localMapper.OptionsToPb(&a.options)
				f.client.EXPECT().Get(a.ctx, &surveyv1.GetRequest{
					Id:      a.id.String(),
					Options: optionsMapped,
				}).Return(nil, status.Error(codes.InvalidArgument, testErr.Error()))
				f.mapper.EXPECT().OptionsToPb(&a.options).Return(optionsMapped)
				return nil, diterrors.NewValidationError(testErr)
			},
		},
		{
			name: "get survey not found error",
			args: args{
				ctx:     ctx,
				id:      surveyID,
				options: surveys.SurveyFilterOptions{true, true, false, false},
			},
			want: func(a args, f fields) (*surveys.Survey, error) {
				optionsMapped := localMapper.OptionsToPb(&a.options)
				f.client.EXPECT().Get(a.ctx, &surveyv1.GetRequest{
					Id:      a.id.String(),
					Options: optionsMapped,
				}).Return(nil, status.Error(codes.NotFound, testErr.Error()))
				f.mapper.EXPECT().OptionsToPb(&a.options).Return(optionsMapped)
				return nil, diterrors.ErrNotFound
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				client: surveyv1.NewMockSurveyAPIClient(ctrl),
				mapper: NewMockSurveyMapper(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			sr := NewSurveyRepository(f.client, f.mapper)
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
