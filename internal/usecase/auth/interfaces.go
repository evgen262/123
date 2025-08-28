package auth

import (
	"context"
	"net"

	entityAuth "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/auth"
)

//go:generate mockgen -source=interfaces.go -destination=./auth_mock.go -package=auth

type Repository interface {
	GetRedirectURL(ctx context.Context, callbackURI string) (string, error)
	Auth(ctx context.Context, code, state, callbackURI string) (*entityAuth.AuthSudir, error)
	AuthPortal(ctx context.Context, params entityAuth.AuthPortalParams) (*entityAuth.Auth1C, error)
	// CreateSession метод для создания сессии. На вход принимает пользователя СУДИР, его девайс и идентификатор сессии портала 1С. Возвращает Access и Refresh токены
	CreateSession(ctx context.Context, user *entityAuth.UserSudir, clientIP net.IP, device *entityAuth.Device, auth *entityAuth.Auth1C) (entityAuth.TokensPair, error)
	GetSession(ctx context.Context, accessToken string) (*entityAuth.Session, error)
	Logout(ctx context.Context, session *entityAuth.Session, accessToken, refreshToken string) error
	// ChangePortal метод для смены активного портала. На вход принимает идентификатор выбранного портала 1С и сессию. Возвращает порталы и сессию портала 1С
	ChangePortal(ctx context.Context, portalID int, session *entityAuth.Session) ([]*entityAuth.Portal, string, error)
	RefreshTokensPair(ctx context.Context, accessToken, refreshToken string) (*entityAuth.TokensPair, error)
}

type RedirectSessionRepository interface {
	CreateSession(ctx context.Context, userInfo *entityAuth.RedirectSessionUserInfo) (string, error)
}
