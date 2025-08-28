package http

import (
	"bytes"
	"errors"
	"net/http"
	"testing"

	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/gintest.git"
	"github.com/gin-gonic/gin"
	"go.uber.org/mock/gomock"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view"
	viewAnalytics "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/analytics"
	entityAnalytics "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/analytics"
	interactorAnalytics "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/usecase/analytics"
)

func Test_analyticsHandlers_addMetrics(t *testing.T) {
	type fields struct {
		interactor *MockAnalyticsInteractor
	}

	tests := []struct {
		name     string
		testCase func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase
	}{
		{
			name: "correct",
			testCase: func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase {
				metricID := ""
				userAgent := "DeviceID=f78f5813-4e7a-4093-99d0-8c434fa38e22;DeviceType=web"
				requestBody := []byte(`{"key":"value"}`)

				f.interactor.EXPECT().AddMetrics(gomock.Any(), entityAnalytics.XCFCUserAgentHeader(userAgent), requestBody).Return(metricID, nil)

				return &gintest.HandlerTestCase{
					Request: &gintest.Request{
						Method:  http.MethodPost,
						Path:    "/analytics/metrics",
						Headers: map[string][]string{"X-CFC-Useragent": {userAgent}},
						Body:    bytes.NewBuffer(requestBody),
					},
					Response: gintest.NewResponse(
						http.StatusOK,
						map[string][]string{
							"Content-Type": {"application/json; charset=utf-8"},
						},
						nil,
						nil,
					).JsonBody(view.NewSuccessResponse(viewAnalytics.AddMetricsResponse{
						MetricID: metricID,
					})),
					HandlerFunc: handler,
				}
			},
		},
		{
			name: "invalid user agent",
			testCase: func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase {
				userAgent := "invalid user agent"
				requestBody := []byte(`{"key":"value"}`)

				return &gintest.HandlerTestCase{
					Request: &gintest.Request{
						Method:  http.MethodPost,
						Path:    "/analytics/metrics",
						Headers: map[string][]string{"X-CFC-Useragent": {userAgent}},
						Body:    bytes.NewBuffer(requestBody),
					},
					Response: gintest.NewResponse(
						http.StatusBadRequest,
						map[string][]string{
							"Content-Type":   {"application/json; charset=utf-8"},
							StatusCodeHeader: {"AM_01"},
						},
						nil,
						nil,
					).JsonBody(view.NewErrorResponse(view.ErrMessageInvalidRequest)),
					HandlerFunc: handler,
				}
			},
		},
		{
			name: "read body error",
			testCase: func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase {
				userAgent := "DeviceID=f78f5813-4e7a-4093-99d0-8c434fa38e22;DeviceType=web"
				// Create a reader that returns an error when reading.
				er := &errorReader{}

				return &gintest.HandlerTestCase{
					Request: &gintest.Request{
						Method:  http.MethodPost,
						Path:    "/analytics/metrics",
						Headers: map[string][]string{"X-CFC-Useragent": {userAgent}},
						Body:    er,
					},
					Response: gintest.NewResponse(
						http.StatusBadRequest,
						map[string][]string{
							"Content-Type":   {"application/json; charset=utf-8"},
							StatusCodeHeader: {"AM_02"},
						},
						nil,
						nil,
					).JsonBody(view.NewErrorResponse(view.ErrMessageInvalidRequest)),
					HandlerFunc: handler,
				}
			},
		},
		{
			name: "unauthenticated error",
			testCase: func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase {
				userAgent := "DeviceID=f78f5813-4e7a-4093-99d0-8c434fa38e22;DeviceType=web"
				requestBody := []byte(`{"key":"value"}`)

				f.interactor.EXPECT().AddMetrics(gomock.Any(), entityAnalytics.XCFCUserAgentHeader(userAgent), requestBody).Return("", interactorAnalytics.ErrUnauthenticated)

				return &gintest.HandlerTestCase{
					Request: &gintest.Request{
						Method:  http.MethodPost,
						Path:    "/analytics/metrics",
						Headers: map[string][]string{"X-CFC-Useragent": {userAgent}},
						Body:    bytes.NewBuffer(requestBody),
					},
					Response: gintest.NewResponse(
						http.StatusUnauthorized,
						map[string][]string{
							"Content-Type":   {"application/json; charset=utf-8"},
							StatusCodeHeader: {"AM_03"},
						},
						nil,
						nil,
					).JsonBody(view.NewErrorResponse(view.ErrMessageUnauthenticated)),
					HandlerFunc: handler,
				}
			},
		},
		{
			name: "validation error",
			testCase: func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase {
				userAgent := "DeviceID=f78f5813-4e7a-4093-99d0-8c434fa38e22;DeviceType=web"
				requestBody := []byte(`{"key":"value"}`)
				validationError := diterrors.NewValidationError(errors.New("validation error"))

				f.interactor.EXPECT().AddMetrics(gomock.Any(), entityAnalytics.XCFCUserAgentHeader(userAgent), requestBody).Return("", validationError)

				return &gintest.HandlerTestCase{
					Request: &gintest.Request{
						Method:  http.MethodPost,
						Path:    "/analytics/metrics",
						Headers: map[string][]string{"X-CFC-Useragent": {userAgent}},
						Body:    bytes.NewBuffer(requestBody),
					},
					Response: gintest.NewResponse(
						http.StatusBadRequest,
						map[string][]string{
							"Content-Type":   {"application/json; charset=utf-8"},
							StatusCodeHeader: {"AM_04"},
						},
						nil,
						nil,
					).JsonBody(view.NewErrorResponse(view.ErrMessageInvalidRequest)),
					HandlerFunc: handler,
				}
			},
		},
		{
			name: "internal server error",
			testCase: func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase {
				userAgent := "DeviceID=f78f5813-4e7a-4093-99d0-8c434fa38e22;DeviceType=web"
				requestBody := []byte(`{"key":"value"}`)
				internalError := errors.New("internal server error")

				f.interactor.EXPECT().AddMetrics(gomock.Any(), entityAnalytics.XCFCUserAgentHeader(userAgent), requestBody).Return("", internalError)

				return &gintest.HandlerTestCase{
					Request: &gintest.Request{
						Method:  http.MethodPost,
						Path:    "/analytics/metrics",
						Headers: map[string][]string{"X-CFC-Useragent": {userAgent}},
						Body:    bytes.NewBuffer(requestBody),
					},
					Response: gintest.NewResponse(
						http.StatusInternalServerError,
						map[string][]string{
							"Content-Type":   {"application/json; charset=utf-8"},
							StatusCodeHeader: {"AM_05"},
						},
						nil,
						nil,
					).JsonBody(view.NewErrorResponse(view.ErrMessageInternalError)),
					HandlerFunc: handler,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			f := fields{
				interactor: NewMockAnalyticsInteractor(ctrl),
			}

			h := NewAnalyticsHandlers(
				f.interactor,
			)
			ginTest := gintest.NewGinTest()
			ginTest.TestHandler(t, tt.testCase(h.addMetrics, f))
		})
	}
}

// errorReader is a reader that returns an error when Read is called.
type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("Simulated read error")
}
