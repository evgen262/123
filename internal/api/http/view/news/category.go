package news

import (
	"time"

	"github.com/google/uuid"
)

//go:generate ditgen -source=category.go -zero=true -all-fields=true

type NewCategory struct {
	Name       string              `json:"title"`
	AuthorID   uuid.UUID           `json:"author"`
	Visibility *CategoryVisibility `json:"visibilityConfig"`
}

type UpdateCategory struct {
	ID         uuid.UUID           `json:"id"`
	Name       string              `json:"title"`
	UpdatedAt  *time.Time          `json:"updatedAt,omitempty"`
	AuthorID   uuid.UUID           `json:"authorID"`
	Visibility *CategoryVisibility `json:"visibilityConfig"`
}

type CategoryResult struct {
	ID        uuid.UUID  `json:"id"`
	UpdatedAt *time.Time `json:"updatedAt"`
}

type Category struct {
	ID         uuid.UUID           `json:"id"`
	Name       string              `json:"title"`
	PortalID   int                 `json:"portalID"`
	UpdatedAt  *time.Time          `json:"updatedAt"`
	Visibility *CategoryVisibility `json:"visibilityConfig"`
}

type CategoryVisibility struct {
	Condition        string      `json:"condition"`
	ComplexIDs       []int       `json:"complexesIds"`
	OIVs             []int       `json:"oivsIds"`
	OrgIDs           []uuid.UUID `json:"organizationsIds"`
	ProductIDs       []uuid.UUID `json:"productsIds"`
	SubdivisionNames []string    `json:"subdivisionsNames"`
	PositionNames    []string    `json:"positionsNames"`
	EmployeeIDs      []uuid.UUID `json:"personsIds"`
	RoleNames        []string    `json:"rolesUshrNames"`
}

type SearchCategoryRequest struct {
	Page  int    `json:"page"`
	Limit *int   `json:"limit"`
	Query string `json:"query"`
}

type SearchCategoryResult struct {
	TotalCount int                        `json:"total"`
	Categories []SearchCategoryResultItem `json:"data"`
}

type SearchCategoryResultItem struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
