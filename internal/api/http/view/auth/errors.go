package auth

import (
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
)

// Ошибки от и для Фронтов
const (
	ErrMessageFrontInternal diterrors.StringError = "Ошибка сервера"

	/*
		Ошибки валидации
	*/

	ErrMessageFrontInvalidParams        diterrors.StringError = "Переданы неверные данные"
	ErrMessageFrontInvalidSUDIRRedirect diterrors.StringError = "Неверный код авторизации"
	ErrMessageFrontSamePortals          diterrors.StringError = "Выбранный и активный портал совпадают"
	ErrMessageFrontEmptyPortalURL       diterrors.StringError = "У выбранного портала отсутствует URL"

	/*
		Ошибки наличия данных или связи данных
	*/

	ErrMessageFrontSUDIRUserNotFound     diterrors.StringError = "Нет данных о сотруднике"
	ErrMessageFrontSKSEmployeeNotFound   diterrors.StringError = "В кадровой службе данные не найдены"
	ErrMessageFrontPortalForUserNotFound diterrors.StringError = "Нет портала, связанного с организацией сотрудника"
	ErrMessageFrontUserNotIntoPortal     diterrors.StringError = "Сотрудник не найден на портале"

	/*
		Ошибки связанные с ограничением доступов
	*/

	ErrMessageFrontUnauthenticated   diterrors.StringError = "Необходима аутентификация"
	ErrMessageFrontUserAccessDenied  diterrors.StringError = "Доступ к порталу ограничен"
	ErrMessageFrontUnavailablePortal diterrors.StringError = "Выбранный портал недоступен пользователю"
)
