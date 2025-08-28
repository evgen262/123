package analytics

import (
	metricsv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/analytics/metrics/v1"

	entityAnalytics "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/analytics"
)

type mapperAnalytics struct {
}

func NewMetricsMapper() *mapperAnalytics {
	return &mapperAnalytics{}
}

func (m mapperAnalytics) MetricsCFCHeadersToPb(request entityAnalytics.CFCHeaders) *metricsv1.AddRequest_CFCHeaders {
	return &metricsv1.AddRequest_CFCHeaders{
		CfcHeader:        string(request.Header),
		CfcPortal:        request.Portal,
		CfcAuthorization: request.Authorization,
	}
}
