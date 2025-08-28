package http

import (
	"context"
	"errors"
	"net/http"

	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view"
)

type filesHandlers struct {
	filesInteractor FilesInteractor
}

func NewFilesHandlers(filesInteractor FilesInteractor) *filesHandlers {
	return &filesHandlers{
		filesInteractor: filesInteractor,
	}
}

// @Summary Получение публичного файла по идентификатору
// @Description Возвращает файл на основе переданного идентификатора файла
// @Tags     Файлы
// @Produce  octet-stream
// @Param    file_id path string true "Идентификатор файла"
// @Router   /{file_id} [get]
// @Success  200 {file} octet-stream "Файл"
// @Failure  400,500 {object} ErrorResponse
// @Security ApiKeyAuth
func (fh *filesHandlers) get(c *gin.Context) {
	ctx, cancelCtx := context.WithTimeout(c.Request.Context(), RequestTimeOut)
	defer cancelCtx()

	fileId, err := uuid.Parse(c.Param("file_id"))
	if err != nil {
		c.Header(StatusCodeHeader, "WAF_01")
		c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		return
	}

	result, err := fh.filesInteractor.Get(ctx, fileId)
	if err != nil {
		switch {
		case errors.Is(err, diterrors.ErrUnauthenticated):
			c.Header(StatusCodeHeader, "WAF_02")
			c.JSON(http.StatusUnauthorized, view.NewErrorResponse(view.ErrMessageUnauthenticated))
			return
		case errors.Is(err, diterrors.ErrPermissionDenied):
			c.Header(StatusCodeHeader, "WAF_03")
			c.JSON(http.StatusForbidden, view.NewErrorResponse(view.ErrPermissionDenied))
			return
		case errors.Is(err, diterrors.ErrNotFound):
			c.Header(StatusCodeHeader, "WAF_04")
			c.JSON(http.StatusNotFound, view.NewErrorResponse(view.ErrMessageNotFound))
			return
		case errors.As(err, new(diterrors.ValidationError)):
			c.Header(StatusCodeHeader, "WAF_05")
			c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
			return
		default:
			c.Header(StatusCodeHeader, "WAF_06")
			c.JSON(http.StatusInternalServerError, view.NewErrorResponse(view.ErrMessageInternalError))
			return
		}
	}

	c.Writer.Header().Set("X-File-Name", result.GetFileName())
	c.Data(http.StatusOK, result.Metadata.ContentType, result.Payload)
}
