package news

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"

	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"github.com/google/uuid"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/news"
)

//go:generate ditgen -source=./category.go -zero=true

const CATEGORY_NAME_MAX_LENGTH = 255

type NewCategory struct {
	Name       string
	AuthorID   uuid.UUID
	PortalID   int
	Visibility *CategoryVisibility
}

type CategoryVisibility struct {
	Condition        string
	ComplexIDs       []int
	OIVs             []int
	OrgIDs           []uuid.UUID
	ProductIDs       []uuid.UUID
	SubdivisionNames []string
	PositionNames    []string
	EmployeeIDs      []uuid.UUID
	RoleNames        []string
}

type UpdateCategory struct {
	ID         uuid.UUID
	Name       string
	UpdatedAt  *time.Time
	AuthorID   uuid.UUID
	PortalID   int
	Visibility *CategoryVisibility
}

func (c *UpdateCategory) Validate() error {
	if c == nil {
		return nil
	}

	c.Name = strings.TrimSpace(c.Name)

	if c.Name == "" {
		return diterrors.NewValidationError(errors.New("UpdateCategory.Validate: category name is required"))
	}

	if utf8.RuneCount([]byte(c.Name)) > CATEGORY_NAME_MAX_LENGTH {
		return diterrors.NewValidationError(errors.New("UpdateCategory.Validate: category name too long"))
	}
	return nil
}

type SearchCategory struct {
	Query      string
	PortalID   int
	Pagination CategoryPagination
}

type FilterCategoryBy int

const (
	FilterCategoryByUnknown FilterCategoryBy = iota
	FilterCategoryByIDs
	FilterCategoryByName
	FilterCategoryByPortalIDs
)

type FilterCategory struct {
	By         FilterCategoryBy
	Name       *string
	PortalIDs  []int
	IDs        []uuid.UUID
	Pagination *CategoryPagination
	Visitor    *news.Visitor
}

type CategoryPagination struct {
	Page  int
	Limit *int
	Total int
}

type CategoryNameAndPortalID struct {
	Name     string
	PortalID int
}

func (p *CategoryPagination) Validate() error {
	if p == nil {
		return diterrors.NewValidationError(errors.New("CategoryPagination.Validate: category page is required"))
	}

	if p.Limit != nil && *p.Limit < 0 {
		return diterrors.NewValidationError(errors.New("CategoryPagination.Validate: category limit is required"))
	}

	if p.Page < 1 {
		return diterrors.NewValidationError(errors.New("CategoryPagination.Validate: category page is invalid"))
	}

	return nil
}
