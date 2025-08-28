package auth

import (
	"errors"
	"fmt"
	"slices"
	"time"

	"go.uber.org/zap/zapcore"
)

//go:generate ditgen -source=auth.go

type Token struct {
	Value     string
	ExpiredAt *time.Time
}

// TokensPair пара токенов
type TokensPair struct {
	// Access токен
	AccessToken Token
	// Refresh токен
	RefreshToken Token
}

type Portal struct {
	ID         int
	Name       string
	URL        string
	Image      string
	IsSelected bool
}

func (p *Portal) GetID() int {
	if p == nil {
		return 0
	}
	return p.ID
}

// MarshalLogObject маршаллер для добавления портала в логи
func (p *Portal) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if p == nil {
		return errors.New("portal is nil")
	}

	enc.AddInt("id", p.ID)
	enc.AddString("name", p.Name)
	enc.AddString("url", p.URL)
	enc.AddString("image", p.Image)
	enc.AddBool("is_selected", p.IsSelected)

	return nil
}

// ActivePortal активный портал пользователя
type ActivePortal struct {
	// Portal портал пользователя
	Portal Portal
	// SID идентификатор сессии 1С
	SID string
}

func (ap *ActivePortal) GetPortal() Portal {
	if ap == nil {
		return Portal{}
	}
	return ap.Portal
}

func (ap *ActivePortal) GetPortalID() int {
	if ap == nil {
		return 0
	}
	return ap.Portal.ID
}

// MarshalLogObject маршаллер для добавления активного портала в логи
func (ap *ActivePortal) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if ap == nil {
		return errors.New("active portal is nil")
	}

	if err := enc.AddObject("portal", &ap.Portal); err != nil {
		enc.AddString("portal", fmt.Sprintf("can't marshal portal: %s", err.Error()))
	}
	enc.AddString("sid", ap.SID)

	return nil
}

type Auth struct {
	JWTToken      Token
	RefreshToken  Token
	Portals       []*Portal
	PortalSession string
}

type Auth1C struct {
	PortalSession string
	EmployeeID    string
	PersonID      string
}

func (a *Auth1C) GetPortalSession() string {
	if a == nil {
		return ""
	}
	return a.PortalSession
}

func (a *Auth1C) GetEmployeeID() string {
	if a == nil {
		return ""
	}
	return a.EmployeeID
}

func (a *Auth1C) GetPersonID() string {
	if a == nil {
		return ""
	}
	return a.PersonID
}

type AuthSudir struct {
	AccessToken  string
	RefreshToken string
	User         *UserSudir
}

type AuthPortalParams struct {
	PortalURL     string
	PortalSession string
	Device        *Device
	User          User1C
}

type AccessList []string

func (al AccessList) Have(s string) bool {
	return slices.Contains(al, s)
}
