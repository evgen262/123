package portalv2

import "time"

// Комплекс.
type Complex struct {
	// Идентификатор.
	ID int
	// Наименование.
	Name string
	// Описание.
	Description *string
	// Идентификатор изображения.
	ImageID *string
	// Группа комплекса.
	ComplexGroup int
	// Сортировка внутри группы.
	Sort int
	// Флаг неактивности.
	IsDisabled bool
	// Ответственный за комплекс.
	Responsible *ComplexResponsible
	// Порталы
	Portals []*ComplexPortal
	// Дата/время создания.
	CreatedAt time.Time
	// Дата/время изменения.
	UpdatedAt *time.Time
	// Дата/время удаления.
	DeletedAt *time.Time
}

// Портал в комплексе.
type ComplexPortal struct {
	// Идентификатор портала.
	ID int
	// Сортировка.
	Sort int
}

// Ответственный за комплекс.
type ComplexResponsible struct {
	// Имя ответственного.
	FirstName string
	// Фамилия ответственного.
	LastName string
	// Отчество ответственного.
	MiddleName *string
	// Идентификатор изображения руководителя.
	ImageID *string
	// Описание должности ответственного.
	Description string
}

type GetComplexOptions struct {
	// Учитывать удаленные комплексы.
	IsDeleted bool
	// Учитывать неактивные комплексы.
	IsDisabled bool
}

type FilterComplexesFilters struct {
	// Идентификаторы комплексов.
	IDs []int
	// Идентификаторы порталов.
	PortalIDs []int
}

type FilterComplexesOptions struct {
	// Добавить в выдачу удаленные комплексы.
	WithDeleted bool
	// Добавить в выдачу неактивные комплексы.
	WithDisabled bool
}
