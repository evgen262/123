package news

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"

	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"github.com/google/uuid"
)

const CATEGORY_NAME_MAX_LENGTH = 255

type Category struct {
	ID         uuid.UUID
	Name       string
	UpdatedAt  *time.Time
	Visibility CategoryVisibility
}

func (c *Category) Validate() error {
	if c == nil {
		return nil
	}

	c.Name = strings.TrimSpace(c.Name)

	if c.Name == "" {
		return diterrors.NewValidationError(errors.New("Category.Validate: category name is required"))
	}

	if utf8.RuneCount([]byte(c.Name)) > CATEGORY_NAME_MAX_LENGTH {
		return diterrors.NewValidationError(errors.New("Category.Validate: category name too long"))
	}
	return nil
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

type CategoryPagination struct {
	Page  int
	Limit *int
	Total int
}

type CategoriesWithPagination struct {
	Categories []*Category
	Pagination CategoryPagination
}
