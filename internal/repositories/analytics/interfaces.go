package analytics

import (
	metricsv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/analytics/metrics/v1"

	entityAnalytics "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/analytics"
)

//go:generate mockgen -source=interfaces.go -destination=./analytics_mock.go -package=analytics

type MetricsMapper interface {
	MetricsCFCHeadersToPb(headers entityAnalytics.CFCHeaders) *metricsv1.AddRequest_CFCHeaders
}
