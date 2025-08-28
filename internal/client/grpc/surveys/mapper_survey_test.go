package surveys

import (
	"fmt"
	"testing"
	"time"

	answervariantv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/surveys/answervariant/v1"
	questionv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/surveys/question/v1"
	respondentv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/surveys/respondent/v1"
	sharedv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/surveys/shared/v1"
	surveyv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/surveys/survey/v1"
	surveys "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/survey"
	timeUtils "git.mos.ru/buch-cloud/moscow-team-2.0/build/time-utils.git"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestSurveyMapper_PaginationToEntity(t *testing.T) {
	type fields struct {
		timeUtils *timeUtils.MockTimeUtils
	}
	type args struct {
		response *sharedv1.PaginationResponse
	}

	testTime := time.Date(2023, 9, 5, 0, 0, 0, 0, time.UTC)
	testId := uuid.New()
	lastId := surveys.SurveyID(testId)

	tests := []struct {
		name string
		args args
		want func(a args, f fields) (*surveys.Pagination, error)
	}{
		{
			name: "correct",
			args: args{
				response: &sharedv1.PaginationResponse{
					Limit:           10,
					LastId:          wrapperspb.String(testId.String()),
					LastCreatedTime: timestamppb.New(testTime),
					Total:           250,
				},
			},
			want: func(a args, f fields) (*surveys.Pagination, error) {
				f.timeUtils.EXPECT().TimestampToTime(a.response.LastCreatedTime).Return(&testTime)
				return &surveys.Pagination{
					Limit:    a.response.Limit,
					LastId:   &lastId,
					LastDate: &testTime,
					Total:    250,
				}, nil
			},
		},
		{
			name: "error parse last survey id",
			args: args{
				response: &sharedv1.PaginationResponse{
					Limit:           10,
					LastId:          wrapperspb.String("test uuid"),
					LastCreatedTime: timestamppb.New(testTime),
					Total:           250,
				},
			},
			want: func(a args, f fields) (*surveys.Pagination, error) {
				f.timeUtils.EXPECT().TimestampToTime(a.response.LastCreatedTime).Return(&testTime)
				return nil, fmt.Errorf("can't parse last survey id in pagination options: invalid UUID length: 9")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				timeUtils: timeUtils.NewMockTimeUtils(ctrl),
			}
			p := NewSurveyMapper(f.timeUtils)
			want, wantErr := tt.want(tt.args, f)
			got, err := p.PaginationToEntity(tt.args.response)
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

func TestSurveyMapper_PaginationToPb(t *testing.T) {
	type fields struct {
		timeUtils *timeUtils.MockTimeUtils
	}

	type args struct {
		pagination *surveys.Pagination
	}
	testTime := time.Date(2023, 9, 5, 0, 0, 0, 0, time.UTC)
	testId := uuid.New()
	lastId := surveys.SurveyID(testId)

	tests := []struct {
		name string
		args args
		want func(a args, f fields) *sharedv1.PaginationRequest
	}{
		{
			name: "correct",
			args: args{pagination: &surveys.Pagination{
				Limit:    10,
				LastId:   &lastId,
				LastDate: &testTime,
				Total:    0,
			}},
			want: func(a args, f fields) *sharedv1.PaginationRequest {
				t := timestamppb.New(testTime)
				f.timeUtils.EXPECT().TimeToTimestamp(a.pagination.LastDate).Return(t)
				return &sharedv1.PaginationRequest{
					Limit:           a.pagination.Limit,
					LastId:          wrapperspb.String(testId.String()),
					LastCreatedTime: t,
				}
			},
		},
		{
			name: "nil input",
			args: args{pagination: nil},
			want: func(a args, f fields) *sharedv1.PaginationRequest {
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				timeUtils: timeUtils.NewMockTimeUtils(ctrl),
			}
			p := NewSurveyMapper(f.timeUtils)
			want := tt.want(tt.args, f)
			got := p.PaginationToPb(tt.args.pagination)
			assert.Equal(t, want, got)
		})
	}
}

func TestSurveyMapper_RespondentIDsToEntity(t *testing.T) {
	type fields struct {
		timeUtils *timeUtils.MockTimeUtils
	}
	type args struct {
		respondentIDs []string
	}

	testID := "aaa7b814-9cad-11d1-80b4-00c04fd410c8"
	testUUID, _ := uuid.Parse(testID)
	tests := []struct {
		name string
		args args
		want func(a args, f fields) (surveys.RespondentIDs, error)
	}{
		{
			name: "correct",
			args: args{
				respondentIDs: []string{testID},
			},
			want: func(a args, f fields) (surveys.RespondentIDs, error) {
				repIDS := []surveys.RespondentID{surveys.RespondentID(testUUID)}
				return repIDS, nil
			},
		},
		{
			name: "uuid parse error",
			args: args{
				respondentIDs: []string{"testID"},
			},
			want: func(a args, f fields) (surveys.RespondentIDs, error) {
				return nil, fmt.Errorf("can't parse uuid: invalid UUID length: 6")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				timeUtils: timeUtils.NewMockTimeUtils(ctrl),
			}
			ap := NewSurveyMapper(f.timeUtils)
			want, wantErr := tt.want(tt.args, f)
			got, err := ap.RespondentIDsToEntity(tt.args.respondentIDs)
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

func TestSurveyMapper_RespondentToEntity(t *testing.T) {
	type fields struct {
		timeUtils *timeUtils.MockTimeUtils
	}
	type args struct {
		respondent *respondentv1.Respondent
	}

	testID := "aaa7b814-9cad-11d1-80b4-00c04fd410c8"
	testUUID, _ := uuid.Parse(testID)
	tests := []struct {
		name string
		args args
		want func(a args, f fields) (*surveys.SurveyRespondent, error)
	}{
		{
			name: "correct with unknown",
			args: args{
				respondent: &respondentv1.Respondent{
					Respondent: &respondentv1.Respondent_Unknown_{},
				},
			},
			want: func(a args, f fields) (*surveys.SurveyRespondent, error) {
				return &surveys.SurveyRespondent{Type: surveys.RespondentTypeUnknown}, nil
			},
		},
		{
			name: "correct with anonymous",
			args: args{
				respondent: &respondentv1.Respondent{
					Respondent: &respondentv1.Respondent_Anonymous_{},
				},
			},
			want: func(a args, f fields) (*surveys.SurveyRespondent, error) {
				return &surveys.SurveyRespondent{Type: surveys.RespondentTypeAnonymous}, nil
			},
		},
		{
			name: "correct with user",
			args: args{
				respondent: &respondentv1.Respondent{
					Respondent: &respondentv1.Respondent_User_{User: &respondentv1.Respondent_User{
						Ids: []string{testID},
					}},
				},
			},
			want: func(a args, f fields) (*surveys.SurveyRespondent, error) {
				return &surveys.SurveyRespondent{
						Type: surveys.RespondentTypeUser,
						Ids:  surveys.RespondentIDs{surveys.RespondentID(testUUID)}},
					nil
			},
		},
		{
			name: "uuid parse error",
			args: args{
				respondent: &respondentv1.Respondent{
					Respondent: &respondentv1.Respondent_User_{User: &respondentv1.Respondent_User{
						Ids: []string{"testID"},
					}},
				},
			},
			want: func(a args, f fields) (*surveys.SurveyRespondent, error) {
				return nil, fmt.Errorf("can't parse uuid: invalid UUID length: 6")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				timeUtils: timeUtils.NewMockTimeUtils(ctrl),
			}
			ap := NewSurveyMapper(f.timeUtils)
			want, wantErr := tt.want(tt.args, f)
			got, err := ap.RespondentToEntity(tt.args.respondent)
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

func TestSurveyMapper_OptionsToPb(t *testing.T) {
	type fields struct {
		timeUtils *timeUtils.MockTimeUtils
	}
	type args struct {
		options *surveys.SurveyFilterOptions
	}

	tests := []struct {
		name string
		args args
		want func(a args, f fields) *sharedv1.Options
	}{
		{
			name: "nil options correct",
			args: args{
				options: nil,
			},
			want: func(a args, f fields) *sharedv1.Options {
				return nil
			},
		},
		{
			name: "correct",
			args: args{
				options: &surveys.SurveyFilterOptions{
					WithQuestions:         true,
					WithAnswers:           true,
					WithInactiveQuestions: true,
				},
			},
			want: func(a args, f fields) *sharedv1.Options {
				return &sharedv1.Options{
					WithQuestions:         true,
					WithAnswers:           true,
					WithInactiveQuestions: true,
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				timeUtils: timeUtils.NewMockTimeUtils(ctrl),
			}
			ap := NewSurveyMapper(f.timeUtils)
			want := tt.want(tt.args, f)
			got := ap.OptionsToPb(tt.args.options)
			assert.Equal(t, want, got)
		})
	}
}

func TestSurveyMapper_SurveyToPb(t *testing.T) {
	type fields struct {
		timeUtils *timeUtils.MockTimeUtils
	}
	type args struct {
		survey *surveys.Survey
	}
	testINT := 5
	testID := "aaa7b814-9cad-11d1-80b4-00c04fd410c8"
	testUUID, err := uuid.Parse(testID)
	assert.NoError(t, err)
	answerID := surveys.AnswerID(testUUID)
	questionID := surveys.QuestionID(testUUID)
	surveyID := surveys.SurveyID(testUUID)

	tests := []struct {
		name string
		args args
		want func(a args, f fields) *surveyv1.Survey
	}{
		{
			name: "correct with unknown",
			args: args{
				survey: &surveys.Survey{
					ID:    &surveyID,
					Title: "survey",
					Respondent: &surveys.SurveyRespondent{
						Type: surveys.RespondentTypeUnknown,
					},
					Questions: []*surveys.Question{{
						ID:       &questionID,
						SurveyID: &surveyID,
						Text:     "question",
						Rules: &surveys.QuestionRules{
							PickMaxCount: &testINT,
							PickMinCount: &testINT,
						},
						Answers: []*surveys.AnswerVariant{{
							ID:          &answerID,
							QuestionID:  &questionID,
							Text:        "answer variant",
							WithContent: true,
							Weight:      10,
							Image: &surveys.Image{
								ID: "image",
								ExternalImageInfo: &surveys.ExternalProperties{
									ID:       testUUID,
									FileName: "filename",
									URL:      "url",
									Size:     15,
								},
							},
							Rules: &surveys.AnswerRules{
								ContentMinDigit:  &testINT,
								ContentMaxDigit:  &testINT,
								ContentMinLength: &testINT,
								ContentMaxLength: &testINT,
							},
						}},
					}},
				},
			},
			want: func(a args, f fields) *surveyv1.Survey {
				f.timeUtils.EXPECT().TimeToTimestamp(&a.survey.ActivePeriodStart).Return(nil)
				f.timeUtils.EXPECT().TimeToTimestamp(&a.survey.ActivePeriodEnd).Return(nil)
				f.timeUtils.EXPECT().TimeToTimestamp(a.survey.CreatedAt).Return(nil)
				f.timeUtils.EXPECT().TimeToTimestamp(a.survey.UpdatedAt).Return(nil)
				f.timeUtils.EXPECT().TimeToTimestamp(a.survey.DeletedAt).Return(nil)
				f.timeUtils.EXPECT().TimeToTimestamp(a.survey.Questions[0].DeletedAt).Return(nil)
				f.timeUtils.EXPECT().TimeToTimestamp(a.survey.Questions[0].Answers[0].DeletedAt).Return(nil)
				result := &surveyv1.Survey{
					Id:    testUUID.String(),
					Title: "survey",
					Respondent: &respondentv1.Respondent{
						Respondent: &respondentv1.Respondent_Unknown_{},
					},
					Questions: []*questionv1.Question{{
						Id:       testUUID.String(),
						SurveyId: testUUID.String(),
						Text:     "question",
						Rules: &questionv1.Question_Rules{
							PickMinCount: &wrapperspb.Int64Value{Value: int64(testINT)},
							PickMaxCount: &wrapperspb.Int64Value{Value: int64(testINT)},
						},
						AnswerVariants: []*answervariantv1.AnswerVariant{{
							Id:          testUUID.String(),
							QuestionId:  testUUID.String(),
							Text:        "answer variant",
							WithContent: true,
							Weight:      10,
							Image: &answervariantv1.AnswerVariant_Image{
								Id: "image",
								ExternalImageInfo: &answervariantv1.AnswerVariant_Image_ExternalInfo{
									Id:       testUUID.String(),
									Filename: "filename",
									Url:      "url",
									Size:     15,
								},
							},
							Rules: &answervariantv1.AnswerVariant_Rules{
								ContentMinDigit:  &wrapperspb.Int64Value{Value: int64(testINT)},
								ContentMaxDigit:  &wrapperspb.Int64Value{Value: int64(testINT)},
								ContentMinLength: &wrapperspb.Int64Value{Value: int64(testINT)},
								ContentMaxLength: &wrapperspb.Int64Value{Value: int64(testINT)},
							},
						}},
					}},
				}
				return result
			},
		},
		{
			name: "correct with user",
			args: args{
				survey: &surveys.Survey{
					ID:    &surveyID,
					Title: "survey",
					Respondent: &surveys.SurveyRespondent{
						Type: surveys.RespondentTypeUser,
						Ids:  surveys.RespondentIDs{surveys.RespondentID(testUUID)},
					},
					Questions: []*surveys.Question{{
						ID:       &questionID,
						SurveyID: &surveyID,
						Text:     "question",
						Answers: []*surveys.AnswerVariant{{
							ID:          &answerID,
							QuestionID:  &questionID,
							Text:        "answer variant",
							WithContent: true,
							Image: &surveys.Image{
								ID: "image",
								ExternalImageInfo: &surveys.ExternalProperties{
									ID:       testUUID,
									FileName: "filename",
									URL:      "url",
									Size:     15,
								},
							},
							Rules: &surveys.AnswerRules{
								ContentMinDigit:  &testINT,
								ContentMaxDigit:  &testINT,
								ContentMinLength: &testINT,
								ContentMaxLength: &testINT,
							},
						}},
					}},
				},
			},
			want: func(a args, f fields) *surveyv1.Survey {
				f.timeUtils.EXPECT().TimeToTimestamp(&a.survey.ActivePeriodStart).Return(nil)
				f.timeUtils.EXPECT().TimeToTimestamp(&a.survey.ActivePeriodEnd).Return(nil)
				f.timeUtils.EXPECT().TimeToTimestamp(a.survey.CreatedAt).Return(nil)
				f.timeUtils.EXPECT().TimeToTimestamp(a.survey.UpdatedAt).Return(nil)
				f.timeUtils.EXPECT().TimeToTimestamp(a.survey.DeletedAt).Return(nil)
				f.timeUtils.EXPECT().TimeToTimestamp(a.survey.Questions[0].DeletedAt).Return(nil)
				f.timeUtils.EXPECT().TimeToTimestamp(a.survey.Questions[0].Answers[0].DeletedAt).Return(nil)
				result := &surveyv1.Survey{
					Id:    testUUID.String(),
					Title: "survey",
					Respondent: &respondentv1.Respondent{
						Respondent: &respondentv1.Respondent_User_{
							User: &respondentv1.Respondent_User{
								Ids: []string{testID},
							},
						},
					},
					Questions: []*questionv1.Question{{
						Id:       testUUID.String(),
						SurveyId: testUUID.String(),
						Text:     "question",
						AnswerVariants: []*answervariantv1.AnswerVariant{{
							Id:          testUUID.String(),
							QuestionId:  testUUID.String(),
							Text:        "answer variant",
							WithContent: true,
							Image: &answervariantv1.AnswerVariant_Image{
								Id: "image",
								ExternalImageInfo: &answervariantv1.AnswerVariant_Image_ExternalInfo{
									Id:       testUUID.String(),
									Filename: "filename",
									Url:      "url",
									Size:     15,
								},
							},
							Rules: &answervariantv1.AnswerVariant_Rules{
								ContentMinDigit:  &wrapperspb.Int64Value{Value: int64(testINT)},
								ContentMaxDigit:  &wrapperspb.Int64Value{Value: int64(testINT)},
								ContentMinLength: &wrapperspb.Int64Value{Value: int64(testINT)},
								ContentMaxLength: &wrapperspb.Int64Value{Value: int64(testINT)},
							},
						}},
					}},
				}
				return result
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				timeUtils: timeUtils.NewMockTimeUtils(ctrl),
			}
			ap := NewSurveyMapper(f.timeUtils)
			want := tt.want(tt.args, f)
			got := ap.SurveyToPb(tt.args.survey)
			assert.Equal(t, want, got)
		})
	}
}

func TestSurveyMapper_NewSurveyToPb(t *testing.T) {
	type args struct {
		survey *surveys.Survey
	}

	type fields struct {
		timeUtils *timeUtils.MockTimeUtils
	}

	testINT := 5
	testUUID := uuid.New()

	tests := []struct {
		name string
		args args
		want func(a args, f fields) *surveyv1.AddRequest_Survey
	}{
		{
			name: "nil add request",
			args: args{survey: nil},
			want: func(a args, f fields) *surveyv1.AddRequest_Survey {
				return nil
			},
		},
		{
			name: "correct with user",
			args: args{
				survey: &surveys.Survey{
					Title: "test title",
					Questions: []*surveys.Question{
						{
							Text: "test question 1",
							Rules: &surveys.QuestionRules{
								PickMaxCount: &testINT,
								PickMinCount: &testINT,
							},
							Answers: []*surveys.AnswerVariant{
								{
									Text:        "test answer 1",
									WithContent: true,
									Image: &surveys.Image{
										ID: "test",
										ExternalImageInfo: &surveys.ExternalProperties{
											ID: testUUID,
										},
									},
									Weight: 10,
									Rules: &surveys.AnswerRules{
										ContentMinLength: &testINT,
										ContentMaxLength: &testINT,
										ContentMinDigit:  &testINT,
										ContentMaxDigit:  &testINT,
									},
								},
							},
						},
					},
					Respondent: &surveys.SurveyRespondent{
						Type: surveys.RespondentTypeUser,
						Ids:  surveys.RespondentIDs{surveys.RespondentID(testUUID)},
					},
				},
			},
			want: func(a args, f fields) *surveyv1.AddRequest_Survey {
				f.timeUtils.EXPECT().TimeToTimestamp(&a.survey.ActivePeriodStart).Return(nil)
				f.timeUtils.EXPECT().TimeToTimestamp(&a.survey.ActivePeriodEnd).Return(nil)
				return &surveyv1.AddRequest_Survey{
					Title: "test title",
					Questions: []*surveyv1.AddRequest_Survey_Question{
						{
							Text: "test question 1",
							Rules: &surveyv1.AddRequest_Survey_Question_Rules{
								PickMinCount: &wrapperspb.Int64Value{Value: int64(testINT)},
								PickMaxCount: &wrapperspb.Int64Value{Value: int64(testINT)},
							},
							AnswerVariants: []*surveyv1.AddRequest_Survey_Question_AnswerVariant{
								{
									Text:        "test answer 1",
									WithContent: true,
									Image: &surveyv1.AddRequest_Survey_Question_AnswerVariant_Image{
										Id: "test",
										ExternalImageInfo: &surveyv1.AddRequest_Survey_Question_AnswerVariant_Image_ExternalInfo{
											Id: testUUID.String(),
										},
									},
									Weight: 10,
									Rules: &surveyv1.AddRequest_Survey_Question_AnswerVariant_Rules{
										ContentMinLength: wrapperspb.Int64(int64(testINT)),
										ContentMaxLength: wrapperspb.Int64(int64(testINT)),
										ContentMinDigit:  wrapperspb.Int64(int64(testINT)),
										ContentMaxDigit:  wrapperspb.Int64(int64(testINT)),
									},
								},
							},
						},
					},
					Respondent: &respondentv1.Respondent{
						Respondent: &respondentv1.Respondent_User_{
							User: &respondentv1.Respondent_User{
								Ids: []string{testUUID.String()},
							},
						},
					},
				}
			},
		},
		{
			name: "correct with anonymous",
			args: args{
				survey: &surveys.Survey{
					Title: "test title",
					Questions: []*surveys.Question{
						{
							Text:  "test question 1",
							Rules: &surveys.QuestionRules{},
							Answers: []*surveys.AnswerVariant{
								{
									Text:        "test answer 1",
									WithContent: true,
									Image: &surveys.Image{
										ID: "test",
										ExternalImageInfo: &surveys.ExternalProperties{
											ID: testUUID,
										},
									},
									Rules: &surveys.AnswerRules{
										ContentMinLength: &testINT,
									},
								},
							},
						},
					},
					Respondent: &surveys.SurveyRespondent{
						Type: surveys.RespondentTypeAnonymous,
					},
				},
			},
			want: func(a args, f fields) *surveyv1.AddRequest_Survey {
				f.timeUtils.EXPECT().TimeToTimestamp(&a.survey.ActivePeriodStart).Return(nil)
				f.timeUtils.EXPECT().TimeToTimestamp(&a.survey.ActivePeriodEnd).Return(nil)
				return &surveyv1.AddRequest_Survey{
					Title: "test title",
					Questions: []*surveyv1.AddRequest_Survey_Question{
						{
							Text:  "test question 1",
							Rules: &surveyv1.AddRequest_Survey_Question_Rules{},
							AnswerVariants: []*surveyv1.AddRequest_Survey_Question_AnswerVariant{
								{
									Text:        "test answer 1",
									WithContent: true,
									Image: &surveyv1.AddRequest_Survey_Question_AnswerVariant_Image{
										Id: "test",
										ExternalImageInfo: &surveyv1.AddRequest_Survey_Question_AnswerVariant_Image_ExternalInfo{
											Id: testUUID.String(),
										},
									},
									Rules: &surveyv1.AddRequest_Survey_Question_AnswerVariant_Rules{
										ContentMinLength: wrapperspb.Int64(int64(testINT)),
									},
								},
							},
						},
					},
					Respondent: &respondentv1.Respondent{
						Respondent: &respondentv1.Respondent_Anonymous_{},
					},
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				timeUtils: timeUtils.NewMockTimeUtils(ctrl),
			}
			want := tt.want(tt.args, f)
			s := NewSurveyMapper(f.timeUtils)
			got := s.NewSurveyToPb(tt.args.survey)
			assert.Equal(t, want, got)
		})
	}
}

func TestSurveyMapper_SurveyToEntity(t *testing.T) {
	type fields struct {
		timeUtils *timeUtils.MockTimeUtils
	}

	type args struct {
		survey *surveyv1.Survey
	}

	testUUID := uuid.New()
	testINT := 1
	answerID := surveys.AnswerID(testUUID)
	questionID := surveys.QuestionID(testUUID)
	surveyID := surveys.SurveyID(testUUID)

	tests := []struct {
		name string
		args args
		want func(a args, f fields) (*surveys.Survey, error)
	}{
		{
			name: "correct with user",
			args: args{survey: &surveyv1.Survey{
				Id: testUUID.String(),
				Questions: []*questionv1.Question{
					{
						Id:           testUUID.String(),
						SurveyId:     testUUID.String(),
						QuestionType: questionv1.Question_QUESTION_TYPE_TEXT,
						Rules: &questionv1.Question_Rules{
							PickMinCount: &wrapperspb.Int64Value{Value: int64(testINT)},
							PickMaxCount: &wrapperspb.Int64Value{Value: int64(testINT)},
						},
						AnswerVariants: []*answervariantv1.AnswerVariant{
							{
								Id:          testUUID.String(),
								QuestionId:  testUUID.String(),
								WithContent: true,
								ContentType: answervariantv1.AnswerVariant_CONTENT_TYPE_DIGIT,
								Image: &answervariantv1.AnswerVariant_Image{
									Id: "test id",
									ExternalImageInfo: &answervariantv1.AnswerVariant_Image_ExternalInfo{
										Id:   testUUID.String(),
										Size: 1,
									},
								},
								Weight: 10,
								Rules: &answervariantv1.AnswerVariant_Rules{
									ContentMinLength: &wrapperspb.Int64Value{Value: int64(testINT)},
									ContentMaxLength: &wrapperspb.Int64Value{Value: int64(testINT)},
									ContentMinDigit:  &wrapperspb.Int64Value{Value: int64(testINT)},
									ContentMaxDigit:  &wrapperspb.Int64Value{Value: int64(testINT)},
								},
							},
						},
					},
				},
				Respondent: &respondentv1.Respondent{
					Respondent: &respondentv1.Respondent_User_{
						User: &respondentv1.Respondent_User{
							Ids: []string{testUUID.String()},
						},
					},
				},
			}},
			want: func(a args, f fields) (*surveys.Survey, error) {
				f.timeUtils.EXPECT().TimestampToTime(a.survey.GetQuestions()[0].GetDeletedTime()).Return(nil).AnyTimes()
				return &surveys.Survey{
					ID:                &surveyID,
					ActivePeriodStart: a.survey.GetActivePeriodStartTime().AsTime(),
					ActivePeriodEnd:   a.survey.GetActivePeriodEndTime().AsTime(),
					Questions: []*surveys.Question{
						{
							ID:       &questionID,
							SurveyID: &surveyID,
							Type:     surveys.QuestionTypeText,
							Rules: &surveys.QuestionRules{
								PickMaxCount: &testINT,
								PickMinCount: &testINT,
							},
							Answers: []*surveys.AnswerVariant{
								{
									ID:         &answerID,
									QuestionID: &questionID,
									Text:       a.survey.GetQuestions()[0].GetAnswerVariants()[0].GetText(),
									Image: &surveys.Image{
										ID: "test id",
										ExternalImageInfo: &surveys.ExternalProperties{
											ID:   testUUID,
											Size: 1,
										},
									},
									Weight: 10,
									Rules: &surveys.AnswerRules{
										ContentMinLength: &testINT,
										ContentMaxLength: &testINT,
										ContentMinDigit:  &testINT,
										ContentMaxDigit:  &testINT,
									},
									WithContent: true,
									ContentType: surveys.ContentTypeDigit,
								},
							},
						},
					},
					Respondent: &surveys.SurveyRespondent{
						Type: surveys.RespondentTypeUser,
						Ids:  surveys.RespondentIDs{surveys.RespondentID(testUUID)},
					},
				}, nil
			},
		},
		{
			name: "respondent nil",
			args: args{survey: &surveyv1.Survey{
				Id:         testUUID.String(),
				Respondent: nil,
			}},
			want: func(a args, f fields) (*surveys.Survey, error) {
				f.timeUtils.EXPECT().TimestampToTime(nil).Return(nil).AnyTimes()
				return &surveys.Survey{
					ID:                &surveyID,
					ActivePeriodStart: a.survey.GetActivePeriodStartTime().AsTime(),
					ActivePeriodEnd:   a.survey.GetActivePeriodEndTime().AsTime(),
					Questions:         []*surveys.Question{},
					Respondent:        nil,
				}, nil
			},
		},
		{
			name: "correct with anonymous",
			args: args{survey: &surveyv1.Survey{
				Id: testUUID.String(),
				Respondent: &respondentv1.Respondent{
					Respondent: &respondentv1.Respondent_Anonymous_{},
				},
			}},
			want: func(a args, f fields) (*surveys.Survey, error) {
				f.timeUtils.EXPECT().TimestampToTime(nil).Return(nil).AnyTimes()
				return &surveys.Survey{
					ID:                &surveyID,
					ActivePeriodStart: a.survey.GetActivePeriodStartTime().AsTime(),
					ActivePeriodEnd:   a.survey.GetActivePeriodEndTime().AsTime(),
					Questions:         []*surveys.Question{},
					Respondent: &surveys.SurveyRespondent{
						Type: surveys.RespondentTypeAnonymous,
						Ids:  nil,
					},
				}, nil
			},
		},
		{
			name: "incorrect uuid user",
			args: args{survey: &surveyv1.Survey{
				Id: testUUID.String(),
				Respondent: &respondentv1.Respondent{
					Respondent: &respondentv1.Respondent_User_{
						User: &respondentv1.Respondent_User{
							Ids: []string{"test uuid"},
						},
					},
				},
			}},
			want: func(a args, f fields) (*surveys.Survey, error) {
				f.timeUtils.EXPECT().TimestampToTime(nil).Return(nil).AnyTimes()
				return nil, fmt.Errorf("can't convert respondent to entity: can't parse uuid: invalid UUID length: 9")
			},
		},
		{
			name: "incorrect uuid survey",
			args: args{survey: &surveyv1.Survey{
				Id: "test uuid",
			}},
			want: func(a args, f fields) (*surveys.Survey, error) {
				return nil, fmt.Errorf("can't parse survey id: invalid UUID length: 9")
			},
		},
		{
			name: "incorrect uuid question",
			args: args{survey: &surveyv1.Survey{
				Id: testUUID.String(),
				Questions: []*questionv1.Question{
					{
						Id: "test uuid",
					},
				},
			}},
			want: func(a args, f fields) (*surveys.Survey, error) {
				f.timeUtils.EXPECT().TimestampToTime(nil).Return(nil).AnyTimes()
				return nil, fmt.Errorf("can't convert questions to entity: can't parse question id: invalid UUID length: 9")
			},
		},
		{
			name: "incorrect uuid variant",
			args: args{survey: &surveyv1.Survey{
				Id: testUUID.String(),
				Questions: []*questionv1.Question{
					{
						Id:       testUUID.String(),
						SurveyId: testUUID.String(),
						AnswerVariants: []*answervariantv1.AnswerVariant{
							{
								Id: "test uuid",
							},
						},
					},
				},
			}},
			want: func(a args, f fields) (*surveys.Survey, error) {
				f.timeUtils.EXPECT().TimestampToTime(nil).Return(nil).AnyTimes()
				return nil, fmt.Errorf("can't convert questions to entity: can't convert answers variants to entity: can't parse variant id: invalid UUID length: 9")
			},
		},
		{
			name: "incorrect uuid image",
			args: args{survey: &surveyv1.Survey{
				Id: testUUID.String(),
				Questions: []*questionv1.Question{
					{
						Id:       testUUID.String(),
						SurveyId: testUUID.String(),
						AnswerVariants: []*answervariantv1.AnswerVariant{
							{
								Id:         testUUID.String(),
								QuestionId: testUUID.String(),
								Image: &answervariantv1.AnswerVariant_Image{
									Id: "some id",
									ExternalImageInfo: &answervariantv1.AnswerVariant_Image_ExternalInfo{
										Id: "test uuid",
									},
								},
							},
						},
					},
				},
			}},
			want: func(a args, f fields) (*surveys.Survey, error) {
				f.timeUtils.EXPECT().TimestampToTime(nil).Return(nil).AnyTimes()
				return nil, fmt.Errorf("can't convert questions to entity: can't convert answers variants to entity: can't parse image id: invalid UUID length: 9")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				timeUtils: timeUtils.NewMockTimeUtils(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			s := NewSurveyMapper(f.timeUtils)
			got, err := s.SurveyToEntity(tt.args.survey)
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

func TestSurveyMapper_SurveysToEntity(t *testing.T) {
	type fields struct {
		timeUtils *timeUtils.MockTimeUtils
	}

	type args struct {
		surveys []*surveyv1.Survey
	}

	testUUID := uuid.New()
	testMinLength := 1
	answerID := surveys.AnswerID(testUUID)
	questionID := surveys.QuestionID(testUUID)
	surveyID := surveys.SurveyID(testUUID)

	tests := []struct {
		name string
		args args
		want func(a args, f fields) ([]*surveys.Survey, error)
	}{
		{
			name: "correct",
			args: args{surveys: []*surveyv1.Survey{{
				Id: testUUID.String(),
				Questions: []*questionv1.Question{
					{
						Id:           testUUID.String(),
						SurveyId:     testUUID.String(),
						QuestionType: questionv1.Question_QUESTION_TYPE_TEXT,
						AnswerVariants: []*answervariantv1.AnswerVariant{
							{
								Id:          testUUID.String(),
								QuestionId:  testUUID.String(),
								WithContent: true,
								ContentType: answervariantv1.AnswerVariant_CONTENT_TYPE_DIGIT,
								Weight:      10,
								Image: &answervariantv1.AnswerVariant_Image{
									Id: "test id",
									ExternalImageInfo: &answervariantv1.AnswerVariant_Image_ExternalInfo{
										Id:   testUUID.String(),
										Size: 1,
									},
								},
								Rules: &answervariantv1.AnswerVariant_Rules{
									ContentMinLength: &wrapperspb.Int64Value{Value: int64(testMinLength)},
								},
							},
						},
					},
				},
				Respondent: &respondentv1.Respondent{
					Respondent: &respondentv1.Respondent_User_{
						User: &respondentv1.Respondent_User{
							Ids: []string{testUUID.String()},
						},
					},
				},
			}}},
			want: func(a args, f fields) ([]*surveys.Survey, error) {
				f.timeUtils.EXPECT().TimestampToTime(a.surveys[0].GetActivePeriodStartTime()).Return(nil)
				f.timeUtils.EXPECT().TimestampToTime(a.surveys[0].GetActivePeriodEndTime()).Return(nil)
				f.timeUtils.EXPECT().TimestampToTime(a.surveys[0].GetCreatedTime()).Return(nil)
				f.timeUtils.EXPECT().TimestampToTime(a.surveys[0].GetUpdatedTime()).Return(nil)
				return []*surveys.Survey{{
					ID:                &surveyID,
					ActivePeriodStart: a.surveys[0].GetActivePeriodStartTime().AsTime(),
					ActivePeriodEnd:   a.surveys[0].GetActivePeriodEndTime().AsTime(),
					Questions: []*surveys.Question{
						{
							ID:       &questionID,
							SurveyID: &surveyID,
							Type:     surveys.QuestionTypeText,
							Answers: []*surveys.AnswerVariant{
								{
									ID:         &answerID,
									QuestionID: &questionID,
									Weight:     10,
									Text:       a.surveys[0].GetQuestions()[0].GetAnswerVariants()[0].GetText(),
									Image: &surveys.Image{
										ID: "test id",
										ExternalImageInfo: &surveys.ExternalProperties{
											ID:   testUUID,
											Size: 1,
										},
									},
									Rules: &surveys.AnswerRules{
										ContentMinLength: &testMinLength,
									},
									WithContent: true,
									ContentType: surveys.ContentTypeDigit,
								},
							},
						},
					},
					Respondent: &surveys.SurveyRespondent{
						Type: surveys.RespondentTypeUser,
						Ids:  surveys.RespondentIDs{surveys.RespondentID(testUUID)},
					},
				}}, nil
			},
		},
		{
			name: "err",
			args: args{surveys: []*surveyv1.Survey{{Id: "test"}}},
			want: func(a args, f fields) ([]*surveys.Survey, error) {
				return nil, fmt.Errorf("can't convert survey to entity: can't parse survey id: invalid UUID length: 4")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				timeUtils: timeUtils.NewMockTimeUtils(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			s := NewSurveyMapper(f.timeUtils)
			got, err := s.SurveysToEntities(tt.args.surveys)
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
