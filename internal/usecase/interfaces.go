package usecase

import (
	"context"

	"git.mos.ru/buch-cloud/moscow-team-2.0/pud/auth.git/internal/client/http/kadry"
	"git.mos.ru/buch-cloud/moscow-team-2.0/pud/auth.git/internal/client/http/sudir"
	"git.mos.ru/buch-cloud/moscow-team-2.0/pud/auth.git/internal/entity"
	"github.com/google/uuid"
)

//go:generate mockgen -source=interfaces.go -destination=./usecases_mock.go -package=usecase

// KadryClient интерфейс клиента СКС
type KadryClient interface {
	// GetEmployeesInfo
	//  возвращает информацию о сотрудниках в СКС
	GetEmployeesInfo(ctx context.Context, cloudId string, attributes ...kadry.AttributeName) ([]entity.EmployeeInfo, error)
}

// SudirClient интерфейс клиента СУДИР
type SudirClient interface {
	// AuthURL
	//  URL для авторизации пользователя в СУДИР
	AuthURL(options sudir.AuthURLOptions) string
	// CodeExchange
	//  проверка кода авторизации и обмен его на токены
	CodeExchange(ctx context.Context, code string, options sudir.CodeExchangeOptions) (*sudir.OAuthResponse, error)
	// LoginCredentials
	//  авторизация по client_id client_secret мобильного приложения
	LoginCredentials(ctx context.Context, options sudir.LoginOptions) (*sudir.OAuthResponse, error)
	// RefreshToken
	//  проверка валидности токена и обмен его на новый
	RefreshToken(ctx context.Context, refreshToken string) (*sudir.OAuthResponse, error)
	// ParseToken
	//  получение payload из jwt токена
	ParseToken(token string) (*sudir.JWTPayload, error)

	GetUserInfo(ctx context.Context, accessToken string) (*sudir.UserInfo, error)
	ValidateToken(ctx context.Context, accessToken string) (*sudir.ValidationInfo, error)
	Logout(ctx context.Context, clientID, registrationToken string) error
}

// StateRepository репозиторий state
type StateRepository interface {
	New(ctx context.Context, options *entity.StateOptions) (*entity.State, error)
	IsExists(ctx context.Context, stateID string) error
	Get(ctx context.Context, stateID string) (*entity.State, error)
	Delete(ctx context.Context, stateID string)
}

// TokenRepository репозиторий refresh_token
type TokenRepository interface {
	Save(ctx context.Context, id, token string) error
	Get(ctx context.Context, id string) (string, error)
	Delete(ctx context.Context, id string)
}

// EmployeeRepository репозиторий employees
type EmployeeRepository interface {
	Save(ctx context.Context, key string, employees []entity.EmployeeInfo) error
	Get(ctx context.Context, key string) ([]entity.EmployeeInfo, error)
	GetPersonIDByEmployeeEmail(ctx context.Context, email string) (uuid.UUID, error)
	GetEmployeesInfoByPersonID(ctx context.Context, personID uuid.UUID) ([]entity.EmployeeInfo, error)
}
