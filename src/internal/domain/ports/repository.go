package ports

import (
	"context"

	"redistore/internal/domain"
)

// Repository is an interface to be implemented for some
// operation related to Product & Card entity
type Repository interface {

	// Insert creates a new record in db and returns the stored item
	// or an error if there was problem
	InsertProduct(ctx context.Context, product domain.Product) (*domain.Product, error)

	// SearchProductsByTitle make a full-text search and returns matched products.
	SearchProductsByTitle(ctx context.Context, titleKeywords string) ([]domain.Product, error)

	// GetProductByID gets an id and , find related product in the database and return it.
	GetProductByID(ctx context.Context, id string) (*domain.Product, error)

	// GetProductList  find  all products in the database and return them.
	GetProductList(ctx context.Context) ([]domain.Product, error)

	// GetProductByID gets an id and , find related card in the database and return it.
	GetCardByID(ctx context.Context, id string) (*domain.Card, error)

	// InsertCard creates a new record in db and returns the stored item
	// or an error if there was problem
	InsertCard(ctx context.Context, card domain.Card) (*domain.Card, error)

	// UpdateCard gets an Card entity, find it in the database and update it.
	UpdateCard(ctx context.Context, card domain.Card) error
}
