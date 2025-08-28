package analytics

import (
	"context"
	"fmt"

	metricsv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/analytics/metrics/v1"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"google.golang.org/grpc/codes"

	entityAnalytics "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/analytics"
)

type analyticsRepository struct {
	client metricsv1.MetricsAPIClient

	mapper MetricsMapper
}

func NewAnalyticsRepository(client metricsv1.MetricsAPIClient, mapper MetricsMapper) *analyticsRepository {
	return &analyticsRepository{
		client: client,
		mapper: mapper,
	}
}

// AddMetrics - отправляет метрики в сервис analytics
func (r analyticsRepository) AddMetrics(ctx context.Context, headers entityAnalytics.CFCHeaders, body []byte) (metricID string, err error) {
	// Отправляем запрос в analytics
	resp, err := r.client.Add(ctx, &metricsv1.AddRequest{
		Body:       body,
		CfcHeaders: r.mapper.MetricsCFCHeadersToPb(headers),
	})
	if err != nil {
		localErr := diterrors.NewLocalizedError("ru-RU", err)
		switch localErr.Code() {
		case codes.InvalidArgument:
			return "", fmt.Errorf("client.Add: %w", diterrors.NewValidationError(localErr, diterrors.ErrValidationFields{
				Field:   "in",
				Message: localErr.Error(),
			}))
		default:
			return "", fmt.Errorf("client.Add: %w", err)
		}
	}

	// Возвращаем результат
	return resp.GetMetricId(), nil
}
