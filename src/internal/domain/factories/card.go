package factories

import (
	"math/rand"
	"redistore/internal/domain"
	"time"
)

// CardFactory is contract for Card factory
type CardFactory interface {
	Create() domain.Card
	CreateMany(int) []domain.Card
}

var Card CardFactory = &cardFactory{}

type cardFactory struct {
}

func (cf *cardFactory) Create() domain.Card {
	return domain.Card{
		ID:        uint(rand.Uint32()),
		Price:     0,
		CreatedAt: time.Now().UTC().Unix(),
	}
}

func (cf *cardFactory) CreateMany(count int) []domain.Card {
	var cards []domain.Card
	for i := 0; i < count; i++ {
		cards = append(cards, cf.Create())
	}
	return cards
}
