package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"git.mos.ru/buch-cloud/moscow-team-2.0/pud/auth.git/internal/entity"
	"git.mos.ru/buch-cloud/moscow-team-2.0/pud/auth.git/internal/util/converter"
)

const (
	statePrefix = "state:"
	stateTTL    = 5 * time.Minute
)

var (
	ErrStateNotFound = errors.New("state not found")
	ErrNotFound      = errors.New("not found")
	ErrInvalidState  = errors.New("invalid state")
)

type stateRepository struct {
	source CacheSource
	logger ditzap.Logger

	basePrefix string
}

func NewStateRepository(basePrefix string, cacheSource CacheSource, logger ditzap.Logger) *stateRepository {
	return &stateRepository{
		basePrefix: basePrefix,
		source:     cacheSource,
		logger:     logger,
	}
}

func (sr *stateRepository) getKey(id string) string {
	return sr.basePrefix + statePrefix + id
}

// New генерация нового state
func (sr *stateRepository) New(ctx context.Context, options *entity.StateOptions) (*entity.State, error) {
	stateID := uuid.NewString()

	state := &entity.State{
		ID: stateID,
	}

	if options != nil {
		if options.CallbackURL != "" {
			state.CallbackURL = options.CallbackURL
		}

		if options.ClientID != "" {
			state.ClientID = options.ClientID
		}

		if options.ClientSecret != "" {
			state.ClientSecret = options.ClientSecret
		}

		if options.CodeVerifier != "" {
			state.CodeVerifier = options.CodeVerifier
		}

		if options.DeviceID != "" {
			state.DeviceID = options.DeviceID
		}

		if options.UserAgent != "" {
			state.UserAgent = options.UserAgent
		}
	}

	data, err := json.Marshal(state)
	if err != nil {
		return nil, err
	}

	if err := sr.source.SetEx(ctx, sr.getKey(stateID), data, stateTTL); err != nil {
		sr.logger.Error("не удалось сохранить state",
			zap.Error(err),
		)
		return nil, err
	}
	return state, nil
}

// IsExists проверка существования идентификатора
func (sr *stateRepository) IsExists(ctx context.Context, stateID string) error {
	if _, err := uuid.Parse(stateID); err != nil {
		return ErrInvalidState
	}

	exists, err := sr.source.Exists(ctx, sr.getKey(stateID))
	if err != nil {
		sr.logger.Error("не удалось получить state",
			zap.Error(err),
		)
		return err
	}

	if !exists {
		return ErrStateNotFound
	}
	return nil
}

// Delete удаление идентификатора
func (sr *stateRepository) Delete(ctx context.Context, stateID string) {
	if err := sr.source.Delete(ctx, sr.getKey(stateID)); err != nil {
		sr.logger.Error("не удалось удалить state",
			zap.Error(err),
		)
	}
}

// Get получение state
func (sr *stateRepository) Get(ctx context.Context, stateID string) (*entity.State, error) {
	dataStr, err := sr.source.Get(ctx, sr.getKey(stateID))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrNotFound
		}
		sr.logger.Error("не удалось получить state",
			zap.Error(err),
		)
		return nil, err
	}

	data := converter.StringToBytes(dataStr)

	var state entity.State
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}

	return &state, nil
}
