package postgres

import (
	"gorm.io/gorm"
	"redistore/internal/domain"
)

type Product struct {
	gorm.Model
	Title       string `gorm:"size:128;column:title;index:title"`
	Description string `gorm:"size:256;column:description"`
	Price       uint   `gorm:"column:price"`
	Category    string `gorm:"size:256;column:category"`
}

func NewRepoProduct(product domain.Product) *Product {
	return &Product{
		Title:       product.Title,
		Description: product.Description,
		Price:       product.Price,
		Category:    string(product.Category),
	}
}

func NewDomainProduct(p Product) domain.Product {
	return domain.Product{
		ID:          p.Model.ID,
		Title:       p.Title,
		Description: p.Description,
		Price:       p.Price,
		Category:    domain.Category(p.Category),
		CreatedAt:   p.CreatedAt.Unix(),
		UpdatedAt:   p.UpdatedAt.Unix(),
	}
}
