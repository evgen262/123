package http

import (
	"net/http"
	"testing"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/file"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/gintest.git"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_filesHandlers_get(t *testing.T) {
	type fields struct {
		interactor *MockFilesInteractor
	}

	testUUID := uuid.New()

	tests := []struct {
		name     string
		testCase func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase
	}{
		{
			name: "invalid file id err",
			testCase: func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase {
				return &gintest.HandlerTestCase{
					Request: &gintest.Request{
						Method: http.MethodGet,
						Path:   "/:file_id",
					},
					Response: gintest.NewResponse(
						http.StatusBadRequest,
						map[string][]string{
							"Content-Type":   {"application/json; charset=utf-8"},
							StatusCodeHeader: {"WAF_01"},
						},
						nil,
						nil,
					).JsonBody(view.NewErrorResponse(view.ErrMessageInvalidRequest)),
					Params: []gintest.Param{
						{
							Key:   "file_id",
							Value: "invalid",
						},
					},
					HandlerFunc: handler,
				}
			},
		},
		{
			name: "unauthenticated err",
			testCase: func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase {
				f.interactor.EXPECT().Get(gomock.Any(), testUUID).Return(nil, diterrors.ErrUnauthenticated)
				return &gintest.HandlerTestCase{
					Request: &gintest.Request{
						Method: http.MethodGet,
						Path:   "/:file_id",
					},
					Response: gintest.NewResponse(
						http.StatusUnauthorized,
						map[string][]string{
							"Content-Type":   {"application/json; charset=utf-8"},
							StatusCodeHeader: {"WAF_02"},
						},
						nil,
						nil,
					).JsonBody(view.NewErrorResponse(view.ErrMessageUnauthenticated)),
					Params: []gintest.Param{
						{
							Key:   "file_id",
							Value: testUUID.String(),
						},
					},
					HandlerFunc: handler,
				}
			},
		},
		{
			name: "permission denied err",
			testCase: func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase {
				f.interactor.EXPECT().Get(gomock.Any(), testUUID).Return(nil, diterrors.ErrPermissionDenied)
				return &gintest.HandlerTestCase{
					Request: &gintest.Request{
						Method: http.MethodGet,
						Path:   "/:file_id",
					},
					Response: gintest.NewResponse(
						http.StatusForbidden,
						map[string][]string{
							"Content-Type":   {"application/json; charset=utf-8"},
							StatusCodeHeader: {"WAF_03"},
						},
						nil,
						nil,
					).JsonBody(view.NewErrorResponse(view.ErrPermissionDenied)),
					Params: []gintest.Param{
						{
							Key:   "file_id",
							Value: testUUID.String(),
						},
					},
					HandlerFunc: handler,
				}
			},
		},
		{
			name: "not found err",
			testCase: func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase {
				f.interactor.EXPECT().Get(gomock.Any(), testUUID).Return(nil, diterrors.ErrNotFound)
				return &gintest.HandlerTestCase{
					Request: &gintest.Request{
						Method: http.MethodGet,
						Path:   "/:file_id",
					},
					Response: gintest.NewResponse(
						http.StatusNotFound,
						map[string][]string{
							"Content-Type":   {"application/json; charset=utf-8"},
							StatusCodeHeader: {"WAF_04"},
						},
						nil,
						nil,
					).JsonBody(view.NewErrorResponse(view.ErrMessageNotFound)),
					Params: []gintest.Param{
						{
							Key:   "file_id",
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
				testValidationErr := diterrors.NewValidationError(assert.AnError)
				f.interactor.EXPECT().Get(gomock.Any(), testUUID).Return(nil, testValidationErr)
				return &gintest.HandlerTestCase{
					Request: &gintest.Request{
						Method: http.MethodGet,
						Path:   "/:file_id",
					},
					Response: gintest.NewResponse(
						http.StatusBadRequest,
						map[string][]string{
							"Content-Type":   {"application/json; charset=utf-8"},
							StatusCodeHeader: {"WAF_05"},
						},
						nil,
						nil,
					).JsonBody(view.NewErrorResponse(view.ErrMessageInvalidRequest)),
					Params: []gintest.Param{
						{
							Key:   "file_id",
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
				f.interactor.EXPECT().Get(gomock.Any(), testUUID).Return(nil, assert.AnError)
				return &gintest.HandlerTestCase{
					Request: &gintest.Request{
						Method: http.MethodGet,
						Path:   "/:file_id",
					},
					Response: gintest.NewResponse(
						http.StatusInternalServerError,
						map[string][]string{
							"Content-Type":   {"application/json; charset=utf-8"},
							StatusCodeHeader: {"WAF_06"},
						},
						nil,
						nil,
					).JsonBody(view.NewErrorResponse(view.ErrMessageInternalError)),
					Params: []gintest.Param{
						{
							Key:   "file_id",
							Value: testUUID.String(),
						},
					},
					HandlerFunc: handler,
				}
			},
		},
		{
			name: "correct",
			testCase: func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase {
				testFile := &file.File{
					Name: "testFile.txt",
					Metadata: file.Metadata{
						ContentType: "text/plain",
					},
					Payload: []byte("test file content"),
				}
				f.interactor.EXPECT().Get(gomock.Any(), testUUID).Return(testFile, nil)
				return &gintest.HandlerTestCase{
					Request: &gintest.Request{
						Method: http.MethodGet,
						Path:   "/:file_id",
					},
					Response: gintest.NewResponse(http.StatusOK, map[string][]string{
						"X-File-Name":  {testFile.GetFileName()},
						"Content-Type": {testFile.Metadata.ContentType},
					}, testFile.Payload, nil),
					Params: []gintest.Param{
						{
							Key:   "file_id",
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
				interactor: NewMockFilesInteractor(ctrl),
			}

			fh := NewFilesHandlers(f.interactor)
			ginTest := gintest.NewGinTest()
			ginTest.TestHandler(t, tt.testCase(fh.get, f))
		})
	}
}
