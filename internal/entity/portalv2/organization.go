package portalv2

import "time"

// Организация.
type Organization struct {
	// Идентификатор организации
	ID int
	// Идентификатор организации в сервисе такси
	TaxiID int
	// Внешний идентификатор из 1С.
	ExtID string
	// Идентификатор портала
	PortalID int
	// Наименование
	Name string
	// Полное наименование
	FullName *string
	// Короткое наименование
	ShortName *string
	// ИНН
	INN string
	// Адрес
	Address *string
	// Телефон
	Phone *string
	// Флаг неактивности
	IsDisabled bool
	// Руководитель организации
	Manager *OrganizationManager
	// Идентификатор изображения
	ImageID *string
	// Дата/время создания
	CreatedAt time.Time
	// Дата/время обновления
	UpdatedAt *time.Time
	// Дата/время удаления
	DeletedAt *time.Time
}

// Руководитель организации.
type OrganizationManager struct {
	// Имя
	FirstName string
	// Фамилия
	LastName string
	// Отчество
	MiddleName *string
	// Идентификатор изображения
	ImageID *string
}

// Параметры получения организации.
type GetOrganizationOptions struct {
	// Флаг удаленности
	IsDeleted bool
	// Флаг неактивности
	IsDisabled bool
}

// Параметры фильтрации организаций.
type FilterOrganizationsFilters struct {
	// Список идентификаторов организаций
	IDs []int
	// Список наименований организаций
	Names []string
	// Список ИНН
	INNs []string
	// Список идентификаторов организаций в сервисе такси
	TaxiIDs []int
	// Список внешних идентификаторов из 1С
	ExtIDs []string
	// Список идентификаторов порталов
	PortalIDs []int
	// Список URL-адресов порталов
	PortalURLs []string
}

// Параметры фильтрации организаций.
type FilterOrganizationsOptions struct {
	// Флаг удаленности
	WithDeleted bool
	// Флаг неактивности
	WithDisabled bool
}
