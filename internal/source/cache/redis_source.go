package cache

import (
	"context"
	"time"
)

type cacheSource struct {
	redis Redis
}

func NewCacheSource(redis Redis) *cacheSource {
	return &cacheSource{
		redis: redis,
	}
}

func (s *cacheSource) Get(ctx context.Context, key string) (string, error) {
	return s.redis.Get(ctx, key).Result()
}

func (s *cacheSource) SetEx(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return s.redis.SetEx(ctx, key, value, ttl).Err()
}

func (s *cacheSource) Exists(ctx context.Context, keys ...string) (bool, error) {
	cnt, err := s.redis.Exists(ctx, keys...).Result()
	if err != nil {
		return false, err
	}
	return int64(len(keys)) == cnt, nil
}

func (s *cacheSource) Delete(ctx context.Context, keys ...string) error {
	if err := s.redis.Del(ctx, keys...).Err(); err != nil {
		return err
	}
	return nil
}
