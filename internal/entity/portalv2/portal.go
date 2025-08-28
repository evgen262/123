package portalv2

import (
	"time"
)

// Статусы портала.
type PortalStatus int

const (
	PortalStatusInvalid PortalStatus = iota
	// Портал активен
	PortalStatusActive
	// Портал неактивен
	PortalStatusInactive
	// Портал в режиме обслуживания
	PortalStatusMaintenance
)

// Портал.
type Portal struct {
	// Идентификатор
	ID int
	// Полное наименование
	Name string
	// Короткое наименование
	ShortName string
	// Адрес портала, без схемы
	Url string
	// Идентификатор изображения
	ImageID *string
	// Статус
	Status PortalStatus
	// Дата/время обновления статуса
	StatusUpdatedTime *time.Time
	// Флаг неактивности
	IsDisabled bool
	// Сортировка
	Sort int
	// Руководитель портала
	Manager *PortalManager
	// Организации
	Organizations []*Organization
	// Дата/время создания
	CreatedAt time.Time
	// Дата/время обновления
	UpdatedAt *time.Time
	// Дата/время удаления
	DeletedAt *time.Time
}

// Портал с подсчетом количества организаций и сотрудников.
type PortalWithCounts struct {
	// Портал
	Portal *Portal
	// Количество организаций
	OrgsCount int
	// Количество сотрудников
	EmployeesCount int
}

// Руководитель портала.
type PortalManager struct {
	// Имя
	FirstName string
	// Фамилия
	LastName string
	// Отчество
	MiddleName *string
	// Должность
	Prosition string
	// Идентификатор изображения
	ImageID *string
}

// Фильтры для получения списка порталов
type FilterPortalsFilters struct {
	// Идентификаторы порталов
	IDs []int
}

// Параметры для получения списка порталов
type FilterPortalsOptions struct {
	// Флаг включения удаленных порталов
	WithDeleted bool
	// Флаг включения неактивных порталов
	WithDisabled bool
	// Флаг включения подсчета количества сотрудников
	WithEmployeesCount bool
}
