package analytics

import (
	"context"

	entityAnalytics "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/analytics"
)

//go:generate mockgen -source=interfaces.go -destination=./analytics_mock.go -package=analytics

type AnalyticsRepository interface {
	AddMetrics(ctx context.Context, headers entityAnalytics.CFCHeaders, body []byte) (metricID string, err error)
}
