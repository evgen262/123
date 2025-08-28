package repositories

import (
	"context"
	"time"
)

//go:generate mockgen -source=interfaces.go -destination=./repositories_mock.go -package=repositories
type CacheSource interface {
	Get(ctx context.Context, key string) (string, error)
	SetEx(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Exists(ctx context.Context, keys ...string) (bool, error)
	Delete(ctx context.Context, keys ...string) error
}
