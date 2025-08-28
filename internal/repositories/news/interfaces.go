package news

import (
	commentsv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/newsfacade/comment/v1"
	newsv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/newsfacade/news/v1"
	dtoNews "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/dto/news"
	"github.com/google/uuid"

	entityNews "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/news"
)

//go:generate mockgen -source=interfaces.go -destination=./news_mock.go -package=news

type NewsMapper interface {
	NewsFullToEntity(newsPb []*newsv1.News) []*entityNews.NewsFull
	OnceNewsFullToEntity(newsPb *newsv1.News) *entityNews.NewsFull
	NewsToEntity(newsPb []*newsv1.News) []*entityNews.News
	OnceNewsToEntity(newsPb *newsv1.News) *entityNews.News
	ParticipantsIdsToPb(participants []*uuid.UUID) []string
	StatusToPb(status entityNews.NewsStatus) newsv1.NewsStatus
	StatusToEntity(statusPb newsv1.NewsStatus) entityNews.NewsStatus
	VisibilityToPb(vis *entityNews.NewsVisibility) *newsv1.Visibility
	CommentsToEntity(commentsPb []*newsv1.Comment) []*entityNews.NewsComment
	CommentToEntity(comment *newsv1.Comment) *entityNews.NewsComment
	NewCommentToPb(comment *dtoNews.NewComment) *newsv1.Comment
}

type CommentMapper interface {
	CommentsToEntity(commentsPb []*commentsv1.Comment) []*entityNews.NewsComment
	CommentToEntity(comment *commentsv1.Comment) *entityNews.NewsComment
}
