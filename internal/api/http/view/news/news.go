package news

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

//go:generate ditgen -source=news.go -zero=true

type NewNews struct {
	Title                  string          `json:"title"`
	Slug                   string          `json:"slug"`
	ImageID                *uuid.UUID      `json:"titleImageId"`
	CategoryID             uuid.UUID       `json:"categoryId"`
	ProviderOrganizationID *uuid.UUID      `json:"providerOrganizationId"`
	ProviderProductID      *uuid.UUID      `json:"providerProductId"`
	Properties             NewsProperties  `json:"flags"`
	Content                json.RawMessage `json:"content"`
	ParticipantsIDs        []*uuid.UUID    `json:"participantsIds"`
	PublishDate            *time.Time      `json:"publishDate"`
}

type UpdateNews struct {
	Slug                   string          `json:"slug"`
	Title                  string          `json:"title"`
	ImageID                *string         `json:"titleImageId"`
	CategoryID             *string         `json:"categoryId"`
	ProviderOrganizationId *string         `json:"providerOrganizationId"`
	Properties             NewsProperties  `json:"flags"`
	ProviderProductId      *string         `json:"providerProductId"`
	Visibility             NewsVisibility  `json:"visibilityConfig"`
	Content                json.RawMessage `json:"content"`
	ParticipantsIDs        []*uuid.UUID    `json:"participantsIds"`
	PublishDate            *time.Time      `json:"publishDate"`
	UpdatedAt              *time.Time      `json:"updatedAt"`
}

type News struct {
	Id                   uuid.UUID           `json:"id"`
	Slug                 string              `json:"slug"`
	Title                string              `json:"title"`
	ImageID              *uuid.UUID          `json:"titleImageId"`
	Participants         []*NewsParticipants `json:"participants"`
	Category             *NewsCategory       `json:"category"`
	Author               Author              `json:"author"`
	Status               NewsStatus          `json:"status"`
	Properties           NewsProperties      `json:"flags"`
	ProviderOrganization *NewsOrganization   `json:"providerOrganization"`
	ProviderProduct      *NewsProduct        `json:"providerProduct"`
	Content              json.RawMessage     `json:"content"`
	UpdatedAt            *time.Time          `json:"updatedAt"`
	CreateDate           *time.Time          `json:"createDate"`
	PublishDate          *time.Time          `json:"publishDate"`
}

type NewsProperties struct {
	ViewsEnabled    bool `json:"isViewsEnabled"`
	LikesEnabled    bool `json:"isLikesEnabled"`
	CommentsEnabled bool `json:"isCommentsEnabled"`
	OnMainPage      bool `json:"onMainPage"`
	MainPagePinned  bool `json:"isMainPagePinned"`
}

type NewsStatus string

const (
	NewsStatusDraft          NewsStatus = "DRAFT"
	NewsStatusPublished      NewsStatus = "PUBLISHED"
	NewsStatusUnpublished    NewsStatus = "UNPUBLISHED"
	NewsStatusWaitingPublish NewsStatus = "WAITING_PUBLISH"
)

type NewsVisibility struct {
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

type SearchNewsRequest struct {
	OrderBy   string                   `json:"orderBy"`
	OrderType string                   `json:"orderType"`
	Page      int                      `json:"page"`
	Limit     int                      `json:"limit"`
	Query     string                   `json:"query"`
	Filters   SearchNewsRequestFilters `json:"filters"`
}

type SearchNewsRequestFilters struct {
	Status                   NewsStatus   `json:"status"`
	ProviderOrganizationsIds []*uuid.UUID `json:"providerOrganizationsIds"`
	ProviderProductsNames    []string     `json:"providerProductsNames"`
	CategoriesNames          []string     `json:"categoriesNames"`
	AuthorsNames             []string     `json:"authorsNames"`
	OnMainPage               bool         `json:"onMainPage"`
	IsPinnedOnMainPage       bool         `json:"isMainPagePinned"`
}

type SearchNewsResponse struct {
	Total int                       `json:"total"`
	Data  []*SearchNewsResponseItem `json:"data"`
}

type SearchNewsResponseItem struct {
	Id                   uuid.UUID               `json:"id"`
	Slug                 string                  `json:"slug"`
	Title                string                  `json:"title"`
	ImageID              *uuid.UUID              `json:"titleImageId"`
	Category             *NewsCategory           `json:"category"`
	ProviderOrganization *NewsOrganization       `json:"providerOrganization"`
	ProviderProduct      *NewsProduct            `json:"providerProduct"`
	Author               Author                  `json:"author"`
	Flags                SearchNewsResponseFlags `json:"flags"`
	Status               NewsStatus              `json:"status"`
	CreateAt             *time.Time              `json:"createDate"`
	PublishedAt          *time.Time              `json:"publishDate"`
}

type NewsCategory struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type NewsOrganization struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type NewsProduct struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type Author struct {
	ID      *uuid.UUID `json:"id"`
	Name    string     `json:"name"`
	ImageID *uuid.UUID `json:"imageId"`
}

type NewsParticipants struct {
	Id      *uuid.UUID `json:"id"`
	Name    string     `json:"name"`
	ImageID *uuid.UUID `json:"imageId"`
}

type SearchNewsResponseFlags struct {
	OnMainPage         bool `json:"onMainPage"`
	IsPinnedOnMainPage bool `json:"isMainPagePinned"`
}

type UpdateNewsFlags struct {
	OnMainPage         *bool      `json:"onMainPage"`
	IsPinnedOnMainPage *bool      `json:"isMainPagePinned"`
	UpdatedAt          *time.Time `json:"updatedAt"`
}
