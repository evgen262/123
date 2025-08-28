package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

//go:generate mockgen -source=interfaces.go -destination=./cache_mock.go -package=cache
type Redis interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	SetEx(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	HGet(ctx context.Context, key, field string) *redis.StringCmd
	HGetAll(ctx context.Context, key string) *redis.MapStringStringCmd
	HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	Exists(ctx context.Context, keys ...string) *redis.IntCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Ping(ctx context.Context) *redis.StatusCmd
	Info(ctx context.Context, sections ...string) *redis.StringCmd
	Close() error
}
