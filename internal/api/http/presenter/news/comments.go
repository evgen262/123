package news

import (
	"strings"

	viewNews "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/news"
	dtoNews "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/dto/news"
	entityNews "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/news"
	"github.com/google/uuid"
)

type commentsPresenter struct {
}

func NewCommentsPresenter() *commentsPresenter {
	return &commentsPresenter{}
}

func (cp *commentsPresenter) NewCommentToDTO(newsID uuid.UUID, v *viewNews.NewNewsComment) dtoNews.NewComment {
	comment := dtoNews.NewComment{
		NewsID: newsID,
	}
	if v != nil {
		comment.Text = strings.TrimSpace(v.Text)
	}
	return comment
}

func (cp *commentsPresenter) CommentsToView(list []*entityNews.NewsComment) []*viewNews.NewsComment {
	comments := make([]*viewNews.NewsComment, 0, len(list))

	if len(list) == 0 {
		return comments
	}

	for _, l := range list {
		comment := cp.CommentToView(l)
		if comment == nil {
			continue
		}

		comments = append(comments, comment)
	}

	return comments
}

func (cp *commentsPresenter) CommentToView(comment *entityNews.NewsComment) *viewNews.NewsComment {
	if comment == nil {
		return nil
	}

	fio := make([]string, 0, 3)
	fio = append(fio,
		comment.Author.LastName,
		comment.Author.FirstName,
	)
	if comment.Author.GetMiddleNamePtr() != nil {
		fio = append(fio, comment.Author.GetMiddleName())
	}
	name := strings.Join(fio, " ")

	var imgID string
	if comment.Author.GetImageIDPtr() != nil && comment.Author.GetImageID() != uuid.Nil {
		imgID = comment.Author.GetImageID().String()
	}

	view := &viewNews.NewsComment{
		ID:         comment.ID,
		CreateAt:   comment.GetCreatedAtPtr(),
		Message:    comment.Message,
		IsUserMade: false,
		Author: viewNews.CommentAuthor{
			ID:      comment.Author.GetID(),
			Name:    name,
			ImageID: imgID,
			// TODO: Уточнить у СА требование к этому флагу
			IsActive: true,
		},
	}

	if !comment.GetDeletedAt().IsZero() {
		view.IsDeleted = true
		// TODO: Уточнить у СА требование к тексту комментария
		view.Message = ""
	}

	return view
}
