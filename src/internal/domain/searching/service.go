package searching

import (
	"context"
	"redistore/internal/domain"
)

type Service interface {
	SearchProductsByTitle(ctx context.Context, titleKeywords string) ([]domain.Product, error)
}
