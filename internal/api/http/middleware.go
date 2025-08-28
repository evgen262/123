package http

import (
	"context"
	"errors"
	"net"
	"net/http"
	"strings"
	"time"

	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	timeUtils "git.mos.ru/buch-cloud/moscow-team-2.0/build/time-utils.git"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view"
	authView "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/auth"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity"
	entityAuth "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/auth"
)

const onlyAuthOptKey = "onlyAuthCheck"

const (
	RefreshTokenCookie = "komanda-tkn"
	JWTTokenHeader     = "Authorization"
	UserAgentHeader    = "User-Agent"
)

const (
	XRequestId = "x-request-id"
	ClientIP   = "client-ip"
	UserAgent  = "http-user-agent"
)

func NewAuthSessionMiddleware(ai AuthInteractor, tu timeUtils.TimeUtils, opts MiddlewareOptions) gin.HandlerFunc {
	var onlyAuth bool
	if onlyAuthOpt := opts.Get(onlyAuthOptKey); onlyAuthOpt != nil {
		onlyAuth, _ = onlyAuthOpt.GetValue().(bool)
	}

	return func(c *gin.Context) {
		startTime := tu.New()
		defer func() {
			httpDurationMetric.WithLabelValues(c.Request.Method, c.Request.URL.EscapedPath()).Observe(time.Since(*startTime).Seconds())
		}()

		tokenHeader := c.GetHeader(JWTTokenHeader)
		if tokenHeader == "" {
			// Нет заголовка Authorization
			c.Header(StatusCodeHeader, "WASM_01")
			c.AbortWithStatusJSON(http.StatusUnauthorized, view.NewErrorResponse(view.ErrMessageUnauthenticated))
			return
		}

		if !strings.HasPrefix(tokenHeader, "Bearer ") {
			// В заголовке Authorization отсутствует Bearer
			c.Header(StatusCodeHeader, "WASM_02")
			c.AbortWithStatusJSON(http.StatusUnauthorized, view.NewErrorResponse(view.ErrMessageUnauthenticated))
			return
		}

		accessToken := strings.TrimPrefix(tokenHeader, "Bearer ")

		// TODO опция ингорируется, поскольку пока отсутствует ручка для проверки существования сессии, всегда выполняется получение сессии
		_ = onlyAuth
		/*
			if onlyAuth {
				// Проверяем существование сессии (токен валидный и не протух, сессия по id из payload найдена)
				if !uc.SessionExists(c, token) {
					c.AbortWithStatusJSON(http.StatusForbidden, view.NewErrorResponse(view.ErrMessageUnauthorized))
					return
				}
				c.Next()
				return
			}
		*/

		// Получаем сессию по токену (токен валидный и не протух, сессия по id найдена и не протухла)
		session, err := ai.GetSession(c, accessToken)
		if err != nil {
			switch {
			case errors.As(err, new(diterrors.ValidationError)):
				// Невалидный JWT токен или идентификатор сессии
				c.Header(StatusCodeHeader, "WASM_03")
				c.AbortWithStatusJSON(http.StatusForbidden, view.NewErrorResponse(view.ErrMessageInvalidRequest))
			case errors.Is(err, diterrors.ErrFailedPrecondition):
				// JWT токен истёк
				c.Header(StatusCodeHeader, "WASM_04")
				c.AbortWithStatusJSON(http.StatusUnauthorized, view.NewErrorResponse(view.ErrMessageUnauthenticated))
			case errors.Is(err, diterrors.ErrNotFound):
				// Сессия не найдена
				c.Header(StatusCodeHeader, "WASM_05")
				c.AbortWithStatusJSON(http.StatusUnauthorized, view.NewErrorResponse(view.ErrMessageUnauthorized))
			case errors.Is(err, diterrors.ErrUnauthenticated):
				// Сессия истекла
				c.Header(StatusCodeHeader, "WASM_06")
				c.AbortWithStatusJSON(http.StatusUnauthorized, view.NewErrorResponse(authView.ErrMessageFrontUnauthenticated))
			default:
				// Внутренняя ошибка сервиса
				c.Header(StatusCodeHeader, "WASM_07")
				c.AbortWithStatusJSON(http.StatusInternalServerError, view.NewErrorResponse(view.ErrMessageInternalError))

			}
			return
		}

		if session == nil {
			c.Header(StatusCodeHeader, "WASM_08")
			c.AbortWithStatusJSON(http.StatusUnauthorized, view.NewErrorResponse(view.ErrMessageUnauthorized))
			return
		}

		if SharedFields.AccessList != nil {
			if !SharedFields.AccessList.Have(session.GetUser().GetEmail()) {
				c.Header(StatusCodeHeader, "WASM_09")
				c.AbortWithStatusJSON(http.StatusForbidden, view.NewErrorResponse(view.ErrMessageUserAccessDenied))
				return
			}
		}

		c.Request = c.Request.WithContext(entity.WithSession(c.Request.Context(), session))

		/*
			TODO: подумать об изменении активности сессии
			session.LastActive = startTime
			c.Request = c.Request.WithContext(entity.WithSession(c.Request.Context(), sessionRepository.Update(c, session)))
		*/

		c.Next()
	}
}

func NewRequestIDMiddleware(_ MiddlewareOptions) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqId := uuid.NewString()
		var mdCtx context.Context
		if _, ok := metadata.FromOutgoingContext(c.Request.Context()); ok {
			mdCtx = metadata.AppendToOutgoingContext(c.Request.Context(), XRequestId, reqId)
		} else {
			mdCtx = metadata.NewOutgoingContext(c.Request.Context(), metadata.New(map[string]string{XRequestId: reqId}))
		}
		c.Header(XRequestId, reqId)
		c.Request = c.Request.WithContext(mdCtx)

		c.Next()
	}
}

func NewHeadersMiddleware(_ MiddlewareOptions) gin.HandlerFunc {
	return func(c *gin.Context) {
		userAgent := c.GetHeader(UserAgentHeader)
		clientIP := c.ClientIP()

		md := map[string]string{
			ClientIP:  clientIP,
			UserAgent: userAgent,
		}
		device := &entityAuth.Device{
			UserAgent: userAgent,
		}

		var mdCtx context.Context
		if _, ok := metadata.FromOutgoingContext(c.Request.Context()); ok {
			mdKV := make([]string, 0, len(md))
			for k, v := range md {
				mdKV = append(mdKV, k, v)
			}
			mdCtx = metadata.AppendToOutgoingContext(c.Request.Context(), mdKV...)
		} else {
			mdCtx = metadata.NewOutgoingContext(c.Request.Context(), metadata.New(md))
		}

		c.Request = c.Request.WithContext(mdCtx)
		c.Request = c.Request.WithContext(entity.WithClientIP(c.Request.Context(), net.ParseIP(clientIP)))
		c.Request = c.Request.WithContext(entity.WithDevice(c.Request.Context(), device))

		c.Next()
	}
}
