package employee

import (
	"time"

	timeUtils "git.mos.ru/buch-cloud/moscow-team-2.0/build/time-utils.git"
	"github.com/google/uuid"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity"
)

//go:generate ditgen -source=person.go

// Person информация о пользователе как о физическом лице
type Person struct {
	// ID идентификатор
	ID uuid.UUID
	// ExtID Внешний идентификатор (в базах 1С)
	ExtID string
	// CloudID идентификатор в СУДИР
	CloudID string
	// LastName фамилия
	LastName string
	// FirstName имя
	FirstName string
	// MiddleName отчество
	MiddleName string
	// INN ИНН
	INN string
	// SNILS СНИЛС
	SNILS string
	// Gender пол
	Gender entity.Gender
	// Birthday дата рождения
	Birthday timeUtils.Date
	// Phone личный телефон (мобильный)
	Phone string
	// Socials контакты
	Socials *Socials
	// ImageID фото (идентификатор)
	ImageID string
	// IsActive флаг активности
	IsActive bool
	// CreatedAt дата создания
	CreatedAt *time.Time
	// CreatedAt дата обновления
	UpdatedAt *time.Time
}

// Socials информация о контактах физического лица
type Socials struct {
	// Telegram телеграм
	Telegram string
	// ChatID идентификатор пользователя чата
	ChatID *string
}
