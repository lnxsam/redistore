package postgres

import (
	"context"
	"errors"
	"redistore/internal/data"
	"redistore/internal/domain"
	"redistore/pkg/yerror"

	"gorm.io/gorm"
)

func NewDBDataSource(db *gorm.DB) data.DBDataSource {
	return &postgres{
		db: db,
	}
}

type postgres struct {
	db *gorm.DB
}

func (p *postgres) GetCardByID(ctx context.Context, id string) (*domain.Card, error) {
	const op yerror.Op = "postgres.GetCardByID"
	repoCard := new(Card)

	err := p.db.WithContext(ctx).Where("id = ?", id).First(&repoCard).Error
	if err != nil {
		return nil, yerror.E(op, errors.New("no card found"), yerror.LevelError, yerror.KindInternal)
	}

	card := NewDomainCard(*repoCard)
	return card, nil
}

func (p *postgres) GetProductList(ctx context.Context) ([]domain.Product, error) {
	const op yerror.Op = "postgres.GetWithdraws"
	var repoProductList []Product

	err := p.db.WithContext(ctx).
		Find(&repoProductList).Error
	if err != nil {
		return nil, yerror.E(op, errors.New("no withdraw found"), yerror.LevelError, yerror.KindInternal)
	}

	var domainProducts = make([]domain.Product, len(repoProductList))
	for i, repoProduct := range repoProductList {
		domainProducts[i] = NewDomainProduct(repoProduct)
	}
	return domainProducts, nil
}

func (p *postgres) SearchProductsByTitle(ctx context.Context, titleKeywords string) ([]domain.Product, error) {
	const op yerror.Op = "postgres.SearchProductsByTitle"
	var repoProductList []Product

	err := p.db.WithContext(ctx).
		Where("title LIKE ?", "%"+titleKeywords+"%").
		Find(&repoProductList).Error
	if err != nil {
		return nil, yerror.E(op, errors.New("no withdraw found"), yerror.LevelError, yerror.KindInternal)
	}

	var domainProducts = make([]domain.Product, len(repoProductList))
	for i, repoProduct := range repoProductList {
		domainProducts[i] = NewDomainProduct(repoProduct)
	}
	return domainProducts, nil
}

func (p *postgres) InsertCard(ctx context.Context, domainCard domain.Card) (*domain.Card, error) {
	const op yerror.Op = "postgres.InsertCard"

	repoCard := NewRepoCard(domainCard)

	err := p.db.WithContext(ctx).Create(&repoCard).Error
	if err != nil {
		return nil, yerror.E(op, err, yerror.KindInternal, yerror.LevelError)
	}

	return NewDomainCard(*repoCard), nil
}

func (p *postgres) AutoMigrate() error {
	const op yerror.Op = "data_sources.AutoMigrate"

	err := p.db.AutoMigrate(&Product{}, Card{})
	if err != nil {
		panic("initialize db failed")
	}

	return nil
}

func (p *postgres) InsertProduct(ctx context.Context, domainProduct domain.Product) (*domain.Product, error) {
	const op yerror.Op = "postgres.InsertProduct"

	repoProduct := NewRepoProduct(domainProduct)

	err := p.db.WithContext(ctx).Create(&repoProduct).Error
	if err != nil {
		return nil, yerror.E(op, err, yerror.KindInternal, yerror.LevelError)
	}

	product := NewDomainProduct(*repoProduct)
	return &product, nil
}

func (p *postgres) UpdateCard(ctx context.Context, domainCard domain.Card) error {
	const op yerror.Op = "postgres.UpdateCard"

	repoCard := NewRepoCard(domainCard)

	err := p.db.WithContext(ctx).Save(&repoCard).Error
	if err != nil {
		return yerror.E(op, err, yerror.KindInternal, yerror.LevelError)
	}

	return nil
}

func (p *postgres) GetProductByID(ctx context.Context, id string) (*domain.Product, error) {
	const op yerror.Op = "postgres.GetProductByID"
	repoProduct := new(Product)

	err := p.db.WithContext(ctx).Where("id = ?", id).First(&repoProduct).Error
	if err != nil {
		return nil, yerror.E(op, errors.New("no product found"), yerror.LevelError, yerror.KindInternal)
	}

	product := NewDomainProduct(*repoProduct)
	return &product, nil
}
