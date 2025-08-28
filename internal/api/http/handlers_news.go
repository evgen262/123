package http

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view"
	viewNews "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/news"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/dto"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/dto/news"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type newsHandlers struct {
	categoryInteractor NewsCategoryInteractor
	commentsInteractor NewsCommentsInteractor
	commentsPresenter  NewsCommentsPresenter
	newsInteractor     NewsInteractor
	newsPresenter      NewsAdminPresenter
	logger             ditzap.Logger
}

func NewNewsHandlers(
	categoryInteractor NewsCategoryInteractor,
	commentsInteractor NewsCommentsInteractor,
	commentsPresenter NewsCommentsPresenter,
	newsInteractor NewsInteractor,
	newsPresenter NewsAdminPresenter,
	logger ditzap.Logger,
) *newsHandlers {
	return &newsHandlers{
		categoryInteractor: categoryInteractor,
		commentsInteractor: commentsInteractor,
		commentsPresenter:  commentsPresenter,
		newsInteractor:     newsInteractor,
		newsPresenter:      newsPresenter,
		logger:             logger,
	}
}

func (n *newsHandlers) getNews(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c.Request.Context(), RequestTimeOut)
	defer cancelCtx()

	slug := c.Param(SLUG_PARAM_KEY)
	if slug == "" {
		c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		return
	}

	news, err := n.newsInteractor.Get(ctx, slug)
	if err != nil {
		type errResponse struct {
			code     int
			response *view.Response
		}
		resp := &errResponse{}
		switch {
		case errors.Is(err, diterrors.ErrNotFound):
			resp.code = http.StatusNotFound
			resp.response = view.NewErrorResponse(viewNews.ErrNewsNotFound)
		case errors.Is(err, diterrors.ErrFailedPrecondition):
			fallthrough
		case errors.As(err, new(diterrors.ValidationError)):
			resp.code = http.StatusBadRequest
			resp.response = view.NewErrorResponse(err)
		default:
			resp.code = http.StatusInternalServerError
			resp.response = view.NewErrorResponse(view.ErrMessageInternalError)
		}
		c.JSON(resp.code, resp.response)
		return
	}

	c.JSON(http.StatusOK, view.NewSuccessResponse(n.newsPresenter.FullNewsToView(news)))
}

func (n *newsHandlers) searchNews(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c.Request.Context(), RequestTimeOut)
	defer cancelCtx()

	search := new(viewNews.SearchNewsRequest)
	if err := c.BindJSON(search); err != nil {
		c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		return
	}

	req := n.newsPresenter.SearchNewsToDTO(search)

	news, err := n.newsInteractor.Search(ctx, req)
	if err != nil {
		type errResponse struct {
			code     int
			response *view.Response
		}
		resp := &errResponse{}
		switch {
		case errors.Is(err, diterrors.ErrFailedPrecondition):
			fallthrough
		case errors.As(err, new(diterrors.ValidationError)):
			resp.code = http.StatusBadRequest
			resp.response = view.NewErrorResponse(err)
		default:
			resp.code = http.StatusInternalServerError
			resp.response = view.NewErrorResponse(view.ErrMessageInternalError)
		}
		c.JSON(resp.code, resp.response)
		return
	}

	c.JSON(http.StatusOK, view.NewSuccessResponse(&viewNews.SearchNewsResponse{
		Total: news.Total,
		Data:  n.newsPresenter.FullNewsToSearchItems(news.News),
	}))
}

func (n *newsHandlers) createComment(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), RequestTimeOut)
	defer cancel()

	newsID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.Header(StatusCodeHeader, "ncc_01")
		c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		return
	}

	newComment := new(viewNews.NewNewsComment)
	if err := c.BindJSON(newComment); err != nil {
		c.Header(StatusCodeHeader, "ncc_02")
		c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		return
	}

	dto := n.newsPresenter.NewCommentToDTO(newsID, newComment)

	_, count, err := n.commentsInteractor.Create(ctx, dto)
	if err != nil {
		type errResponse struct {
			code     int
			response *view.Response
		}
		resp := &errResponse{}

		var (
			alreadyExists diterrors.AlreadyExistsError
			validationErr diterrors.ValidationError
		)

		switch {
		case errors.Is(err, diterrors.ErrFailedPrecondition):
			c.Header(StatusCodeHeader, "ncc_03")
			resp.code = http.StatusPreconditionFailed
			resp.response = view.NewErrorResponse(err)
		case errors.As(err, &validationErr):
			c.Header(StatusCodeHeader, "ncc_04")
			resp.code = http.StatusBadRequest
			resp.response = view.NewErrorResponse(err)
		case errors.As(err, &alreadyExists):
			c.Header(StatusCodeHeader, "ncc_05")
			resp.code = http.StatusConflict
			resp.response = view.NewErrorResponse(alreadyExists)
		default:
			c.Header(StatusCodeHeader, "ncc_06")
			n.logger.Error("create comment failed", zap.Error(err))
			resp.code = http.StatusInternalServerError
			resp.response = view.NewErrorResponse(view.ErrMessageInternalError)
		}
		c.JSON(resp.code, resp.response)
		return
	}

	// 5) Ответ целевого вида
	type commentsPayload struct {
		Count      int  `json:"count"`
		IsUserMade bool `json:"isUserMade"`
	}
	type dataPayload struct {
		Comments commentsPayload `json:"comments"`
	}
	resp := &dataPayload{
		Comments: commentsPayload{
			Count:      count,
			IsUserMade: true,
		},
	}

	c.JSON(http.StatusOK, view.NewSuccessResponse(resp))
}

func (n *newsHandlers) listComments(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), RequestTimeOut)
	defer cancel()

	newsID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.Header(StatusCodeHeader, "nlc_01")
		c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		return
	}

	// TODO: переделать и убрать в маппер
	// POKAZ: получение query-параметров
	sortField, ok := c.GetQuery("sortBy")
	if !ok || sortField != "date" {
		c.Header(StatusCodeHeader, "nlc_02")
		c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		return
	}
	// на текущий момент сортировка только по дате создания комментария
	sort := dto.SortFieldDateCreate

	// orderType
	orderType, ok := c.GetQuery("orderType")
	if !ok {
		c.Header(StatusCodeHeader, "nlc_03")
		c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		return
	}
	order := dto.OrderDirectionAsc
	if orderType == "DESC" {
		order = dto.OrderDirectionDesc
	}

	// lastCommentId
	lastID, ok := c.GetQuery("lastCommentId")
	var AfterID *uuid.UUID
	if ok && lastID != "" {
		aID, err := uuid.Parse(lastID)
		if err != nil {
			c.Header(StatusCodeHeader, "nlc_04")
			c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
			return
		}
		AfterID = &aID
	}

	// limit
	limit, ok := c.GetQuery("limit")
	if !ok {
		c.Header(StatusCodeHeader, "nlc_05")
		c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		return
	}
	l, err := strconv.Atoi(limit)
	if err != nil || l < 0 {
		c.Header(StatusCodeHeader, "nlc_06")
		c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		return
	}
	// POKAZ: end

	params := &news.FilterComments{
		NewsID:    newsID,
		SortField: sort,
		Order:     order,
		AfterID:   AfterID,
		Limit:     l,
	}

	// пока total нет в контракте
	list, _, err := n.commentsInteractor.List(ctx, params)
	if err != nil {
		c.Header(StatusCodeHeader, "nlc_07")
		c.JSON(http.StatusInternalServerError, view.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, view.NewSuccessResponse(n.commentsPresenter.CommentsToView(list)))
	return
}
