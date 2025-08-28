package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view"
)

var AppInstance appInstance

type appInstance string

const (
	AppInstanceDevelop appInstance = "develop"
	AppInstanceProd    appInstance = "prod"
)

const StatusCodeHeader = "X-Status-Code"

func NotImplementedHandler(c *gin.Context) {
	if AppInstance == AppInstanceDevelop {
		c.Header("req-method", c.Request.Method)
		c.Header("req-url", c.Request.URL.String())
		c.Header("req-uri", c.Request.RequestURI)
	}
	c.JSON(http.StatusMethodNotAllowed, view.NewErrorResponse(view.ErrMessageMethodNotAllowed))
}
