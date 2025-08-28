package auth

import (
	"context"
	"errors"
	"net"
	"net/url"
	"testing"
	"time"

	authv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/authfacade/auth/v1"
	authErrorsv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/authfacade/errors/v1"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	timeUtils "git.mos.ru/buch-cloud/moscow-team-2.0/build/time-utils.git"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	entityAuth "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/auth"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/repositories"
)

func Test_authRepository_GetRedirectURL(t *testing.T) {
	type fields struct {
		client *authv1.MockAuthAPIClient
		mapper *MockMapperAuth
		tu     *timeUtils.MockTimeUtils
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
			name: "client err",
			args: args{
				ctx:         ctx,
				callbackURI: "https://www.test.com/?testQuery=test",
			},
			want: func(a args, f fields) (string, error) {
				testRequest := &authv1.GetURLRequest{
					CallbackUrl: "/?testQuery=test",
				}
				f.client.EXPECT().GetURL(a.ctx, testRequest).Return(nil, testErr)
				return "", diterrors.GrpcErrorToError(testErr)
			},
		},
		{
			name: "correct",
			args: args{
				ctx:         ctx,
				callbackURI: "https://www.test.com/?testQuery=test",
			},
			want: func(a args, f fields) (string, error) {
				testRequest := &authv1.GetURLRequest{
					CallbackUrl: "/?testQuery=test",
				}
				testResponse := &authv1.GetURLResponse{}
				f.client.EXPECT().GetURL(a.ctx, testRequest).Return(testResponse, nil)
				return "", nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				client: authv1.NewMockAuthAPIClient(ctrl),
				mapper: NewMockMapperAuth(ctrl),
				tu:     timeUtils.NewMockTimeUtils(ctrl),
				logger: ditzap.NewMockLogger(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			ar := NewAuthRepository(f.client, f.mapper, url.URL{}, "test-app-name", time.Second, time.Second, f.tu, f.logger)
			got, err := ar.GetRedirectURL(tt.args.ctx, tt.args.callbackURI)
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

func Test_authRepository_Auth(t *testing.T) {
	type fields struct {
		client *authv1.MockAuthAPIClient
		mapper *MockMapperAuth
		tu     *timeUtils.MockTimeUtils
		logger *ditzap.MockLogger
	}
	type args struct {
		ctx         context.Context
		code        string
		state       string
		callbackURI string
	}

	ctx := context.TODO()
	testErr := errors.New("test error")

	tests := []struct {
		name string
		args args
		want func(a args, f fields) (*entityAuth.AuthSudir, error)
	}{
		{
			name: "client err",
			args: args{
				ctx:         ctx,
				code:        "",
				state:       "",
				callbackURI: "",
			},
			want: func(a args, f fields) (*entityAuth.AuthSudir, error) {
				testRequest := &authv1.AuthRequest{}
				f.client.EXPECT().Auth(a.ctx, testRequest).Return(nil, testErr)
				return nil, diterrors.GrpcErrorToError(testErr)
			},
		},
		{
			name: "correct",
			args: args{
				ctx:         ctx,
				code:        "",
				state:       "",
				callbackURI: "",
			},
			want: func(a args, f fields) (*entityAuth.AuthSudir, error) {
				testRequest := &authv1.AuthRequest{}
				testResponse := &authv1.AuthResponse{}
				f.client.EXPECT().Auth(a.ctx, testRequest).Return(testResponse, nil)
				testUserSudir := &entityAuth.UserSudir{}
				f.mapper.EXPECT().UserToEntity(testResponse.GetUser()).Return(testUserSudir)
				return &entityAuth.AuthSudir{
					AccessToken:  testResponse.GetAccessToken(),
					RefreshToken: testResponse.GetRefreshToken(),
					User:         testUserSudir,
				}, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				client: authv1.NewMockAuthAPIClient(ctrl),
				mapper: NewMockMapperAuth(ctrl),
				tu:     timeUtils.NewMockTimeUtils(ctrl),
				logger: ditzap.NewMockLogger(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			ar := NewAuthRepository(f.client, f.mapper, url.URL{}, "test-app-name", time.Second, time.Second, f.tu, f.logger)
			got, err := ar.Auth(tt.args.ctx, tt.args.code, tt.args.state, tt.args.callbackURI)
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

func Test_authRepository_AuthPortal(t *testing.T) {
	type fields struct {
		client *authv1.MockAuthAPIClient
		mapper *MockMapperAuth
		tu     *timeUtils.MockTimeUtils
		logger *ditzap.MockLogger
	}
	type args struct {
		ctx    context.Context
		params entityAuth.AuthPortalParams
	}

	ctx := context.TODO()
	testErr := errors.New("test error")

	tests := []struct {
		name string
		args args
		want func(a args, f fields) (*entityAuth.Auth1C, error)
	}{
		{
			name: "client err",
			args: args{
				ctx:    ctx,
				params: entityAuth.AuthPortalParams{},
			},
			want: func(a args, f fields) (*entityAuth.Auth1C, error) {
				testRequest := &authv1.Auth1CRequest{
					User: &authv1.User1C{
						UserType: authv1.User1C_USER_TYPE_WEB,
					},
				}
				f.client.EXPECT().Auth1C(a.ctx, testRequest).Return(nil, testErr)
				return nil, diterrors.GrpcErrorToError(testErr)
			},
		},
		{
			name: "correct",
			args: args{
				ctx:    ctx,
				params: entityAuth.AuthPortalParams{},
			},
			want: func(a args, f fields) (*entityAuth.Auth1C, error) {
				testRequest := &authv1.Auth1CRequest{
					User: &authv1.User1C{
						UserType: authv1.User1C_USER_TYPE_WEB,
					},
				}
				testResponse := &authv1.Auth1CResponse{
					SessionId: "testSID",
				}
				f.client.EXPECT().Auth1C(a.ctx, testRequest).Return(testResponse, nil)
				return &entityAuth.Auth1C{
					PortalSession: testResponse.GetSessionId(),
					EmployeeID:    testResponse.GetEmployeeId(),
					PersonID:      testResponse.GetPersonId(),
				}, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				client: authv1.NewMockAuthAPIClient(ctrl),
				mapper: NewMockMapperAuth(ctrl),
				tu:     timeUtils.NewMockTimeUtils(ctrl),
				logger: ditzap.NewMockLogger(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			ar := NewAuthRepository(f.client, f.mapper, url.URL{}, "test-app-name", time.Second, time.Second, f.tu, f.logger)
			got, err := ar.AuthPortal(tt.args.ctx, tt.args.params)
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

func Test_authRepository_CreateSession(t *testing.T) {
	type fields struct {
		client *authv1.MockAuthAPIClient
		mapper *MockMapperAuth
		tu     *timeUtils.MockTimeUtils
		logger *ditzap.MockLogger
	}
	type args struct {
		ctx      context.Context
		user     *entityAuth.UserSudir
		clientIP net.IP
		device   *entityAuth.Device
		auth     *entityAuth.Auth1C
	}

	ctx := context.TODO()
	testErr := errors.New("test error")

	tests := []struct {
		name string
		args args
		want func(a args, f fields) (entityAuth.TokensPair, error)
	}{
		{
			name: "no user portals err",
			args: args{
				ctx:  ctx,
				user: &entityAuth.UserSudir{},
			},
			want: func(a args, f fields) (entityAuth.TokensPair, error) {
				return entityAuth.TokensPair{}, repositories.ErrNoUserPortals
			},
		},
		{
			name: "client err",
			args: args{
				ctx:  ctx,
				user: &entityAuth.UserSudir{Portals: []*entityAuth.Portal{{}}},
			},
			want: func(a args, f fields) (entityAuth.TokensPair, error) {
				testRequest := &authv1.CreateSessionRequest{
					User: &authv1.CreateSessionRequest_User{
						Id:       uuid.Nil.String(),
						Portal:   &authv1.CreateSessionRequest_UserPortal{},
						Employee: &authv1.CreateSessionRequest_Employee{},
						Person:   &authv1.CreateSessionRequest_Person{},
					},
					UserType:           authv1.UserType_USER_TYPE_AUTH,
					UserIp:             "<nil>",
					Device:             &authv1.CreateSessionRequest_Device{},
					SudirInfo:          &authv1.CreateSessionRequest_SudirInfo{},
					Issuer:             "test-app-name",
					AccessTtlDuration:  durationpb.New(time.Second),
					RefreshTtlDuration: durationpb.New(time.Second),
				}
				f.client.EXPECT().CreateSession(a.ctx, testRequest).Return(nil, testErr)
				return entityAuth.TokensPair{}, diterrors.GrpcErrorToError(testErr)
			},
		},
		{
			name: "correct",
			args: args{
				ctx:    ctx,
				user:   &entityAuth.UserSudir{Portals: []*entityAuth.Portal{{}}},
				device: &entityAuth.Device{UserAgent: "testUserAgent"},
			},
			want: func(a args, f fields) (entityAuth.TokensPair, error) {
				testRequest := &authv1.CreateSessionRequest{
					User: &authv1.CreateSessionRequest_User{
						Id:       uuid.Nil.String(),
						Portal:   &authv1.CreateSessionRequest_UserPortal{},
						Employee: &authv1.CreateSessionRequest_Employee{},
						Person:   &authv1.CreateSessionRequest_Person{},
					},
					UserType: authv1.UserType_USER_TYPE_AUTH,
					UserIp:   "<nil>",
					Device: &authv1.CreateSessionRequest_Device{
						UserAgent: "testUserAgent",
					},
					SudirInfo:          &authv1.CreateSessionRequest_SudirInfo{},
					Issuer:             "test-app-name",
					AccessTtlDuration:  durationpb.New(time.Second),
					RefreshTtlDuration: durationpb.New(time.Second),
				}
				testTime := time.Time{}
				testResponse := &authv1.CreateSessionResponse{
					Tokens: &authv1.TokensPair{
						AccessToken:  &authv1.Token{Value: "testAccessToken", ExpiredTime: timestamppb.New(testTime)},
						RefreshToken: &authv1.Token{Value: "testRefreshToken", ExpiredTime: timestamppb.New(testTime)},
					},
				}
				f.client.EXPECT().CreateSession(a.ctx, testRequest).Return(testResponse, nil)
				return entityAuth.TokensPair{
					AccessToken: entityAuth.Token{
						Value:     testResponse.GetTokens().GetAccessToken().GetValue(),
						ExpiredAt: &testTime,
					},
					RefreshToken: entityAuth.Token{
						Value:     testResponse.GetTokens().GetRefreshToken().GetValue(),
						ExpiredAt: &testTime,
					},
				}, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				client: authv1.NewMockAuthAPIClient(ctrl),
				mapper: NewMockMapperAuth(ctrl),
				tu:     timeUtils.NewMockTimeUtils(ctrl),
				logger: ditzap.NewMockLogger(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			ar := NewAuthRepository(f.client, f.mapper, url.URL{}, "test-app-name", time.Second, time.Second, f.tu, f.logger)
			got, err := ar.CreateSession(tt.args.ctx, tt.args.user, tt.args.clientIP, tt.args.device, tt.args.auth)
			if wantErr != nil {
				assert.EqualError(t, err, wantErr.Error())
				assert.Equal(t, entityAuth.TokensPair{}, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, want, got)
			}
		})
	}
}

func Test_authRepository_GetSession(t *testing.T) {
	type fields struct {
		client *authv1.MockAuthAPIClient
		mapper *MockMapperAuth
		tu     *timeUtils.MockTimeUtils
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
			name: "nil session err",
			args: args{
				ctx: ctx,
			},
			want: func(a args, f fields) (*entityAuth.Session, error) {
				return nil, diterrors.NewValidationError(repositories.ErrEmptyAccessToken)
			},
		},
		{
			name: "get session err",
			args: args{
				ctx:         ctx,
				accessToken: "testAccessToken",
			},
			want: func(a args, f fields) (*entityAuth.Session, error) {
				testRequest := &authv1.GetSessionRequest{AccessToken: a.accessToken}
				f.client.EXPECT().GetSession(a.ctx, testRequest).Return(nil, testErr)
				return nil, diterrors.GrpcErrorToError(testErr)
			},
		},
		{
			name: "map session err",
			args: args{
				ctx:         ctx,
				accessToken: "testAccessToken",
			},
			want: func(a args, f fields) (*entityAuth.Session, error) {
				testRequest := &authv1.GetSessionRequest{AccessToken: a.accessToken}
				testResponse := &authv1.GetSessionResponse{}
				f.client.EXPECT().GetSession(a.ctx, testRequest).Return(testResponse, nil)
				f.mapper.EXPECT().SessionToEntity(testResponse.GetSession()).Return(nil)
				return nil, diterrors.ErrNotFound
			},
		},
		{
			name: "correct",
			args: args{
				ctx:         ctx,
				accessToken: "testAccessToken",
			},
			want: func(a args, f fields) (*entityAuth.Session, error) {
				testRequest := &authv1.GetSessionRequest{AccessToken: a.accessToken}
				testResponse := &authv1.GetSessionResponse{}
				f.client.EXPECT().GetSession(a.ctx, testRequest).Return(testResponse, nil)
				testSession := &entityAuth.Session{}
				f.mapper.EXPECT().SessionToEntity(testResponse.GetSession()).Return(testSession)
				return testSession, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				client: authv1.NewMockAuthAPIClient(ctrl),
				mapper: NewMockMapperAuth(ctrl),
				tu:     timeUtils.NewMockTimeUtils(ctrl),
				logger: ditzap.NewMockLogger(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			ar := NewAuthRepository(f.client, f.mapper, url.URL{}, "test-app-name", time.Second, time.Second, f.tu, f.logger)
			got, err := ar.GetSession(tt.args.ctx, tt.args.accessToken)
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

func Test_authRepository_Logout(t *testing.T) {
	type fields struct {
		client *authv1.MockAuthAPIClient
		mapper *MockMapperAuth
		tu     *timeUtils.MockTimeUtils
		logger *ditzap.MockLogger
	}
	type args struct {
		ctx                       context.Context
		session                   *entityAuth.Session
		accessToken, refreshToken string
	}

	ctx := context.TODO()
	testErr := errors.New("test error")

	tests := []struct {
		name string
		args args
		want func(a args, f fields) error
	}{
		{
			name: "nil session err",
			args: args{
				ctx: ctx,
			},
			want: func(a args, f fields) error {
				return diterrors.NewValidationError(repositories.ErrNilSession)
			},
		},
		{
			name: "empty access token err",
			args: args{
				ctx:     ctx,
				session: &entityAuth.Session{},
			},
			want: func(a args, f fields) error {
				return diterrors.NewValidationError(repositories.ErrEmptyAccessToken)
			},
		},
		{
			name: "empty refresh token err",
			args: args{
				ctx:         ctx,
				session:     &entityAuth.Session{},
				accessToken: "testAccessToken",
			},
			want: func(a args, f fields) error {
				return diterrors.NewValidationError(repositories.ErrEmptyRefreshToken)
			},
		},
		{
			name: "nil active portal err",
			args: args{
				ctx:          ctx,
				session:      &entityAuth.Session{},
				accessToken:  "testAccessToken",
				refreshToken: "testRefreshToken",
			},
			want: func(a args, f fields) error {
				return repositories.ErrNilSessionActivePortal
			},
		},
		{
			name: "nil user err",
			args: args{
				ctx: ctx,
				session: &entityAuth.Session{
					ActivePortal: &entityAuth.ActivePortal{
						Portal: entityAuth.Portal{
							URL: "testURL",
						},
						SID: "testSID",
					}},
				accessToken:  "testAccessToken",
				refreshToken: "testRefreshToken",
			},
			want: func(a args, f fields) error {
				return repositories.ErrNilSessionUser
			},
		},
		{
			name: "nil sudir info err",
			args: args{
				ctx: ctx,
				session: &entityAuth.Session{
					User: &entityAuth.User{
						CloudID: "testCloudID",
					},
					ActivePortal: &entityAuth.ActivePortal{
						Portal: entityAuth.Portal{
							URL: "testURL",
						},
						SID: "testSID",
					}},
				accessToken:  "testAccessToken",
				refreshToken: "testRefreshToken",
			},
			want: func(a args, f fields) error {
				return repositories.ErrNilSessionSudirInfo
			},
		},
		{
			name: "logout err",
			args: args{
				ctx: ctx,
				session: &entityAuth.Session{
					User: &entityAuth.User{
						CloudID: "testCloudID",
					},
					Device: &entityAuth.Device{
						SudirInfo: &entityAuth.SudirInfo{
							SID: "testSudirID",
						},
					},
					ActivePortal: &entityAuth.ActivePortal{
						Portal: entityAuth.Portal{
							URL: "testURL",
						},
						SID: "testSID",
					}},
				accessToken:  "testAccessToken",
				refreshToken: "testRefreshToken",
			},
			want: func(a args, f fields) error {
				testRequest := &authv1.LogoutRequest{
					PortalUrl: a.session.GetActivePortal().Portal.URL,
					SessionId: a.session.GetID().String(),
					LogoutBy: &authv1.LogoutRequest_CloudId{
						CloudId: a.session.GetUser().CloudID,
					},
					SudirSid:      a.session.GetDevice().GetSudirInfo().SID,
					PortalSession: a.session.GetActivePortal().SID,
					Tokens: &authv1.LogoutRequest_TokensPair{
						AccessToken:  a.accessToken,
						RefreshToken: a.refreshToken,
					},
				}
				f.client.EXPECT().Logout(a.ctx, testRequest).Return(nil, testErr)
				return diterrors.GrpcErrorToError(testErr)
			},
		},
		{
			name: "correct",
			args: args{
				ctx: ctx,
				session: &entityAuth.Session{
					User: &entityAuth.User{
						CloudID: "testCloudID",
					},
					Device: &entityAuth.Device{
						SudirInfo: &entityAuth.SudirInfo{
							SID: "testSudirID",
						},
					},
					ActivePortal: &entityAuth.ActivePortal{
						Portal: entityAuth.Portal{
							URL: "testURL",
						},
						SID: "testSID",
					}},
				accessToken:  "testAccessToken",
				refreshToken: "testRefreshToken",
			},
			want: func(a args, f fields) error {
				testRequest := &authv1.LogoutRequest{
					PortalUrl: a.session.GetActivePortal().Portal.URL,
					SessionId: a.session.GetID().String(),
					LogoutBy: &authv1.LogoutRequest_CloudId{
						CloudId: a.session.GetUser().CloudID,
					},
					SudirSid:      a.session.GetDevice().GetSudirInfo().SID,
					PortalSession: a.session.GetActivePortal().SID,
					Tokens: &authv1.LogoutRequest_TokensPair{
						AccessToken:  a.accessToken,
						RefreshToken: a.refreshToken,
					},
				}
				f.client.EXPECT().Logout(a.ctx, testRequest).Return(&authv1.LogoutResponse{}, nil)
				return nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				client: authv1.NewMockAuthAPIClient(ctrl),
				mapper: NewMockMapperAuth(ctrl),
				tu:     timeUtils.NewMockTimeUtils(ctrl),
				logger: ditzap.NewMockLogger(ctrl),
			}
			wantErr := tt.want(tt.args, f)
			ar := NewAuthRepository(f.client, f.mapper, url.URL{}, "test-app-name", time.Second, time.Second, f.tu, f.logger)
			err := ar.Logout(tt.args.ctx, tt.args.session, tt.args.accessToken, tt.args.refreshToken)
			if wantErr != nil {
				assert.EqualError(t, err, wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_authRepository_ChangePortal(t *testing.T) {
	type fields struct {
		client *authv1.MockAuthAPIClient
		mapper *MockMapperAuth
		tu     *timeUtils.MockTimeUtils
		logger *ditzap.MockLogger
	}
	type args struct {
		ctx      context.Context
		portalID int
		session  *entityAuth.Session
	}

	ctx := context.TODO()
	testErr := errors.New("test error")

	tests := []struct {
		name string
		args args
		want func(a args, f fields) ([]*entityAuth.Portal, string, error)
	}{
		{
			name: "nil session err",
			args: args{
				ctx: ctx,
			},
			want: func(a args, f fields) ([]*entityAuth.Portal, string, error) {
				return nil, "", repositories.ErrNilSession
			},
		},
		{
			name: "change portal err",
			args: args{
				ctx:      ctx,
				portalID: 1,
				session:  &entityAuth.Session{},
			},
			want: func(a args, f fields) ([]*entityAuth.Portal, string, error) {
				testSession := &authv1.Session{}
				f.mapper.EXPECT().SessionToPb(a.session).Return(testSession)
				testRequest := &authv1.ChangePortalRequest{
					SelectedPortalId: int32(a.portalID),
					Session:          testSession,
				}
				f.client.EXPECT().ChangePortal(a.ctx, testRequest).Return(nil, testErr)
				return nil, "", diterrors.GrpcErrorToError(testErr)
			},
		},
		{
			name: "correct",
			args: args{
				ctx:      ctx,
				portalID: 1,
				session:  &entityAuth.Session{},
			},
			want: func(a args, f fields) ([]*entityAuth.Portal, string, error) {
				testSession := &authv1.Session{}
				f.mapper.EXPECT().SessionToPb(a.session).Return(testSession)
				testRequest := &authv1.ChangePortalRequest{
					SelectedPortalId: int32(a.portalID),
					Session:          testSession,
				}
				testResponse := &authv1.ChangePortalResponse{}
				testPortals := []*entityAuth.Portal{{}}
				f.client.EXPECT().ChangePortal(a.ctx, testRequest).Return(testResponse, nil)
				f.mapper.EXPECT().PortalsToEntity(testResponse.GetPortals()).Return(testPortals)
				return testPortals, testResponse.GetPortalSid(), nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				client: authv1.NewMockAuthAPIClient(ctrl),
				mapper: NewMockMapperAuth(ctrl),
				tu:     timeUtils.NewMockTimeUtils(ctrl),
				logger: ditzap.NewMockLogger(ctrl),
			}
			wantP, wantS, wantErr := tt.want(tt.args, f)
			ar := NewAuthRepository(f.client, f.mapper, url.URL{}, "test-app-name", time.Second, time.Second, f.tu, f.logger)
			gotP, gotS, err := ar.ChangePortal(tt.args.ctx, tt.args.portalID, tt.args.session)
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

func Test_authRepository_RefreshTokensPair(t *testing.T) {
	type fields struct {
		client *authv1.MockAuthAPIClient
		mapper *MockMapperAuth
		tu     *timeUtils.MockTimeUtils
		logger *ditzap.MockLogger
	}
	type args struct {
		ctx                       context.Context
		accessToken, refreshToken string
	}

	ctx := context.TODO()
	testErr := errors.New("test error")
	testTime := time.Now()
	testTimePb := timestamppb.New(testTime)

	tests := []struct {
		name string
		args args
		want func(a args, f fields) (*entityAuth.TokensPair, error)
	}{
		{
			name: "empty access token err",
			args: args{
				ctx:          ctx,
				accessToken:  "",
				refreshToken: "testRefreshToken",
			},
			want: func(a args, f fields) (*entityAuth.TokensPair, error) {
				return nil, repositories.NewDetailsError("access_token", repositories.ErrEmptyAccessToken.Error(), false)
			},
		},
		{
			name: "empty refresh token err",
			args: args{
				ctx:          ctx,
				accessToken:  "testAccessToken",
				refreshToken: "",
			},
			want: func(a args, f fields) (*entityAuth.TokensPair, error) {
				return nil, repositories.NewDetailsError("refresh_token", repositories.ErrEmptyRefreshToken.Error(), false)
			},
		},
		{
			name: "details err",
			args: args{
				ctx:          ctx,
				accessToken:  "testAccessToken",
				refreshToken: "testRefreshToken",
			},
			want: func(a args, f fields) (*entityAuth.TokensPair, error) {
				testRequest := &authv1.RefreshTokensPairRequest{
					Tokens: &authv1.RefreshTokensPairRequest_TokensPair{
						AccessToken:  a.accessToken,
						RefreshToken: a.refreshToken,
					}}
				testStatus := status.New(codes.InvalidArgument, testErr.Error())
				testDetails := &authErrorsv1.AuthInvalidArgument{
					Field:          "access_token",
					Message:        "test message",
					ReauthRequired: true,
				}
				testStatusWithDetails, _ := testStatus.WithDetails(testDetails)
				f.client.EXPECT().RefreshTokensPair(a.ctx, testRequest).Return(nil, testStatusWithDetails.Err())
				return nil, repositories.NewDetailsError(testDetails.GetField(), testDetails.GetMessage(), testDetails.GetReauthRequired())
			},
		},
		{
			name: "err",
			args: args{
				ctx:          ctx,
				accessToken:  "testAccessToken",
				refreshToken: "testRefreshToken",
			},
			want: func(a args, f fields) (*entityAuth.TokensPair, error) {
				testRequest := &authv1.RefreshTokensPairRequest{
					Tokens: &authv1.RefreshTokensPairRequest_TokensPair{
						AccessToken:  a.accessToken,
						RefreshToken: a.refreshToken,
					}}
				f.client.EXPECT().RefreshTokensPair(a.ctx, testRequest).Return(nil, testErr)
				return nil, diterrors.GrpcErrorToError(testErr)
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
				testRequest := &authv1.RefreshTokensPairRequest{
					Tokens: &authv1.RefreshTokensPairRequest_TokensPair{
						AccessToken:  a.accessToken,
						RefreshToken: a.refreshToken,
					}}
				testResponse := &authv1.RefreshTokensPairResponse{
					Tokens: &authv1.TokensPair{
						AccessToken: &authv1.Token{
							Value:       "testAccessToken",
							ExpiredTime: testTimePb,
						},
						RefreshToken: &authv1.Token{
							Value:       "testRefreshToken",
							ExpiredTime: testTimePb,
						},
					},
				}
				f.client.EXPECT().RefreshTokensPair(a.ctx, testRequest).Return(testResponse, nil)
				f.tu.EXPECT().TimestampToTime(testTimePb).Return(&testTime).Times(2)
				return &entityAuth.TokensPair{
					AccessToken: entityAuth.Token{
						Value:     testResponse.GetTokens().GetAccessToken().GetValue(),
						ExpiredAt: &testTime,
					},
					RefreshToken: entityAuth.Token{
						Value:     testResponse.GetTokens().GetRefreshToken().GetValue(),
						ExpiredAt: &testTime,
					},
				}, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				client: authv1.NewMockAuthAPIClient(ctrl),
				mapper: NewMockMapperAuth(ctrl),
				tu:     timeUtils.NewMockTimeUtils(ctrl),
				logger: ditzap.NewMockLogger(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			ar := NewAuthRepository(f.client, f.mapper, url.URL{}, "test-app-name", time.Second, time.Second, f.tu, f.logger)
			got, err := ar.RefreshTokensPair(tt.args.ctx, tt.args.accessToken, tt.args.refreshToken)
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
