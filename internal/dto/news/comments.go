package news

import (
	"fmt"
	"strings"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/dto"
	entityNews "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/news"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"github.com/google/uuid"
)

//go:generate ditgen -source=comments.go
const (
	ErrNewsIDRequired      diterrors.StringError = "news id required"
	ErrCommentTextRequired diterrors.StringError = "comment text required"
	ErrAuthorRequired      diterrors.StringError = "author required"
)

type NewComment struct {
	NewsID   uuid.UUID  // ID новости (из path {id})
	Text     string     // Текст комментария (из body)
	AuthorID *uuid.UUID // Автор из контекста аутентификации
}

func (nc *NewComment) Validate() error {
	if nc == nil {
		return diterrors.NewValidationError(diterrors.ErrInputEmpty)
	}

	// Проверяем news_id
	if nc.NewsID == uuid.Nil {
		return diterrors.NewValidationError(ErrNewsIDRequired,
			diterrors.ErrValidationFields{
				Field:   "news_id",
				Message: fmt.Sprintf("invalid news_id"),
			},
		)
	}

	// Проверяем текст
	if strings.TrimSpace(nc.Text) == "" {
		return diterrors.NewValidationError(ErrCommentTextRequired, diterrors.ErrValidationFields{
			Field:   "text",
			Message: "text is empty",
		})
	}

	// Проверяем author_id (из контекста аутентификации)
	if nc.AuthorID == nil || *nc.AuthorID == uuid.Nil {
		return diterrors.NewValidationError(ErrAuthorRequired,
			diterrors.ErrValidationFields{
				Field:   "author_id",
				Message: fmt.Sprintf("invalid author_id"),
			},
		)
	}

	return nil
}

type FilterComments struct {
	NewsID    uuid.UUID
	SortField dto.SortField
	Order     dto.OrderDirection
	AfterID   *uuid.UUID
	Limit     int
	Visitor   *entityNews.Visitor
}
