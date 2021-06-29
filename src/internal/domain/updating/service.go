package updating

import (
	"context"
)

type Service interface {
	AddProductToCard(ctx context.Context, cardID, productID string, count uint) error
	RemoveProductFromCard(ctx context.Context, cardID, productID string) error
}