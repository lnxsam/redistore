package domain

import "strconv"

type Card struct {
	ID        uint
	UserID    string
	CardItems map[string]*CardItem
	Price     uint
	CreatedAt int64
	UpdatedAt int64
}

type CardItem struct {
	Product *Product
	Count   uint
}

func NewCardItem(Count uint, product *Product) *CardItem {
	return &CardItem{Count: Count, Product: product}
}

func (c *Card) AddProduct(product *Product, count uint) {
	if c.CardItems == nil {
		c.CardItems = make(map[string]*CardItem)
	}
	if cardItem, ok := c.CardItems[strconv.FormatUint(uint64(product.ID), 10)]; ok {
		cardItem.Count += count
	} else {
		c.CardItems[strconv.FormatUint(uint64(product.ID), 10)] = &CardItem{Product: product, Count: count}
	}
	c.Price += count * product.Price
}

func (c *Card) RemoveCardItem(id string) {
	if cardItem, ok := c.CardItems[id]; ok {
		c.Price -= cardItem.Count * cardItem.Product.Price
		delete(c.CardItems, id)
	}
}
