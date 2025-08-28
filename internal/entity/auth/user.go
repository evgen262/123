package auth

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap/zapcore"
)

//go:generate ditgen --source=user.go

// User пользователь сессии
type User struct {
	// ID идентификатор пользователя
	ID uuid.UUID
	// Идентификатор СУДИР.
	CloudID string
	// Login имя пользователя
	Login string
	// Email пользователя
	Email string
	// SNILS СНИЛС пользователя
	SNILS string
	// Employee информация о пользователе как о сотруднике
	Employee *Employee
	// Person информация о пользователе как о физическом лице
	Person *Person
}

func (u *User) GetLogin() string {
	if u == nil {
		return ""
	}
	return u.Login
}

func (u *User) GetEmail() string {
	if u == nil {
		return ""
	}
	return u.Email
}

// MarshalLogObject маршаллер для добавления пользователя в логи
func (u *User) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if u == nil {
		return errors.New("user is nil")
	}

	enc.AddString("id", u.ID.String())
	enc.AddString("cloud_id", u.CloudID)
	enc.AddString("login", u.Login)
	enc.AddString("email", u.Email)
	enc.AddString("SNILS", u.SNILS)
	if employee := u.GetEmployee(); employee != nil {
		if err := enc.AddObject("employee", employee); err != nil {
			enc.AddString("employee", fmt.Sprintf("can't marshal employee: %s", err.Error()))
		}
	}
	if person := u.GetPerson(); person != nil {
		if err := enc.AddObject("person", person); err != nil {
			enc.AddString("person", fmt.Sprintf("can't marshal person: %s", err.Error()))
		}
	}

	return nil
}

type UserSudir struct {
	// Идентификатор СУДИР.
	CloudID string
	// Имя для входа.
	Login string
	// Email сотрудника.
	Email string
	// ФИО сотрудника
	FIO string
	// Фамилия сотрудника
	LastName string
	// Имя сотрудника
	FirstName string
	// Отчество сотрудника
	MiddleName string
	// СНИЛС.
	SNILS string
	// Идентификатор сессии СУДИР
	SID string
	// Идентификатор в домене.
	//  Заполняется при вторичном входе
	Sub string
	// Порталы доступные пользователю.
	Portals []*Portal
	// Информация о сотруднике
	Employees []*EmployeeInfo
}

type User1C struct {
	// CloudID идентификатор пользователя в СУДИР.
	CloudID string
	// SNILS снилс сотрудника.
	SNILS string
	// email пользователя.
	Email string
}

type RedirectSessionUserInfo struct {
	SessionID string
	Email     string
	SNILS     string
	PortalURL string
	TargetURL string
	UserAgent string
	IP        string
}
