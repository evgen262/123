package auth

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap/zapcore"
)

//go:generate ditgen --source=session.go

// SessionID идентификатор сессии
type SessionID uuid.UUID

// UUID возвращает идентификатор сессии в виде uuid.UUID
func (sid *SessionID) UUID() uuid.UUID {
	if sid == nil {
		return uuid.Nil
	}
	return (uuid.UUID)(*sid)
}

// String возвращает строковую форму идентификатора сессии или пустую строку
func (sid *SessionID) String() string {
	if sid == nil {
		return ""
	}
	return sid.UUID().String()
}

// UserAuthType тип аутентификации пользователя
type UserAuthType int

const (
	// UserAuthTypeInvalid невалидный тип пользователя
	UserAuthTypeInvalid = iota
	// UserAuthTypeAnon анонимный тип пользователя
	UserAuthTypeAnon
	// UserAuthTypeAuth тип пользователя после универсального входа
	UserAuthTypeAuth
	// UserAuthTypeOldAuth тип пользователя после старой аутентификации
	UserAuthTypeOldAuth
	// UserAuthTypeService тип пользователя сервис
	UserAuthTypeService
)

// Session данные для создания сессии и получения access и refresh токенов
type Session struct {
	// ID идентификатор сессии
	ID *SessionID
	// User пользователь
	User *User
	// UserAuthType тип аутентификации пользователя
	UserAuthType UserAuthType
	// UserIP IP адрес пользователя
	UserIP net.IP
	// Device устройство пользователя
	Device *Device
	// ActivePortal активный портал пользователя
	ActivePortal *ActivePortal
	// Issuer сервис запрашивающий токены
	Issuer string
	// Subject субъект, которому выдан токен
	Subject string
	// LastActiveTime время последней активности сессии
	LastActiveTime *time.Time
	// AccessExpiredTime время истечения срока действия токена
	AccessExpiredTime time.Time
	// RefreshExpiredTime время истечения срока для обновления сессии
	RefreshExpiredTime time.Time
	// CreatedTime время создания сессии
	CreatedTime time.Time
	// RefreshedTime время последнего обновления сессии
	RefreshedTime time.Time
	// IsActive активна ли сессия
	IsActive bool
}

// MarshalLogObject маршаллер для добавления сессии в логи
func (s *Session) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if s == nil {
		return errors.New("session is nil")
	}

	enc.AddString("id", s.GetID().String())
	if user := s.GetUser(); user != nil {
		if err := enc.AddObject("user", user); err != nil {
			enc.AddString("user", fmt.Sprintf("can't marshal user: %s", err.Error()))
		}
	}
	enc.AddInt("user_type", int(s.UserAuthType))
	enc.AddString("user_ip", s.UserIP.String())
	if device := s.GetDevice(); device != nil {
		if err := enc.AddObject("device", device); err != nil {
			enc.AddString("device", fmt.Sprintf("can't marshal device: %s", err.Error()))
		}
	}
	if portal := s.GetActivePortal(); portal != nil {
		if err := enc.AddObject("active_portal", portal); err != nil {
			enc.AddString("active_portal", fmt.Sprintf("can't marshal active portal: %s", err.Error()))
		}
	}
	enc.AddString("issuer", s.Issuer)
	enc.AddString("subject", s.Subject)
	if s.LastActiveTime != nil {
		enc.AddTime("last_active_time", *s.LastActiveTime)
	}
	enc.AddTime("access_expired_time", s.AccessExpiredTime)
	enc.AddTime("refresh_expired_time", s.RefreshExpiredTime)
	enc.AddTime("created_time", s.CreatedTime)
	enc.AddTime("refreshed_time", s.RefreshedTime)
	enc.AddBool("is_active", s.IsActive)
	return nil
}
