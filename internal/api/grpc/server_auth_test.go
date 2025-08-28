package grpc

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"

	authv1 "git.mos.ru/buch-cloud/moscow-team-2.0/pud/auth.git/internal/api/grpc/gen/infogorod/auth/auth/v1"
	"git.mos.ru/buch-cloud/moscow-team-2.0/pud/auth.git/internal/api/grpc/view"
	"git.mos.ru/buch-cloud/moscow-team-2.0/pud/auth.git/internal/client/http/kadry"
	"git.mos.ru/buch-cloud/moscow-team-2.0/pud/auth.git/internal/client/http/sudir"
	"git.mos.ru/buch-cloud/moscow-team-2.0/pud/auth.git/internal/entity"
	"git.mos.ru/buch-cloud/moscow-team-2.0/pud/auth.git/internal/repositories"
	"git.mos.ru/buch-cloud/moscow-team-2.0/pud/auth.git/internal/usecase"
)

func Test_authServer_GetURL(t *testing.T) {
	type fields struct {
		interactor *MockAuthInteractor
		presenter  *MockAuthPresenter
	}
	type args struct {
		ctx     context.Context
		request *authv1.GetURLRequest
	}

	ctx := context.TODO()
	testErr := errors.New("some test error")
	mdCtx := metadata.NewIncomingContext(ctx, metadata.Pairs("deviceid", "testDeviceID", "x-cfc-useragent", "testUserAgent"))

	tests := []struct {
		name string
		args args
		want func(a args, f fields) (*authv1.GetURLResponse, error)
	}{
		{
			name: "empty credentials error",
			args: args{
				ctx: ctx,
				request: &authv1.GetURLRequest{
					Credentials: &authv1.Credentials{},
				},
			},
			want: func(a args, f fields) (*authv1.GetURLResponse, error) {
				return nil, NewApiError(codes.InvalidArgument, "bad credentials provided").
					WithLocalizedMessage(view.ErrStatusTextBadCredentialsProvided)
			},
		},
		{
			name: "error",
			args: args{
				ctx:     ctx,
				request: &authv1.GetURLRequest{CallbackUrl: "https://some/url/path"},
			},
			want: func(a args, f fields) (*authv1.GetURLResponse, error) {
				f.interactor.EXPECT().GetAuthURL(a.ctx, a.request.GetCallbackUrl(), "", "").Return("", testErr)
				return nil, &apiError{
					Code:    codes.Internal,
					Message: fmt.Errorf("auth service error: %w", testErr).Error(),
					localizeMessage: &errdetails.LocalizedMessage{
						Locale:  "ru-RU",
						Message: view.ErrStatusTextServiceError,
					},
				}
			},
		},
		{
			name: "correct",
			args: args{
				ctx: mdCtx,
				request: &authv1.GetURLRequest{
					CallbackUrl: "https://some/url/path",
					Credentials: &authv1.Credentials{
						ClientId:     "testClientID",
						ClientSecret: "testClientSecret",
					},
				},
			},
			want: func(a args, f fields) (*authv1.GetURLResponse, error) {
				testUrl := "https://service/path/to/login"
				a.ctx = context.WithValue(a.ctx, "deviceid", "testDeviceID")
				a.ctx = context.WithValue(a.ctx, "x-cfc-useragent", "testUserAgent")
				f.interactor.EXPECT().GetAuthURL(a.ctx, a.request.GetCallbackUrl(), a.request.GetCredentials().GetClientId(), a.request.GetCredentials().GetClientSecret()).Return(testUrl, nil)
				return &authv1.GetURLResponse{
					RedirectUrl: testUrl,
				}, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				interactor: NewMockAuthInteractor(ctrl),
				presenter:  NewMockAuthPresenter(ctrl),
			}

			want, wantErr := tt.want(tt.args, f)

			a := NewAuthServer(f.interactor, f.presenter)
			got, err := a.GetURL(tt.args.ctx, tt.args.request)

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

func Test_authServer_Auth(t *testing.T) {
	type fields struct {
		interactor *MockAuthInteractor
		presenter  *MockAuthPresenter
	}
	type args struct {
		ctx     context.Context
		request *authv1.AuthRequest
	}

	ctx := context.TODO()
	testTime := time.Now()
	testErr := errors.New("some test error")

	tests := []struct {
		name string
		args args
		want func(a args, f fields) (*authv1.AuthResponse, error)
	}{
		{
			name: "empty code",
			args: args{
				ctx: ctx,
				request: &authv1.AuthRequest{
					State:       "71bac977-ea7f-4156-9504-60f7d443ab62",
					Code:        "",
					CallbackUrl: "",
				},
			},
			want: func(a args, f fields) (*authv1.AuthResponse, error) {
				return nil, &apiError{
					Code:    codes.InvalidArgument,
					Message: "code not provided",
					localizeMessage: &errdetails.LocalizedMessage{
						Locale:  "ru-RU",
						Message: view.ErrStatusTextCodeNotProvided,
					},
				}
			},
		},
		{
			name: "empty callback url",
			args: args{
				ctx: ctx,
				request: &authv1.AuthRequest{
					State:       "71bac977-ea7f-4156-9504-60f7d443ab62",
					Code:        "XXijKLnoayN8xX9ap3JUMLBATcv2m8V9irtIvrKEtofDfkH0eGe2XAb9qzvKJXXuuJHtVokGL_iW5urIGNtQgWuhV1MK46C4D4jPEMl77-Q",
					CallbackUrl: "",
				},
			},
			want: func(a args, f fields) (*authv1.AuthResponse, error) {
				return nil, &apiError{
					Code:    codes.InvalidArgument,
					Message: "callback url not provided",
					localizeMessage: &errdetails.LocalizedMessage{
						Locale:  "ru-RU",
						Message: view.ErrStatusTextCallbackURLNotProvided,
					},
				}
			},
		},
		{
			name: "sudir invalid grant",
			args: args{
				ctx: ctx,
				request: &authv1.AuthRequest{
					State:       "71bac977-ea7f-4156-9504-60f7d443ab62",
					Code:        "XXijKLnoayN8xX9ap3JUMLBATcv2m8V9irtIvrKEtofDfkH0eGe2XAb9qzvKJXXuuJHtVokGL_iW5urIGNtQgWuhV1MK46C4D4jPEMl77-Q",
					CallbackUrl: "https://some/callback/url",
				},
			},
			want: func(a args, f fields) (*authv1.AuthResponse, error) {
				f.interactor.EXPECT().Auth(a.ctx, a.request.GetCode(), a.request.GetState(), a.request.GetCallbackUrl()).Return(nil,
					sudir.ErrInvalidGrant)
				return nil, &apiError{
					Code:    codes.InvalidArgument,
					Message: "invalid code provided",
					localizeMessage: &errdetails.LocalizedMessage{
						Locale:  "ru-RU",
						Message: view.ErrStatusTextBadCodeProvided,
					},
				}
			},
		},
		{
			name: "sudir invalid state",
			args: args{
				ctx: ctx,
				request: &authv1.AuthRequest{
					State:       "71bac977-ea7f-4156-9504-60f7d443ab62",
					Code:        "XXijKLnoayN8xX9ap3JUMLBATcv2m8V9irtIvrKEtofDfkH0eGe2XAb9qzvKJXXuuJHtVokGL_iW5urIGNtQgWuhV1MK46C4D4jPEMl77-Q",
					CallbackUrl: "https://some/callback/url",
				},
			},
			want: func(a args, f fields) (*authv1.AuthResponse, error) {
				f.interactor.EXPECT().Auth(a.ctx, a.request.GetCode(), a.request.GetState(), a.request.GetCallbackUrl()).Return(nil, repositories.ErrInvalidState)
				return nil, &apiError{
					Code:    codes.InvalidArgument,
					Message: "invalid state provided",
					localizeMessage: &errdetails.LocalizedMessage{
						Locale:  "ru-RU",
						Message: view.ErrStatusTextStateInvalid,
					},
				}
			},
		},
		{
			name: "sudir state not found",
			args: args{
				ctx: ctx,
				request: &authv1.AuthRequest{
					State:       "71bac977-ea7f-4156-9504-60f7d443ab62",
					Code:        "XXijKLnoayN8xX9ap3JUMLBATcv2m8V9irtIvrKEtofDfkH0eGe2XAb9qzvKJXXuuJHtVokGL_iW5urIGNtQgWuhV1MK46C4D4jPEMl77-Q",
					CallbackUrl: "https://some/callback/url",
				},
			},
			want: func(a args, f fields) (*authv1.AuthResponse, error) {
				f.interactor.EXPECT().Auth(a.ctx, a.request.GetCode(), a.request.GetState(), a.request.GetCallbackUrl()).Return(nil, repositories.ErrStateNotFound)
				return nil, &apiError{
					Code:    codes.InvalidArgument,
					Message: "invalid state provided",
					localizeMessage: &errdetails.LocalizedMessage{
						Locale:  "ru-RU",
						Message: view.ErrStatusTextStateInvalid,
					},
				}
			},
		},
		{
			name: "usecase error",
			args: args{
				ctx: ctx,
				request: &authv1.AuthRequest{
					State:       "71bac977-ea7f-4156-9504-60f7d443ab62",
					Code:        "XXijKLnoayN8xX9ap3JUMLBATcv2m8V9irtIvrKEtofDfkH0eGe2XAb9qzvKJXXuuJHtVokGL_iW5urIGNtQgWuhV1MK46C4D4jPEMl77-Q",
					CallbackUrl: "https://some/callback/url",
				},
			},
			want: func(a args, f fields) (*authv1.AuthResponse, error) {
				f.interactor.EXPECT().Auth(a.ctx, a.request.GetCode(), a.request.GetState(), a.request.GetCallbackUrl()).Return(nil, testErr)
				return nil, &apiError{
					Code:    codes.Internal,
					Message: fmt.Errorf("auth service error: %w", testErr).Error(),
					localizeMessage: &errdetails.LocalizedMessage{
						Locale:  "ru-RU",
						Message: view.ErrStatusTextServiceError,
					},
				}
			},
		},
		{
			name: "url mismatch error",
			args: args{
				ctx: ctx,
				request: &authv1.AuthRequest{
					State:       "71bac977-ea7f-4156-9504-60f7d443ab62",
					Code:        "XXijKLnoayN8xX9ap3JUMLBATcv2m8V9irtIvrKEtofDfkH0eGe2XAb9qzvKJXXuuJHtVokGL_iW5urIGNtQgWuhV1MK46C4D4jPEMl77-Q",
					CallbackUrl: "https://some/callback/url",
				},
			},
			want: func(a args, f fields) (*authv1.AuthResponse, error) {
				f.interactor.EXPECT().Auth(a.ctx, a.request.GetCode(), a.request.GetState(), a.request.GetCallbackUrl()).Return(nil, usecase.ErrCallbackURLMismatch)
				return nil, NewApiError(codes.InvalidArgument, "callback url mismatch").
					WithLocalizedMessage(view.ErrStatusTextCallbackURLMismatch)
			},
		},
		{
			name: "employees error",
			args: args{
				ctx: ctx,
				request: &authv1.AuthRequest{
					State:       "71bac977-ea7f-4156-9504-60f7d443ab62",
					Code:        "XXijKLnoayN8xX9ap3JUMLBATcv2m8V9irtIvrKEtofDfkH0eGe2XAb9qzvKJXXuuJHtVokGL_iW5urIGNtQgWuhV1MK46C4D4jPEMl77-Q",
					CallbackUrl: "https://some/callback/url",
				},
			},
			want: func(a args, f fields) (*authv1.AuthResponse, error) {
				authEntity := &entity.AuthInfo{
					OAuth: &entity.OAuth{
						AccessToken:  "some_test_access_token",
						RefreshToken: "some_test_refresh_token",
						Expiry:       &testTime,
					},
					User: &entity.User{
						CloudID: "12-some-uuid-345",
						Info: &entity.UserInfo{
							Email:     "test_email@mail.local",
							LogonName: "TestUN",
						},
					},
					Device: &entity.Device{
						ID:        "testDeviceID",
						ClientID:  "testClientID",
						UserAgent: "testUserAgent",
					},
				}
				testViewUser := &authv1.User{
					CloudId: "12-some-uuid-345",
					Info: &authv1.UserInfo{
						Email:     "test_email@mail.local",
						LogonName: "TestUN",
					},
				}
				testViewDevice := &authv1.AuthResponse_Device{
					Id:        "testDeviceID",
					ClientId:  "testClientID",
					UserAgent: "testUserAgent",
				}

				f.interactor.EXPECT().Auth(a.ctx, a.request.GetCode(), a.request.GetState(), a.request.GetCallbackUrl()).Return(authEntity, kadry.ErrEmployeesService)
				f.presenter.EXPECT().UserToPb(authEntity.User).Return(testViewUser)
				f.presenter.EXPECT().DeviceToPb(authEntity.Device).Return(testViewDevice)
				return &authv1.AuthResponse{
					User:        testViewUser,
					AccessToken: authEntity.OAuth.AccessToken,
					Device:      testViewDevice,
				}, nil
			},
		},
		{
			name: "correct",
			args: args{
				ctx: ctx,
				request: &authv1.AuthRequest{
					State:       "71bac977-ea7f-4156-9504-60f7d443ab62",
					Code:        "XXijKLnoayN8xX9ap3JUMLBATcv2m8V9irtIvrKEtofDfkH0eGe2XAb9qzvKJXXuuJHtVokGL_iW5urIGNtQgWuhV1MK46C4D4jPEMl77-Q",
					CallbackUrl: "https://some/callback/url",
				},
			},
			want: func(a args, f fields) (*authv1.AuthResponse, error) {
				authEntity := &entity.AuthInfo{
					OAuth: &entity.OAuth{
						AccessToken:  "some_test_access_token",
						RefreshToken: "some_test_refresh_token",
						Expiry:       &testTime,
					},
					User: &entity.User{
						CloudID: "12-some-uuid-345",
						Info: &entity.UserInfo{
							Email:     "test_email@mail.local",
							LogonName: "TestUN",
						},
						Employees: []entity.EmployeeInfo{
							{
								CloudID: "12-some-uuid-345",
								Inn:     "770987654321",
								OrgID:   "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
								FIO:     "Иванов Иван Иванович",
							},
							{
								CloudID: "12-some-uuid-345",
								Inn:     "771234567890",
								OrgID:   "71bac977-ea7f-4156-9504-60f7d443ab62",
								FIO:     "Иванов Иван Иванович",
							},
						},
					},
				}
				testViewUser := &authv1.User{
					CloudId: "12-some-uuid-345",
					Info: &authv1.UserInfo{
						Email:     "test_email@mail.local",
						LogonName: "TestUN",
					},
					Employees: []*authv1.Employee{
						{
							Inn:   "770987654321",
							OrgId: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
							Fio:   "Иванов Иван Иванович",
						},
						{
							Inn:   "771234567890",
							OrgId: "71bac977-ea7f-4156-9504-60f7d443ab62",
							Fio:   "Иванов Иван Иванович",
						},
					},
				}
				testViewDevice := &authv1.AuthResponse_Device{
					Id:        "testDeviceID",
					ClientId:  "testClientID",
					UserAgent: "testUserAgent",
				}

				f.interactor.EXPECT().Auth(a.ctx, a.request.GetCode(), a.request.GetState(), a.request.GetCallbackUrl()).Return(authEntity, nil)
				f.presenter.EXPECT().UserToPb(authEntity.User).Return(testViewUser)
				f.presenter.EXPECT().DeviceToPb(authEntity.Device).Return(testViewDevice)
				return &authv1.AuthResponse{
					User:        testViewUser,
					AccessToken: authEntity.OAuth.AccessToken,
					Device:      testViewDevice,
				}, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				interactor: NewMockAuthInteractor(ctrl),
				presenter:  NewMockAuthPresenter(ctrl),
			}

			want, wantErr := tt.want(tt.args, f)

			a := NewAuthServer(f.interactor, f.presenter)
			got, err := a.Auth(tt.args.ctx, tt.args.request)

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

func Test_authServer_Login(t *testing.T) {
	type fields struct {
		interactor *MockAuthInteractor
		presenter  *MockAuthPresenter
	}
	type args struct {
		ctx     context.Context
		request *authv1.LoginRequest
	}

	ctx := context.TODO()
	testErr := errors.New("some test error")

	tests := []struct {
		name   string
		fields fields
		args   args
		want   func(a args, f fields) (*authv1.LoginResponse, error)
	}{
		{
			name: "invalid login by err",
			args: args{
				ctx:     ctx,
				request: &authv1.LoginRequest{},
			},
			want: func(a args, f fields) (*authv1.LoginResponse, error) {
				return nil, NewApiError(codes.InvalidArgument, "credentials not provided").
					WithLocalizedMessage(view.ErrStatusTextCredentialsNotProvided)
			},
		},
		{
			name: "credentials nil err",
			args: args{
				ctx: ctx,
				request: &authv1.LoginRequest{
					LoginBy: &authv1.LoginRequest_Credentials{
						Credentials: nil,
					},
				},
			},
			want: func(a args, f fields) (*authv1.LoginResponse, error) {
				return nil, NewApiError(codes.InvalidArgument, "credentials not provided").
					WithLocalizedMessage(view.ErrStatusTextCredentialsNotProvided)
			},
		},
		{
			name: "invalid client err",
			args: args{
				ctx: ctx,
				request: &authv1.LoginRequest{
					LoginBy: &authv1.LoginRequest_Credentials{
						Credentials: &authv1.Credentials{
							ClientId:     "testClientID",
							ClientSecret: "testClientSecret",
						},
					},
				},
			},
			want: func(a args, f fields) (*authv1.LoginResponse, error) {
				testCredentials := a.request.GetCredentials()
				f.interactor.EXPECT().LoginByCredentials(a.ctx, testCredentials.GetClientId(), testCredentials.GetClientSecret()).Return(nil, sudir.ErrInvalidClient)
				return nil, NewApiError(codes.Unauthenticated, "bad credentials provided").
					WithLocalizedMessage(view.ErrStatusTextBadCredentialsProvided)
			},
		},
		{
			name: "invalid client err",
			args: args{
				ctx: ctx,
				request: &authv1.LoginRequest{
					LoginBy: &authv1.LoginRequest_Credentials{
						Credentials: &authv1.Credentials{
							ClientId:     "testClientID",
							ClientSecret: "testClientSecret",
						},
					},
				},
			},
			want: func(a args, f fields) (*authv1.LoginResponse, error) {
				testCredentials := a.request.GetCredentials()
				f.interactor.EXPECT().LoginByCredentials(a.ctx, testCredentials.GetClientId(), testCredentials.GetClientSecret()).Return(nil, sudir.ErrAccessDenied)
				return nil, NewApiError(codes.PermissionDenied, "bad credentials provided").
					WithLocalizedMessage(view.ErrStatusTextBadCredentialsProvided)
			},
		},
		{
			name: "login by credentials err",
			args: args{
				ctx: ctx,
				request: &authv1.LoginRequest{
					LoginBy: &authv1.LoginRequest_Credentials{
						Credentials: &authv1.Credentials{
							ClientId:     "testClientID",
							ClientSecret: "testClientSecret",
						},
					},
				},
			},
			want: func(a args, f fields) (*authv1.LoginResponse, error) {
				f.interactor.EXPECT().LoginByCredentials(a.ctx, a.request.GetCredentials().GetClientId(), a.request.GetCredentials().GetClientSecret()).Return(nil, testErr)
				return nil, NewApiError(codes.Internal, "auth service error", testErr).
					WithLocalizedMessage(view.ErrStatusTextServiceError)
			},
		},
		{
			name: "correct",
			args: args{
				ctx: ctx,
				request: &authv1.LoginRequest{
					LoginBy: &authv1.LoginRequest_Credentials{
						Credentials: &authv1.Credentials{
							ClientId:     "testClientID",
							ClientSecret: "testClientSecret",
						},
					},
				},
			},
			want: func(a args, f fields) (*authv1.LoginResponse, error) {
				testInfo := &entity.AuthInfo{OAuth: &entity.OAuth{AccessToken: "testAccessToken"}}
				f.interactor.EXPECT().LoginByCredentials(a.ctx, a.request.GetCredentials().GetClientId(), a.request.GetCredentials().GetClientSecret()).Return(testInfo, nil)
				return &authv1.LoginResponse{AccessToken: testInfo.OAuth.AccessToken}, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				interactor: NewMockAuthInteractor(ctrl),
				presenter:  NewMockAuthPresenter(ctrl),
			}

			want, wantErr := tt.want(tt.args, f)

			a := NewAuthServer(f.interactor, f.presenter)
			got, err := a.Login(tt.args.ctx, tt.args.request)

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

func Test_authServer_RefreshToken(t *testing.T) {
	type fields struct {
		interactor *MockAuthInteractor
		presenter  *MockAuthPresenter
	}
	type args struct {
		ctx     context.Context
		request *authv1.RefreshTokenRequest
	}

	ctx := context.TODO()
	testErr := errors.New("some test error")

	tests := []struct {
		name string
		args args
		want func(a args, f fields) (*authv1.RefreshTokenResponse, error)
	}{
		{
			name: "id not provided err",
			args: args{
				ctx:     ctx,
				request: &authv1.RefreshTokenRequest{},
			},
			want: func(a args, f fields) (*authv1.RefreshTokenResponse, error) {
				return nil, NewApiError(codes.InvalidArgument, "id not provided").
					WithLocalizedMessage(view.ErrStatusTextCloudIDNotProvided)
			},
		},
		{
			name: "invalid grant err",
			args: args{
				ctx: ctx,
				request: &authv1.RefreshTokenRequest{
					Id: &authv1.RefreshTokenRequest_CloudId{
						CloudId: "testCloudID",
					},
				},
			},
			want: func(a args, f fields) (*authv1.RefreshTokenResponse, error) {
				f.interactor.EXPECT().RefreshToken(a.ctx, a.request.GetCloudId(), a.request.GetSid()).Return("",
					sudir.ErrInvalidGrant)
				return nil, NewApiError(codes.Internal, "sudir service error").
					WithLocalizedMessage(view.ErrStatusTextServiceError)
			},
		},
		{
			name: "not found err",
			args: args{
				ctx: ctx,
				request: &authv1.RefreshTokenRequest{
					Id: &authv1.RefreshTokenRequest_CloudId{
						CloudId: "testCloudID",
					},
				},
			},
			want: func(a args, f fields) (*authv1.RefreshTokenResponse, error) {
				f.interactor.EXPECT().RefreshToken(a.ctx, a.request.GetCloudId(), a.request.GetSid()).Return("", repositories.ErrNotFound)
				return nil, NewApiError(codes.NotFound, "refresh token not found").
					WithLocalizedMessage(view.ErrStatusTextRefreshTokenNotFound)
			},
		},
		{
			name: "usecase refresh token err",
			args: args{
				ctx: ctx,
				request: &authv1.RefreshTokenRequest{
					Id: &authv1.RefreshTokenRequest_ClientId{
						ClientId: "testClientID",
					},
				},
			},
			want: func(a args, f fields) (*authv1.RefreshTokenResponse, error) {
				f.interactor.EXPECT().RefreshToken(a.ctx, a.request.GetClientId(), a.request.GetSid()).Return("", testErr)
				return nil, NewApiError(codes.Internal, "auth service error", testErr).WithLocalizedMessage(view.ErrStatusTextServiceError)
			},
		},
		{
			name: "correct",
			args: args{
				ctx: ctx,
				request: &authv1.RefreshTokenRequest{
					Id: &authv1.RefreshTokenRequest_ClientId{
						ClientId: "testClientID",
					},
				},
			},
			want: func(a args, f fields) (*authv1.RefreshTokenResponse, error) {
				testToken := "testToken"
				f.interactor.EXPECT().RefreshToken(a.ctx, a.request.GetClientId(), a.request.GetSid()).Return(testToken, nil)
				return &authv1.RefreshTokenResponse{AccessToken: testToken}, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				interactor: NewMockAuthInteractor(ctrl),
				presenter:  NewMockAuthPresenter(ctrl),
			}

			want, wantErr := tt.want(tt.args, f)

			a := NewAuthServer(f.interactor, f.presenter)
			got, err := a.RefreshToken(tt.args.ctx, tt.args.request)

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

func Test_authServer_Logout(t *testing.T) {
	type fields struct {
		interactor *MockAuthInteractor
		presenter  *MockAuthPresenter
	}
	type args struct {
		ctx     context.Context
		request *authv1.LogoutRequest
	}

	ctx := context.TODO()
	testErr := errors.New("some test error")

	tests := []struct {
		name   string
		fields fields
		args   args
		want   func(a args, f fields) (*authv1.LogoutResponse, error)
	}{
		{
			name: "credentials nil err",
			args: args{
				ctx: ctx,
				request: &authv1.LogoutRequest{
					Id: &authv1.LogoutRequest_Credentials{},
				},
			},
			want: func(a args, f fields) (*authv1.LogoutResponse, error) {
				return nil, NewApiError(codes.InvalidArgument, "credentials not provided").
					WithLocalizedMessage(view.ErrStatusTextBadCredentialsProvided)
			},
		},
		{
			name: "token empty err",
			args: args{
				ctx: ctx,
				request: &authv1.LogoutRequest{
					Id: &authv1.LogoutRequest_Credentials{
						Credentials: &authv1.LogoutRequest_ClientCredentials{},
					},
				},
			},
			want: func(a args, f fields) (*authv1.LogoutResponse, error) {
				return nil, NewApiError(codes.InvalidArgument, "registration token not provided").
					WithLocalizedMessage(view.ErrStatusTextAccessTokenNotProvided)
			},
		},
		{
			name: "id not provided err",
			args: args{
				ctx: ctx,
			},
			want: func(a args, f fields) (*authv1.LogoutResponse, error) {
				return nil, NewApiError(codes.InvalidArgument, "id not provided").
					WithLocalizedMessage(view.ErrStatusTextBadCredentialsProvided)
			},
		},
		{
			name: "access denied err",
			args: args{
				ctx: ctx,
				request: &authv1.LogoutRequest{
					Id: &authv1.LogoutRequest_CloudId{
						CloudId: "testCloudID",
					},
				},
			},
			want: func(a args, f fields) (*authv1.LogoutResponse, error) {
				f.interactor.EXPECT().Logout(a.ctx, a.request.GetCloudId(), "", "").Return(sudir.ErrAccessDenied)
				return nil, NewApiError(codes.PermissionDenied, "access denied").
					WithLocalizedMessage(view.ErrStatusTextAccessDenied)
			},
		},
		{
			name: "usecase logout err",
			args: args{
				ctx: ctx,
				request: &authv1.LogoutRequest{
					Id: &authv1.LogoutRequest_CloudId{
						CloudId: "testCloudID",
					},
				},
			},
			want: func(a args, f fields) (*authv1.LogoutResponse, error) {
				f.interactor.EXPECT().Logout(a.ctx, a.request.GetCloudId(), "", "").Return(testErr)
				return nil, NewApiError(codes.Internal, "auth service error").
					WithLocalizedMessage(view.ErrStatusTextServiceError)
			},
		},
		{
			name: "correct",
			args: args{
				ctx: ctx,
				request: &authv1.LogoutRequest{
					Id: &authv1.LogoutRequest_CloudId{
						CloudId: "testCloudID",
					},
				},
			},
			want: func(a args, f fields) (*authv1.LogoutResponse, error) {
				f.interactor.EXPECT().Logout(a.ctx, a.request.GetCloudId(), "", "").Return(nil)
				return &authv1.LogoutResponse{}, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				interactor: NewMockAuthInteractor(ctrl),
				presenter:  NewMockAuthPresenter(ctrl),
			}

			want, wantErr := tt.want(tt.args, f)

			a := NewAuthServer(f.interactor, f.presenter)
			got, err := a.Logout(tt.args.ctx, tt.args.request)

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

func Test_authServer_GetUser(t *testing.T) {
	type fields struct {
		interactor *MockAuthInteractor
		presenter  *MockAuthPresenter
	}
	type args struct {
		ctx     context.Context
		request *authv1.GetUserRequest
	}

	ctx := context.TODO()
	testErr := errors.New("some test error")

	tests := []struct {
		name   string
		fields fields
		args   args
		want   func(a args, f fields) (*authv1.GetUserResponse, error)
	}{
		{
			name: "credentials nil err",
			args: args{
				ctx:     ctx,
				request: &authv1.GetUserRequest{},
			},
			want: func(a args, f fields) (*authv1.GetUserResponse, error) {
				return nil, NewApiError(codes.InvalidArgument, "access_token not provided").
					WithLocalizedMessage(view.ErrStatusTextAccessTokenNotProvided)
			},
		},
		{
			name: "get user info err",
			args: args{
				ctx: ctx,
				request: &authv1.GetUserRequest{
					AccessToken: "testAccessToken",
				},
			},
			want: func(a args, f fields) (*authv1.GetUserResponse, error) {
				f.interactor.EXPECT().GetUserInfo(a.ctx, a.request.GetAccessToken()).Return(nil, testErr)
				return nil, testErr
			},
		},
		{
			name: "correct",
			args: args{
				ctx: ctx,
				request: &authv1.GetUserRequest{
					AccessToken: "testAccessToken",
				},
			},
			want: func(a args, f fields) (*authv1.GetUserResponse, error) {
				testUserInfo := &entity.UserInfo{}
				f.interactor.EXPECT().GetUserInfo(a.ctx, a.request.GetAccessToken()).Return(testUserInfo, nil)
				testUserInfoPb := &authv1.UserInfo{}
				f.presenter.EXPECT().UserInfoToPb(testUserInfo).Return(testUserInfoPb)
				return &authv1.GetUserResponse{UserInfo: testUserInfoPb}, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				interactor: NewMockAuthInteractor(ctrl),
				presenter:  NewMockAuthPresenter(ctrl),
			}

			want, wantErr := tt.want(tt.args, f)

			a := NewAuthServer(f.interactor, f.presenter)
			got, err := a.GetUser(tt.args.ctx, tt.args.request)

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

func Test_authServer_GetEmployees(t *testing.T) {
	type fields struct {
		interactor *MockAuthInteractor
		presenter  *MockAuthPresenter
	}
	type args struct {
		ctx     context.Context
		request *authv1.GetEmployeesRequest
	}

	ctx := context.TODO()
	testErr := errors.New("some test error")

	tests := []struct {
		name   string
		fields fields
		args   args
		want   func(a args, f fields) (*authv1.GetEmployeesResponse, error)
	}{
		{
			name: "invalid key type err",
			args: args{
				ctx:     ctx,
				request: &authv1.GetEmployeesRequest{},
			},
			want: func(a args, f fields) (*authv1.GetEmployeesResponse, error) {
				return nil, diterrors.NewApiError(codes.InvalidArgument, "invalid key type").
					WithLocalizedMessage("Передан некорректный ключ")
			},
		},
		{
			name: "get employees err",
			args: args{
				ctx: ctx,
				request: &authv1.GetEmployeesRequest{
					Key: &authv1.GetEmployeesRequest_CloudId{
						CloudId: "testCloudID",
					},
				},
			},
			want: func(a args, f fields) (*authv1.GetEmployeesResponse, error) {
				testParams := entity.EmployeeGetParams{
					Key:     entity.EmployeeGetKeyCloudID,
					CloudID: a.request.GetCloudId(),
				}
				f.interactor.EXPECT().GetEmployees(a.ctx, testParams).Return(nil, testErr)
				return nil, diterrors.NewApiError(codes.Internal, "auth service error", testErr).
					WithLocalizedMessage(view.ErrStatusTextServiceError)
			},
		},
		{
			name: "get employees validation err",
			args: args{
				ctx: ctx,
				request: &authv1.GetEmployeesRequest{
					Key: &authv1.GetEmployeesRequest_CloudId{
						CloudId: "testCloudID",
					},
				},
			},
			want: func(a args, f fields) (*authv1.GetEmployeesResponse, error) {
				testParams := entity.EmployeeGetParams{
					Key:     entity.EmployeeGetKeyCloudID,
					CloudID: a.request.GetCloudId(),
				}
				testValidationErr := diterrors.NewValidationError(testErr)
				f.interactor.EXPECT().GetEmployees(a.ctx, testParams).Return(nil, testValidationErr)
				return nil, diterrors.NewApiError(codes.InvalidArgument, "invalid request", testValidationErr).
					WithLocalizedMessage("Неверные параметры запроса")
			},
		},
		{
			name: "get employees not found err",
			args: args{
				ctx: ctx,
				request: &authv1.GetEmployeesRequest{
					Key: &authv1.GetEmployeesRequest_CloudId{
						CloudId: "testCloudID",
					},
				},
			},
			want: func(a args, f fields) (*authv1.GetEmployeesResponse, error) {
				testParams := entity.EmployeeGetParams{
					Key:     entity.EmployeeGetKeyCloudID,
					CloudID: a.request.GetCloudId(),
				}
				f.interactor.EXPECT().GetEmployees(a.ctx, testParams).Return(nil, diterrors.ErrNotFound)
				return nil, diterrors.NewApiError(codes.NotFound, "employee not found", diterrors.ErrNotFound).
					WithLocalizedMessage("Сотрудник не найден")
			},
		},
		{
			name: "correct",
			args: args{
				ctx: ctx,
				request: &authv1.GetEmployeesRequest{
					Key: &authv1.GetEmployeesRequest_Email{
						Email: "testEmail",
					},
				},
			},
			want: func(a args, f fields) (*authv1.GetEmployeesResponse, error) {
				testParams := entity.EmployeeGetParams{
					Key:   entity.EmployeeGetKeyEmail,
					Email: a.request.GetEmail(),
				}
				testEmployees := make([]entity.EmployeeInfo, 0)
				f.interactor.EXPECT().GetEmployees(a.ctx, testParams).Return(testEmployees, nil)
				testEmployeesPb := make([]*authv1.Employee, 0)
				f.presenter.EXPECT().EmployeesToPb(testEmployees).Return(testEmployeesPb)
				return &authv1.GetEmployeesResponse{Employees: testEmployeesPb}, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				interactor: NewMockAuthInteractor(ctrl),
				presenter:  NewMockAuthPresenter(ctrl),
			}

			want, wantErr := tt.want(tt.args, f)

			a := NewAuthServer(f.interactor, f.presenter)
			got, err := a.GetEmployees(tt.args.ctx, tt.args.request)

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

func Test_authServer_Validate(t *testing.T) {
	type fields struct {
		interactor *MockAuthInteractor
		presenter  *MockAuthPresenter
	}
	type args struct {
		ctx     context.Context
		request *authv1.ValidateRequest
	}

	ctx := context.TODO()
	testErr := errors.New("some test error")

	tests := []struct {
		name   string
		fields fields
		args   args
		want   func(a args, f fields) (*authv1.ValidateResponse, error)
	}{
		{
			name: "empty token err",
			args: args{
				ctx:     ctx,
				request: &authv1.ValidateRequest{},
			},
			want: func(a args, f fields) (*authv1.ValidateResponse, error) {
				return nil, NewApiError(codes.InvalidArgument, "access_token not provided").
					WithLocalizedMessage(view.ErrStatusTextAccessTokenNotProvided)
			},
		},
		{
			name: "token expired err",
			args: args{
				ctx: ctx,
				request: &authv1.ValidateRequest{
					AccessToken: "testAccessToken",
				},
			},
			want: func(a args, f fields) (*authv1.ValidateResponse, error) {
				f.interactor.EXPECT().IsValidToken(a.ctx, a.request.GetAccessToken()).Return(nil, usecase.ErrTokenIsExpire)
				return nil, NewApiError(codes.PermissionDenied, "token expired").
					WithLocalizedMessage(view.ErrStatusTextAccessTokenExpired)
			},
		},
		{
			name: "token expired err",
			args: args{
				ctx: ctx,
				request: &authv1.ValidateRequest{
					AccessToken: "testAccessToken",
				},
			},
			want: func(a args, f fields) (*authv1.ValidateResponse, error) {
				f.interactor.EXPECT().IsValidToken(a.ctx, a.request.GetAccessToken()).Return(nil, testErr)
				return nil, NewApiError(codes.Internal, "auth service error").
					WithLocalizedMessage(view.ErrStatusTextServiceError)
			},
		},
		{
			name: "correct",
			args: args{
				ctx: ctx,
				request: &authv1.ValidateRequest{
					AccessToken: "testAccessToken",
				},
			},
			want: func(a args, f fields) (*authv1.ValidateResponse, error) {
				testTokenInfo := &entity.TokenInfo{}
				f.interactor.EXPECT().IsValidToken(a.ctx, a.request.GetAccessToken()).Return(testTokenInfo, nil)
				testTokenInfoPb := &authv1.TokenInfo{}
				f.presenter.EXPECT().TokenInfoToPb(testTokenInfo).Return(testTokenInfoPb)
				return &authv1.ValidateResponse{Info: testTokenInfoPb}, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				interactor: NewMockAuthInteractor(ctrl),
				presenter:  NewMockAuthPresenter(ctrl),
			}

			want, wantErr := tt.want(tt.args, f)

			a := NewAuthServer(f.interactor, f.presenter)
			got, err := a.Validate(tt.args.ctx, tt.args.request)

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
