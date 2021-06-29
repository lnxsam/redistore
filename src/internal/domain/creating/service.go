package creating

import (
	"context"
	"redistore/internal/domain"
)

type Service interface {
	CreateProduct(ctx context.Context, Title, Description string, Price uint, Category string) (*domain.Product, error)
	CreateCard(ctx context.Context, userID string) (*domain.Card, error)
}
