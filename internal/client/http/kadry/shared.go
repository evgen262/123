package kadry

import (
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
)

var ErrEmployeesService diterrors.StringError = "ошибка при работе с системой кадров"

const (
	ErrStatusTextBadRequest   = "неверные параметры для запроса"
	ErrStatusTextErrorRequest = "ошибка запроса к серверу СКС"
	ErrStatusTextServiceError = "неверный ответ от сервера СКС"
)
