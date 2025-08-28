package news

import (
	commentsv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/newsfacade/comment/v1"
	entityNews "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/news"
)

func NewCommentsMapper(sharedMapper SharedMapper) *commentsMapper {
	return &commentsMapper{
		sharedMapper: sharedMapper,
	}
}

type commentsMapper struct {
	sharedMapper SharedMapper
}

func (cm *commentsMapper) CommentsToEntity(commentsPb []*commentsv1.Comment) []*entityNews.NewsComment {
	if len(commentsPb) == 0 {
		return nil
	}

	comments := make([]*entityNews.NewsComment, 0, len(commentsPb))
	for _, commentPb := range commentsPb {
		comment := cm.CommentToEntity(commentPb)
		if comment == nil {
			continue
		}

		comments = append(comments, comment)
	}
	return comments
}

func (cm *commentsMapper) CommentToEntity(comment *commentsv1.Comment) *entityNews.NewsComment {
	if comment == nil {
		return nil
	}

	id := cm.sharedMapper.StringToUUIDPtr(comment.GetId(), false)
	return &entityNews.NewsComment{
		ID:      *id,
		Message: comment.GetText(),
		Author: entityNews.Author{
			ID:         cm.sharedMapper.StringToUUIDPtr(comment.GetAuthor().GetId(), true),
			LastName:   comment.GetAuthor().GetLastName(),
			FirstName:  comment.GetAuthor().GetFirstName(),
			MiddleName: cm.sharedMapper.StringValueToStringPtr(comment.GetAuthor().GetMiddleName()),
			ImageID:    cm.sharedMapper.StringValueToUUIDPtr(comment.GetAuthor().GetImageId(), true),
		},
		CreatedAt: cm.sharedMapper.TimestampToTime(comment.GetCreateTime()),
		UpdatedAt: cm.sharedMapper.TimestampToTime(comment.GetUpdateTime()),
		DeletedAt: cm.sharedMapper.TimestampToTime(comment.GetDeleteTime()),
	}
}
