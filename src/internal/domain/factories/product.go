package factories

import (
	"math/rand"
	"redistore/internal/domain"
	"time"
)

// ProductFactory is contract for Product factory
type ProductFactory interface {
	Create() domain.Product
	CreateMany(int) []domain.Product
}

var Product ProductFactory = &productFactory{}

type productFactory struct {
}

func (af *productFactory) Create() domain.Product {
	return domain.Product{
		ID:          uint(rand.Uint32()),
		Title:       "Product Title",
		Category:    domain.Car,
		Price:       1000,
		Description: "Description",
		CreatedAt:   time.Now().UTC().Unix(),
	}
}

func (af *productFactory) CreateMany(count int) []domain.Product {
	var products []domain.Product
	for i := 0; i < count; i++ {
		products = append(products, af.Create())
	}
	return products
}
