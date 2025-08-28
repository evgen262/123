package news

import (
	"time"

	"github.com/google/uuid"
)

//go:generate ditgen -source=./news.go -zero=true

type News struct {
	ID              uuid.UUID
	Title           string
	Slug            string
	ImageID         *uuid.UUID
	Portals         []int
	CategoryID      uuid.UUID
	OrganizationID  *uuid.UUID
	ProductID       *uuid.UUID
	Participants    []*uuid.UUID
	Status          NewsStatus
	Body            []byte
	Author          Author
	OnMain          bool
	Pinned          bool
	CanDisplayViews bool
	CanReacts       bool
	CanCommented    bool
	PublicationAt   *time.Time
	Visibility      *NewsVisibility
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
}

func (n *News) GetStatus() NewsStatus {
	switch n.Status {
	case NewsStatusDraft:
		return NewsStatusDraft
	case NewsStatusWaitingPublish:
		return NewsStatusWaitingPublish
	case NewsStatusPublished:
		return NewsStatusPublished
	case NewsStatusUnpublished:
		return NewsStatusUnpublished
	default:
		return NewsStatusInvalid
	}
}

type NewsStatus int

const (
	NewsStatusInvalid NewsStatus = iota
	NewsStatusDraft
	NewsStatusWaitingPublish
	NewsStatusPublished
	NewsStatusUnpublished
)

type Author struct {
	ID         *uuid.UUID
	LastName   string
	FirstName  string
	MiddleName *string
	ImageID    *uuid.UUID
}

type NewsFull struct {
	ID              uuid.UUID
	Title           string
	Slug            string
	ImageID         *uuid.UUID
	Category        *Category
	Organization    *NewsOrganization
	Product         *NewsProduct
	Participants    []*Participant
	Status          NewsStatus
	Body            []byte
	Author          Author
	OnMain          bool
	Pinned          bool
	CanDisplayViews bool
	Views           int
	CanReacts       bool
	Likes           int
	CanCommented    bool
	Comments        []*NewsComment
	PublicationAt   *time.Time
	Visibility      *NewsNamedVisibility
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
}

func (n *NewsFull) GetStatus() NewsStatus {
	switch n.Status {
	case NewsStatusDraft:
		return NewsStatusDraft
	case NewsStatusWaitingPublish:
		return NewsStatusWaitingPublish
	case NewsStatusPublished:
		return NewsStatusPublished
	case NewsStatusUnpublished:
		return NewsStatusUnpublished
	default:
		return NewsStatusInvalid
	}
}

type NewsPortal struct {
	ID   int
	Name string
}

type NewsOrganization struct {
	ID   uuid.UUID
	Name string
}

type NewsProduct struct {
	ID   uuid.UUID
	Name string
}

type Participant struct {
	ID         *uuid.UUID
	LastName   string
	FirstName  string
	MiddleName *string
	ImageID    *uuid.UUID
}

type NewsComment struct {
	ID        uuid.UUID
	Message   string
	Author    Author
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

type NewsVisibility struct {
	ComplexIDs []int
	PortalsIDs []int
}

type NewsNamedVisibility struct {
	Portals []*NewsPortal
}
