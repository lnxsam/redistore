package listing

import (
	"context"
	"redistore/internal/domain"
	"redistore/internal/domain/ports"
	"redistore/pkg/yerror"
)

type Service interface {
	GetProductList(ctx context.Context) ([]domain.Product, error)
}

func New(repo ports.Repository) Service {
	return service{
		repo: repo,
	}
}

type service struct {
	repo ports.Repository
}

func (s service) GetProductList(ctx context.Context) ([]domain.Product, error) {
	const op yerror.Op = "domain.listing.service.GetProductList"

	products, err := s.repo.GetProductList(ctx)

	if err != nil {
		return nil, yerror.E(op, err)
	}
	return products, nil
}
