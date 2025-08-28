package http

import (
	"net/http"
	"testing"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view"

	"git.mos.ru/buch-cloud/moscow-team-2.0/build/gintest.git"
	"github.com/gin-gonic/gin"
)

func TestNotImplementedHandler(t *testing.T) {
	type fields struct {
	}

	tests := []struct {
		name     string
		testCase func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase
	}{
		{
			name: "correct",
			testCase: func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase {
				return &gintest.HandlerTestCase{
					Request: &gintest.Request{
						Method: http.MethodPost,
						Path:   "/test",
					},
					Response: gintest.NewResponse(http.StatusMethodNotAllowed, nil, nil, nil).
						JsonBody(view.NewErrorResponse(view.ErrMessageMethodNotAllowed)),
					HandlerFunc: handler,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fields{}

			ginTest := gintest.NewGinTest()
			ginTest.TestHandler(t, tt.testCase(NotImplementedHandler, f))
		})
	}
}
