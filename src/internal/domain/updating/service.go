package updating

import (
	"context"
	"errors"
	"redistore/internal/domain/ports"
	"redistore/pkg/yerror"
)

type Service interface {
	AddProductToCard(ctx context.Context, cardID, productID string, count uint) error
	RemoveProductFromCard(ctx context.Context, cardID, productID string) error
}

func New(repo ports.Repository) Service {
	return service{
		repo: repo,
	}
}

type service struct {
	repo ports.Repository
}

func (s service) AddProductToCard(ctx context.Context, cardID, productID string, count uint) error {
	const op yerror.Op = "domain.updating.service.AddProductToCard"

	if cardID == "" {
		return yerror.E(op, yerror.KindInvalidArgument, errors.New("the cardID is empty"))
	}

	if productID == "" {
		return yerror.E(op, yerror.KindInvalidArgument, errors.New("the productID is empty"))
	}

	if count == 0 {
		return yerror.E(op, yerror.KindInvalidArgument, errors.New("the count is invalid"))
	}

	card, err := s.repo.GetCardByID(ctx, cardID)
	if err != nil {
		return yerror.E(op, err)
	}
	product, err := s.repo.GetProductByID(ctx, productID)
	if err != nil {
		return yerror.E(op, err)
	}

	card.AddProduct(product, count)

	err = s.repo.UpdateCard(ctx, *card)
	if err != nil {
		return yerror.E(op, err)
	}
	return nil
}

func (s service) RemoveProductFromCard(ctx context.Context, cardID, productID string) error {
	const op yerror.Op = "domain.updating.service.RemoveProductFromCard"

	if cardID == "" {
		return yerror.E(op, yerror.KindInvalidArgument, errors.New("the cardID is empty"))
	}

	if productID == "" {
		return yerror.E(op, yerror.KindInvalidArgument, errors.New("the productID is empty"))
	}

	card, err := s.repo.GetCardByID(ctx, cardID)
	if err != nil {
		return yerror.E(op, err)
	}

	card.RemoveCardItem(productID)

	err = s.repo.UpdateCard(ctx, *card)
	if err != nil {
		return yerror.E(op, err)
	}
	return nil
}
