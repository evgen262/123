package news

import (
	"fmt"
	"time"

	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"github.com/google/uuid"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/dto"
	entityNews "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/news"
)

//go:generate ditgen -source=./news.go -zero=true  -all-fields=true

const (
	ErrNewsSlugRequired     diterrors.StringError = "news slug required"
	ErrNewsTitleRequired    diterrors.StringError = "news title required"
	ErrCategoryRequired     diterrors.StringError = "category required"
	ErrNewsStatus           diterrors.StringError = "incorrect news status"
	ErrPublishTimeBeforeNow diterrors.StringError = "publish time must be in the future"
)

type NewNews struct {
	Title           string
	Slug            string
	ImageID         *uuid.UUID
	CategoryID      uuid.UUID
	OrganizationID  *uuid.UUID
	ProductID       *uuid.UUID
	Participants    []*uuid.UUID
	Status          entityNews.NewsStatus
	Body            []byte
	Author          Author
	OnMain          bool
	Pinned          bool
	CanDisplayViews bool
	CanReacts       bool
	CanCommented    bool
	PublicationAt   *time.Time
	Visibility      *entityNews.NewsVisibility
}

func (nn *NewNews) GetValidStatus() entityNews.NewsStatus {
	switch nn.Status {
	case entityNews.NewsStatusDraft:
		return entityNews.NewsStatusDraft
	case entityNews.NewsStatusWaitingPublish:
		return entityNews.NewsStatusWaitingPublish
	case entityNews.NewsStatusPublished:
		return entityNews.NewsStatusPublished
	case entityNews.NewsStatusUnpublished:
		return entityNews.NewsStatusUnpublished
	default:
		return entityNews.NewsStatusInvalid
	}
}

func (nn *NewNews) Validate() error {
	if nn == nil {
		return diterrors.NewValidationError(diterrors.ErrInputEmpty)
	}

	if nn.Slug == "" {
		return diterrors.NewValidationError(ErrNewsSlugRequired)
	}

	if nn.Title == "" {
		return diterrors.NewValidationError(ErrNewsTitleRequired)
	}

	if nn.CategoryID == uuid.Nil {
		return diterrors.NewValidationError(ErrCategoryRequired,
			diterrors.ErrValidationFields{
				Field:   "category_id",
				Message: fmt.Sprintf("invalid category_id <%s>", nn.CategoryID),
			},
		)
	}

	if nn.GetStatus() == entityNews.NewsStatusInvalid {
		return diterrors.NewValidationError(ErrNewsStatus)
	}

	// Делаем проверку на дату публикации.
	//  За текущее время берется 1 минута назад для устранения проблем с лагом времени.
	nowTime := time.Now().Add(-1 * time.Minute)
	if nn.PublicationAt != nil && nn.PublicationAt.Before(nowTime) {
		return diterrors.NewValidationError(ErrPublishTimeBeforeNow, diterrors.ErrValidationFields{
			Field:   "publication_at",
			Message: fmt.Sprintf("incorrect publication_at time: %s", nn.PublicationAt),
		})
	}

	return nil
}

type Author struct {
	ID         uuid.UUID
	LastName   string
	FirstName  string
	MiddleName *string
	ImageID    *uuid.UUID
}

type UpdateNews struct {
	Title           *string
	Slug            *string
	ImageID         *uuid.UUID
	CategoryID      *uuid.UUID
	OrganizationID  *uuid.UUID
	ProductID       *uuid.UUID
	Participants    []*uuid.UUID
	Status          entityNews.NewsStatus
	Body            []byte
	OnMain          *bool
	Pinned          *bool
	CanDisplayViews *bool
	CanReacts       *bool
	CanCommented    *bool
	PublicationAt   *time.Time
	Visibility      *entityNews.NewsVisibility
	UpdatedAt       *time.Time
}

type UpdateFlags struct {
	OnMain    *bool
	Pinned    *bool
	UpdatedAt *time.Time
}

func (un *UpdateNews) GetValidStatus() entityNews.NewsStatus {
	switch un.Status {
	case entityNews.NewsStatusDraft:
		return entityNews.NewsStatusDraft
	case entityNews.NewsStatusWaitingPublish:
		return entityNews.NewsStatusWaitingPublish
	case entityNews.NewsStatusPublished:
		return entityNews.NewsStatusPublished
	case entityNews.NewsStatusUnpublished:
		return entityNews.NewsStatusUnpublished
	default:
		return entityNews.NewsStatusInvalid
	}
}

func (un *UpdateNews) Validate() error {
	if un == nil {
		return diterrors.NewValidationError(diterrors.ErrInputEmpty)
	}

	if un.CategoryID != nil && *un.CategoryID == uuid.Nil {
		return diterrors.NewValidationError(ErrCategoryRequired,
			diterrors.ErrValidationFields{
				Field:   "category_id",
				Message: fmt.Sprintf("invalid category_id <%s>", un.CategoryID),
			},
		)
	}

	if un.Title == nil || *un.Title == "" {
		return diterrors.NewValidationError(ErrNewsTitleRequired)
	}

	if un.Slug == nil || *un.Slug == "" {
		return diterrors.NewValidationError(ErrNewsSlugRequired)
	}

	// Делаем проверку на дату публикации.
	//  За текущее время берется 1 минута назад для устранения проблем с лагом времени.
	nowTime := time.Now().Add(-1 * time.Minute)
	if un.PublicationAt != nil && un.PublicationAt.Before(nowTime) {
		return diterrors.NewValidationError(ErrPublishTimeBeforeNow, diterrors.ErrValidationFields{
			Field:   "publication_at",
			Message: fmt.Sprintf("incorrect publication_at time: %s", un.PublicationAt),
		})
	}

	return nil
}

type SearchNews struct {
	Query      string
	Filter     *SearchNewsFilter
	Pagination SearchNewsPagination
	Order      SearchNewsOrder
	Visitor    *entityNews.Visitor
}

type SearchNewsFilter struct {
	Status                   entityNews.NewsStatus
	ProviderOrganizationsIds []*uuid.UUID
	ProviderProductsNames    []string
	CategoriesNames          []string
	AuthorsNames             []string
	OnMainPage               bool
	IsPinnedOnMainPage       bool
}

type SearchNewsScroll struct {
	LastID    *uuid.UUID
	CreatedAt *time.Time
	Limit     int
}

type SearchNewsPagination struct {
	Page  int
	Limit int
}

type SearchNewsOrderBy int

const (
	SearchNewsOrderByTitle SearchNewsOrderBy = iota
	SearchNewsOrderByCreatedAt
)

type SearchNewsOrder struct {
	By        SearchNewsOrderBy
	Direction dto.OrderDirection
}

func (sno SearchNewsOrder) GetValidDirection() dto.OrderDirection {
	if sno.Direction == dto.OrderDirectionAsc {
		return dto.OrderDirectionAsc
	}
	return dto.OrderDirectionDesc
}

type SearchNewsResult struct {
	News  []*entityNews.NewsFull
	Total int
}
