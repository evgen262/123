package usecase

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"git.mos.ru/buch-cloud/moscow-team-2.0/pud/auth.git/internal/client/http/kadry"
	"git.mos.ru/buch-cloud/moscow-team-2.0/pud/auth.git/internal/client/http/sudir"
	"git.mos.ru/buch-cloud/moscow-team-2.0/pud/auth.git/internal/entity"
)

/*
TODO: исправить тест для существующих изменений
func Test_authUseCase_Auth(t *testing.T) {
	ctx := context.TODO()
	ctxCancelLess := context.WithoutCancel(ctx)
	testErr := fmt.Errorf("testErr")
	testTime := time.Date(2023, 7, 14, 10, 10, 10, 0, time.UTC)
	testToken := &sudir.OAuthResponse{
		IDToken:      "",
		AccessToken:  "",
		RefreshToken: "",
		Expiry:       &testTime,
	}
	testJWTPayload := &sudir.JWTPayload{
		RegisteredClaims: jwt.RegisteredClaims{},
		UserClaims: sudir.UserClaims{
			CloudGUID:  "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
			Company:    "ГКУ \"Инфогород\"",
			Department: "Отдел мобильной разработки",
			Email:      "test@it.mos.ru",
			LogonName:  "test",
			Position:   "Программист",
		},
	}
	type fields struct {
		sudirClient        *MockSudirClient
		kadryClient        *MockKadryClient
		stateRepository    *MockStateRepository
		tokenRepository    *MockTokenRepository
		employeeRepository *MockEmployeeRepository
		logger             *ditzap.MockLogger
	}
	type args struct {
		ctx           context.Context
		ctxCancelLess context.Context
		code          string
		state         string
		callbackUrl   string
	}
	tests := []struct {
		name string
		args args
		want func(a args, f fields) (*entity.AuthInfo, error)
	}{
		{
			name: "err not valid state",
			args: args{
				ctx:   ctx,
				state: "33ebcb16-fail-4d54-b1e5-f6727a0090ac",
			},
			want: func(a args, f fields) (*entity.AuthInfo, error) {
				f.stateRepository.EXPECT().IsValid(a.ctx, a.state).Return(repositories.ErrInvalidState)
				return nil, repositories.ErrInvalidState
			},
		},
		{
			name: "err not found state",
			args: args{
				ctx:   ctx,
				state: "33ebcb16-fail-4d54-b1e5-f6727a0090ac",
			},
			want: func(a args, f fields) (*entity.AuthInfo, error) {
				f.stateRepository.EXPECT().IsValid(a.ctx, a.state).Return(repositories.ErrNotFound)
				return nil, repositories.ErrNotFound
			},
		},
		{
			name: "code not valid or sudir error",
			args: args{
				ctx:           ctx,
				ctxCancelLess: ctxCancelLess,
				code:          "",
				state:         "33ebcb16-3ff5-4d54-b1e5-f6727a0090ac",
			},
			want: func(a args, f fields) (*entity.AuthInfo, error) {
				f.stateRepository.EXPECT().IsValid(a.ctx, a.state).Return(nil)
				f.sudirClient.EXPECT().CodeExchange(a.ctx, a.code, a.callbackUrl).Return(nil, testErr)
				f.stateRepository.EXPECT().Delete(a.ctxCancelLess, a.state)
				return nil, testErr
			},
		},
		{
			name: "no jwt token",
			args: args{
				ctx:           ctx,
				ctxCancelLess: ctxCancelLess,
				code:          "9HA9yRq2NNPjh8UDBMrvGa9JF6dGRK_r-qTw0Uo16tDwJ6PbQiLo8kBKv7za-YOGP2Vu-xu9StSTUyp0_BvXVel8Eho7SMfJwyjmMtPu2LI",
				state:         "33ebcb16-3ff5-4d54-b1e5-f6727a0090ac",
			},
			want: func(a args, f fields) (*entity.AuthInfo, error) {
				f.stateRepository.EXPECT().IsValid(a.ctx, a.state).Return(nil)
				f.sudirClient.EXPECT().CodeExchange(a.ctx, a.code, a.callbackUrl).Return(testToken, sudir.ErrNoJWTToken)
				f.logger.EXPECT().Error("в ответе СУДИР отсутствует IDToken", zap.Error(sudir.ErrNoJWTToken))
				f.stateRepository.EXPECT().Delete(a.ctxCancelLess, a.state)
				return nil, sudir.ErrNoJWTToken
			},
		},
		{
			name: "parse token error",
			args: args{
				ctx:           ctx,
				ctxCancelLess: ctxCancelLess,
				code:          "9HA9yRq2NNPjh8UDBMrvGa9JF6dGRK_r-qTw0Uo16tDwJ6PbQiLo8kBKv7za-YOGP2Vu-xu9StSTUyp0_BvXVel8Eho7SMfJwyjmMtPu2LI",
				state:         "33ebcb16-3ff5-4d54-b1e5-f6727a0090ac",
			},
			want: func(a args, f fields) (*entity.AuthInfo, error) {
				f.stateRepository.EXPECT().IsValid(a.ctx, a.state).Return(nil)
				f.sudirClient.EXPECT().CodeExchange(a.ctx, a.code, a.callbackUrl).Return(testToken, nil)
				f.sudirClient.EXPECT().ParseToken(testToken.IDToken).Return(nil, testErr)
				f.logger.EXPECT().Error("ошибка парсинга IDToken", zap.Error(testErr), zap.String("id_token", testToken.IDToken))
				f.stateRepository.EXPECT().Delete(a.ctxCancelLess, a.state)
				return nil, testErr
			},
		},
		{
			name: "empty cloudID",
			args: args{
				ctx:           ctx,
				ctxCancelLess: ctxCancelLess,
				code:          "9HA9yRq2NNPjh8UDBMrvGa9JF6dGRK_r-qTw0Uo16tDwJ6PbQiLo8kBKv7za-YOGP2Vu-xu9StSTUyp0_BvXVel8Eho7SMfJwyjmMtPu2LI",
				state:         "33ebcb16-3ff5-4d54-b1e5-f6727a0090ac",
			},
			want: func(a args, f fields) (*entity.AuthInfo, error) {
				badJWTPayload := &sudir.JWTPayload{
					RegisteredClaims: jwt.RegisteredClaims{},
					UserClaims: sudir.UserClaims{
						Company:    "ГКУ \"Инфогород\"",
						Department: "Отдел мобильной разработки",
						Email:      "test@it.mos.ru",
						LogonName:  "test",
						Position:   "Программист",
					},
				}

				f.stateRepository.EXPECT().IsValid(a.ctx, a.state).Return(nil)
				f.sudirClient.EXPECT().CodeExchange(a.ctx, a.code, a.callbackUrl).Return(testToken, nil)
				f.sudirClient.EXPECT().ParseToken(testToken.IDToken).Return(badJWTPayload, nil)
				f.logger.EXPECT().Warn("в payload JWT отсутствует CloudGUID",
					zap.String("logon_name", "test"),
					zap.String("email", "test@it.mos.ru"),
				)

				f.stateRepository.EXPECT().Delete(a.ctxCancelLess, a.state)
				return &entity.AuthInfo{
					OAuth: &entity.OAuth{
						AccessToken:  testToken.AccessToken,
						RefreshToken: testToken.RefreshToken,
						Expiry:       testToken.Expiry,
					},
					User: &entity.User{
						CloudID:   "",
						Email:     "test@it.mos.ru",
						LogonName: "test",
					},
				}, nil
			},
		},
		{
			name: "employees repo ok",
			args: args{
				ctx:           ctx,
				ctxCancelLess: ctxCancelLess,
				code:          "9HA9yRq2NNPjh8UDBMrvGa9JF6dGRK_r-qTw0Uo16tDwJ6PbQiLo8kBKv7za-YOGP2Vu-xu9StSTUyp0_BvXVel8Eho7SMfJwyjmMtPu2LI",
				state:         "33ebcb16-3ff5-4d54-b1e5-f6727a0090ac",
			},
			want: func(a args, f fields) (*entity.AuthInfo, error) {
				f.stateRepository.EXPECT().IsValid(a.ctx, a.state).Return(nil)
				f.sudirClient.EXPECT().CodeExchange(a.ctx, a.code, a.callbackUrl).Return(testToken, nil)
				f.sudirClient.EXPECT().ParseToken(testToken.IDToken).Return(testJWTPayload, nil)

				f.tokenRepository.EXPECT().Save(a.ctx, testJWTPayload.CloudGUID, testToken.RefreshToken).Return(nil)

				f.employeeRepository.EXPECT().Get(a.ctx, testJWTPayload.CloudGUID).Return([]entity.EmployeeInfo{
					{
						CloudID: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
						Inn:     "770123456789",
					},
					{
						CloudID: "22c2a7dc-34ef-d9d6-c048-76b39bfbaf6a",
						Inn:     "779876543210",
					},
				}, nil)

				f.stateRepository.EXPECT().Delete(a.ctxCancelLess, a.state)
				return &entity.AuthInfo{
					OAuth: &entity.OAuth{
						AccessToken:  testToken.AccessToken,
						RefreshToken: testToken.RefreshToken,
						Expiry:       testToken.Expiry,
					},
					User: &entity.User{
						CloudID:   entity.CloudID(testJWTPayload.CloudGUID),
						Email:     testJWTPayload.Email,
						LogonName: testJWTPayload.LogonName,
						Employees: []entity.EmployeeInfo{
							{
								CloudID: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
								Inn:     "770123456789",
							},
							{
								CloudID: "22c2a7dc-34ef-d9d6-c048-76b39bfbaf6a",
								Inn:     "779876543210",
							},
						},
					},
				}, nil
			},
		},
		{
			name: "employees service error",
			args: args{
				ctx:           ctx,
				ctxCancelLess: ctxCancelLess,
				code:          "9HA9yRq2NNPjh8UDBMrvGa9JF6dGRK_r-qTw0Uo16tDwJ6PbQiLo8kBKv7za-YOGP2Vu-xu9StSTUyp0_BvXVel8Eho7SMfJwyjmMtPu2LI",
				state:         "33ebcb16-3ff5-4d54-b1e5-f6727a0090ac",
			},
			want: func(a args, f fields) (*entity.AuthInfo, error) {
				f.stateRepository.EXPECT().IsValid(a.ctx, a.state).Return(nil)
				f.sudirClient.EXPECT().CodeExchange(a.ctx, a.code, a.callbackUrl).Return(testToken, nil)
				f.sudirClient.EXPECT().ParseToken(testToken.IDToken).Return(testJWTPayload, nil)

				f.tokenRepository.EXPECT().Save(a.ctx, testJWTPayload.CloudGUID, testToken.RefreshToken).Return(nil)
				f.employeeRepository.EXPECT().Get(a.ctx, testJWTPayload.CloudGUID).Return(nil, repositories.ErrNotFound)

				f.kadryClient.EXPECT().GetEmployeesInfo(a.ctx, testJWTPayload.CloudGUID, kadry.PersonID, kadry.InnOrg).Return(nil, testErr)

				f.stateRepository.EXPECT().Delete(a.ctxCancelLess, a.state)
				return &entity.AuthInfo{
					OAuth: &entity.OAuth{
						AccessToken:  testToken.AccessToken,
						RefreshToken: testToken.RefreshToken,
						Expiry:       testToken.Expiry,
					},
					User: &entity.User{
						CloudID:   entity.CloudID(testJWTPayload.CloudGUID),
						Email:     testJWTPayload.Email,
						LogonName: testJWTPayload.LogonName,
					},
				}, kadry.ErrEmployeesService
			},
		},
		{
			name: "employees service wrong data",
			args: args{
				ctx:           ctx,
				ctxCancelLess: ctxCancelLess,
				code:          "9HA9yRq2NNPjh8UDBMrvGa9JF6dGRK_r-qTw0Uo16tDwJ6PbQiLo8kBKv7za-YOGP2Vu-xu9StSTUyp0_BvXVel8Eho7SMfJwyjmMtPu2LI",
				state:         "33ebcb16-3ff5-4d54-b1e5-f6727a0090ac",
			},
			want: func(a args, f fields) (*entity.AuthInfo, error) {
				wrongEmployees := []entity.EmployeeInfo{
					{
						CloudID: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
						Inn:     "770123456789",
					},
					{
						CloudID: "22c2a7dc-34ef-d9d6-c048-76b39bfbaf6a",
						Inn:     "779876543210",
					},
				}

				f.stateRepository.EXPECT().IsValid(a.ctx, a.state).Return(nil)
				f.sudirClient.EXPECT().CodeExchange(a.ctx, a.code, a.callbackUrl).Return(testToken, nil)
				f.sudirClient.EXPECT().ParseToken(testToken.IDToken).Return(testJWTPayload, nil)

				f.tokenRepository.EXPECT().Save(a.ctx, testJWTPayload.CloudGUID, testToken.RefreshToken).Return(nil)
				f.employeeRepository.EXPECT().Get(a.ctx, testJWTPayload.CloudGUID).Return(nil, repositories.ErrNotFound)

				f.kadryClient.EXPECT().GetEmployeesInfo(a.ctx, testJWTPayload.CloudGUID, kadry.PersonID, kadry.InnOrg).Return(wrongEmployees, nil)
				f.logger.EXPECT().Warn("Неверный PersonID из системы кадров",
					zap.String("sudirID", "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f"),
					zap.String("PersonID", "22c2a7dc-34ef-d9d6-c048-76b39bfbaf6a"),
				)

				f.employeeRepository.EXPECT().Save(a.ctx, testJWTPayload.CloudGUID, []entity.EmployeeInfo{{CloudID: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f", Inn: "770123456789"}})

				f.stateRepository.EXPECT().Delete(a.ctxCancelLess, a.state)
				return &entity.AuthInfo{
					OAuth: &entity.OAuth{
						AccessToken:  testToken.AccessToken,
						RefreshToken: testToken.RefreshToken,
						Expiry:       testToken.Expiry,
					},
					User: &entity.User{
						CloudID:   entity.CloudID(testJWTPayload.CloudGUID),
						Email:     testJWTPayload.Email,
						LogonName: testJWTPayload.LogonName,
						Employees: []entity.EmployeeInfo{
							{
								CloudID: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
								Inn:     "770123456789",
							},
						},
					},
				}, nil
			},
		},
		{
			name: "correct",
			args: args{
				ctx:           ctx,
				ctxCancelLess: ctxCancelLess,
				code:          "9HA9yRq2NNPjh8UDBMrvGa9JF6dGRK_r-qTw0Uo16tDwJ6PbQiLo8kBKv7za-YOGP2Vu-xu9StSTUyp0_BvXVel8Eho7SMfJwyjmMtPu2LI",
				state:         "33ebcb16-3ff5-4d54-b1e5-f6727a0090ac",
			},
			want: func(a args, f fields) (*entity.AuthInfo, error) {
				testEmployees := []entity.EmployeeInfo{
					{
						CloudID: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
						Inn:     "770123456789",
						OrgID:   "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
						FIO:     "Иванов Иван Иванович",
					},
					{
						CloudID: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
						Inn:     "779876543210",
						OrgID:   "71bac977-ea7f-4156-9504-60f7d443ab62",
						FIO:     "Иванов Иван Иванович",
					},
				}

				f.stateRepository.EXPECT().IsValid(a.ctx, a.state).Return(nil)
				f.sudirClient.EXPECT().CodeExchange(a.ctx, a.code, a.callbackUrl).Return(testToken, nil)
				f.sudirClient.EXPECT().ParseToken(testToken.IDToken).Return(testJWTPayload, nil)

				f.tokenRepository.EXPECT().Save(a.ctx, testJWTPayload.CloudGUID, testToken.RefreshToken).Return(nil)
				f.employeeRepository.EXPECT().Get(a.ctx, testJWTPayload.CloudGUID).Return(nil, repositories.ErrNotFound)

				f.kadryClient.EXPECT().GetEmployeesInfo(a.ctx, testJWTPayload.CloudGUID, kadry.PersonID, kadry.InnOrg).Return(testEmployees, nil)

				f.employeeRepository.EXPECT().Save(
					a.ctx,
					testJWTPayload.CloudGUID,
					testEmployees,
				)

				f.stateRepository.EXPECT().Delete(a.ctxCancelLess, a.state)
				return &entity.AuthInfo{
					OAuth: &entity.OAuth{
						AccessToken:  testToken.AccessToken,
						RefreshToken: testToken.RefreshToken,
						Expiry:       testToken.Expiry,
					},
					User: &entity.User{
						CloudID:   entity.CloudID(testJWTPayload.CloudGUID),
						Email:     testJWTPayload.Email,
						LogonName: testJWTPayload.LogonName,
						Employees: []entity.EmployeeInfo{
							{
								CloudID: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
								Inn:     "770123456789",
								OrgID:   "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
								FIO:     "Иванов Иван Иванович",
							},
							{
								CloudID: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
								Inn:     "779876543210",
								OrgID:   "71bac977-ea7f-4156-9504-60f7d443ab62",
								FIO:     "Иванов Иван Иванович",
							},
						},
					},
				}, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			f := fields{
				sudirClient:        NewMockSudirClient(ctrl),
				kadryClient:        NewMockKadryClient(ctrl),
				stateRepository:    NewMockStateRepository(ctrl),
				tokenRepository:    NewMockTokenRepository(ctrl),
				employeeRepository: NewMockEmployeeRepository(ctrl),
				logger:             ditzap.NewMockLogger(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			uc := NewAuthUseCase(
				f.sudirClient,
				f.kadryClient,
				f.stateRepository,
				f.tokenRepository,
				f.employeeRepository,
				f.logger,
			)
			got, err := uc.Auth(tt.args.ctx, tt.args.code, tt.args.state, tt.args.callbackUrl)
			if wantErr != nil {
				assert.EqualError(t, err, wantErr.Error())
				switch true {
				case errors.Is(wantErr, kadry.ErrEmployeesService):
					assert.Equal(t, want, got)
					return
				default:
					assert.Nil(t, got)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, want, got)
			}
		})
	}
}
*/

func Test_authUseCase_GetAuthURL(t *testing.T) {
	type fields struct {
		sudirClient        *MockSudirClient
		kadryClient        *MockKadryClient
		stateRepository    *MockStateRepository
		tokenRepository    *MockTokenRepository
		employeeRepository *MockEmployeeRepository
		ctx                context.Context
		logger             *ditzap.MockLogger
	}
	type args struct {
		ctx                                 context.Context
		callbackURL, clientID, clientSecret string
	}

	ctx := context.TODO()
	ctxValue := context.WithValue(ctx, "deviceid", "testDeviceID")
	ctxValue = context.WithValue(ctxValue, "x-cfc-useragent", "testUserAgent")
	testErr := errors.New("some test error")

	tests := []struct {
		name string
		args args
		want func(a args, f fields) (string, error)
	}{
		{
			name: "error",
			args: args{
				ctx:         ctxValue,
				callbackURL: "testCallbackURL",
			},
			want: func(a args, f fields) (string, error) {
				testStateOptions := &entity.StateOptions{
					CallbackURL: "testCallbackURL",
				}
				f.stateRepository.EXPECT().New(a.ctx, testStateOptions).Return(nil, testErr)
				return "", testErr
			},
		},
		{
			name: "correct",
			args: args{
				ctx:         ctxValue,
				callbackURL: "testCallbackURL",
			},
			want: func(a args, f fields) (string, error) {
				testStateOptions := &entity.StateOptions{
					CallbackURL: "testCallbackURL",
				}
				testOptions := sudir.AuthURLOptions{
					IsOffline:   true,
					RedirectURI: "testCallbackURL",
					State:       "testID",
				}
				f.stateRepository.EXPECT().New(a.ctx, testStateOptions).Return(&entity.State{ID: "testID"}, nil)
				testAuthURL := "testAuthURL"
				f.sudirClient.EXPECT().AuthURL(testOptions).Return(testAuthURL)
				return testAuthURL, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			f := fields{
				sudirClient:        NewMockSudirClient(ctrl),
				kadryClient:        NewMockKadryClient(ctrl),
				stateRepository:    NewMockStateRepository(ctrl),
				tokenRepository:    NewMockTokenRepository(ctrl),
				employeeRepository: NewMockEmployeeRepository(ctrl),
				ctx:                context.TODO(),
				logger:             ditzap.NewMockLogger(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			uc := NewAuthUseCase(
				f.sudirClient,
				f.kadryClient,
				f.stateRepository,
				f.tokenRepository,
				f.employeeRepository,
				f.logger,
			)
			got, err := uc.GetAuthURL(tt.args.ctx, tt.args.callbackURL, tt.args.clientID, tt.args.clientSecret)

			if wantErr != nil {
				assert.EqualError(t, err, wantErr.Error())
				assert.Empty(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, want, got)
			}
		})
	}
}

/*
func Test_authUseCase_RefreshToken(t *testing.T) {
	ctx := context.TODO()
	testTime := time.Date(2023, 7, 14, 10, 10, 10, 0, time.UTC)
	testErr := fmt.Errorf("testErr")
	type fields struct {
		sudirClient        *MockSudirClient
		kadryClient        *MockKadryClient
		stateRepository    *MockStateRepository
		tokenRepository    *MockTokenRepository
		employeeRepository *MockEmployeeRepository
		callbackURL        *url.URL
		logger             *ditzap.MockLogger
	}
	type args struct {
		ctx     context.Context
		cloudID string
	}
	tests := []struct {
		name string
		args args
		want func(a args, f fields) (string, error)
	}{
		{
			name: "error token repo",
			args: args{
				ctx:     ctx,
				cloudID: "some wrong token",
			},
			want: func(a args, f fields) (string, error) {
				f.tokenRepository.EXPECT().Get(a.ctx, a.cloudID).Return("", testErr)
				return "", testErr
			},
		},
		{
			name: "error refresh token",
			args: args{
				ctx:     ctx,
				cloudID: "some wrong token",
			},
			want: func(a args, f fields) (string, error) {
				testToken := "WNjZilQWNRjK9yFMg3YGmSAp-KTJFYySKT64TEiKGuBX_kRaTRImpztssfRetWT81sl98llHXwpCl9C1HyHALg"

				f.tokenRepository.EXPECT().Get(a.ctx, a.cloudID).Return(testToken, nil)
				f.sudirClient.EXPECT().RefreshToken(a.ctx, testToken).Return(nil, testErr)
				return "", testErr
			},
		},
		{
			name: "correct",
			args: args{
				ctx:     ctx,
				cloudID: "some correct token",
			},
			want: func(a args, f fields) (string, error) {
				testOauth := &sudir.OAuthResponse{
					AccessToken:  "new access token",
					RefreshToken: "new refresh token",
					Expiry:       &testTime,
				}
				testToken := "WNjZilQWNRjK9yFMg3YGmSAp-KTJFYySKT64TEiKGuBX_kRaTRImpztssfRetWT81sl98llHXwpCl9C1HyHALg"

				f.tokenRepository.EXPECT().Get(a.ctx, a.cloudID).Return(testToken, nil)
				f.sudirClient.EXPECT().RefreshToken(a.ctx, testToken).Return(testOauth, nil)
				f.tokenRepository.EXPECT().Save(a.ctx, a.cloudID, testOauth.RefreshToken)

				return testOauth.AccessToken, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			f := fields{
				sudirClient:        NewMockSudirClient(ctrl),
				kadryClient:        NewMockKadryClient(ctrl),
				stateRepository:    NewMockStateRepository(ctrl),
				tokenRepository:    NewMockTokenRepository(ctrl),
				employeeRepository: NewMockEmployeeRepository(ctrl),
				logger:             ditzap.NewMockLogger(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)

			uc := NewAuthUseCase(
				f.sudirClient,
				f.kadryClient,
				f.stateRepository,
				f.tokenRepository,
				f.employeeRepository,
				f.logger,
			)

			got, err := uc.RefreshToken(tt.args.ctx, tt.args.cloudID)
			if wantErr != nil {
				assert.EqualError(t, err, wantErr.Error())
				assert.Empty(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, want, got)
			}
		})
	}
}

func Test_authUseCase_Logout(t *testing.T) {
	ctx := context.TODO()

	type fields struct {
		sudirClient        *MockSudirClient
		kadryClient        *MockKadryClient
		stateRepository    *MockStateRepository
		tokenRepository    *MockTokenRepository
		employeeRepository *MockEmployeeRepository
		logger             *ditzap.MockLogger
	}
	type args struct {
		ctx     context.Context
		cloudID string
	}
	tests := []struct {
		name string
		args args
		want func(a args, f fields)
	}{
		{
			name: "correct",
			args: args{
				ctx:     ctx,
				cloudID: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
			},
			want: func(a args, f fields) {
				f.tokenRepository.EXPECT().Delete(a.ctx, a.cloudID)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			f := fields{
				sudirClient:        NewMockSudirClient(ctrl),
				kadryClient:        NewMockKadryClient(ctrl),
				stateRepository:    NewMockStateRepository(ctrl),
				tokenRepository:    NewMockTokenRepository(ctrl),
				employeeRepository: NewMockEmployeeRepository(ctrl),
				logger:             ditzap.NewMockLogger(ctrl),
			}

			tt.want(tt.args, f)

			uc := NewAuthUseCase(
				f.sudirClient,
				f.kadryClient,
				f.stateRepository,
				f.tokenRepository,
				f.employeeRepository,
				f.logger,
			)

			uc.Logout(tt.args.ctx, tt.args.cloudID)
		})
	}
}
*/

func Test_authUseCase_GetEmployees(t *testing.T) {
	type fields struct {
		sudirClient        *MockSudirClient
		kadryClient        *MockKadryClient
		stateRepository    *MockStateRepository
		tokenRepository    *MockTokenRepository
		employeeRepository *MockEmployeeRepository
		logger             *ditzap.MockLogger
	}
	type args struct {
		ctx    context.Context
		params entity.EmployeeGetParams
	}

	ctx := context.TODO()
	testPersonID := uuid.New()
	testErr := errors.New("some test error")

	tests := []struct {
		name string
		args args
		want func(a args, f fields) ([]entity.EmployeeInfo, error)
	}{
		{
			name: "invalid key type",
			args: args{
				ctx:    ctx,
				params: entity.EmployeeGetParams{},
			},
			want: func(a args, f fields) ([]entity.EmployeeInfo, error) {
				f.employeeRepository.EXPECT().Get(a.ctx, a.params.ParamByKey()).Return(nil, testErr)
				return nil, diterrors.NewValidationError(ErrInvalidKeyType)
			},
		},
		{
			name: "correct get from cache by cloud id",
			args: args{
				ctx: ctx,
				params: entity.EmployeeGetParams{
					Key:     entity.EmployeeGetKeyCloudID,
					CloudID: "testCloudID",
				},
			},
			want: func(a args, f fields) ([]entity.EmployeeInfo, error) {
				testEmployees := []entity.EmployeeInfo{{}}
				f.employeeRepository.EXPECT().Get(a.ctx, a.params.CloudID).Return([]entity.EmployeeInfo{{}}, nil)
				return testEmployees, nil
			},
		},
		{
			name: "correct get from cache by email",
			args: args{
				ctx: ctx,
				params: entity.EmployeeGetParams{
					Key:     entity.EmployeeGetKeyEmail,
					CloudID: "testEmail",
				},
			},
			want: func(a args, f fields) ([]entity.EmployeeInfo, error) {
				testEmployees := []entity.EmployeeInfo{{}}
				f.employeeRepository.EXPECT().Get(a.ctx, a.params.Email).Return([]entity.EmployeeInfo{{}}, nil)
				return testEmployees, nil
			},
		},
		{
			name: "get employees info by cloud id err",
			args: args{
				ctx: ctx,
				params: entity.EmployeeGetParams{
					Key:     entity.EmployeeGetKeyCloudID,
					CloudID: "testCloudID",
				},
			},
			want: func(a args, f fields) ([]entity.EmployeeInfo, error) {
				f.employeeRepository.EXPECT().Get(a.ctx, a.params.CloudID).Return([]entity.EmployeeInfo{}, nil)
				f.kadryClient.EXPECT().GetEmployeesInfo(a.ctx, a.params.CloudID, kadry.PersonID, kadry.InnOrg, kadry.SNILS).Return(nil, testErr)
				f.logger.EXPECT().Error("ошибка запроса к серверу СКС",
					entity.LogModuleUE,
					entity.LogCode("UE_028"),
					zap.String("cloud-id", a.params.CloudID),
					zap.Error(testErr),
				)
				return nil, kadry.ErrEmployeesService
			},
		},
		{
			name: "get person by employee email err",
			args: args{
				ctx: ctx,
				params: entity.EmployeeGetParams{
					Key:     entity.EmployeeGetKeyEmail,
					CloudID: "testEmail",
				},
			},
			want: func(a args, f fields) ([]entity.EmployeeInfo, error) {
				f.employeeRepository.EXPECT().Get(a.ctx, a.params.Email).Return(nil, nil)
				f.employeeRepository.EXPECT().GetPersonIDByEmployeeEmail(a.ctx, a.params.Email).Return(uuid.Nil, testErr)
				f.logger.EXPECT().Error("ошибка получения personID из сервиса сотрудников",
					entity.LogModuleUE,
					entity.LogCode("UE_029"),
					zap.String("email", a.params.Email),
					zap.Error(testErr),
				)
				return nil, fmt.Errorf("ошибка получения personID из сервиса сотрудников: %w", testErr)
			},
		},
		{
			name: "get employees info by person id err",
			args: args{
				ctx: ctx,
				params: entity.EmployeeGetParams{
					Key:     entity.EmployeeGetKeyEmail,
					CloudID: "testEmail",
				},
			},
			want: func(a args, f fields) ([]entity.EmployeeInfo, error) {
				f.employeeRepository.EXPECT().Get(a.ctx, a.params.Email).Return(nil, nil)
				f.employeeRepository.EXPECT().GetPersonIDByEmployeeEmail(a.ctx, a.params.Email).Return(testPersonID, nil)
				f.employeeRepository.EXPECT().GetEmployeesInfoByPersonID(a.ctx, testPersonID).Return(nil, testErr)
				f.logger.EXPECT().Error("ошибка получения данных о сотрудниках из сервиса сотрудников",
					entity.LogModuleUE,
					entity.LogCode("UE_030"),
					ditzap.UUID("person-id", testPersonID),
					zap.String("email", a.params.Email),
					zap.Error(testErr),
				)
				return nil, fmt.Errorf("ошибка получения данных о сотрудниках из сервиса сотрудников: %w", testErr)
			},
		},
		{
			name: "correct by cloud id with employees",
			args: args{
				ctx: ctx,
				params: entity.EmployeeGetParams{
					Key:     entity.EmployeeGetKeyCloudID,
					CloudID: "testCloudID",
				},
			},
			want: func(a args, f fields) ([]entity.EmployeeInfo, error) {
				f.employeeRepository.EXPECT().Get(a.ctx, a.params.CloudID).Return([]entity.EmployeeInfo{}, nil)
				testEmployees := []entity.EmployeeInfo{
					{CloudID: "testCloudID"},
					{CloudID: "invalidTestCloudID"},
				}
				f.kadryClient.EXPECT().GetEmployeesInfo(a.ctx, a.params.CloudID, kadry.PersonID, kadry.InnOrg, kadry.SNILS).Return(testEmployees, nil)
				f.logger.EXPECT().Warn("неверный CloudID из системы кадров",
					zap.String("sudir-id", a.params.CloudID),
					zap.String("cloud-id", testEmployees[1].CloudID),
				)
				testUserEmployees := []entity.EmployeeInfo{{CloudID: "testCloudID"}}
				f.employeeRepository.EXPECT().Save(a.ctx, a.params.CloudID, testUserEmployees).Return(testErr)
				f.logger.EXPECT().Warn("ошибка при сохранении данных о сотрудниках",
					zap.String(a.params.Key.String(), a.params.ParamByKey()),
					zap.Error(testErr),
				)
				return testUserEmployees, nil
			},
		},
		{
			name: "correct by email with employees",
			args: args{
				ctx: ctx,
				params: entity.EmployeeGetParams{
					Key:     entity.EmployeeGetKeyEmail,
					CloudID: "testEmail",
				},
			},
			want: func(a args, f fields) ([]entity.EmployeeInfo, error) {
				f.employeeRepository.EXPECT().Get(a.ctx, a.params.Email).Return(nil, nil)
				f.employeeRepository.EXPECT().GetPersonIDByEmployeeEmail(a.ctx, a.params.Email).Return(testPersonID, nil)
				testEmployees := []entity.EmployeeInfo{
					{CloudID: "testCloudID"},
					{CloudID: "TestCloudID"},
				}
				f.employeeRepository.EXPECT().GetEmployeesInfoByPersonID(a.ctx, testPersonID).Return(testEmployees, nil)
				f.employeeRepository.EXPECT().Save(a.ctx, a.params.Email, testEmployees).Return(testErr)
				f.logger.EXPECT().Warn("ошибка при сохранении данных о сотрудниках",
					zap.String(a.params.Key.String(), a.params.ParamByKey()),
					zap.Error(testErr),
				)
				return testEmployees, nil
			},
		},
		{
			name: "correct by cloud id without employees",
			args: args{
				ctx: ctx,
				params: entity.EmployeeGetParams{
					Key:     entity.EmployeeGetKeyCloudID,
					CloudID: "testCloudID",
				},
			},
			want: func(a args, f fields) ([]entity.EmployeeInfo, error) {
				f.employeeRepository.EXPECT().Get(a.ctx, a.params.CloudID).Return(nil, nil)
				f.kadryClient.EXPECT().GetEmployeesInfo(a.ctx, a.params.CloudID, kadry.PersonID, kadry.InnOrg, kadry.SNILS).Return(nil, nil)
				f.logger.EXPECT().Warn("данные о сотрудниках не найдены в СКС",
					entity.LogModuleUE,
					entity.LogCode("UE_031"),
					zap.String(a.params.Key.String(), a.params.ParamByKey()),
				)
				return make([]entity.EmployeeInfo, 0), nil
			},
		},
		{
			name: "correct by email without employees",
			args: args{
				ctx: ctx,
				params: entity.EmployeeGetParams{
					Key:     entity.EmployeeGetKeyEmail,
					CloudID: "testEmail",
				},
			},
			want: func(a args, f fields) ([]entity.EmployeeInfo, error) {
				f.employeeRepository.EXPECT().Get(a.ctx, a.params.Email).Return(nil, nil)
				f.employeeRepository.EXPECT().GetPersonIDByEmployeeEmail(a.ctx, a.params.Email).Return(testPersonID, nil)
				f.employeeRepository.EXPECT().GetEmployeesInfoByPersonID(a.ctx, testPersonID).Return([]entity.EmployeeInfo{}, nil)
				f.logger.EXPECT().Warn("данные о сотрудниках не найдены в сервисе сотрудников",
					entity.LogModuleUE,
					entity.LogCode("UE_032"),
					zap.String("email", a.params.Email),
				)
				return make([]entity.EmployeeInfo, 0), nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				sudirClient:        NewMockSudirClient(ctrl),
				kadryClient:        NewMockKadryClient(ctrl),
				stateRepository:    NewMockStateRepository(ctrl),
				tokenRepository:    NewMockTokenRepository(ctrl),
				employeeRepository: NewMockEmployeeRepository(ctrl),
				logger:             ditzap.NewMockLogger(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			uc := NewAuthUseCase(
				f.sudirClient,
				f.kadryClient,
				f.stateRepository,
				f.tokenRepository,
				f.employeeRepository,
				f.logger,
			)
			got, err := uc.GetEmployees(tt.args.ctx, tt.args.params)

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
