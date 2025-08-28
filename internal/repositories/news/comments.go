package news

import (
	"context"
	"fmt"

	newsv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/newsfacade/comment/v1"
	sharedv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/newsfacade/shared/v1"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/dto"
	dtoNews "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/dto/news"
	entityNews "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/news"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"github.com/google/uuid"
)

type commentsRepository struct {
	newsApi    newsv1.CommentAPIClient
	mapper     CommentMapper
	newsMapper NewsMapper
}

func NewCommentsRepository(newsFacadeApi newsv1.CommentAPIClient, mapper CommentMapper, newsMapper NewsMapper) *commentsRepository {
	return &commentsRepository{
		newsApi:    newsFacadeApi,
		mapper:     mapper,
		newsMapper: newsMapper,
	}
}

func (c *commentsRepository) Create(ctx context.Context, in dtoNews.NewComment) (uuid.UUID, int, error) {
	req := &newsv1.CreateRequest{
		NewsId: in.NewsID.String(),
		Text:   in.Text,
	}

	if in.GetAuthorID() != nil {
		req.AuthorId = in.GetAuthorID().String()
	}

	resp, err := c.newsApi.Create(ctx, req)
	if err != nil {
		return uuid.Nil, 0, fmt.Errorf("newsFacadeApi.Create: %w", diterrors.GrpcErrorToError(err))
	}
	// Разбор ответа
	commentID := uuid.Nil
	if idStr := resp.GetId(); idStr != "" {
		if parsedID, parseErr := uuid.Parse(idStr); parseErr == nil {
			commentID = parsedID
		}
	}

	// Возвращаем uuid.Nil если не удалось распарсить, и реально вернувшееся количество комментариев
	return commentID, int(resp.GetCount()), nil
}

func (c *commentsRepository) List(ctx context.Context, params *dtoNews.FilterComments) ([]*entityNews.NewsComment, int, error) {

	// TODO: В задаче на рефакторинг убрать в маппер
	sort := &newsv1.ListByNewsRequest_Sort{
		By: newsv1.ListByNewsRequest_Sort_SORT_FIELDS_CREATED_DATE,
	}

	if params.Order == dto.OrderDirectionDesc {
		sort.DirectionAsk = false
	}

	pagination := &sharedv1.ScrollPaginationRequest{
		Limit: int32(params.Limit),
	}

	if params.AfterID != nil {
		pagination.LastId = params.AfterID.String()
	}

	req := &newsv1.ListByNewsRequest{
		NewsId:     params.NewsID.String(),
		Sort:       sort,
		Pagination: pagination,
	}

	resp, err := c.newsApi.ListByNews(ctx, req)
	if err != nil {
		return nil, 0, fmt.Errorf("newsFacadeApi.Create: %w", diterrors.GrpcErrorToError(err))
	}

	return c.mapper.CommentsToEntity(resp.GetComments()), int(resp.GetPagination().GetTotal()), nil
}
