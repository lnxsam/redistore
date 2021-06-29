package redis

import (
	"context"
	"time"

	"redistore/internal/data"
	"redistore/pkg/yerror"

	redisPkg "github.com/go-redis/redis/v8"
)

func NewCacheDataSource(redis *redisPkg.Client) data.CacheDataSource {
	return &cacheDataSource{
		redis: redis,
	}
}

type cacheDataSource struct {
	redis *redisPkg.Client
}

func (c *cacheDataSource) Set(ctx context.Context, key string, data []byte, time time.Duration) error {
	const op yerror.Op = "cache_data_source.Set"
	err := c.redis.Set(ctx, key, data, time).Err()
	if err != nil {
		return yerror.E(op, err)
	}
	return nil
}

func (c *cacheDataSource) Get(ctx context.Context, key string) (string, error) {
	const op yerror.Op = "cache_data_source.Get"
	redisValue, err := c.redis.Get(ctx, key).Result()
	if err == redisPkg.Nil {
		return "", nil
	} else if err != nil {
		return "", yerror.E(op, err)
	}
	return redisValue, nil
}

func (c *cacheDataSource) FlushKey(ctx context.Context, key string) error {
	const op yerror.Op = "cache_data_source.FlushKey"

	err := c.redis.Del(ctx, key).Err()
	if err != nil {
		return yerror.E(op, err)
	}

	return nil
}

func (c *cacheDataSource) FlushAll(ctx context.Context) error {
	const op yerror.Op = "cache_data_source.FlushAll"
	err := c.redis.FlushAll(ctx).Err()
	if err != nil {
		return yerror.E(op, err)
	}
	return nil
}
