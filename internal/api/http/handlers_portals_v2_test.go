package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/portalsv2"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portalv2"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/gintest.git"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_portalsV2Handlers_getPortals(t *testing.T) {
	type fields struct {
		portalsV2Interactor   *MockPortalsV2Interactor
		portalsV2Presenter    *MockPortalsV2Presenter
		complexesV2Interactor *MockComplexesV2Interactor
		complexesV2Presenter  *MockComplexesV2Presenter
		logger                *ditzap.MockLogger
	}

	tests := []struct {
		name     string
		testCase func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase
	}{
		{
			name: "invalid json body err",
			testCase: func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase {
				invalidJSONBody := []byte(`{"ids": [1, 2, "extra": "field"`)

				f.logger.EXPECT().Debug(gomock.Any(), gomock.Any()).Times(1)

				return &gintest.HandlerTestCase{
					Request: &gintest.Request{
						Method: http.MethodPost,
						Path:   "/oivs",
						Body:   bytes.NewReader(invalidJSONBody),
					},
					Response: gintest.NewResponse(
						http.StatusBadRequest,
						map[string][]string{
							"Content-Type": {"application/json; charset=utf-8"},
						},
						nil,
						nil,
					).JsonBody(view.NewErrorResponse(view.ErrMessageInvalidRequest)),
					HandlerFunc: handler,
				}
			},
		},
		{
			name: "validation error from interactor",
			testCase: func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase {
				validFilterRequest := portalsv2.PortalsFilterRequest{PortalIDs: []int{1, 2}}
				validFilterEntity := portalv2.FilterPortalsFilters{IDs: []int{1, 2}}
				testValidationError := diterrors.NewValidationError(assert.AnError)

				f.portalsV2Presenter.EXPECT().PortalsFilterToEntity(&validFilterRequest).Return(&validFilterEntity).Times(1)

				f.portalsV2Interactor.EXPECT().Filter(gomock.Any(), &validFilterEntity, &portalv2.FilterPortalsOptions{WithEmployeesCount: true}).Return(nil, testValidationError).Times(1)

				body, _ := json.Marshal(validFilterRequest)

				return &gintest.HandlerTestCase{
					Request: &gintest.Request{
						Method: http.MethodPost,
						Path:   "/oivs",
						Body:   bytes.NewReader(body),
					},
					Response: gintest.NewResponse(
						http.StatusBadRequest,
						map[string][]string{
							"Content-Type": {"application/json; charset=utf-8"},
						},
						nil,
						nil,
					).JsonBody(view.NewErrorResponse(view.ErrMessageInvalidRequest)),
					HandlerFunc: handler,
				}
			},
		},
		{
			name: "internal error from interactor",
			testCase: func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase {
				validFilterRequest := portalsv2.PortalsFilterRequest{PortalIDs: []int{1, 2}}
				validFilterEntity := portalv2.FilterPortalsFilters{IDs: []int{1, 2}}
				testInternalError := assert.AnError

				f.portalsV2Presenter.EXPECT().PortalsFilterToEntity(&validFilterRequest).Return(&validFilterEntity).Times(1)

				f.portalsV2Interactor.EXPECT().Filter(gomock.Any(), &validFilterEntity, &portalv2.FilterPortalsOptions{WithEmployeesCount: true}).Return(nil, testInternalError).Times(1)

				body, _ := json.Marshal(validFilterRequest)

				return &gintest.HandlerTestCase{
					Request: &gintest.Request{
						Method: http.MethodPost,
						Path:   "/oivs",
						Body:   bytes.NewReader(body),
					},
					Response: gintest.NewResponse(
						http.StatusInternalServerError,
						map[string][]string{
							"Content-Type": {"application/json; charset=utf-8"},
						},
						nil,
						nil,
					).JsonBody(view.NewErrorResponse(testInternalError)),
					HandlerFunc: handler,
				}
			},
		},
		{
			name: "success case",
			testCase: func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase {
				validFilterRequest := portalsv2.PortalsFilterRequest{PortalIDs: []int{1, 2}}
				validFilterEntity := portalv2.FilterPortalsFilters{IDs: []int{1, 2}}
				interactorResult := []*portalv2.PortalWithCounts{
					&portalv2.PortalWithCounts{Portal: &portalv2.Portal{ID: 1, Name: "OIV 1"}, EmployeesCount: 100},
					&portalv2.PortalWithCounts{Portal: &portalv2.Portal{ID: 2, Name: "OIV 2"}, EmployeesCount: 200},
				}
				presenterResult := []*portalsv2.Portal{
					&portalsv2.Portal{ID: 1, Name: "OIV 1", Count: portalsv2.Count{Employees: 100}},
					&portalsv2.Portal{ID: 2, Name: "OIV 2", Count: portalsv2.Count{Employees: 200}},
				}

				f.portalsV2Presenter.EXPECT().PortalsFilterToEntity(&validFilterRequest).Return(&validFilterEntity).Times(1)

				f.portalsV2Interactor.EXPECT().Filter(gomock.Any(), &validFilterEntity, &portalv2.FilterPortalsOptions{WithEmployeesCount: true}).Return(interactorResult, nil).Times(1)

				f.portalsV2Presenter.EXPECT().PortalsWithCountToView(interactorResult).Return(presenterResult).Times(1)

				body, _ := json.Marshal(validFilterRequest)

				return &gintest.HandlerTestCase{
					Request: &gintest.Request{
						Method: http.MethodPost,
						Path:   "/oivs",
						Body:   bytes.NewReader(body),
					},
					Response: gintest.NewResponse(
						http.StatusOK,
						map[string][]string{
							"Content-Type": {"application/json; charset=utf-8"},
						},
						nil,
						nil,
					).JsonBody(presenterResult),
					HandlerFunc: handler,
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			f := fields{
				portalsV2Interactor:   NewMockPortalsV2Interactor(ctrl),
				portalsV2Presenter:    NewMockPortalsV2Presenter(ctrl),
				complexesV2Interactor: NewMockComplexesV2Interactor(ctrl),
				complexesV2Presenter:  NewMockComplexesV2Presenter(ctrl),
				logger:                ditzap.NewMockLogger(ctrl),
			}

			ph := NewPortalsV2Handlers(
				f.portalsV2Interactor,
				f.portalsV2Presenter,
				f.complexesV2Interactor,
				f.complexesV2Presenter,
				f.logger,
			)

			ginTest := gintest.NewGinTest()
			ginTest.TestHandler(t, tt.testCase(ph.getPortals, f))
		})
	}
}

func Test_portalsV2Handlers_getComplexes(t *testing.T) {
	type fields struct {
		portalsV2Interactor   *MockPortalsV2Interactor
		portalsV2Presenter    *MockPortalsV2Presenter
		complexesV2Interactor *MockComplexesV2Interactor
		complexesV2Presenter  *MockComplexesV2Presenter
		logger                *ditzap.MockLogger
	}

	tests := []struct {
		name     string
		testCase func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase
	}{
		{
			name: "validation error from interactor",
			testCase: func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase {
				testValidationError := diterrors.NewValidationError(assert.AnError)

				f.complexesV2Interactor.EXPECT().Filter(gomock.Any(), nil, nil).Return(nil, testValidationError).Times(1)

				return &gintest.HandlerTestCase{
					Request: &gintest.Request{
						Method: http.MethodPost,
						Path:   "/complexes",
						Body:   nil,
					},
					Response: gintest.NewResponse(
						http.StatusBadRequest,
						map[string][]string{
							"Content-Type": {"application/json; charset=utf-8"},
						},
						nil,
						nil,
					).JsonBody(view.NewErrorResponse(testValidationError.Unwrap())),
					HandlerFunc: handler,
				}
			},
		},
		{
			name: "internal error from interactor",
			testCase: func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase {
				testInternalError := assert.AnError

				f.complexesV2Interactor.EXPECT().Filter(gomock.Any(), nil, nil).Return(nil, testInternalError).Times(1)

				return &gintest.HandlerTestCase{
					Request: &gintest.Request{
						Method: http.MethodPost,
						Path:   "/complexes",
						Body:   nil,
					},
					Response: gintest.NewResponse(
						http.StatusInternalServerError,
						map[string][]string{
							"Content-Type": {"application/json; charset=utf-8"},
						},
						nil,
						nil,
					).JsonBody(view.NewErrorResponse(testInternalError)),
					HandlerFunc: handler,
				}
			},
		},
		{
			name: "success case",
			testCase: func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase {
				interactorResult := []*portalv2.Complex{
					{ID: 1},
					{ID: 2},
				}
				presenterResult := []*portalsv2.Complex{
					{ID: 1},
					{ID: 2},
				}

				f.complexesV2Interactor.EXPECT().Filter(gomock.Any(), nil, nil).Return(interactorResult, nil).Times(1)

				f.complexesV2Presenter.EXPECT().ComplexesToView(interactorResult).Return(presenterResult).Times(1)

				return &gintest.HandlerTestCase{
					Request: &gintest.Request{
						Method: http.MethodPost,
						Path:   "/complexes",
						Body:   nil,
					},
					Response: gintest.NewResponse(
						http.StatusOK,
						map[string][]string{
							"Content-Type": {"application/json; charset=utf-8"},
						},
						nil,
						nil,
					).JsonBody(presenterResult),
					HandlerFunc: handler,
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				portalsV2Interactor:   NewMockPortalsV2Interactor(ctrl),
				portalsV2Presenter:    NewMockPortalsV2Presenter(ctrl),
				complexesV2Interactor: NewMockComplexesV2Interactor(ctrl),
				complexesV2Presenter:  NewMockComplexesV2Presenter(ctrl),
				logger:                ditzap.NewMockLogger(ctrl),
			}

			ph := NewPortalsV2Handlers(
				f.portalsV2Interactor,
				f.portalsV2Presenter,
				f.complexesV2Interactor,
				f.complexesV2Presenter,
				f.logger,
			)

			ginTest := gintest.NewGinTest()
			ginTest.TestHandler(t, tt.testCase(ph.getComplexes, f))
		})
	}
}
