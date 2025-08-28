package auth

import (
	"context"
	"errors"
	"fmt"

	redirectsessionv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/authfacade/redirectsession/v1"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"go.uber.org/zap"

	entitySession "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/auth"
)

type redirectSessionRepository struct {
	client redirectsessionv1.RedirectSessionAPIClient
	mapper RedirectSessionMapper

	logger ditzap.Logger
}

func NewRedirectSessionRepository(client redirectsessionv1.RedirectSessionAPIClient, mapper RedirectSessionMapper, logger ditzap.Logger) *redirectSessionRepository {
	return &redirectSessionRepository{
		client: client,
		mapper: mapper,

		logger: logger,
	}
}

func (r redirectSessionRepository) CreateSession(ctx context.Context, userInfo *entitySession.RedirectSessionUserInfo) (string, error) {
	if userInfo == nil {
		return "", fmt.Errorf("redirectSessionRepository.Create: %w", ErrUserInfoRequired)
	}

	result, err := r.client.CreateSession(ctx, &redirectsessionv1.CreateSessionRequest{
		UserInfo: r.mapper.UserInfoToCreateRequestUserUnfoPb(userInfo),
	})

	if err != nil {
		switch {
		case errors.As(err, new(diterrors.ValidationError)):
			r.logger.Warn(
				"invalid session", zap.Error(err),
				zap.String("method", "sr.client.CreateSession(ctx, session)"),
				zap.Any("session", userInfo),
			)
			return "", fmt.Errorf("client.CreateSession: invalid session: %w", err)
		default:
			r.logger.Error(
				"can't create session", zap.Error(err),
				zap.String("method", "sr.client.CreateSession(ctx, session)"),
				zap.Any("session", userInfo),
			)
			return "", fmt.Errorf("client.CreateSession: can't create session: %w", err)
		}
	}

	return result.GetSessionId(), nil
}
