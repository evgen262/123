package surveys

import (
	"testing"
	"time"

	viewSurveys "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/surveys"
	entitySurvey "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/survey"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSurveysPresenter_ToNewEntity(t *testing.T) {
	type fields struct {
		logger *ditzap.MockLogger
	}

	type args struct {
		survey *viewSurveys.NewSurvey
	}

	testUUID := uuid.New()
	testINT := 1
	testINT1 := 15

	tests := []struct {
		name string
		args args
		want func(a args, f fields) *entitySurvey.Survey
	}{
		{
			name: "correct 1",
			args: args{survey: &viewSurveys.NewSurvey{
				Title:          "test",
				RespondentType: 2,
				Respondents:    []*viewSurveys.SurveyRespondent{{ID: testUUID}},
				Questions: []*viewSurveys.NewSurveyQuestion{
					{
						Text: "test",
						Type: "text",
						Rules: &viewSurveys.SurveyQuestionRules{
							PickMinCount: &testINT,
							PickMaxCount: &testINT,
						},
						AnswerVariants: []*viewSurveys.NewSurveyAnswerVariant{
							{
								Text:        "test",
								WithContent: true,
								ContentType: &testINT,
								Image: &viewSurveys.SurveyImage{
									ID: "test id",
								},
								Rules: &viewSurveys.SurveyAnswerVariantRules{
									ContentMinLength: &testINT,
									ContentMaxLength: &testINT,
									ContentMinDigit:  &testINT,
									ContentMaxDigit:  &testINT,
								},
							},
						},
					},
				},
			}},
			want: func(a args, f fields) *entitySurvey.Survey {
				return &entitySurvey.Survey{
					Title: "test",
					Questions: []*entitySurvey.Question{
						{
							Text: "test",
							Type: entitySurvey.QuestionTypeText,
							Rules: &entitySurvey.QuestionRules{
								PickMinCount: &testINT,
								PickMaxCount: &testINT,
							},
							Answers: []*entitySurvey.AnswerVariant{
								{
									Text:        "test",
									WithContent: true,
									ContentType: entitySurvey.ContentTypeText,
									Image: &entitySurvey.Image{
										ID:                "test id",
										ExternalImageInfo: &entitySurvey.ExternalProperties{},
									},
									Rules: &entitySurvey.AnswerRules{
										ContentMinLength: &testINT,
										ContentMaxLength: &testINT,
										ContentMinDigit:  &testINT,
										ContentMaxDigit:  &testINT,
									},
								},
							},
						},
					},
					Respondent: &entitySurvey.SurveyRespondent{
						Type: entitySurvey.RespondentTypeUser,
						Ids:  entitySurvey.RespondentIDs{entitySurvey.RespondentID(testUUID)},
					},
				}
			},
		},
		{
			name: "correct 2",
			args: args{survey: &viewSurveys.NewSurvey{
				Title:          "test",
				RespondentType: 5,
				Questions: []*viewSurveys.NewSurveyQuestion{
					{
						Text: "test",
						Type: "test",
						AnswerVariants: []*viewSurveys.NewSurveyAnswerVariant{
							{
								Text:        "test",
								WithContent: true,
								ContentType: &testINT1,
								Image: &viewSurveys.SurveyImage{
									ID: "test id",
								},
								Rules: &viewSurveys.SurveyAnswerVariantRules{
									ContentMinLength: &testINT,
									ContentMaxLength: &testINT,
									ContentMinDigit:  &testINT,
									ContentMaxDigit:  &testINT,
								},
							},
						},
					},
				},
			}},
			want: func(a args, f fields) *entitySurvey.Survey {
				f.logger.EXPECT().Warn("unknown respondent type, switched to default type")
				return &entitySurvey.Survey{
					Title: "test",
					Questions: []*entitySurvey.Question{
						{
							Text: "test",
							Type: entitySurvey.QuestionTypeInvalid,
							Answers: []*entitySurvey.AnswerVariant{
								{
									Text:        "test",
									WithContent: true,
									ContentType: entitySurvey.ContentTypeInvalid,
									Image: &entitySurvey.Image{
										ID:                "test id",
										ExternalImageInfo: &entitySurvey.ExternalProperties{},
									},
									Rules: &entitySurvey.AnswerRules{
										ContentMinLength: &testINT,
										ContentMaxLength: &testINT,
										ContentMinDigit:  &testINT,
										ContentMaxDigit:  &testINT,
									},
								},
							},
						},
					},
					Respondent: &entitySurvey.SurveyRespondent{
						Type: entitySurvey.RespondentTypeAll,
					},
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				logger: ditzap.NewMockLogger(ctrl),
			}
			want := tt.want(tt.args, f)
			s := NewSurveysPresenter(f.logger)
			got := s.ToNewEntity(tt.args.survey)
			assert.Equal(t, want, got)
		})
	}
}

func TestSurveysPresenter_ToEntity(t *testing.T) {
	type fields struct {
		logger *ditzap.MockLogger
	}

	type args struct {
		survey *viewSurveys.Survey
	}

	testUUID := uuid.New()
	testINT := 1
	testTime := time.Now()
	answerID := entitySurvey.AnswerID(testUUID)
	questionID := entitySurvey.QuestionID(testUUID)
	surveyID := entitySurvey.SurveyID(testUUID)

	tests := []struct {
		name string
		args args
		want func(a args, f fields) *entitySurvey.Survey
	}{
		{
			name: "correct 1",
			args: args{survey: &viewSurveys.Survey{
				ID:             testUUID,
				Title:          "test",
				RespondentType: 2,
				Respondents:    []*viewSurveys.SurveyRespondent{{ID: testUUID}},
				CreatedAt:      testTime,
				UpdatedAt:      &testTime,
				DeletedAt:      &testTime,
				Questions: []*viewSurveys.SurveyQuestion{
					{
						ID:        testUUID,
						SurveyID:  testUUID,
						Text:      "test",
						Type:      "text",
						DeletedAt: &testTime,
						Rules: &viewSurveys.SurveyQuestionRules{
							PickMinCount: &testINT,
							PickMaxCount: &testINT,
						},
						AnswerVariants: []*viewSurveys.SurveyAnswerVariant{
							{
								ID:          testUUID,
								QuestionID:  testUUID,
								Text:        "test",
								WithContent: true,
								ContentType: &testINT,
								DeletedAt:   &testTime,
								Image: &viewSurveys.SurveyImage{
									ID: "test id",
								},
								Rules: &viewSurveys.SurveyAnswerVariantRules{
									ContentMinLength: &testINT,
									ContentMaxLength: &testINT,
									ContentMinDigit:  &testINT,
									ContentMaxDigit:  &testINT,
								},
							},
						},
					},
				},
			}},
			want: func(a args, f fields) *entitySurvey.Survey {
				return &entitySurvey.Survey{
					ID:        &surveyID,
					Title:     "test",
					CreatedAt: &testTime,
					UpdatedAt: &testTime,
					DeletedAt: &testTime,
					Questions: []*entitySurvey.Question{
						{
							ID:        &questionID,
							SurveyID:  &surveyID,
							Text:      "test",
							Type:      entitySurvey.QuestionTypeText,
							DeletedAt: &testTime,
							Rules: &entitySurvey.QuestionRules{
								PickMinCount: &testINT,
								PickMaxCount: &testINT,
							},
							Answers: []*entitySurvey.AnswerVariant{
								{
									ID:          &answerID,
									QuestionID:  &questionID,
									Text:        "test",
									WithContent: true,
									ContentType: entitySurvey.ContentTypeText,
									DeletedAt:   &testTime,
									Image: &entitySurvey.Image{
										ID:                "test id",
										ExternalImageInfo: &entitySurvey.ExternalProperties{},
									},
									Rules: &entitySurvey.AnswerRules{
										ContentMinLength: &testINT,
										ContentMaxLength: &testINT,
										ContentMinDigit:  &testINT,
										ContentMaxDigit:  &testINT,
									},
								},
							},
						},
					},
					Respondent: &entitySurvey.SurveyRespondent{
						Type: entitySurvey.RespondentTypeUser,
						Ids:  entitySurvey.RespondentIDs{entitySurvey.RespondentID(testUUID)},
					},
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				logger: ditzap.NewMockLogger(ctrl),
			}
			want := tt.want(tt.args, f)
			s := NewSurveysPresenter(f.logger)
			got := s.ToEntity(tt.args.survey)
			assert.Equal(t, want, got)
		})
	}
}

func TestSurveysPresenter_ToView(t *testing.T) {
	type fields struct {
		logger *ditzap.MockLogger
	}

	type args struct {
		survey *entitySurvey.Survey
	}

	testUUID := uuid.New()
	testINT := 1
	testTime := time.Now()
	answerID := entitySurvey.AnswerID(testUUID)
	questionID := entitySurvey.QuestionID(testUUID)
	surveyID := entitySurvey.SurveyID(testUUID)

	tests := []struct {
		name string
		args args
		want func(a args, f fields) *viewSurveys.Survey
	}{
		{
			name: "correct 1",
			args: args{survey: &entitySurvey.Survey{
				ID:        &surveyID,
				Title:     "test",
				CreatedAt: &testTime,
				UpdatedAt: &testTime,
				DeletedAt: &testTime,
				Questions: []*entitySurvey.Question{
					{
						ID:       &questionID,
						SurveyID: &surveyID,
						Text:     "test",
						Type:     entitySurvey.QuestionTypeText,
						Rules: &entitySurvey.QuestionRules{
							PickMinCount: &testINT,
							PickMaxCount: &testINT,
						},
						DeletedAt: &testTime,
						Answers: []*entitySurvey.AnswerVariant{
							{
								ID:          &answerID,
								QuestionID:  &questionID,
								Text:        "test",
								WithContent: true,
								ContentType: entitySurvey.ContentTypeText,
								DeletedAt:   &testTime,
								Image: &entitySurvey.Image{
									ID:                "test id",
									ExternalImageInfo: &entitySurvey.ExternalProperties{},
								},
								Rules: &entitySurvey.AnswerRules{
									ContentMinLength: &testINT,
									ContentMaxLength: &testINT,
									ContentMinDigit:  &testINT,
									ContentMaxDigit:  &testINT,
								},
							},
						},
					},
				},
				Respondent: &entitySurvey.SurveyRespondent{
					Type: entitySurvey.RespondentTypeUser,
					Ids:  entitySurvey.RespondentIDs{entitySurvey.RespondentID(testUUID)},
				},
			}},
			want: func(a args, f fields) *viewSurveys.Survey {
				return &viewSurveys.Survey{
					ID:             testUUID,
					Title:          "test",
					RespondentType: 2,
					Respondents:    []*viewSurveys.SurveyRespondent{{ID: testUUID}},
					CreatedAt:      testTime,
					UpdatedAt:      &testTime,
					DeletedAt:      &testTime,
					Questions: []*viewSurveys.SurveyQuestion{
						{
							ID:        testUUID,
							SurveyID:  testUUID,
							Text:      "test",
							Type:      "text",
							DeletedAt: &testTime,
							Rules: &viewSurveys.SurveyQuestionRules{
								PickMinCount: &testINT,
								PickMaxCount: &testINT,
							},
							AnswerVariants: []*viewSurveys.SurveyAnswerVariant{
								{
									ID:          testUUID,
									QuestionID:  testUUID,
									Text:        "test",
									WithContent: true,
									ContentType: &testINT,
									DeletedAt:   &testTime,
									Image: &viewSurveys.SurveyImage{
										ID: "test id",
									},
									Rules: &viewSurveys.SurveyAnswerVariantRules{
										ContentMinLength: &testINT,
										ContentMaxLength: &testINT,
										ContentMinDigit:  &testINT,
										ContentMaxDigit:  &testINT,
									},
								},
							},
						},
					},
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				logger: ditzap.NewMockLogger(ctrl),
			}
			want := tt.want(tt.args, f)
			s := NewSurveysPresenter(f.logger)
			got := s.ToView(tt.args.survey)
			assert.Equal(t, want, got)
		})
	}
}

func TestSurveysPresenter_ToShortView(t *testing.T) {
	type fields struct {
		logger *ditzap.MockLogger
	}

	type args struct {
		survey *entitySurvey.Survey
	}

	testUUID := uuid.New()
	testINT := 1
	testTime := time.Now()
	answerID := entitySurvey.AnswerID(testUUID)
	questionID := entitySurvey.QuestionID(testUUID)
	surveyID := entitySurvey.SurveyID(testUUID)

	tests := []struct {
		name string
		args args
		want func(a args, f fields) *viewSurveys.SurveyInfo
	}{
		{
			name: "correct 1",
			args: args{survey: &entitySurvey.Survey{
				ID:        &surveyID,
				Title:     "test",
				CreatedAt: &testTime,
				UpdatedAt: &testTime,
				DeletedAt: &testTime,
				Questions: []*entitySurvey.Question{
					{
						ID:       &questionID,
						SurveyID: &surveyID,
						Text:     "test",
						Type:     entitySurvey.QuestionTypeText,
						Rules: &entitySurvey.QuestionRules{
							PickMinCount: &testINT,
							PickMaxCount: &testINT,
						},
						DeletedAt: &testTime,
						Answers: []*entitySurvey.AnswerVariant{
							{
								ID:          &answerID,
								QuestionID:  &questionID,
								Text:        "test",
								WithContent: true,
								ContentType: entitySurvey.ContentTypeText,
								DeletedAt:   &testTime,
								Image: &entitySurvey.Image{
									ID:                "test id",
									ExternalImageInfo: &entitySurvey.ExternalProperties{},
								},
								Rules: &entitySurvey.AnswerRules{
									ContentMinLength: &testINT,
									ContentMaxLength: &testINT,
									ContentMinDigit:  &testINT,
									ContentMaxDigit:  &testINT,
								},
							},
						},
					},
				},
				Respondent: &entitySurvey.SurveyRespondent{
					Type: entitySurvey.RespondentTypeUser,
					Ids:  entitySurvey.RespondentIDs{entitySurvey.RespondentID(testUUID)},
				},
			}},
			want: func(a args, f fields) *viewSurveys.SurveyInfo {
				return &viewSurveys.SurveyInfo{
					Title: "test",
					Questions: []*viewSurveys.SurveyQuestionInfo{
						{
							ID:   uuid.UUID(questionID),
							Text: "test",
							Type: "text",
							Rules: &viewSurveys.SurveyQuestionRules{
								PickMinCount: &testINT,
								PickMaxCount: &testINT,
							},
							AnswerVariants: []*viewSurveys.SurveyAnswerVariantInfo{
								{
									ID:          uuid.UUID(answerID),
									Text:        "test",
									WithContent: true,
									ContentType: &testINT,
									Image: &viewSurveys.SurveyImageInfo{
										ID: "test id",
									},
									Rules: &viewSurveys.SurveyAnswerVariantRules{
										ContentMinLength: &testINT,
										ContentMaxLength: &testINT,
										ContentMinDigit:  &testINT,
										ContentMaxDigit:  &testINT,
									},
								},
							},
						},
					},
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				logger: ditzap.NewMockLogger(ctrl),
			}
			want := tt.want(tt.args, f)
			s := NewSurveysPresenter(f.logger)
			got := s.ToShortView(tt.args.survey)
			assert.Equal(t, want, got)
		})
	}
}

func TestSurveysPresenter_IDsToEntities(t *testing.T) {
	type fields struct {
		logger *ditzap.MockLogger
	}

	type args struct {
		ids []uuid.UUID
	}

	testUUID := uuid.New()
	surveyID := entitySurvey.SurveyID(testUUID)

	tests := []struct {
		name string
		args args
		want func(a args, f fields) entitySurvey.SurveyIDs
	}{
		{
			name: "correct",
			args: args{
				ids: []uuid.UUID{testUUID},
			},
			want: func(a args, f fields) entitySurvey.SurveyIDs {
				return entitySurvey.SurveyIDs{surveyID}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				logger: ditzap.NewMockLogger(ctrl),
			}
			want := tt.want(tt.args, f)
			s := NewSurveysPresenter(f.logger)
			got := s.IDsToEntities(tt.args.ids)
			assert.Equal(t, want, got)
		})
	}
}

func TestSurveysPresenter_RespondentToEntity(t *testing.T) {
	type fields struct {
		logger *ditzap.MockLogger
	}

	type args struct {
		respondent *viewSurveys.GetAllSurveysRespondent
	}

	testUUID := uuid.New()
	respID := entitySurvey.RespondentID(testUUID)

	tests := []struct {
		name string
		args args
		want func(a args, f fields) *entitySurvey.SurveyRespondent
	}{
		{
			name: "correct",
			args: args{
				respondent: &viewSurveys.GetAllSurveysRespondent{
					RespondentType: 2,
					RespondentIDs:  []uuid.UUID{testUUID},
				},
			},
			want: func(a args, f fields) *entitySurvey.SurveyRespondent {
				return &entitySurvey.SurveyRespondent{
					Type: entitySurvey.RespondentTypeUser,
					Ids:  entitySurvey.RespondentIDs{respID},
				}
			},
		},
		{
			name: "nil input",
			args: args{
				respondent: nil,
			},
			want: func(a args, f fields) *entitySurvey.SurveyRespondent {
				return &entitySurvey.SurveyRespondent{
					Type: entitySurvey.RespondentTypeAll,
					Ids:  nil,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				logger: ditzap.NewMockLogger(ctrl),
			}
			want := tt.want(tt.args, f)
			s := NewSurveysPresenter(f.logger)
			got := s.RespondentToEntity(tt.args.respondent)
			assert.Equal(t, want, got)
		})
	}
}

func TestSurveysPresenter_OptionsToEntity(t *testing.T) {
	type fields struct {
		logger *ditzap.MockLogger
	}

	type args struct {
		options *viewSurveys.SurveysOptions
	}

	tests := []struct {
		name string
		args args
		want func(a args, f fields) entitySurvey.SurveyFilterOptions
	}{
		{
			name: "correct",
			args: args{
				options: &viewSurveys.SurveysOptions{},
			},
			want: func(a args, f fields) entitySurvey.SurveyFilterOptions {
				return entitySurvey.SurveyFilterOptions{}
			},
		},
		{
			name: "nil",
			args: args{
				options: &viewSurveys.SurveysOptions{},
			},
			want: func(a args, f fields) entitySurvey.SurveyFilterOptions {
				return entitySurvey.SurveyFilterOptions{}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				logger: ditzap.NewMockLogger(ctrl),
			}
			want := tt.want(tt.args, f)
			s := NewSurveysPresenter(f.logger)
			got := s.OptionsToEntity(tt.args.options)
			assert.Equal(t, want, got)
		})
	}
}

func TestSurveysPresenter_PaginationToEntity(t *testing.T) {
	type fields struct {
		logger *ditzap.MockLogger
	}

	type args struct {
		pagination *viewSurveys.GetAllSurveysPagination
	}

	testTime := time.Now()
	testUUID := uuid.New()
	sID := entitySurvey.SurveyID(testUUID)
	testInt := 5

	tests := []struct {
		name string
		args args
		want func(a args, f fields) entitySurvey.Pagination
	}{
		{
			name: "correct",
			args: args{
				pagination: &viewSurveys.GetAllSurveysPagination{
					Limit:    testInt,
					LastId:   &testUUID,
					LastDate: &testTime,
					Total:    &testInt,
				},
			},
			want: func(a args, f fields) entitySurvey.Pagination {
				return entitySurvey.Pagination{
					Limit:    uint32(testInt),
					LastId:   &sID,
					LastDate: &testTime,
					Total:    int64(testInt),
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				logger: ditzap.NewMockLogger(ctrl),
			}
			want := tt.want(tt.args, f)
			s := NewSurveysPresenter(f.logger)
			got := s.PaginationToEntity(tt.args.pagination)
			assert.Equal(t, want, got)
		})
	}
}

func TestSurveysPresenter_SurveysWithPaginationToView(t *testing.T) {
	type fields struct {
		logger *ditzap.MockLogger
	}

	type args struct {
		surveysWithPagination *entitySurvey.SurveysWithPagination
	}

	testUUID := uuid.New()
	testINT := 1
	testTime := time.Now()
	answerID := entitySurvey.AnswerID(testUUID)
	questionID := entitySurvey.QuestionID(testUUID)
	surveyID := entitySurvey.SurveyID(testUUID)

	tests := []struct {
		name string
		args args
		want func(a args, f fields) *viewSurveys.SurveysWithPagination
	}{
		{
			name: "correct",
			args: args{
				surveysWithPagination: &entitySurvey.SurveysWithPagination{
					Surveys: []*entitySurvey.Survey{{
						ID:        &surveyID,
						Title:     "test",
						CreatedAt: &testTime,
						UpdatedAt: &testTime,
						DeletedAt: &testTime,
						Questions: []*entitySurvey.Question{
							{
								ID:        &questionID,
								SurveyID:  &surveyID,
								Text:      "test",
								Type:      entitySurvey.QuestionTypeText,
								DeletedAt: &testTime,
								Answers: []*entitySurvey.AnswerVariant{
									{
										ID:          &answerID,
										QuestionID:  &questionID,
										Text:        "test",
										WithContent: true,
										ContentType: entitySurvey.ContentTypeText,
										DeletedAt:   &testTime,
										Image: &entitySurvey.Image{
											ID:                "test id",
											ExternalImageInfo: &entitySurvey.ExternalProperties{},
										},
										Rules: &entitySurvey.AnswerRules{
											ContentMinLength: &testINT,
											ContentMaxLength: &testINT,
											ContentMinDigit:  &testINT,
											ContentMaxDigit:  &testINT,
										},
									},
								},
							},
						},
						Respondent: &entitySurvey.SurveyRespondent{
							Type: entitySurvey.RespondentTypeUser,
							Ids:  entitySurvey.RespondentIDs{entitySurvey.RespondentID(testUUID)},
						},
					}},
					Pagination: entitySurvey.Pagination{
						LastId: &surveyID,
					},
				},
			},
			want: func(a args, f fields) *viewSurveys.SurveysWithPagination {
				return &viewSurveys.SurveysWithPagination{
					Pagination: viewSurveys.SurveyPagination{
						LastId: &testUUID,
					},
					Surveys: []*viewSurveys.Survey{
						{ID: testUUID,
							Title:          "test",
							RespondentType: 2,
							Respondents:    []*viewSurveys.SurveyRespondent{{ID: testUUID}},
							CreatedAt:      testTime,
							UpdatedAt:      &testTime,
							DeletedAt:      &testTime,
							Questions: []*viewSurveys.SurveyQuestion{
								{
									ID:        testUUID,
									SurveyID:  testUUID,
									Text:      "test",
									Type:      "text",
									DeletedAt: &testTime,
									AnswerVariants: []*viewSurveys.SurveyAnswerVariant{
										{
											ID:          testUUID,
											QuestionID:  testUUID,
											Text:        "test",
											WithContent: true,
											ContentType: &testINT,
											DeletedAt:   &testTime,
											Image: &viewSurveys.SurveyImage{
												ID: "test id",
											},
											Rules: &viewSurveys.SurveyAnswerVariantRules{
												ContentMinLength: &testINT,
												ContentMaxLength: &testINT,
												ContentMinDigit:  &testINT,
												ContentMaxDigit:  &testINT,
											},
										},
									},
								},
							},
						},
					},
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				logger: ditzap.NewMockLogger(ctrl),
			}
			want := tt.want(tt.args, f)
			s := NewSurveysPresenter(f.logger)
			got := s.SurveysWithPaginationToView(tt.args.surveysWithPagination)
			assert.Equal(t, want, got)
		})
	}
}

func TestSurveysPresenter_IDToView(t *testing.T) {
	type fields struct {
		logger *ditzap.MockLogger
	}

	type args struct {
		ID *entitySurvey.SurveyID
	}

	testUUID := uuid.New()
	sID := entitySurvey.SurveyID(testUUID)

	tests := []struct {
		name string
		args args
		want func(a args, f fields) *viewSurveys.IDResponse
	}{
		{
			name: "correct",
			args: args{
				ID: &sID,
			},
			want: func(a args, f fields) *viewSurveys.IDResponse {
				return &viewSurveys.IDResponse{ID: testUUID.String()}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				logger: ditzap.NewMockLogger(ctrl),
			}
			want := tt.want(tt.args, f)
			s := NewSurveysPresenter(f.logger)
			got := s.IDToView(tt.args.ID)
			assert.Equal(t, want, got)
		})
	}
}
