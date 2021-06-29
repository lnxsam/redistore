package rest

type ProductCreateDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       uint   `json:"price"`
	Category    string `json:"category"`
}

type CardCreateDTO struct {
	UserID string `json:"user_id"`
}

type AddProductToCardDTO struct {
	CardID    string `json:"card_id"`
	ProductID string `json:"product_id"`
	Count     uint   `json:"count"`
}

type RemoveCardItemDTO struct {
	CardID    string `json:"card_id"`
	ProductID string `json:"product_id"`
}

type SearchProductDTO struct {
	Title string `json:"title"`
}
