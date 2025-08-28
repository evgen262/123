package auth

import (
	"errors"

	"go.uber.org/zap/zapcore"
)

// Employee сотрудник сессии
type Employee struct {
	// TODO: добавить идентификатор из employees
	ID string
	// ExtID Внешний идентификатор (в базах 1С)
	ExtID string
}

func (e *Employee) GetExtID() string {
	if e == nil {
		return ""
	}
	return e.ExtID
}

// MarshalLogObject маршаллер для добавления сотрудника сессии в логи
func (e *Employee) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if e == nil {
		return errors.New("employee is nil")
	}

	enc.AddString("ext_id", e.ExtID)

	return nil
}

// Person физическое лицо сессии
type Person struct {
	// ExtID Внешний идентификатор (в базах 1С)
	ExtID string
}

// MarshalLogObject маршаллер для добавления физического лица в сессии в логи
func (p *Person) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if p == nil {
		return errors.New("person is nil")
	}

	enc.AddString("ext_id", p.ExtID)

	return nil
}

type EmployeeInfo struct {
	Inn   string
	OrgID string
}
