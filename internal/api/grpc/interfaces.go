package grpc

import (
	"context"

	authv1 "git.mos.ru/buch-cloud/moscow-team-2.0/pud/auth.git/internal/api/grpc/gen/infogorod/auth/auth/v1"
	"git.mos.ru/buch-cloud/moscow-team-2.0/pud/auth.git/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./grpc_mock.go -package=grpc

type AuthInteractor interface {
	// GetAuthURL
	//  URL для перенаправления пользователя для авторизации в СУДИР
	GetAuthURL(ctx context.Context, callbackURL, clientID, clientSecret string) (string, error)
	// Auth аутентификация пользователя
	//  метод возвращает информацию о пользователе в СУДИР
	//  и oauth2 токены
	Auth(ctx context.Context, code, state, callbackURL string) (*entity.AuthInfo, error)
	// LoginByCredentials вход в СУДИР по динамическим реквизитам
	LoginByCredentials(ctx context.Context, clientID, clientSecret string) (*entity.AuthInfo, error)
	// RefreshToken обновление access_token
	RefreshToken(ctx context.Context, id, sessionID string) (string, error)
	// Logout очистка токенов пользователя и выход из СУДИР
	Logout(ctx context.Context, id, sessionID, registrationToken string) error

	GetUserInfo(ctx context.Context, accessToken string) (*entity.UserInfo, error)
	GetEmployees(ctx context.Context, params entity.EmployeeGetParams) ([]entity.EmployeeInfo, error)
	// IsValidToken проверка access token в СУДИР.
	IsValidToken(ctx context.Context, accessToken string) (*entity.TokenInfo, error)
}

type AuthPresenter interface {
	UserToPb(user *entity.User) *authv1.User
	EmployeesToPb(entities []entity.EmployeeInfo) []*authv1.Employee
	UserInfoToPb(user *entity.UserInfo) *authv1.UserInfo
	DeviceToPb(device *entity.Device) *authv1.AuthResponse_Device
	TokenInfoToPb(info *entity.TokenInfo) *authv1.TokenInfo
}
