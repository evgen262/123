package usecase

import (
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
)

var (
	ErrTokenIsExpire       diterrors.StringError = "срок действия токена истек"
	ErrCallbackURLMismatch diterrors.StringError = "callback_url не совпадает с callback_url в запросе"
	ErrInvalidKeyType      diterrors.StringError = "невалидный ключ получения сотрудников"
)
