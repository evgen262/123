package sudir

import (
	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	// SID идентификатор сессии пользователя
	SID string `json:"sid,omitempty"`
	// Cloud ID пользователя
	CloudGUID string `json:"cloudGUID,omitempty"`
	// Оагинизация
	Company string `json:"company,omitempty"`
	// Подразделение
	Department string `json:"department,omitempty"`
	// Email
	Email string `json:"email,omitempty"`
	// Имя для входа
	LogonName string `json:"logonname,omitempty"`
	// Должность
	Position string `json:"position,omitempty"`
	// Фамилия
	FamilyName string `json:"family_name,omitempty"`
	// Имя
	Name string `json:"name,omitempty"`
	// Отчество
	MiddleName string `json:"middle_name,omitempty"`
}

type JWTPayload struct {
	jwt.RegisteredClaims

	UserClaims
}
