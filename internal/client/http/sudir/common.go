package sudir

import "git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"

const (
	ScopeOpenId   = "openid"
	ScopeProfile  = "profile"
	ScopeEmail    = "email"
	ScopeUserInfo = "userinfo"
	ScopeEmployee = "employee"
	ScopeGroups   = "groups"
)

const (
	ErrStatusTextServiceError  = "ошибка при работе с сервисом СУДИР"
	ErrStatusTextNoIdToken     = "отсутствует маркер идентификации"
	ErrStatusTextInvalidGrant  = "предоставлен неверный код авторизации"
	ErrStatusTextInvalidClient = "предоставлен неверные данные для авторизации"
	ErrStatusAccessDenied      = "доступ запрещён"
)

var (
	ErrNoJWTToken    diterrors.StringError = ErrStatusTextServiceError + ": " + ErrStatusTextNoIdToken
	ErrInvalidGrant  diterrors.StringError = ErrStatusTextInvalidGrant
	ErrInvalidClient diterrors.StringError = ErrStatusTextInvalidClient
	ErrAccessDenied  diterrors.StringError = ErrStatusAccessDenied
)
