package user

import (
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity"
)

// UserInfo информация о пользователе
type UserInfo struct {
	User ShortUser
}

// ShortUser краткий пользователя
type ShortUser struct {
	// TODO: ID реализовать когда появится сервис пользователей
	// ID идентификатор пользователя
	// ID string
	// LastName фамилия пользователя
	LastName string
	// FirstName имя пользователя
	FirstName string
	// MiddleName отчество пользователя
	MiddleName string
	// ImageID фото пользователя (идентификатор)
	ImageID string
	// Gender пол пользователя
	Gender     entity.Gender
	PortalData PortalData
}

type PortalData struct {
	// PersonID идентификатор физ-лица на портале
	PersonID string
	// EmployeeID идентификатор сотрудника на портале
	EmployeeID string
}
