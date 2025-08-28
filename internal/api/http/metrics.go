package http

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpDurationMetric *prometheus.HistogramVec
)

func InitHTTPMetrics() {
	if gin.Mode() == gin.TestMode {
		return
	}
	httpDurationMetric = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_handler_handling_seconds",
		Help: "Histogram of response latency (seconds) of the HTTP until it is finished by the application",
	}, []string{"method", "path"})
}
