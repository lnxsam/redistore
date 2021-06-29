package redis_test

import (
	"context"
	"encoding/json"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	caches "redistore/internal/data/datasource/redis"
	"redistore/internal/domain"
	"testing"
	"time"
)

var (
	key = "key"
)

func miniRedis() string {
	mr, err := miniredis.Run()
	if err != nil {
		return ""
	}

	return mr.Addr()
}

func TestNewCacheDataSource(t *testing.T) {
	assert.NotNil(t, caches.NewCacheDataSource(&redis.Client{}), "NewCacheDataSource() should not return nil")
}

func TestSet(t *testing.T) {
	redisAddress := miniRedis()

	require.NotNil(t, redisAddress, "invalid address")

	client := redis.NewClient(&redis.Options{
		Addr: redisAddress,
	})

	setErr := caches.NewCacheDataSource(client).Set(context.Background(), key, []byte(""), 5*time.Second)

	assert.Nil(t, setErr)
}

func TestGet(t *testing.T) {
	redisAddress := miniRedis()

	require.NotNil(t, redisAddress, "invalid address")

	client := redis.NewClient(&redis.Options{
		Addr: redisAddress,
	})
	model := &domain.Product{
		ID: 1,
	}
	cacheNetwork, _ := json.Marshal(model)

	setErr := caches.NewCacheDataSource(client).Set(context.Background(), key, cacheNetwork, 5*time.Second)

	assert.Nil(t, setErr)

	redisValue, err := caches.NewCacheDataSource(client).Get(context.Background(), key)

	assert.Nil(t, err)

	require.NotNil(t, redisValue)
	assert.Equal(t, model.ID, uint(1), "redis values are not same")
}

func TestFlushKey(t *testing.T) {
	redisAddress := miniRedis()

	require.NotNil(t, redisAddress, "invalid address")

	client := redis.NewClient(&redis.Options{
		Addr: redisAddress,
	})

	model := &domain.Product{
		ID: 1,
	}
	cacheNetwork, _ := json.Marshal(model)

	setErr := caches.NewCacheDataSource(client).Set(context.Background(), key, cacheNetwork, 5*time.Second)

	assert.Nil(t, setErr)

	err := caches.NewCacheDataSource(client).FlushKey(context.Background(), key)

	assert.Nil(t, err)

	redisValue, err := caches.NewCacheDataSource(client).Get(context.Background(), key)

	require.Nil(t, err)
	require.NotNil(t, redisValue)
}

func TestFlushAll(t *testing.T) {
	redisAddress := miniRedis()

	require.NotNil(t, redisAddress, "invalid address")

	client := redis.NewClient(&redis.Options{
		Addr: redisAddress,
	})

	model := &domain.Product{
		ID: 1,
	}
	cacheNetwork, _ := json.Marshal(model)

	setErr := caches.NewCacheDataSource(client).Set(context.Background(), key, cacheNetwork, 5*time.Second)

	assert.Nil(t, setErr)

	err := caches.NewCacheDataSource(client).FlushAll(context.Background())
	assert.Nil(t, err)

	redisValue, err := caches.NewCacheDataSource(client).Get(context.Background(), key)

	require.Nil(t, err)
	require.NotNil(t, redisValue)
}
