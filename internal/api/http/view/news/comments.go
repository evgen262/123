package news

import (
	"time"

	"github.com/google/uuid"
)

//go:generate ditgen -source=comments.go -zero=true

// Тело запроса: { "text": "..." }
type NewNewsComment struct {
	Text string `json:"text"`
}

type NewsComment struct {
	ID         uuid.UUID     `json:"id"`
	CreateAt   *time.Time    `json:"date"`
	Message    string        `json:"text"`
	IsUserMade bool          `json:"isUserMade"`
	IsDeleted  bool          `json:"isDeleted"`
	Author     CommentAuthor `json:"author"`
}

type CommentAuthor struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	ImageID  string    `json:"imageId"`
	IsActive bool      `json:"isActive"`
}
