package analytics

import (
	"context"
	"fmt"

	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"go.uber.org/zap"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity"
	entityAnalytics "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/analytics"
)

type analyticsInteractor struct {
	repository AnalyticsRepository

	logger ditzap.Logger
}

func NewAnalyticsInteractor(repository AnalyticsRepository, logger ditzap.Logger) *analyticsInteractor {
	return &analyticsInteractor{
		repository: repository,

		logger: logger,
	}
}

// AddMetrics - отправляет метрики в сервис analytics
func (i analyticsInteractor) AddMetrics(ctx context.Context, headers entityAnalytics.XCFCUserAgentHeader, body []byte) (string, error) {
	if len(body) == 0 {
		i.logger.Error("body is empty", zap.String("call", "analyticsInteractor.AddMetrics"))
		return "", diterrors.NewValidationError(ErrBodyIsEmpty, diterrors.ErrValidationFields{
			Field:   "body",
			Message: ErrBodyIsEmpty.Error(),
		})
	}

	if headers == "" {
		i.logger.Error("headers are empty", zap.String("call", "analyticsInteractor.AddMetrics"))
		return "", diterrors.NewValidationError(ErrHeadersAreAmpty, diterrors.ErrValidationFields{
			Field:   "headers",
			Message: ErrHeadersAreAmpty.Error(),
		})
	}

	i.logger.Debug(
		"add metrics data",
		zap.String("call", "analyticsInteractor.AddMetrics"),
		zap.String("body", string(body)),
		zap.String("headers", string(headers)),
	)

	session, err := entity.SessionFromContext(ctx)
	if err != nil || session == nil {
		i.logger.Error("invalid session", zap.Error(err), zap.String("call", "entity.SessionFromContext"))
		return "", fmt.Errorf("entity.SessionFromContext: %w", ErrUnauthenticated)
	}

	cfcHeaders := entityAnalytics.CFCHeaders{Header: headers}

	activePortal := session.GetActivePortal()
	if activePortal == nil {
		i.logger.Warn("empty active portal", zap.String("call", "session.GetActivePortal"))
		cfcHeaders.Portal = activePortal.GetPortal().URL
	}

	// TODO: User не передается в сессию контекста.
	if user := session.GetUser(); user != nil {
		if person := user.GetPerson(); person != nil {
			cfcHeaders.Authorization = person.ExtID
		} else {
			i.logger.Warn("empty person", zap.String("call", "user.GetPerson"))
		}
	} else {
		i.logger.Warn("empty user", zap.String("call", "session.GetUser"))
	}

	metricID, err := i.repository.AddMetrics(ctx, cfcHeaders, body)
	if err != nil {
		i.logger.Error("can't add metrics in repository", zap.Error(err),
			zap.String("call", "repository.AddMetrics"),
			zap.String("body", string(body)),
			zap.String("headers", string(headers)),
		)
		return "", fmt.Errorf("repository.AddMetrics: %w", err)
	}

	return metricID, nil
}
