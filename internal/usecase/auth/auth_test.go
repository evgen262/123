package auth

import (
	"context"
	"errors"
	"fmt"
	"net"
	"testing"

	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity"
	entityAuth "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/auth"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/repositories"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/usecase"
)

func Test_authUseCase_GetAuthURL(t *testing.T) {
	type fields struct {
		repo   *MockRepository
		logger *ditzap.MockLogger
	}
	type args struct {
		ctx         context.Context
		callbackURI string
	}

	ctx := context.TODO()
	testErr := errors.New("test error")

	tests := []struct {
		name string
		args args
		want func(a args, f fields) (string, error)
	}{
		{
			name: "get url localized err",
			args: args{
				ctx:         ctx,
				callbackURI: "",
			},
			want: func(a args, f fields) (string, error) {
				testLocalizedErr := diterrors.NewLocalizedError("", testErr)
				f.repo.EXPECT().GetRedirectURL(a.ctx, a.callbackURI).Return("", testLocalizedErr)
				f.logger.EXPECT().Error("cant get auth url", zap.Error(errors.Unwrap(testLocalizedErr)))
				return "", fmt.Errorf("cant get auth url: %w", testLocalizedErr)
			},
		},
		{
			name: "get url err",
			args: args{
				ctx:         ctx,
				callbackURI: "",
			},
			want: func(a args, f fields) (string, error) {
				f.repo.EXPECT().GetRedirectURL(a.ctx, a.callbackURI).Return("", testErr)
				f.logger.EXPECT().Error("cant get auth url", zap.Error(testErr))
				return "", fmt.Errorf("cant get auth url: %w", testErr)
			},
		},
		{
			name: "correct",
			args: args{
				ctx:         ctx,
				callbackURI: "",
			},
			want: func(a args, f fields) (string, error) {
				f.repo.EXPECT().GetRedirectURL(a.ctx, a.callbackURI).Return("", nil)
				return "", nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				repo:   NewMockRepository(ctrl),
				logger: ditzap.NewMockLogger(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			au := NewAuthUseCase(f.repo, f.logger, nil)
			got, err := au.GetAuthURL(tt.args.ctx, tt.args.callbackURI)
			if wantErr != nil {
				assert.EqualError(t, err, wantErr.Error())
				assert.Equal(t, "", got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, want, got)
			}
		})
	}
}

func Test_authUseCase_Auth(t *testing.T) {
	type fields struct {
		repo   *MockRepository
		logger *ditzap.MockLogger
	}
	type args struct {
		ctx         context.Context
		code        string
		state       string
		callbackURI string
		device      *entityAuth.Device
		clientIP    net.IP
	}

	ctx := context.TODO()
	testErr := errors.New("test error")
	testDevice := &entityAuth.Device{
		UserAgent: "testUserAgent",
	}

	tests := []struct {
		name string
		args args
		want func(a args, f fields) (*entityAuth.Auth, error)
	}{
		{
			name: "auth not found err",
			args: args{
				ctx:         ctx,
				code:        "",
				state:       "",
				callbackURI: "",
			},
			want: func(a args, f fields) (*entityAuth.Auth, error) {
				f.repo.EXPECT().Auth(a.ctx, a.code, a.state, a.callbackURI).Return(nil, diterrors.ErrNotFound)
				return nil, ErrEmployeesNotFound
			},
		},
		{
			name: "auth localized err",
			args: args{
				ctx:         ctx,
				code:        "",
				state:       "",
				callbackURI: "",
			},
			want: func(a args, f fields) (*entityAuth.Auth, error) {
				testLocalizedErr := diterrors.NewLocalizedError("", testErr)
				f.repo.EXPECT().Auth(a.ctx, a.code, a.state, a.callbackURI).Return(nil, testLocalizedErr)
				f.logger.EXPECT().Error("cant authenticate user",
					zap.String("code", a.code),
					zap.String("state", a.state),
					zap.Error(testLocalizedErr),
				)
				return nil, fmt.Errorf("cant authenticate user: %w", testLocalizedErr)
			},
		},
		{
			name: "auth validation err",
			args: args{
				ctx:         ctx,
				code:        "",
				state:       "",
				callbackURI: "",
			},
			want: func(a args, f fields) (*entityAuth.Auth, error) {
				testValidationErr := diterrors.NewValidationError(testErr)
				f.repo.EXPECT().Auth(a.ctx, a.code, a.state, a.callbackURI).Return(nil, testValidationErr)
				f.logger.EXPECT().Debug("cant authenticate user",
					zap.String("code", a.code),
					zap.String("state", a.state),
					zap.Error(testValidationErr),
				)
				return nil, fmt.Errorf("cant authenticate user: %w", testValidationErr)
			},
		},
		{
			name: "nil user err",
			args: args{
				ctx:         ctx,
				code:        "",
				state:       "",
				callbackURI: "",
			},
			want: func(a args, f fields) (*entityAuth.Auth, error) {
				testAuthSudir := &entityAuth.AuthSudir{}
				f.repo.EXPECT().Auth(a.ctx, a.code, a.state, a.callbackURI).Return(testAuthSudir, nil)
				return nil, diterrors.ErrPermissionDenied
			},
		},
		{
			name: "access list permission denied err",
			args: args{
				ctx: ctx,
			},
			want: func(a args, f fields) (*entityAuth.Auth, error) {
				testAuthSudir := &entityAuth.AuthSudir{User: &entityAuth.UserSudir{Email: "test2@example.com"}}
				f.repo.EXPECT().Auth(a.ctx, a.code, a.state, a.callbackURI).Return(testAuthSudir, nil)
				return nil, ErrUserAccessDenied
			},
		},
		/* TODO: Отдать этот корнер-кейс на пересмотр аналитику
		{
			name: "no cloudID err",
			args: args{
				ctx:         ctx,
				code:        "",
				state:       "",
				callbackURI: "",
			},
			want: func(a args, f fields) (*entityAuth.Auth, error) {
				testAuthSudir := &entityAuth.AuthSudir{
					User: &entityAuth.UserSudir{
						Email: "test@example.com",
					},
				}
				f.repo.EXPECT().Auth(a.ctx, a.code, a.state, a.callbackURI).Return(testAuthSudir, nil)
				return nil, ErrSUDIRNoCloudID
			},
		},
		*/
		{
			name: "empty portals err",
			args: args{
				ctx:         ctx,
				code:        "",
				state:       "",
				callbackURI: "",
			},
			want: func(a args, f fields) (*entityAuth.Auth, error) {
				testAuthSudir := &entityAuth.AuthSudir{
					User: &entityAuth.UserSudir{
						Login:   "testUserLogin",
						CloudID: "testCloudID",
						Email:   "test@example.com",
						FIO:     "testUserFIO",
					},
				}
				f.repo.EXPECT().Auth(a.ctx, a.code, a.state, a.callbackURI).Return(testAuthSudir, nil)
				f.logger.EXPECT().Error("нет порталов связанных с сотрудником",
					entityAuth.LogModuleUE,
					entityAuth.LogCode("UE_030"),
					zap.String("login", testAuthSudir.GetUser().Login),
					zap.String("email", testAuthSudir.GetUser().Email),
					zap.String("user", testAuthSudir.GetUser().FIO),
					zap.Error(nil),
				)
				return nil, ErrPortalsNotFound
			},
		},
		{
			name: "auth portal not found err",
			args: args{
				ctx:         ctx,
				code:        "",
				state:       "",
				callbackURI: "",
			},
			want: func(a args, f fields) (*entityAuth.Auth, error) {
				testAuthSudir := &entityAuth.AuthSudir{
					User: &entityAuth.UserSudir{
						Login:   "testUserLogin",
						CloudID: "testCloudID",
						Email:   "test@example.com",
						FIO:     "testUserFIO",
						SNILS:   "000-000-000-00",
						Portals: []*entityAuth.Portal{{
							Name: "testName", URL: "testURL"}},
					},
				}
				f.repo.EXPECT().Auth(a.ctx, a.code, a.state, a.callbackURI).Return(testAuthSudir, nil)
				testAuthParams := entityAuth.AuthPortalParams{
					PortalURL: "testURL",
					User: entityAuth.User1C{
						CloudID: "testCloudID",
						SNILS:   "000-000-000-00",
						Email:   "test@example.com",
					},
				}
				f.repo.EXPECT().AuthPortal(a.ctx, testAuthParams).Return(nil, diterrors.ErrNotFound)
				return nil, diterrors.ErrNotFound
			},
		},
		{
			name: "auth portal localized err",
			args: args{
				ctx:         ctx,
				code:        "",
				state:       "",
				callbackURI: "",
			},
			want: func(a args, f fields) (*entityAuth.Auth, error) {
				testAuthSudir := &entityAuth.AuthSudir{
					User: &entityAuth.UserSudir{
						Login:   "testUserLogin",
						CloudID: "testCloudID",
						Email:   "test@example.com",
						FIO:     "testUserFIO",
						SNILS:   "000-000-000-00",
						Portals: []*entityAuth.Portal{{
							Name: "testName", URL: "testURL"}},
					},
				}
				f.repo.EXPECT().Auth(a.ctx, a.code, a.state, a.callbackURI).Return(testAuthSudir, nil)
				testAuthParams := entityAuth.AuthPortalParams{
					PortalURL: "testURL",
					User: entityAuth.User1C{
						CloudID: "testCloudID",
						SNILS:   "000-000-000-00",
						Email:   "test@example.com",
					},
				}
				testLocalizedErr := diterrors.NewLocalizedError("", testErr)
				f.repo.EXPECT().AuthPortal(a.ctx, testAuthParams).Return(nil, testLocalizedErr)
				f.logger.EXPECT().Error("cant authenticate user into portal",
					zap.String("user_id", testAuthSudir.GetUser().CloudID),
					zap.String("portal_url", testAuthSudir.GetUser().Portals[0].URL),
					zap.Error(testLocalizedErr),
				)
				return nil, fmt.Errorf("cant authenticate user: %w", testLocalizedErr)
			},
		},
		{
			name: "auth portal validation err",
			args: args{
				ctx:         ctx,
				code:        "",
				state:       "",
				callbackURI: "",
			},
			want: func(a args, f fields) (*entityAuth.Auth, error) {
				testAuthSudir := &entityAuth.AuthSudir{
					User: &entityAuth.UserSudir{
						Login:   "testUserLogin",
						CloudID: "testCloudID",
						Email:   "test@example.com",
						FIO:     "testUserFIO",
						SNILS:   "000-000-000-00",
						Portals: []*entityAuth.Portal{{
							Name: "testName", URL: "testURL"}},
					},
				}
				f.repo.EXPECT().Auth(a.ctx, a.code, a.state, a.callbackURI).Return(testAuthSudir, nil)
				testAuthParams := entityAuth.AuthPortalParams{
					PortalURL: "testURL",
					User: entityAuth.User1C{
						CloudID: "testCloudID",
						SNILS:   "000-000-000-00",
						Email:   "test@example.com",
					},
				}
				testValidationErr := diterrors.NewValidationError(testErr)
				f.repo.EXPECT().AuthPortal(a.ctx, testAuthParams).Return(nil, testValidationErr)
				f.logger.EXPECT().Debug("cant authenticate user into portal",
					zap.String("user_id", testAuthSudir.GetUser().CloudID),
					zap.String("portal_url", testAuthSudir.GetUser().Portals[0].URL),
					zap.Error(testValidationErr),
				)
				return nil, fmt.Errorf("cant authenticate user: %w", testValidationErr)
			},
		},
		{
			name: "invalid device err",
			args: args{
				ctx:         ctx,
				code:        "",
				state:       "",
				callbackURI: "",
			},
			want: func(a args, f fields) (*entityAuth.Auth, error) {
				testAuthSudir := &entityAuth.AuthSudir{
					User: &entityAuth.UserSudir{
						Login:   "testUserLogin",
						CloudID: "testCloudID",
						Email:   "test@example.com",
						FIO:     "testUserFIO",
						SNILS:   "000-000-000-00",
						Portals: []*entityAuth.Portal{{
							Name: "testName", URL: "testURL"}},
					},
				}
				f.repo.EXPECT().Auth(a.ctx, a.code, a.state, a.callbackURI).Return(testAuthSudir, nil)
				testAuthParams := entityAuth.AuthPortalParams{
					PortalURL: "testURL",
					User: entityAuth.User1C{
						CloudID: "testCloudID",
						SNILS:   "000-000-000-00",
						Email:   "test@example.com",
					},
				}
				testAuth1C := &entityAuth.Auth1C{PortalSession: "testPortalSession"}
				f.repo.EXPECT().AuthPortal(a.ctx, testAuthParams).Return(testAuth1C, nil)
				return nil, ErrInvalidDevice
			},
		},
		{
			name: "create session no user portals err",
			args: args{
				ctx:         entity.WithDevice(ctx, testDevice),
				code:        "",
				state:       "",
				callbackURI: "",
			},
			want: func(a args, f fields) (*entityAuth.Auth, error) {
				testAuthSudir := &entityAuth.AuthSudir{
					User: &entityAuth.UserSudir{
						Login:   "testUserLogin",
						CloudID: "testCloudID",
						Email:   "test@example.com",
						FIO:     "testUserFIO",
						SNILS:   "000-000-000-00",
						Portals: []*entityAuth.Portal{{
							Name: "testName", URL: "testURL"}},
					},
				}
				f.repo.EXPECT().Auth(a.ctx, a.code, a.state, a.callbackURI).Return(testAuthSudir, nil)
				testAuthParams := entityAuth.AuthPortalParams{
					PortalURL: "testURL",
					User: entityAuth.User1C{
						CloudID: "testCloudID",
						SNILS:   "000-000-000-00",
						Email:   "test@example.com",
					},
				}
				testAuth1C := &entityAuth.Auth1C{PortalSession: "testPortalSession"}
				f.repo.EXPECT().AuthPortal(a.ctx, testAuthParams).Return(testAuth1C, nil)
				f.repo.EXPECT().CreateSession(a.ctx, testAuthSudir.GetUser(), nil, testDevice, testAuth1C).Return(entityAuth.TokensPair{}, repositories.ErrNoUserPortals)
				return nil, ErrPortalsNotFound
			},
		},
		{
			name: "create session err",
			args: args{
				ctx:         entity.WithDevice(ctx, testDevice),
				code:        "",
				state:       "",
				callbackURI: "",
			},
			want: func(a args, f fields) (*entityAuth.Auth, error) {
				testAuthSudir := &entityAuth.AuthSudir{
					User: &entityAuth.UserSudir{
						Login:   "testUserLogin",
						CloudID: "testCloudID",
						Email:   "test@example.com",
						FIO:     "testUserFIO",
						SNILS:   "000-000-000-00",
						Portals: []*entityAuth.Portal{{
							Name: "testName", URL: "testURL"}},
					},
				}
				f.repo.EXPECT().Auth(a.ctx, a.code, a.state, a.callbackURI).Return(testAuthSudir, nil)
				testAuthParams := entityAuth.AuthPortalParams{
					PortalURL: "testURL",
					User: entityAuth.User1C{
						CloudID: "testCloudID",
						SNILS:   "000-000-000-00",
						Email:   "test@example.com",
					},
				}
				testAuth1C := &entityAuth.Auth1C{PortalSession: "testPortalSession"}
				f.repo.EXPECT().AuthPortal(a.ctx, testAuthParams).Return(testAuth1C, nil)
				f.repo.EXPECT().CreateSession(a.ctx, testAuthSudir.GetUser(), nil, testDevice, testAuth1C).Return(entityAuth.TokensPair{}, testErr)
				f.logger.EXPECT().Error("can't create session",
					zap.String("user_id", testAuthSudir.GetUser().CloudID),
					zap.Error(testErr),
				)
				return nil, fmt.Errorf("can't create session: %w", testErr)
			},
		},
		{
			name: "correct",
			args: args{
				ctx:         entity.WithDevice(ctx, testDevice),
				code:        "",
				state:       "",
				callbackURI: "",
				device: &entityAuth.Device{
					UserAgent: "testUserAgent",
				},
			},
			want: func(a args, f fields) (*entityAuth.Auth, error) {
				testAuthSudir := &entityAuth.AuthSudir{
					User: &entityAuth.UserSudir{
						Login:   "testUserLogin",
						CloudID: "testCloudID",
						Email:   "test@example.com",
						FIO:     "testUserFIO",
						SNILS:   "000-000-000-00",
						Portals: []*entityAuth.Portal{{
							Name: "testName", URL: "testURL"}},
					},
				}
				f.repo.EXPECT().Auth(a.ctx, a.code, a.state, a.callbackURI).Return(testAuthSudir, nil)
				testAuthParams := entityAuth.AuthPortalParams{
					PortalURL: "testURL",
					User: entityAuth.User1C{
						CloudID: "testCloudID",
						SNILS:   "000-000-000-00",
						Email:   "test@example.com",
					},
				}
				testAuth1C := &entityAuth.Auth1C{PortalSession: "testPortalSession"}
				f.repo.EXPECT().AuthPortal(a.ctx, testAuthParams).Return(testAuth1C, nil)
				testTokensPair := entityAuth.TokensPair{AccessToken: entityAuth.Token{Value: "testToken"}, RefreshToken: entityAuth.Token{Value: "testToken"}}
				f.repo.EXPECT().CreateSession(a.ctx, testAuthSudir.GetUser(), nil, testDevice, testAuth1C).Return(testTokensPair, nil)
				return &entityAuth.Auth{
					JWTToken:      testTokensPair.AccessToken,
					RefreshToken:  testTokensPair.RefreshToken,
					PortalSession: testAuth1C.PortalSession,
					Portals:       testAuthSudir.GetUser().Portals,
				}, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				repo:   NewMockRepository(ctrl),
				logger: ditzap.NewMockLogger(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			accessList := entityAuth.AccessList{"test@example.com"}
			au := NewAuthUseCase(f.repo, f.logger, accessList)
			got, err := au.Auth(tt.args.ctx, tt.args.code, tt.args.state, tt.args.callbackURI)
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

func Test_authUseCase_GetSession(t *testing.T) {
	type fields struct {
		repo   *MockRepository
		logger *ditzap.MockLogger
	}
	type args struct {
		ctx         context.Context
		accessToken string
	}

	ctx := context.TODO()
	testErr := errors.New("test error")

	tests := []struct {
		name string
		args args
		want func(a args, f fields) (*entityAuth.Session, error)
	}{
		{
			name: "validation err",
			args: args{
				ctx:         ctx,
				accessToken: "testAccessToken",
			},
			want: func(a args, f fields) (*entityAuth.Session, error) {
				testValidationErr := diterrors.NewValidationError(testErr)
				f.repo.EXPECT().GetSession(a.ctx, a.accessToken).Return(nil, testValidationErr)
				return nil, testValidationErr
			},
		},
		{
			name: "failed precondition err",
			args: args{
				ctx:         ctx,
				accessToken: "testAccessToken",
			},
			want: func(a args, f fields) (*entityAuth.Session, error) {
				f.repo.EXPECT().GetSession(a.ctx, a.accessToken).Return(nil, diterrors.ErrFailedPrecondition)
				return nil, diterrors.ErrFailedPrecondition
			},
		},
		{
			name: "not found err",
			args: args{
				ctx:         ctx,
				accessToken: "testAccessToken",
			},
			want: func(a args, f fields) (*entityAuth.Session, error) {
				f.repo.EXPECT().GetSession(a.ctx, a.accessToken).Return(nil, diterrors.ErrNotFound)
				return nil, diterrors.ErrNotFound
			},
		},
		{
			name: "unauthenticated err",
			args: args{
				ctx:         ctx,
				accessToken: "testAccessToken",
			},
			want: func(a args, f fields) (*entityAuth.Session, error) {
				f.repo.EXPECT().GetSession(a.ctx, a.accessToken).Return(nil, diterrors.ErrUnauthenticated)
				return nil, diterrors.ErrUnauthenticated
			},
		},
		{
			name: "internal err",
			args: args{
				ctx:         ctx,
				accessToken: "testAccessToken",
			},
			want: func(a args, f fields) (*entityAuth.Session, error) {
				f.repo.EXPECT().GetSession(a.ctx, a.accessToken).Return(nil, testErr)
				f.logger.EXPECT().Warn("can't get session in repository",
					ditzap.JWTField("access_token", a.accessToken),
					zap.Error(testErr))
				return nil, fmt.Errorf("can't get session in repository: %w", testErr)
			},
		},
		{
			name: "correct",
			args: args{
				ctx:         ctx,
				accessToken: "testAccessToken",
			},
			want: func(a args, f fields) (*entityAuth.Session, error) {
				testSession := &entityAuth.Session{}
				f.repo.EXPECT().GetSession(a.ctx, a.accessToken).Return(testSession, nil)
				return testSession, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				repo:   NewMockRepository(ctrl),
				logger: ditzap.NewMockLogger(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			au := NewAuthUseCase(f.repo, f.logger, nil)
			got, err := au.GetSession(tt.args.ctx, tt.args.accessToken)
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

func Test_authUseCase_Logout(t *testing.T) {
	type fields struct {
		repo   *MockRepository
		logger *ditzap.MockLogger
	}
	type args struct {
		ctx                       context.Context
		accessToken, refreshToken string
	}

	ctx := context.TODO()
	testErr := errors.New("test error")
	testSession := &entityAuth.Session{
		User: &entityAuth.User{
			Login: "testLogin",
		},
		ActivePortal: &entityAuth.ActivePortal{
			Portal: entityAuth.Portal{
				Name: "testPortalName",
			},
		},
	}

	tests := []struct {
		name string
		args args
		want func(a args, f fields) error
	}{
		{
			name: "nil session err",
			args: args{
				ctx:         entity.WithSession(ctx, nil),
				accessToken: "testAccessToken",
			},
			want: func(a args, f fields) error {
				f.logger.EXPECT().Error("can't get session from context", zap.Error(errors.New("session not found")))
				return usecase.ErrGetSessionFromContext
			},
		},
		{
			name: "invalid err",
			args: args{
				ctx:         entity.WithSession(ctx, testSession),
				accessToken: "testAccessToken",
			},
			want: func(a args, f fields) error {
				f.repo.EXPECT().Logout(a.ctx, testSession, a.accessToken, a.refreshToken).Return(repositories.ErrNilSession)
				f.logger.EXPECT().Warn("invalid session",
					zap.String("session_id", testSession.GetID().String()),
					zap.String("login", testSession.GetUser().GetLogin()),
					zap.String("portal_url", testSession.GetActivePortal().GetPortal().URL),
					zap.Error(repositories.ErrNilSession),
				)
				return ErrInvalidSession
			},
		},
		{
			name: "session validation err",
			args: args{
				ctx:         entity.WithSession(ctx, testSession),
				accessToken: "testAccessToken",
			},
			want: func(a args, f fields) error {
				testValidationErr := diterrors.NewValidationError(testErr)
				f.repo.EXPECT().Logout(a.ctx, testSession, a.accessToken, a.refreshToken).Return(testValidationErr)
				f.logger.EXPECT().Warn("cant Logout",
					zap.String("session_id", testSession.GetID().String()),
					zap.String("login", testSession.GetUser().GetLogin()),
					zap.String("portal_url", testSession.GetActivePortal().GetPortal().URL),
					zap.Error(testValidationErr),
				)
				return testValidationErr
			},
		},
		{
			name: "session not found err",
			args: args{
				ctx:         entity.WithSession(ctx, testSession),
				accessToken: "testAccessToken",
			},
			want: func(a args, f fields) error {
				f.repo.EXPECT().Logout(a.ctx, testSession, a.accessToken, a.refreshToken).Return(diterrors.ErrNotFound)
				f.logger.EXPECT().Warn("cant Logout",
					zap.String("session_id", testSession.GetID().String()),
					zap.String("login", testSession.GetUser().GetLogin()),
					zap.String("portal_url", testSession.GetActivePortal().GetPortal().URL),
					zap.Error(diterrors.ErrNotFound),
				)
				return diterrors.ErrNotFound
			},
		},
		{
			name: "logout err",
			args: args{
				ctx:         entity.WithSession(ctx, testSession),
				accessToken: "testAccessToken",
			},
			want: func(a args, f fields) error {
				f.repo.EXPECT().Logout(a.ctx, testSession, a.accessToken, a.refreshToken).Return(testErr)
				f.logger.EXPECT().Warn("can't logout in repository",
					zap.String("session_id", testSession.GetID().String()),
					zap.String("login", testSession.GetUser().GetLogin()),
					zap.String("portal_url", testSession.GetActivePortal().GetPortal().URL),
					zap.Error(testErr),
				)
				return fmt.Errorf("can't logout in repository: %w", testErr)
			},
		},
		{
			name: "correct",
			args: args{
				ctx:         entity.WithSession(ctx, testSession),
				accessToken: "testAccessToken",
			},
			want: func(a args, f fields) error {
				f.repo.EXPECT().Logout(a.ctx, testSession, a.accessToken, a.refreshToken).Return(nil)
				return nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				repo:   NewMockRepository(ctrl),
				logger: ditzap.NewMockLogger(ctrl),
			}
			wantErr := tt.want(tt.args, f)
			au := NewAuthUseCase(f.repo, f.logger, nil)
			err := au.Logout(tt.args.ctx, tt.args.accessToken, tt.args.refreshToken)
			if wantErr != nil {
				assert.EqualError(t, err, wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_authUseCase_ChangePortal(t *testing.T) {
	type fields struct {
		repo   *MockRepository
		logger *ditzap.MockLogger
	}
	type args struct {
		ctx      context.Context
		portalID int
	}

	ctx := context.TODO()
	testErr := errors.New("test error")
	testSession := &entityAuth.Session{
		User: &entityAuth.User{
			Login: "testLogin",
		},
		ActivePortal: &entityAuth.ActivePortal{
			Portal: entityAuth.Portal{
				Name: "testPortalName",
			},
		},
	}

	tests := []struct {
		name string
		args args
		want func(a args, f fields) ([]*entityAuth.Portal, string, error)
	}{
		{
			name: "nil session err",
			args: args{
				ctx:      entity.WithSession(ctx, nil),
				portalID: 1,
			},
			want: func(a args, f fields) ([]*entityAuth.Portal, string, error) {
				f.logger.EXPECT().Error("can't get session from context", zap.Error(errors.New("session not found")))
				return nil, "", usecase.ErrGetSessionFromContext
			},
		},
		{
			name: "change portal err",
			args: args{
				ctx:      entity.WithSession(ctx, testSession),
				portalID: 1,
			},
			want: func(a args, f fields) ([]*entityAuth.Portal, string, error) {
				f.repo.EXPECT().ChangePortal(a.ctx, a.portalID, testSession).Return(nil, "", testErr)
				return nil, "", fmt.Errorf("can't change portal in repository: %w", testErr)
			},
		},
		{
			name: "portals not found err",
			args: args{
				ctx:      entity.WithSession(ctx, testSession),
				portalID: 1,
			},
			want: func(a args, f fields) ([]*entityAuth.Portal, string, error) {
				f.repo.EXPECT().ChangePortal(a.ctx, a.portalID, testSession).Return(nil, "", nil)
				return nil, "", ErrPortalsNotFound
			},
		},
		{
			name: "portals not found err",
			args: args{
				ctx:      entity.WithSession(ctx, testSession),
				portalID: 1,
			},
			want: func(a args, f fields) ([]*entityAuth.Portal, string, error) {
				testPortals := []*entityAuth.Portal{{ID: 2}}
				f.repo.EXPECT().ChangePortal(a.ctx, a.portalID, testSession).Return(testPortals, "", nil)
				return nil, "", ErrUnavailablePortal
			},
		},
		{
			name: "empty portal url err",
			args: args{
				ctx:      entity.WithSession(ctx, testSession),
				portalID: 1,
			},
			want: func(a args, f fields) ([]*entityAuth.Portal, string, error) {
				testPortals := []*entityAuth.Portal{{ID: 1}}
				f.repo.EXPECT().ChangePortal(a.ctx, a.portalID, testSession).Return(testPortals, "", nil)
				return nil, "", ErrEmptyPortalURL
			},
		},
		{
			name: "correct",
			args: args{
				ctx:      entity.WithSession(ctx, testSession),
				portalID: 1,
			},
			want: func(a args, f fields) ([]*entityAuth.Portal, string, error) {
				testPortalSID := "testPortalSID"
				testPortals := []*entityAuth.Portal{{ID: 1, URL: "testPortalURL"}}
				f.repo.EXPECT().ChangePortal(a.ctx, a.portalID, testSession).Return(testPortals, testPortalSID, nil)
				return testPortals, testPortalSID, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				repo:   NewMockRepository(ctrl),
				logger: ditzap.NewMockLogger(ctrl),
			}
			wantP, wantS, wantErr := tt.want(tt.args, f)
			au := NewAuthUseCase(f.repo, f.logger, nil)
			gotP, gotS, err := au.ChangePortal(tt.args.ctx, tt.args.portalID)
			if wantErr != nil {
				assert.EqualError(t, err, wantErr.Error())
				assert.Nil(t, gotP)
				assert.Empty(t, gotS)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, wantP, gotP)
				assert.Equal(t, wantS, gotS)
			}
		})
	}
}

func Test_authUseCase_RefreshTokensPair(t *testing.T) {
	type fields struct {
		repo   *MockRepository
		logger *ditzap.MockLogger
	}
	type args struct {
		ctx                       context.Context
		accessToken, refreshToken string
	}

	ctx := context.TODO()
	testErr := errors.New("test error")

	tests := []struct {
		name string
		args args
		want func(a args, f fields) (*entityAuth.TokensPair, error)
	}{
		{
			name: "err",
			args: args{
				ctx:          ctx,
				accessToken:  "testAccessToken",
				refreshToken: "testRefreshToken",
			},
			want: func(a args, f fields) (*entityAuth.TokensPair, error) {
				f.repo.EXPECT().RefreshTokensPair(a.ctx, a.accessToken, a.refreshToken).Return(nil, testErr)
				return nil, fmt.Errorf("can't change portal in repository: %w", testErr)
			},
		},
		{
			name: "correct",
			args: args{
				ctx:          ctx,
				accessToken:  "testAccessToken",
				refreshToken: "testRefreshToken",
			},
			want: func(a args, f fields) (*entityAuth.TokensPair, error) {
				testTokenPair := &entityAuth.TokensPair{}
				f.repo.EXPECT().RefreshTokensPair(a.ctx, a.accessToken, a.refreshToken).Return(testTokenPair, nil)
				return testTokenPair, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				repo:   NewMockRepository(ctrl),
				logger: ditzap.NewMockLogger(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			au := NewAuthUseCase(f.repo, f.logger, nil)
			got, err := au.RefreshTokensPair(tt.args.ctx, tt.args.accessToken, tt.args.refreshToken)
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
