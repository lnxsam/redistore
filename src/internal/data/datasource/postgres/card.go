package postgres

import (
	"encoding/json"
	"gorm.io/gorm"
	"redistore/internal/domain"
)

type Card struct {
	gorm.Model
	CardItems string
	UserID    string
	Price     uint `gorm:"column:price"`
}

func NewRepoCard(card domain.Card) *Card {
	cardItemsString, _ := json.Marshal(card.CardItems)
	repoCard := &Card{
		UserID:    card.UserID,
		CardItems: string(cardItemsString),
		Price:     card.Price,
	}
	repoCard.Model.ID = card.ID
	return repoCard
}

func NewDomainCard(c Card) *domain.Card {
	cardItems := map[string]*domain.CardItem{}
	json.Unmarshal([]byte(c.CardItems), &cardItems)
	return &domain.Card{
		ID:        c.ID,
		UserID:    c.UserID,
		CardItems: cardItems,
		Price:     c.Price,
		CreatedAt: c.CreatedAt.Unix(),
		UpdatedAt: c.UpdatedAt.Unix(),
	}
}
