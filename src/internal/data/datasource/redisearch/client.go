package redisearch

import (
	"context"
	"github.com/RediSearch/redisearch-go/redisearch"
	"log"
	"redistore/internal/data"
	"redistore/internal/domain"
	"strconv"
)

func NewSearchDataSource(redisearch *redisearch.Client) data.SearchDataSource {
	return &cacheDataSource{
		redisearch: redisearch,
	}
}

type cacheDataSource struct {
	redisearch *redisearch.Client
}

func (c cacheDataSource) Set(ctx context.Context, ID uint, Title string, Description string, Price uint, Category domain.Category, CreatedAt int64, UpdatedAt int64) error {
	docID := "rs:product:" + strconv.FormatUint(uint64(ID), 10)
	currentDoc, err := c.redisearch.Get(docID)

	if err != nil {
		return err
	}
	if currentDoc != nil {
		err = c.redisearch.DeleteDocument(docID)
		if err != nil {
			return err
		}
	}
	// Create a document with an id and given score
	doc := redisearch.NewDocument(docID, 1.0)
	doc.Set("ID", ID).Set("Title", Title).Set("Description", Description).Set("Price", Price).
		Set("Category", string(Category)).Set("CreatedAt", CreatedAt).
		Set("UpdatedAt", UpdatedAt)

	// Index the document. The API accepts multiple documents at a time
	if err := c.redisearch.IndexOptions(redisearch.DefaultIndexingOptions, doc); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (c cacheDataSource) Get(ctx context.Context, keywords string) ([]domain.Product, error) {
	// Searching with limit and sorting
	//docs, total, err := c.redisearch.Search(redisearch.NewQuery(keywords).
	docs, _, err := c.redisearch.Search(redisearch.NewQuery(keywords))
	if err != nil {
		return nil, err
	}
	products := make([]domain.Product, len(docs))
	for i := 0; i < len(docs); i++ {
		id, err := strconv.Atoi(docs[i].Properties["ID"].(string))
		if err != nil {
			continue
		}

		price, err := strconv.Atoi(docs[i].Properties["Price"].(string))
		if err != nil {
			continue
		}

		createdAt, err := strconv.Atoi(docs[i].Properties["CreatedAt"].(string))
		if err != nil {
			continue
		}
		updatedAt, err := strconv.Atoi(docs[i].Properties["UpdatedAt"].(string))
		if err != nil {
			continue
		}
		product := domain.Product{
			ID:          uint(id),
			Title:       docs[i].Properties["Title"].(string),
			Description: docs[i].Properties["Description"].(string),
			Price:       uint(price),
			Category:    domain.Category(docs[i].Properties["Category"].(string)),
			CreatedAt:   int64(createdAt),
			UpdatedAt:   int64(updatedAt),
		}
		products[i] = product
	}
	return products, nil
}
