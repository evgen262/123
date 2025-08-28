package auth

import (
	"errors"
	"fmt"

	"go.uber.org/zap/zapcore"
)

//go:generate ditgen --source=device.go

type Device struct {
	// UserAgent устройства
	UserAgent string
	UserLogin string
	// SudirInfo служебные поля СУДИР
	SudirInfo *SudirInfo
}

func (d *Device) GetUserAgent() string {
	if d == nil {
		return ""
	}
	return d.UserAgent
}

func (d *Device) GetUserLogin() string {
	if d == nil {
		return ""
	}
	return d.UserLogin
}

// MarshalLogObject маршаллер для добавления устройства в логи
func (d *Device) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if d == nil {
		return errors.New("device is nil")
	}

	enc.AddString("user_agent", d.UserAgent)
	enc.AddString("user_login", d.UserLogin)
	if info := d.GetSudirInfo(); info != nil {
		if err := enc.AddObject("sudir_info", info); err != nil {
			enc.AddString("sudir_info", fmt.Sprintf("can't marshal sudir info: %s", err.Error()))
		}
	}

	return nil
}

// SudirInfo служебные поля СУДИР
type SudirInfo struct {
	// SID идентификатор сессии СУДИР
	SID string
	// ClientID инстанса динамических реквизитов входа СУДИР
	ClientID string
}

// MarshalLogObject маршаллер для добавления служебных полей СУДИР в логи
func (s *SudirInfo) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if s == nil {
		return errors.New("sudir info is nil")
	}

	enc.AddString("sid", s.SID)
	enc.AddString("client_id", s.ClientID)

	return nil
}
