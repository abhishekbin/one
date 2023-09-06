package onecache

import (
	"context"
	"errors"
	"time"
)

type RedisCache struct {
	// TODO: Implement.
}

func (c *RedisCache) Get(ctx context.Context, key string) (interface{}, bool, error) {
	return nil, false, errors.New("TODO: Implement")
}

func (c *RedisCache) Set(ctx context.Context, key string, val interface{}, expireAfter time.Duration) error {
	return errors.New("TODO: Implement")
}

func (c *RedisCache) Clear(ctx context.Context, key string) (bool, error) {
	return false, errors.New("TODO: Implement")
}

var _ Cache = (*RedisCache)(nil)
