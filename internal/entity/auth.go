package entity

import "time"

type ScopeType int

const (
	ScopeUnknown ScopeType = iota
	ScopeOpenId
	ScopeProfile
	ScopeEmail
	ScopeUserInfo
	ScopeEmployee
	ScopeGroups
)

type CloudID string

type OAuth struct {
	AccessToken  string
	RefreshToken string
	Expiry       *time.Time
}

type AuthInfo struct {
	OAuth  *OAuth
	User   *User
	Device *Device
}

type Device struct {
	ID        string
	ClientID  string
	UserAgent string
}

type User struct {
	CloudID   CloudID
	Info      *UserInfo
	Employees []EmployeeInfo
}

type UserInfo struct {
	SessionID  string
	Sub        string
	LastName   string
	FirstName  string
	MiddleName string
	LogonName  string
	Company    string
	Department string
	Position   string
	Email      string
}

type TokenInfo struct {
	// субъект, которому выдан токен
	Subject string
	// перечень разрешений токена
	Scopes []ScopeType
	// тип токена
	TokenType string
	// идентификатор системы, которой был выдан токен
	ClientID string
	// статус активности токена
	IsActive bool
	// время, когда токен станет невалидным
	ExpirationTime time.Time
	// время, в которое был выдан токен
	IssuedAt time.Time
}
