package data

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"redistore/internal/domain"
	"redistore/internal/domain/ports"
	"redistore/pkg/yerror"
	"time"
)

const (
	getProductByIDKey = "product:"
	getCardByIDKey    = "card:"
	getProducts       = "product:all"

	cacheDurationTime = 100 * time.Hour
)

type DBDataSource interface {
	AutoMigrate() error

	InsertProduct(ctx context.Context, tx domain.Product) (*domain.Product, error)
	GetProductByID(ctx context.Context, id string) (*domain.Product, error)
	GetProductList(ctx context.Context) ([]domain.Product, error)
	SearchProductsByTitle(ctx context.Context, titleKeywords string) ([]domain.Product, error)

	InsertCard(ctx context.Context, tx domain.Card) (*domain.Card, error)
	UpdateCard(ctx context.Context, tx domain.Card) error
	GetCardByID(ctx context.Context, id string) (*domain.Card, error)
}

type CacheDataSource interface {
	Set(ctx context.Context, key string, data []byte, time time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	FlushKey(ctx context.Context, key string) error
	FlushAll(ctx context.Context) error
}
type SearchDataSource interface {
	Set(ctx context.Context, ID uint, Title string, Description string, Price uint,
		Category domain.Category, CreatedAt int64, UpdatedAt int64) error
	Get(ctx context.Context, keywords string) ([]domain.Product, error)
}

func NewRepository(dbDS DBDataSource, chDS CacheDataSource, srchDS SearchDataSource) ports.Repository {
	return repository{
		databaseDS: dbDS,
		cacheDS:    chDS,
		srchDS:     srchDS,
	}
}

type repository struct {
	databaseDS DBDataSource
	cacheDS    CacheDataSource
	srchDS     SearchDataSource
}

func (r repository) GetCardByID(ctx context.Context, id string) (*domain.Card, error) {
	const op yerror.Op = "product_repository.GetCardByID"
	card := new(domain.Card)

	getByIDCacheKey := getCardByIDKey + id

	cache, err := r.cacheDS.Get(ctx, getByIDCacheKey)
	if err != nil {
		return nil, yerror.E(op, err)
	}
	if cache != "" {
		err = json.Unmarshal([]byte(cache), &card)
		if err != nil {
			return nil, yerror.E(op, err)
		}
		return card, nil
	}

	card, err = r.databaseDS.GetCardByID(ctx, id)
	if err != nil {
		return nil, yerror.E(op, err)
	}

	go func() {
		setCache, _ := json.Marshal(card)
		err = r.cacheDS.Set(ctx, getByIDCacheKey, setCache, cacheDurationTime)
		if err != nil {
			log.Print("err while setting redis cache :", err)
		}
	}()

	return card, nil
}

func (r repository) InsertCard(ctx context.Context, card domain.Card) (*domain.Card, error) {
	insertedCard, err := r.databaseDS.InsertCard(ctx, card)
	if err != nil {
		return nil, err
	}
	getByIDCacheKey := fmt.Sprintf("%s%v", getCardByIDKey, insertedCard.ID)
	go func() {
		setCache, _ := json.Marshal(insertedCard)
		err = r.cacheDS.Set(ctx, getByIDCacheKey, setCache, cacheDurationTime)
		if err != nil {
			log.Print("err while setting redis cache :", err)
		}
	}()
	return insertedCard, nil
}

func (r repository) UpdateCard(ctx context.Context, card domain.Card) error {
	err := r.databaseDS.UpdateCard(ctx, card)
	if err != nil {
		return err
	}
	getByIDCacheKey := fmt.Sprintf("%s%v", getCardByIDKey, card.ID)
	go func() {
		setCache, _ := json.Marshal(card)
		err := r.cacheDS.Set(ctx, getByIDCacheKey, setCache, cacheDurationTime)
		if err != nil {
			log.Print("err while flushing redis cache :", err)
		}
	}()
	return nil
}

func (r repository) SearchProductsByTitle(ctx context.Context, titleKeywords string) ([]domain.Product, error) {
	const op yerror.Op = "product_repository.SearchProductsByTitle"
	products, err := r.srchDS.Get(ctx, titleKeywords)
	if err != nil {
		return nil, yerror.E(op, err)
	}
	if len(products) > 0 {
		return products, nil
	}
	products, err = r.databaseDS.SearchProductsByTitle(ctx, titleKeywords)
	if err != nil {
		return nil, yerror.E(op, err)
	}

	return products, nil

}

func (r repository) InsertProduct(ctx context.Context, product domain.Product) (*domain.Product, error) {
	insertedProduct, err := r.databaseDS.InsertProduct(ctx, product)
	if err != nil {
		return nil, err
	}
	getByIDCacheKey := fmt.Sprintf("%s%v", getProductByIDKey, insertedProduct.ID)
	go func() {
		setCache, _ := json.Marshal(insertedProduct)
		err = r.cacheDS.Set(ctx, getByIDCacheKey, setCache, cacheDurationTime)
		if err != nil {
			log.Print("err while setting redis cache :", err)
		}
		err = r.cacheDS.FlushKey(ctx, getProducts)
		if err != nil {
			log.Print("err while deleting key in redis cache :", err)
		}

		err = r.srchDS.Set(ctx, insertedProduct.ID, insertedProduct.Title, insertedProduct.Description, insertedProduct.Price,
			insertedProduct.Category, insertedProduct.UpdatedAt, insertedProduct.CreatedAt)
		if err != nil {
			log.Print("err while setting redis cache :", err)
		}
	}()
	return insertedProduct, nil
}

func (r repository) GetProductByID(ctx context.Context, id string) (*domain.Product, error) {
	const op yerror.Op = "product_repository.GetProductByID"
	product := new(domain.Product)

	getByIDCacheKey := getProductByIDKey + id

	cache, err := r.cacheDS.Get(ctx, getByIDCacheKey)
	if err != nil {
		return nil, yerror.E(op, err)
	}
	if cache != "" {
		err = json.Unmarshal([]byte(cache), &product)
		if err != nil {
			return nil, yerror.E(op, err)
		}
		return product, nil
	}

	product, err = r.databaseDS.GetProductByID(ctx, id)
	if err != nil {
		return nil, yerror.E(op, err)
	}

	go func() {
		setCache, _ := json.Marshal(product)
		err = r.cacheDS.Set(ctx, getByIDCacheKey, setCache, cacheDurationTime)
		if err != nil {
			log.Print("err while setting redis cache :", err)
		}
	}()

	return product, nil
}

func (r repository) GetProductList(ctx context.Context) ([]domain.Product, error) {
	const op yerror.Op = "product_repository.GetProductList"
	products := make([]domain.Product, 0)

	listCacheKey := getProducts

	cache, err := r.cacheDS.Get(ctx, listCacheKey)
	if err != nil {
		return nil, yerror.E(op, err)
	}
	go func() {
		for _, product := range products {
			err = r.srchDS.Set(ctx, product.ID, product.Title, product.Description, product.Price,
				product.Category, product.UpdatedAt, product.CreatedAt)
			if err != nil {
				log.Print("err while setting redis cache :", err)
			}
		}
	}()
	if cache != "" {
		err = json.Unmarshal([]byte(cache), &products)
		if err != nil {
			return nil, yerror.E(op, err)
		}

		go func() {
			for _, product := range products {
				err = r.srchDS.Set(ctx, product.ID, product.Title, product.Description, product.Price, product.Category, product.CreatedAt, product.UpdatedAt)
				if err != nil {
					log.Fatal(err)
				}
			}
		}()

		return products, nil
	}

	products, err = r.databaseDS.GetProductList(ctx)
	if err != nil {
		return nil, yerror.E(op, err)
	}

	go func() {
		setCache, _ := json.Marshal(products)
		err = r.cacheDS.Set(ctx, listCacheKey, setCache, cacheDurationTime)
		if err != nil {
			log.Print("err while setting redis cache :", err)
		}
	}()

	return products, nil
}
