package http

import (
	"errors"
	"net/http"
	"testing"

	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/gintest.git"
	"github.com/gin-gonic/gin"
	"go.uber.org/mock/gomock"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view"
	authView "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/auth"
	usersView "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/users"
	entityAuth "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/auth"
	entityUser "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/user"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/usecase"
	authUseCase "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/usecase/auth"
	usersUseCase "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/usecase/users"
)

func Test_usersHandlers_getMe(t *testing.T) {
	type fields struct {
		interactor *MockUsersInteractor
		presenter  *MockUsersPresenter
	}

	testErr := errors.New("test error")

	tests := []struct {
		name     string
		testCase func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase
	}{
		{
			name: "get session err",
			testCase: func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase {
				f.interactor.EXPECT().GetMe(gomock.Any()).Return(nil, usecase.ErrGetSessionFromContext)
				return &gintest.HandlerTestCase{
					Request: &gintest.Request{
						Method: http.MethodGet,
						Path:   "/users/me",
					},
					Response: gintest.NewResponse(http.StatusUnauthorized, nil, nil, nil).
						JsonBody(view.NewErrorResponse(authView.ErrMessageFrontUnauthenticated)),
					HandlerFunc: handler,
				}
			},
		},
		{
			name: "nil session person id err",
			testCase: func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase {
				f.interactor.EXPECT().GetMe(gomock.Any()).Return(nil, usersUseCase.ErrEmptySessionPersonID)
				return &gintest.HandlerTestCase{
					Request: &gintest.Request{
						Method: http.MethodGet,
						Path:   "/users/me",
					},
					Response: gintest.NewResponse(http.StatusUnauthorized, nil, nil, nil).
						JsonBody(view.NewErrorResponse(authView.ErrMessageFrontUnauthenticated)),
					HandlerFunc: handler,
				}
			},
		},
		{
			name: "internal err",
			testCase: func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase {
				f.interactor.EXPECT().GetMe(gomock.Any()).Return(nil, testErr)
				return &gintest.HandlerTestCase{
					Request: &gintest.Request{
						Method: http.MethodGet,
						Path:   "/users/me",
					},
					Response: gintest.NewResponse(http.StatusInternalServerError, nil, nil, nil).
						JsonBody(view.NewErrorResponse(authView.ErrMessageFrontInternal)),
					HandlerFunc: handler,
				}
			},
		},
		{
			name: "correct",
			testCase: func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase {
				testUser := &entityUser.UserInfo{}
				f.interactor.EXPECT().GetMe(gomock.Any()).Return(testUser, nil)
				testUserView := &usersView.ShortUser{}
				f.presenter.EXPECT().ShortUserToView(&testUser.User).Return(testUserView)
				return &gintest.HandlerTestCase{
					Request: &gintest.Request{
						Method: http.MethodGet,
						Path:   "/users/me",
					},
					Response: gintest.NewResponse(http.StatusOK, nil, nil, nil).
						JsonBody(view.NewSuccessResponse(testUserView)),
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
				interactor: NewMockUsersInteractor(ctrl),
				presenter:  NewMockUsersPresenter(ctrl),
			}

			ph := NewUsersHandlers(nil, f.interactor, nil, f.presenter)
			ginTest := gintest.NewGinTest()
			ginTest.TestHandler(t, tt.testCase(ph.getMe, f))
		})
	}
}

func Test_usersHandlers_changePortal(t *testing.T) {
	type fields struct {
		interactor *MockAuthInteractor
		presenter  *MockAuthPresenter
	}

	testErr := errors.New("test error")

	tests := []struct {
		name     string
		testCase func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase
	}{
		{
			name: "invalid param type err",
			testCase: func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase {
				return &gintest.HandlerTestCase{
					Request: &gintest.Request{
						Method: http.MethodGet,
						Path:   "/changeportal/:id",
					},
					Response: gintest.NewResponse(http.StatusBadRequest, nil, nil, nil).
						JsonBody(view.NewErrorResponse(view.ErrMessageInvalidRequest)),
					Params: []gintest.Param{
						{
							Key:   "id",
							Value: "invalid",
						},
					},
					HandlerFunc: handler,
				}
			},
		},
		{
			name: "invalid param err",
			testCase: func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase {
				return &gintest.HandlerTestCase{
					Request: &gintest.Request{
						Method: http.MethodGet,
						Path:   "/changeportal/:id",
					},
					Response: gintest.NewResponse(http.StatusBadRequest, nil, nil, nil).
						JsonBody(view.NewErrorResponse(view.ErrMessageInvalidId)),
					Params: []gintest.Param{
						{
							Key:   "id",
							Value: "0",
						},
					},
					HandlerFunc: handler,
				}
			},
		},
		{
			name: "change portal validation err",
			testCase: func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase {
				testValidationErr := diterrors.NewValidationError(testErr)
				f.interactor.EXPECT().ChangePortal(gomock.Any(), 1).Return(nil, "", testValidationErr)
				return &gintest.HandlerTestCase{
					Request: &gintest.Request{
						Method: http.MethodGet,
						Path:   "/changeportal/:id",
					},
					Response: gintest.NewResponse(http.StatusBadRequest, nil, nil, nil).
						JsonBody(view.NewErrorResponse(authView.ErrMessageFrontSamePortals)),
					Params: []gintest.Param{
						{
							Key:   "id",
							Value: "1",
						},
					},
					HandlerFunc: handler,
				}
			},
		},
		{
			name: "get session err",
			testCase: func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase {
				f.interactor.EXPECT().ChangePortal(gomock.Any(), 1).Return(nil, "", usecase.ErrGetSessionFromContext)
				return &gintest.HandlerTestCase{
					Request: &gintest.Request{
						Method: http.MethodGet,
						Path:   "/changeportal/:id",
					},
					Response: gintest.NewResponse(http.StatusUnauthorized, nil, nil, nil).
						JsonBody(view.NewErrorResponse(authView.ErrMessageFrontUnauthenticated)),
					Params: []gintest.Param{
						{
							Key:   "id",
							Value: "1",
						},
					},
					HandlerFunc: handler,
				}
			},
		},
		{
			name: "portals not found err",
			testCase: func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase {
				f.interactor.EXPECT().ChangePortal(gomock.Any(), 1).Return(nil, "", authUseCase.ErrPortalsNotFound)
				return &gintest.HandlerTestCase{
					Request: &gintest.Request{
						Method: http.MethodGet,
						Path:   "/changeportal/:id",
					},
					Response: gintest.NewResponse(http.StatusBadRequest, nil, nil, nil).
						JsonBody(view.NewErrorResponse(authView.ErrMessageFrontPortalForUserNotFound)),
					Params: []gintest.Param{
						{
							Key:   "id",
							Value: "1",
						},
					},
					HandlerFunc: handler,
				}
			},
		},
		{
			name: "unavailable portal err",
			testCase: func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase {
				f.interactor.EXPECT().ChangePortal(gomock.Any(), 1).Return(nil, "", authUseCase.ErrUnavailablePortal)
				return &gintest.HandlerTestCase{
					Request: &gintest.Request{
						Method: http.MethodGet,
						Path:   "/changeportal/:id",
					},
					Response: gintest.NewResponse(http.StatusBadRequest, nil, nil, nil).
						JsonBody(view.NewErrorResponse(authView.ErrMessageFrontUnavailablePortal)),
					Params: []gintest.Param{
						{
							Key:   "id",
							Value: "1",
						},
					},
					HandlerFunc: handler,
				}
			},
		},
		{
			name: "empty portal url err",
			testCase: func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase {
				f.interactor.EXPECT().ChangePortal(gomock.Any(), 1).Return(nil, "", authUseCase.ErrEmptyPortalURL)
				return &gintest.HandlerTestCase{
					Request: &gintest.Request{
						Method: http.MethodGet,
						Path:   "/changeportal/:id",
					},
					Response: gintest.NewResponse(http.StatusInternalServerError, nil, nil, nil).
						JsonBody(view.NewErrorResponse(authView.ErrMessageFrontEmptyPortalURL)),
					Params: []gintest.Param{
						{
							Key:   "id",
							Value: "1",
						},
					},
					HandlerFunc: handler,
				}
			},
		},
		{
			name: "unauthenticated err",
			testCase: func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase {
				f.interactor.EXPECT().ChangePortal(gomock.Any(), 1).Return(nil, "", diterrors.ErrUnauthenticated)
				return &gintest.HandlerTestCase{
					Request: &gintest.Request{
						Method: http.MethodGet,
						Path:   "/changeportal/:id",
					},
					Response: gintest.NewResponse(http.StatusUnauthorized, nil, nil, nil).
						JsonBody(view.NewErrorResponse(authView.ErrMessageFrontUnauthenticated)),
					Params: []gintest.Param{
						{
							Key:   "id",
							Value: "1",
						},
					},
					HandlerFunc: handler,
				}
			},
		},
		{
			name: "not found err",
			testCase: func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase {
				f.interactor.EXPECT().ChangePortal(gomock.Any(), 1).Return(nil, "", diterrors.ErrNotFound)
				return &gintest.HandlerTestCase{
					Request: &gintest.Request{
						Method: http.MethodGet,
						Path:   "/changeportal/:id",
					},
					Response: gintest.NewResponse(http.StatusNotFound, nil, nil, nil).
						JsonBody(view.NewErrorResponse(authView.ErrMessageFrontUserNotIntoPortal)),
					Params: []gintest.Param{
						{
							Key:   "id",
							Value: "1",
						},
					},
					HandlerFunc: handler,
				}
			},
		},
		{
			name: "internal err",
			testCase: func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase {
				f.interactor.EXPECT().ChangePortal(gomock.Any(), 1).Return(nil, "", testErr)
				return &gintest.HandlerTestCase{
					Request: &gintest.Request{
						Method: http.MethodGet,
						Path:   "/changeportal/:id",
					},
					Response: gintest.NewResponse(http.StatusInternalServerError, nil, nil, nil).
						JsonBody(view.NewErrorResponse(authView.ErrMessageFrontInternal)),
					Params: []gintest.Param{
						{
							Key:   "id",
							Value: "1",
						},
					},
					HandlerFunc: handler,
				}
			},
		},
		{
			name: "correct",
			testCase: func(handler gin.HandlerFunc, f fields) *gintest.HandlerTestCase {
				testPortals := make([]*entityAuth.Portal, 0)
				testPortalSID := "testPortalSID"
				f.interactor.EXPECT().ChangePortal(gomock.Any(), 1).Return(testPortals, testPortalSID, nil)
				testResponse := &authView.AuthResponse{PortalSession: testPortalSID, Portals: []*authView.Portal{{}}}
				f.presenter.EXPECT().AuthToView(&entityAuth.Auth{PortalSession: testPortalSID, Portals: testPortals}).Return(testResponse)
				return &gintest.HandlerTestCase{
					Request: &gintest.Request{
						Method: http.MethodGet,
						Path:   "/changeportal/:id",
					},
					Response: gintest.NewResponse(http.StatusOK, nil, nil, nil).
						JsonBody(view.NewSuccessResponse(testResponse)),
					Params: []gintest.Param{
						{
							Key:   "id",
							Value: "1",
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
			defer ctrl.Finish()
			f := fields{
				interactor: NewMockAuthInteractor(ctrl),
				presenter:  NewMockAuthPresenter(ctrl),
			}

			ph := NewUsersHandlers(f.interactor, nil, f.presenter, nil)
			ginTest := gintest.NewGinTest()
			ginTest.TestHandler(t, tt.testCase(ph.changePortal, f))
		})
	}
}
