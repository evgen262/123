package view

import (
	"time"

	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
)

type Response struct {
	Error *ErrorResponse `json:"error"`
	Data  any            `json:"data"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	// TODO: появится businessDescription или код
} // @name ErrorResponse

func NewErrorResponse(err error) *Response {
	return &Response{
		Error: &ErrorResponse{
			Message: err.Error(),
		},
	}
}

func NewSuccessResponse(data any) *Response {
	return &Response{
		Data: data,
	} // @name NewSuccessResponse
}

// StringPagination структура для представления пагинации
type StringPagination struct {
	LastId          string     `json:"last_id,omitempty"`
	Limit           *int       `json:"limit,omitempty"`
	LastCreatedTime *time.Time `json:"last_created_time,omitempty"`
	Total           *int       `json:"total,omitempty"`
} // @name StringPagination

func (p *StringPagination) GetLastCreatedTime() *time.Time {
	if p == nil {
		return nil
	}
	return p.LastCreatedTime
}

func (p *StringPagination) GetLimit() *int {
	if p == nil {
		return nil
	}
	return p.Limit
}

func (p *StringPagination) GetTotal() *int {
	if p == nil {
		return nil
	}
	return p.Total
}

// Общие ошибки
const (
	ErrMessageUnauthenticated  diterrors.StringError = "Вы не авторизованы"
	ErrMessageUnauthorized     diterrors.StringError = "Доступ запрещен"
	ErrMessageMethodNotAllowed diterrors.StringError = "Метод не допустим"
	ErrMessageNotFound         diterrors.StringError = "Элемент не найден"
	ErrMessageInternalError    diterrors.StringError = "Что-то пошло не так, попробуйте позднее"
	ErrMessageInvalidId        diterrors.StringError = "Некорректный ID"
	ErrMessageInvalidRequest   diterrors.StringError = "Некорректный запрос"
	ErrPermissionDenied        diterrors.StringError = "Доступ запрещен"
)

// Ошибки аутентификации
const (
	ErrMessageUserNotFound     diterrors.StringError = "Пользователь не найден"
	ErrMessageCodeNotValid     diterrors.StringError = "Неверный код аутентификации"
	ErrMessageCodeNotProvided  diterrors.StringError = "Не предоставлен код аутентификации"
	ErrMessageUserAccessDenied diterrors.StringError = "Доступ к порталу ограничен"
)

const (
	ErrMessagePortalsNotFound diterrors.StringError = "Порталы не найдены"
)

// Ошибки работы с файлами
const (
	ErrMessageFileNotFound diterrors.StringError = "Файл не найден"
)

// Ошибка при взаимодейтвии с сервисом пользователей
const (
	ErrMessageUserIDExists diterrors.StringError = "Пользователь с таким id уже существует"
)
