package repositories

import (
	"context"
	"errors"
	"time"

	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	tokenPrefix = "tok:"
	tokenTtl    = 365 * 24 * time.Hour
)

var (
	ErrTokenNotFound = errors.New("token not found")
)

type tokenRepository struct {
	source CacheSource
	logger ditzap.Logger

	basePrefix string
}

func NewTokenRepository(basePrefix string, cacheSource CacheSource, logger ditzap.Logger) *tokenRepository {
	return &tokenRepository{
		basePrefix: basePrefix,
		source:     cacheSource,
		logger:     logger,
	}
}

func (or *tokenRepository) getKey(cloudID string) string {
	return or.basePrefix + tokenPrefix + cloudID
}

func (or *tokenRepository) Save(ctx context.Context, id, token string) error {
	err := or.source.SetEx(ctx, or.getKey(id), token, tokenTtl)
	if err != nil {
		or.logger.Error("не удалось сохранить refresh токен",
			zap.Error(err),
		)
		return err
	}
	return nil
}

func (or *tokenRepository) Get(ctx context.Context, id string) (string, error) {
	token, err := or.source.Get(ctx, or.getKey(id))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", ErrNotFound
		}
		or.logger.Error("не удалось получить refresh токен",
			zap.Error(err),
		)
		return "", err
	}
	return token, nil
}

func (or *tokenRepository) Delete(ctx context.Context, id string) {
	if err := or.source.Delete(ctx, or.getKey(id)); err != nil {
		or.logger.Error("не удалось удалить refresh токен",
			zap.Error(err),
		)
	}
}
