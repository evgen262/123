package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/news"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view"
	viewNews "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/news"
	dtoNews "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/dto/news"
	usecaseNews "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/usecase/news"
)

type newsAdminHandlers struct {
	categoryInteractor NewsCategoryInteractor
	newsInteractor     NewsAdminInteractor
	newsPresenter      NewsAdminPresenter
	logger             ditzap.Logger
}

func NewNewsAdminHandlers(
	categoryInteractor NewsCategoryInteractor,
	newsInteractor NewsAdminInteractor,
	newsPresenter NewsAdminPresenter,
	logger ditzap.Logger,
) *newsAdminHandlers {
	return &newsAdminHandlers{
		categoryInteractor: categoryInteractor,
		newsInteractor:     newsInteractor,
		newsPresenter:      newsPresenter,
		logger:             logger,
	}
}

func (h *newsAdminHandlers) createNews(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c.Request.Context(), RequestTimeOut)
	defer cancelCtx()

	newNews := new(viewNews.NewNews)
	if err := c.BindJSON(newNews); err != nil {
		c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		return
	}
	dto := h.newsPresenter.NewNewsToDTO(newNews)

	b, err := json.Marshal(newNews.Content)
	if err != nil {
		h.logger.Error("can't marshal news content", zap.Error(err))
		c.JSON(http.StatusInternalServerError, view.NewErrorResponse(view.ErrMessageInternalError))
	}
	dto.Body = b

	newsID, err := h.newsInteractor.Create(ctx, dto)
	if err != nil {
		type errResponse struct {
			code     int
			response *view.Response
		}
		resp := &errResponse{}
		var alreadyExists diterrors.AlreadyExistsError
		switch {
		case errors.Is(err, dtoNews.ErrNewsSlugRequired):
			resp.code = http.StatusBadRequest
			resp.response = view.NewErrorResponse(viewNews.ErrSlugRequired)
		case errors.Is(err, dtoNews.ErrNewsTitleRequired):
			resp.code = http.StatusBadRequest
			resp.response = view.NewErrorResponse(viewNews.ErrTitleRequired)
		case errors.Is(err, dtoNews.ErrCategoryRequired):
			resp.code = http.StatusBadRequest
			resp.response = view.NewErrorResponse(viewNews.ErrCategoryRequired)
		case errors.Is(err, dtoNews.ErrNewsStatus):
			resp.code = http.StatusBadRequest
			resp.response = view.NewErrorResponse(viewNews.ErrStatus)
		case errors.Is(err, dtoNews.ErrPublishTimeBeforeNow):
			resp.code = http.StatusBadRequest
			resp.response = view.NewErrorResponse(viewNews.ErrPublishTimeBeforeNow)
		case errors.As(err, &alreadyExists):
			resp.code = http.StatusConflict
			resp.response = view.NewErrorResponse(alreadyExists)
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

	type response struct {
		ID        uuid.UUID  `json:"id"`
		UpdatedAt *time.Time `json:"-"`
	}

	resp := &response{
		ID: newsID,
		// UpdatedAt Неоткуда взять. Не предусмотрено.
	}

	c.JSON(http.StatusOK, view.NewSuccessResponse(resp))
}

func (h *newsAdminHandlers) updateNews(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c.Request.Context(), RequestTimeOut)
	defer cancelCtx()

	newNews := new(viewNews.UpdateNews)

	id := c.Param(ID_PARAM_KEY)
	parsedID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		return
	}

	if err := c.BindJSON(newNews); err != nil {
		c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		return
	}
	dto := h.newsPresenter.UpdateNewsToDTO(newNews)

	updatedNews, err := h.newsInteractor.Update(ctx, parsedID, dto)
	if err != nil {
		type errResponse struct {
			code     int
			response *view.Response
		}
		resp := &errResponse{}
		var alreadyExists diterrors.AlreadyExistsError
		switch {
		case errors.Is(err, dtoNews.ErrNewsSlugRequired):
			resp.code = http.StatusBadRequest
			resp.response = view.NewErrorResponse(viewNews.ErrSlugRequired)
		case errors.Is(err, dtoNews.ErrNewsTitleRequired):
			resp.code = http.StatusBadRequest
			resp.response = view.NewErrorResponse(viewNews.ErrTitleRequired)
		case errors.Is(err, dtoNews.ErrCategoryRequired):
			resp.code = http.StatusBadRequest
			resp.response = view.NewErrorResponse(viewNews.ErrCategoryRequired)
		case errors.Is(err, dtoNews.ErrNewsStatus):
			resp.code = http.StatusBadRequest
			resp.response = view.NewErrorResponse(viewNews.ErrStatus)
		case errors.Is(err, dtoNews.ErrPublishTimeBeforeNow):
			resp.code = http.StatusBadRequest
			resp.response = view.NewErrorResponse(viewNews.ErrPublishTimeBeforeNow)
		case errors.As(err, &alreadyExists):
			resp.code = http.StatusConflict
			resp.response = view.NewErrorResponse(alreadyExists)
		case errors.Is(err, diterrors.ErrFailedPrecondition):
			fallthrough
		case errors.As(err, new(diterrors.ValidationError)):
			resp.code = http.StatusBadRequest
			resp.response = view.NewErrorResponse(err)
		case errors.Is(err, diterrors.ErrNotFound):
			resp.code = http.StatusNotFound
			resp.response = view.NewErrorResponse(viewNews.ErrNewsNotFound)
		default:
			resp.code = http.StatusInternalServerError
			resp.response = view.NewErrorResponse(view.ErrMessageInternalError)
		}
		c.JSON(resp.code, resp.response)
		return
	}

	type response struct {
		UpdatedAt *time.Time `json:"updatedAt"`
	}

	resp := &response{
		UpdatedAt: updatedNews.GetUpdatedAtPtr(),
	}

	c.JSON(http.StatusOK, view.NewSuccessResponse(resp))
}

func (h *newsAdminHandlers) setStatusNews(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c.Request.Context(), RequestTimeOut)
	defer cancelCtx()

	idParam := c.Param(ID_PARAM_KEY)
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		return
	}

	var req struct {
		Status viewNews.NewsStatus `json:"status"`
	}
	if err := c.BindJSON(&req); err != nil {
		h.logger.Error("failed to bind JSON for status update", zap.Error(err))
		c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		return
	}

	status := h.newsPresenter.StatusToEntity(req.Status)
	if status == news.NewsStatusInvalid {
		h.logger.Debug("invalid status for news", zap.String("news_id", idParam), zap.String("status", string(req.Status)))
		c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		return
	}

	updatedNews, err := h.newsInteractor.ChangeStatus(ctx, id, status)
	if err != nil {
		type errResponse struct {
			code     int
			response *view.Response
		}
		resp := &errResponse{}
		switch {
		case errors.Is(err, dtoNews.ErrNewsStatus):
			resp.code = http.StatusBadRequest
			resp.response = view.NewErrorResponse(viewNews.ErrStatus)
		case errors.Is(err, dtoNews.ErrPublishTimeBeforeNow):
			resp.code = http.StatusBadRequest
			resp.response = view.NewErrorResponse(viewNews.ErrPublishTimeBeforeNow)
		case errors.Is(err, diterrors.ErrFailedPrecondition):
			fallthrough
		case errors.As(err, new(diterrors.ValidationError)):
			resp.code = http.StatusBadRequest
			resp.response = view.NewErrorResponse(err)
		case errors.Is(err, diterrors.ErrNotFound):
			resp.code = http.StatusNotFound
			resp.response = view.NewErrorResponse(viewNews.ErrNewsNotFound)
		default:
			resp.code = http.StatusInternalServerError
			resp.response = view.NewErrorResponse(view.ErrMessageInternalError)
		}
		c.JSON(resp.code, resp.response)
		return
	}

	type response struct {
		UpdatedAt *time.Time `json:"updatedAt"`
	}

	resp := &response{
		UpdatedAt: updatedNews.GetUpdatedAtPtr(),
	}

	c.JSON(http.StatusOK, view.NewSuccessResponse(resp))
}

func (h *newsAdminHandlers) setFlagsNews(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c.Request.Context(), RequestTimeOut)
	defer cancelCtx()

	idParam := c.Param(ID_PARAM_KEY)
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		return
	}

	flags := new(viewNews.UpdateNewsFlags)
	if err := c.BindJSON(flags); err != nil {
		h.logger.Error("failed to bind JSON for flag update", zap.Error(err))
		c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		return
	}

	dto := h.newsPresenter.UpdateFlagsToDTO(flags)

	updatedNews, err := h.newsInteractor.UpdateFlags(ctx, id, dto)
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
		case errors.As(err, new(diterrors.ValidationError)):
			resp.code = http.StatusBadRequest
			resp.response = view.NewErrorResponse(err)
		default:
			h.logger.Error("failed to update news flags", zap.Error(err))
			resp.code = http.StatusInternalServerError
			resp.response = view.NewErrorResponse(view.ErrMessageInternalError)
		}
		c.JSON(resp.code, resp.response)
		return
	}

	type response struct {
		UpdatedAt *time.Time `json:"updatedAt"`
	}

	resp := &response{
		UpdatedAt: updatedNews.GetUpdatedAtPtr(),
	}

	c.JSON(http.StatusOK, view.NewSuccessResponse(resp))
}

func (h *newsAdminHandlers) getNews(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c.Request.Context(), RequestTimeOut)
	defer cancelCtx()

	id := c.Param(ID_PARAM_KEY)
	parsedID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		return
	}
	news, err := h.newsInteractor.Get(ctx, parsedID)
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

	c.JSON(http.StatusOK, view.NewSuccessResponse(h.newsPresenter.FullNewsToView(news)))
}

func (h *newsAdminHandlers) searchNews(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c.Request.Context(), RequestTimeOut)
	defer cancelCtx()

	search := new(viewNews.SearchNewsRequest)
	if err := c.BindJSON(search); err != nil {
		c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		return
	}

	req := h.newsPresenter.SearchNewsToDTO(search)

	news, err := h.newsInteractor.Search(ctx, req)
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
		Data:  h.newsPresenter.FullNewsToSearchItems(news.News),
	}))
}

func (h *newsAdminHandlers) deleteNews(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c.Request.Context(), RequestTimeOut)
	defer cancelCtx()

	id := c.Param(ID_PARAM_KEY)
	parsedID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		return
	}
	if err := h.newsInteractor.Delete(ctx, parsedID); err != nil {
		switch {
		case errors.Is(err, diterrors.ErrNotFound):
			c.JSON(http.StatusNotFound, view.NewErrorResponse(viewNews.ErrNewsNotFound))
		case errors.Is(err, diterrors.ErrFailedPrecondition):
			fallthrough
		case errors.As(err, new(diterrors.ValidationError)):
			c.JSON(http.StatusBadRequest, view.NewErrorResponse(err))
		}
		return
	}
	c.JSON(http.StatusOK, view.NewSuccessResponse(nil))
}

func (h *newsAdminHandlers) createCategory(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c.Request.Context(), RequestTimeOut)
	defer cancelCtx()

	nCategory := new(viewNews.NewCategory)
	if err := c.BindJSON(nCategory); err != nil {
		c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		return
	}

	req := h.newsPresenter.NewCategoryToDTO(nCategory)

	res, err := h.categoryInteractor.Create(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, usecaseNews.ErrCategoryAlreadyExists):
			c.JSON(http.StatusBadRequest, view.NewErrorResponse(viewNews.ErrCategoryAlreadyExists))
		case errors.As(err, new(diterrors.ValidationError)):
			c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		case errors.Is(err, usecaseNews.ErrAuthorNotFound):
			c.JSON(http.StatusInternalServerError, view.NewErrorResponse(viewNews.ErrAuthorNotFound))
		default:
			c.JSON(http.StatusInternalServerError, view.NewErrorResponse(view.ErrMessageInternalError))
		}
		return
	}

	resp := h.newsPresenter.CategoryToResult(res)

	c.JSON(http.StatusOK, view.NewSuccessResponse(resp))
}

func (h *newsAdminHandlers) updateCategory(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c.Request.Context(), RequestTimeOut)
	defer cancelCtx()

	pid := c.Param(ID_PARAM_KEY)
	id, err := uuid.Parse(pid)
	if err != nil || id == uuid.Nil {
		c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		return
	}

	nCategory := new(viewNews.UpdateCategory)
	if err := c.BindJSON(nCategory); err != nil {
		c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		return
	}

	nCategory.ID = id

	req := h.newsPresenter.UpdateCategoryToDTO(nCategory)
	res, err := h.categoryInteractor.Update(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, usecaseNews.ErrCategoryAlreadyExists):
			c.JSON(http.StatusBadRequest, view.NewErrorResponse(viewNews.ErrCategoryAlreadyExists))
		case errors.As(err, new(diterrors.ValidationError)):
			c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		case errors.Is(err, usecaseNews.ErrAuthorNotFound):
			c.JSON(http.StatusNotFound, view.NewErrorResponse(viewNews.ErrAuthorNotFound))
		case errors.Is(err, usecaseNews.ErrCategoryNotFound):
			c.JSON(http.StatusInternalServerError, view.NewErrorResponse(viewNews.ErrCategoryNotFound))
		default:
			c.JSON(http.StatusInternalServerError, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		}
		return
	}

	resp := h.newsPresenter.CategoryToResult(res)

	c.JSON(http.StatusOK, view.NewSuccessResponse(resp))
}

func (h *newsAdminHandlers) deleteCategory(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c.Request.Context(), RequestTimeOut)
	defer cancelCtx()

	pid := c.Param(ID_PARAM_KEY)
	id, err := uuid.Parse(pid)
	if err != nil || id == uuid.Nil {
		c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		return
	}

	if err := h.categoryInteractor.Delete(ctx, id); err != nil {
		switch {
		case errors.As(err, new(diterrors.ValidationError)):
			c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		case errors.Is(err, usecaseNews.ErrCategoryNotFound):
			c.JSON(http.StatusNotFound, view.NewErrorResponse(view.ErrMessageNotFound))
		default:
			c.JSON(http.StatusInternalServerError, view.NewErrorResponse(view.ErrMessageInternalError))
		}
		return
	}

	c.JSON(http.StatusOK, view.NewSuccessResponse(nil))
}

func (h *newsAdminHandlers) searchCategory(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c.Request.Context(), RequestTimeOut)
	defer cancelCtx()

	req := new(viewNews.SearchCategoryRequest)
	if err := c.BindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		return
	}

	dtoReq := &dtoNews.SearchCategory{
		Query: req.Query,
		Pagination: dtoNews.CategoryPagination{
			Page:  req.Page,
			Limit: req.Limit,
		},
	}

	res, err := h.categoryInteractor.Search(ctx, dtoReq)
	if err != nil {
		switch {
		case errors.As(err, new(diterrors.ValidationError)):
			c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		default:
			c.JSON(http.StatusInternalServerError, view.NewErrorResponse(view.ErrMessageInternalError))
		}
		return
	}

	response := &viewNews.SearchCategoryResult{
		TotalCount: res.Pagination.Total,
		Categories: make([]viewNews.SearchCategoryResultItem, 0, len(res.Categories)),
	}
	for _, cat := range res.Categories {
		response.Categories = append(response.Categories, viewNews.SearchCategoryResultItem{
			ID:   cat.ID,
			Name: cat.Name,
		})
	}

	c.JSON(http.StatusOK, view.NewSuccessResponse(response))
}

func (h *newsAdminHandlers) getCategory(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c.Request.Context(), RequestTimeOut)
	defer cancelCtx()

	pid := c.Param(ID_PARAM_KEY)
	id, err := uuid.Parse(pid)
	if err != nil || id == uuid.Nil {
		c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		return
	}

	res, err := h.categoryInteractor.Get(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, usecaseNews.ErrCategoryNotFound):
			c.JSON(http.StatusNotFound, view.NewErrorResponse(view.ErrMessageNotFound))
		case errors.As(err, new(diterrors.ValidationError)):
			c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		default:
			c.JSON(http.StatusInternalServerError, view.NewErrorResponse(view.ErrMessageInternalError))
		}
		return
	}

	// TODO: когда устаканится контракт переделать на нормальную структуру
	type response struct {
		Id        uuid.UUID  `json:"id"`
		Title     string     `json:"title"`
		UpdatedAt *time.Time `json:"updatedAt"`
	}

	resp := &response{
		Id:    res.ID,
		Title: res.Name,
	}
	c.JSON(http.StatusOK, view.NewSuccessResponse(resp))
}
