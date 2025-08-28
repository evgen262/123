package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view"
	viewSurveys "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/surveys"
	surveys "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/survey"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/gintest.git"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func Test_surveyHandlers_getSurvey(t *testing.T) {
	type fields struct {
		interactor *MockSurveysSurveysInteractor
		presenter  *MockSurveysPresenter
		logger     *ditzap.MockLogger
	}

	testErr := errors.New("test")
	errValidation := diterrors.NewValidationError(testErr)
	testUUID := uuid.New()

	tests := []struct {
		name     string
		testCase func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase
	}{
		{
			name: "correct",
			testCase: func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase {
				surveyOpts := surveys.SurveyFilterOptions{true, true, false, false}
				surveyEntity := &surveys.Survey{}
				surveyView := &viewSurveys.SurveyInfo{}
				f.interactor.EXPECT().Get(
					gomock.Any(),
					surveys.SurveyID(testUUID),
					surveyOpts,
				).Return(surveyEntity, nil)
				f.presenter.EXPECT().ToShortView(surveyEntity).Return(surveyView)

				return &gintest.HandlerTestCase{
					Request: &gintest.Request{
						Method: http.MethodGet,
						Path:   "/survey/:id",
					},
					Response: gintest.NewResponse(http.StatusOK, nil, nil, nil).JsonBody(surveyView),
					Params: []gintest.Param{
						{
							Key:   "id",
							Value: testUUID.String(),
						},
					},
					HandlerFunc: handler,
				}
			},
		},
		{
			name: "bad req",
			testCase: func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase {
				_, err := uuid.Parse("test")
				f.logger.EXPECT().Debug("can't parse param id into uuid", zap.Error(err))

				return &gintest.HandlerTestCase{
					Request: &gintest.Request{
						Method: http.MethodGet,
						Path:   "/survey/:id",
					},
					Response: gintest.NewResponse(http.StatusBadRequest, nil, nil, nil).JsonBody(view.NewErrorResponse(view.ErrMessageInvalidRequest)),
					Params: []gintest.Param{
						{
							Key:   "id",
							Value: "test",
						},
					},
					HandlerFunc: handler,
				}
			},
		},
		{
			name: "not found",
			testCase: func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase {
				surveyOpts := surveys.SurveyFilterOptions{true, true, false, false}
				f.interactor.EXPECT().Get(
					gomock.Any(),
					surveys.SurveyID(testUUID),
					surveyOpts,
				).Return(nil, diterrors.ErrNotFound)

				return &gintest.HandlerTestCase{
					Request: &gintest.Request{
						Method: http.MethodGet,
						Path:   "/survey/:id",
					},
					Response: gintest.NewResponse(http.StatusNotFound, nil, nil,
						nil).JsonBody(view.NewErrorResponse(view.ErrMessageNotFound)),
					Params: []gintest.Param{
						{
							Key:   "id",
							Value: testUUID.String(),
						},
					},
					HandlerFunc: handler,
				}
			},
		},
		{
			name: "validation err",
			testCase: func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase {
				surveyOpts := surveys.SurveyFilterOptions{true, true, false, false}
				f.interactor.EXPECT().Get(
					gomock.Any(),
					surveys.SurveyID(testUUID),
					surveyOpts,
				).Return(nil, errValidation)

				return &gintest.HandlerTestCase{
					Request: &gintest.Request{
						Method: http.MethodGet,
						Path:   "/survey/:id",
					},
					Response: gintest.NewResponse(http.StatusBadRequest, nil, nil,
						nil).JsonBody(view.NewErrorResponse(errValidation)),
					Params: []gintest.Param{
						{
							Key:   "id",
							Value: testUUID.String(),
						},
					},
					HandlerFunc: handler,
				}
			},
		},
		{
			name: "internal err",
			testCase: func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase {
				surveyOpts := surveys.SurveyFilterOptions{true, true, false, false}
				f.interactor.EXPECT().Get(
					gomock.Any(),
					surveys.SurveyID(testUUID),
					surveyOpts,
				).Return(nil, testErr)
				f.logger.EXPECT().Error("can't get survey", zap.Error(testErr))

				return &gintest.HandlerTestCase{
					Request: &gintest.Request{
						Method: http.MethodGet,
						Path:   "/survey/:id",
					},
					Response: gintest.NewResponse(http.StatusInternalServerError, nil, nil,
						nil).JsonBody(view.NewErrorResponse(testErr)),
					Params: []gintest.Param{
						{
							Key:   "id",
							Value: testUUID.String(),
						},
					},
					HandlerFunc: handler,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			f := fields{
				interactor: NewMockSurveysSurveysInteractor(ctrl),
				presenter:  NewMockSurveysPresenter(ctrl),
				logger:     ditzap.NewMockLogger(ctrl),
			}

			ph := NewSurveysHandlers(
				f.interactor,
				f.presenter,
				nil,
				nil,
				nil,
				nil,
				f.logger,
			)
			ginTest := gintest.NewGinTest()
			ginTest.TestHandler(t, tt.testCase(ph.getSurvey, f))
		})
	}
}

func Test_surveyAnswerHandlers_addAnswers(t *testing.T) {
	type fields struct {
		interactor *MockSurveysAnswersInteractor
		presenter  *MockSurveysAnswersPresenter
		logger     *ditzap.MockLogger
	}

	testErr := errors.New("test")
	errValidation := diterrors.NewValidationError(testErr)

	tests := []struct {
		name     string
		testCase func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase
	}{
		{
			name: "correct",
			testCase: func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase {
				answers := viewSurveys.NewSurveyAnswers{}
				b, _ := json.Marshal(answers)
				data := bytes.NewReader(b)
				answersEntity := []*surveys.RespondentAnswer{}
				f.interactor.EXPECT().Add(gomock.Any(), answersEntity).Return([]uuid.UUID{}, nil)
				f.presenter.EXPECT().ToNewEntities(&answers).Return(answersEntity)
				result := []*viewSurveys.SurveyAnswerInfo{}
				f.presenter.EXPECT().ToShortViews([]uuid.UUID{}).Return(result)

				return &gintest.HandlerTestCase{
					Request: &gintest.Request{
						Method: http.MethodPost,
						Path:   "/survey/answers",
						Body:   data,
					},
					Response:    gintest.NewResponse(http.StatusCreated, nil, nil, nil).JsonBody(result),
					HandlerFunc: handler,
				}
			},
		},
		{
			name: "internal err",
			testCase: func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase {
				answers := viewSurveys.NewSurveyAnswers{}
				b, _ := json.Marshal(answers)
				data := bytes.NewReader(b)
				answersEntity := []*surveys.RespondentAnswer{}
				f.interactor.EXPECT().Add(gomock.Any(), answersEntity).Return(nil, testErr)
				f.presenter.EXPECT().ToNewEntities(&answers).Return(answersEntity)
				f.logger.EXPECT().Error("can't add answers", zap.Error(testErr))

				return &gintest.HandlerTestCase{
					Request: &gintest.Request{
						Method: http.MethodPost,
						Path:   "/survey/answers",
						Body:   data,
					},
					Response: gintest.NewResponse(http.StatusInternalServerError, nil, nil,
						nil).JsonBody(view.NewErrorResponse(testErr)),
					HandlerFunc: handler,
				}
			},
		},
		{
			name: "validation err",
			testCase: func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase {
				answers := viewSurveys.NewSurveyAnswers{}
				b, _ := json.Marshal(answers)
				data := bytes.NewReader(b)
				answersEntity := []*surveys.RespondentAnswer{}
				f.interactor.EXPECT().Add(gomock.Any(), answersEntity).Return(nil, errValidation)
				f.presenter.EXPECT().ToNewEntities(&answers).Return(answersEntity)

				return &gintest.HandlerTestCase{
					Request: &gintest.Request{
						Method: http.MethodPost,
						Path:   "/survey/answers",
						Body:   data,
					},
					Response: gintest.NewResponse(http.StatusBadRequest, nil, nil,
						nil).JsonBody(view.NewErrorResponse(errValidation)),
					HandlerFunc: handler,
				}
			},
		},
		{
			name: "bad req",
			testCase: func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase {
				type wrong struct {
					RespondentID *[]int `json:"respondent_id,omitempty"`
				}
				answers := wrong{RespondentID: &[]int{1, 2}}
				b, _ := json.Marshal(answers)
				data := bytes.NewReader(b)
				err := json.Unmarshal(b, &viewSurveys.NewSurveyAnswers{})
				f.logger.EXPECT().Debug("can't unbind NewSurveyAnswers json", zap.Error(err))

				return &gintest.HandlerTestCase{
					Request: &gintest.Request{
						Method: http.MethodPost,
						Path:   "/survey/answers",
						Body:   data,
					},
					Response: gintest.NewResponse(http.StatusBadRequest, nil, nil,
						nil).JsonBody(view.NewErrorResponse(view.ErrMessageInvalidRequest)),
					HandlerFunc: handler,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			f := fields{
				interactor: NewMockSurveysAnswersInteractor(ctrl),
				presenter:  NewMockSurveysAnswersPresenter(ctrl),
				logger:     ditzap.NewMockLogger(ctrl),
			}

			ph := NewSurveysHandlers(
				nil,
				nil,
				f.interactor,
				f.presenter,
				nil,
				nil,
				f.logger,
			)
			ginTest := gintest.NewGinTest()
			ginTest.TestHandler(t, tt.testCase(ph.addAnswers, f))
		})
	}
}

func Test_surveyImageHandlers_getImage(t *testing.T) {
	type fields struct {
		interactor *MockSurveysImagesInteractor
		presenter  *MockSurveysImagesPresenter
		logger     *ditzap.MockLogger
	}

	testErr := errors.New("test")

	tests := []struct {
		name     string
		testCase func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase
	}{
		{
			name: "correct",
			testCase: func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase {
				imageData := []byte{123, 125}
				f.interactor.EXPECT().Get(gomock.Any(), "id").Return(imageData, nil)

				return &gintest.HandlerTestCase{
					Request: &gintest.Request{
						Method: http.MethodGet,
						Path:   "/admin/survey/images",
					},
					Response: gintest.NewResponse(
						http.StatusOK,
						http.Header{
							"Content-Type": {"application/octet-stream"},
						},
						imageData,
						nil),
					Params: []gintest.Param{
						{
							Key:   "id",
							Value: "id",
						},
					},
					HandlerFunc: handler,
				}
			},
		},
		{
			name: "internal err",
			testCase: func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase {
				f.interactor.EXPECT().Get(gomock.Any(), "id").Return(nil, testErr)
				f.logger.EXPECT().Error("can't get image", zap.Error(testErr))

				return &gintest.HandlerTestCase{
					Request: &gintest.Request{
						Method: http.MethodGet,
						Path:   "/admin/survey/images",
					},
					Response: gintest.NewResponse(
						http.StatusInternalServerError,
						nil,
						nil,
						nil).JsonBody(view.NewErrorResponse(testErr)),
					Params: []gintest.Param{
						{
							Key:   "id",
							Value: "id",
						},
					},
					HandlerFunc: handler,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			f := fields{
				interactor: NewMockSurveysImagesInteractor(ctrl),
				presenter:  NewMockSurveysImagesPresenter(ctrl),
				logger:     ditzap.NewMockLogger(ctrl),
			}

			ph := NewSurveysHandlers(
				nil,
				nil,
				nil,
				nil,
				f.interactor,
				f.presenter,
				f.logger,
			)
			ginTest := gintest.NewGinTest()
			ginTest.TestHandler(t, tt.testCase(ph.getImage, f))
		})
	}
}
