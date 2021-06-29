package domain

type Category string

const (
	Car         Category = "Car"
	Electricity Category = "Electricity"
)

type Product struct {
	ID          uint
	Title       string
	Description string
	Price       uint
	Category    Category
	CreatedAt   int64
	UpdatedAt   int64
}

func NewProduct(title, description string, price uint, category Category) *Product {
	return &Product{
		Title:       title,
		Description: description,
		Price:       price,
		Category:    category,
	}
}
