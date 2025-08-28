package repositories

import "git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"

const (
	ErrNotFound               diterrors.StringError = "записей не найдено"
	ErrNilSession             diterrors.StringError = "session is nil"
	ErrEmptyAccessToken       diterrors.StringError = "access token is empty"
	ErrEmptyRefreshToken      diterrors.StringError = "refresh token is empty"
	ErrNilSessionActivePortal diterrors.StringError = "session active portal is nil"
	ErrNilSessionUser         diterrors.StringError = "session user is nil"
	ErrNilSessionSudirInfo    diterrors.StringError = "session sudir info is nil"
	ErrNoUserPortals          diterrors.StringError = "user has no portals"
)

type DetailsError struct {
	// Поле, в котором передано некорректное значение
	field string
	// Текст ошибки
	message string
	// Необходимость пройти аутентификацию
	reauthRequired bool
}

func (e *DetailsError) Error() string {
	if e == nil {
		return ""
	}

	return e.message
}

func (e *DetailsError) GetField() string {
	if e == nil {
		return ""
	}

	return e.field
}

func (e *DetailsError) GetMessage() string {
	return e.Error()
}

func (e *DetailsError) GetReauthRequired() bool {
	if e == nil {
		return false
	}

	return e.reauthRequired
}

func NewDetailsError(field, message string, reauth bool) *DetailsError {
	return &DetailsError{
		field:          field,
		message:        message,
		reauthRequired: reauth,
	}
}
