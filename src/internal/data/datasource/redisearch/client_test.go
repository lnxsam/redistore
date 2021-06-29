package redisearch_test

import (
	"context"
	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	search "redistore/internal/data/datasource/redisearch"
	"redistore/internal/domain"
	"testing"
)

var (
	keyword = "Title"
)

func miniRedis() string {
	mr, err := miniredis.Run()
	if err != nil {
		return ""
	}

	return mr.Addr()
}

func TestNewCacheDataSource(t *testing.T) {
	assert.NotNil(t, search.NewSearchDataSource(&redisearch.Client{}), "NewSearchDataSource() should not return nil")
}

func TestSet(t *testing.T) {
	redisAddress := miniRedis()

	require.NotNil(t, redisAddress, "invalid address")

	client := redisearch.NewClient("localhost:6379", "redistore_index")

	model := &domain.Product{
		ID:          1,
		Title:       "Product Title",
		Price:       1000,
		Description: "Description",
	}
	setErr := search.NewSearchDataSource(client).Set(context.Background(), model.ID, model.Title, model.Description, model.Price, model.Category, model.CreatedAt, model.UpdatedAt)

	assert.Nil(t, setErr)
}

func TestGet(t *testing.T) {
	redisAddress := miniRedis()

	require.NotNil(t, redisAddress, "invalid address")

	client := redisearch.NewClient("localhost:6379", "redistore_index")
	model := &domain.Product{
		ID:          1,
		Title:       "Product" + keyword,
		Price:       1000,
		Description: "Description",
	}

	setErr := search.NewSearchDataSource(client).Set(context.Background(), model.ID, model.Title, model.Description, model.Price, model.Category, model.CreatedAt, model.UpdatedAt)

	assert.Nil(t, setErr)

	redisValue, err := search.NewSearchDataSource(client).Get(context.Background(), keyword)

	assert.Nil(t, err)

	require.NotNil(t, redisValue)
	assert.Equal(t, model.ID, uint(1), "redisearch values are not same")
}
