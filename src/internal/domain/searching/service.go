package searching

import (
	"context"
	"errors"
	"redistore/internal/domain"
	"redistore/internal/domain/ports"
	"redistore/pkg/yerror"
)

type Service interface {
	SearchProductsByTitle(ctx context.Context, titleKeywords string) ([]domain.Product, error)
}

func New(repo ports.Repository) Service {
	return service{
		repo: repo,
	}
}

type service struct {
	repo ports.Repository
}

func (s service) SearchProductsByTitle(ctx context.Context, titleKeywords string) ([]domain.Product, error) {
	const op yerror.Op = "domain.creating.service.SearchProductsByTitle"

	if titleKeywords == "" {
		return nil, yerror.E(op, yerror.KindInvalidArgument, errors.New("the titleKeywords is empty"))
	}

	products, err := s.repo.SearchProductsByTitle(ctx, titleKeywords)

	if err != nil {
		return nil, yerror.E(op, err)
	}
	return products, nil
}
