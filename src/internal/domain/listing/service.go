package listing

import (
	"context"
	"redistore/internal/domain"
)

type Service interface {
	GetProductList(ctx context.Context) ([]domain.Product, error)
}
