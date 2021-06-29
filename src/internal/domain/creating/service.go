package creating

import (
	"context"
	"errors"
	"redistore/internal/domain"
	"redistore/internal/domain/ports"
	"redistore/pkg/yerror"
)

type Service interface {
	CreateProduct(ctx context.Context, Title, Description string, Price uint, Category string) (*domain.Product, error)
	CreateCard(ctx context.Context, userID string) (*domain.Card, error)
}

func New(repo ports.Repository) Service {
	return service{
		repo: repo,
	}
}

type service struct {
	repo ports.Repository
}

func (s service) CreateProduct(ctx context.Context, Title, Description string, Price uint, Category string) (*domain.Product, error) {
	const op yerror.Op = "domain.creating.service.CreateProduct"

	if Title == "" {
		return nil, yerror.E(op, yerror.KindInvalidArgument, errors.New("the Title is empty"))
	}
	if Description == "" {
		return nil, yerror.E(op, yerror.KindInvalidArgument, errors.New("the Description is empty"))
	}
	if Price == 0 {
		return nil, yerror.E(op, yerror.KindInvalidArgument, errors.New("the Price is empty"))
	}
	if Category == "" {
		return nil, yerror.E(op, yerror.KindInvalidArgument, errors.New("the Category is empty"))
	}
	product := domain.Product{
		Title:       Title,
		Description: Description,
		Price:       Price,
		Category:    domain.Category(Category),
	}
	createdProduct, err := s.repo.InsertProduct(ctx, product)

	if err != nil {
		return nil, yerror.E(op, err)
	}
	return createdProduct, nil
}

func (s service) CreateCard(ctx context.Context, userID string) (*domain.Card, error) {
	const op yerror.Op = "domain.creating.service.CreateCard"

	if userID == "" {
		return nil, yerror.E(op, yerror.KindInvalidArgument, errors.New("the userID is empty"))
	}
	card := domain.Card{
		UserID: userID,
	}
	createdCard, err := s.repo.InsertCard(ctx, card)

	if err != nil {
		return nil, yerror.E(op, err)
	}
	return createdCard, nil
}
