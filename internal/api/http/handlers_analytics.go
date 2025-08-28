package http

import (
	"context"
	"errors"
	"io"
	"net/http"

	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"github.com/gin-gonic/gin"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view"
	viewAnalytics "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/analytics"
	entityAnalytics "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/analytics"
	interactorAnalytics "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/usecase/analytics"
)

const headerXCFCUserAgent = "X-CFC-UserAgent"

type analyticsHandlers struct {
	interactor AnalyticsInteractor
}

func NewAnalyticsHandlers(interactor AnalyticsInteractor) analyticsHandlers {
	return analyticsHandlers{
		interactor: interactor,
	}
}

func (h analyticsHandlers) addMetrics(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), RequestTimeOut)
	defer cancel()

	// Валидируем заголовок X-CFC-UserAgent
	cfcUserAgentHeader := entityAnalytics.XCFCUserAgentHeader(c.GetHeader(headerXCFCUserAgent))
	if !cfcUserAgentHeader.IsValid() {
		c.Header(StatusCodeHeader, "AM_01")
		c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		return
	}

	// Валидация тела выполняется в сервисе аналитики.
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.Header(StatusCodeHeader, "AM_02")
		c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		return
	}

	// Отправляем метрики
	metricID, err := h.interactor.AddMetrics(ctx, cfcUserAgentHeader, body)
	// Обрабатываем ошибки
	if err != nil {
		switch {
		case errors.Is(err, interactorAnalytics.ErrUnauthenticated):
			// Не авторизованный пользователь
			c.Header(StatusCodeHeader, "AM_03")
			c.JSON(http.StatusUnauthorized, view.NewErrorResponse(view.ErrMessageUnauthenticated))
		case errors.As(err, new(diterrors.ValidationError)):
			// Не валидный Request
			c.Header(StatusCodeHeader, "AM_04")
			c.JSON(http.StatusBadRequest, view.NewErrorResponse(view.ErrMessageInvalidRequest))
		default:
			// Internal
			c.Header(StatusCodeHeader, "AM_05")
			c.JSON(http.StatusInternalServerError, view.NewErrorResponse(view.ErrMessageInternalError))
		}
		return
	}

	// Все ОК возвращаем результат
	c.JSON(http.StatusOK, view.NewSuccessResponse(viewAnalytics.AddMetricsResponse{
		MetricID: metricID,
	}))
}
