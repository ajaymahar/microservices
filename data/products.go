package data

import (
	"encoding/json"
	"io"
	"time"
)

// Product represent single Product
type Product struct {
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	SKU         string `json:"sku"`
	CreatedAt   string `json:"-"`
	UpdatedAt   string `json:"-"`
	DeletedAt   string `json:"-"`
}

// Products is list of Product
type Products []*Product

func (p *Products) ToJSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(p)
}

// GetProducts returns the list of products
func GetProducts() Products {
	return productList
}

// Dummy data for the product
var productList = []*Product{
	&Product{
		ID:          123,
		Name:        "product1",
		Description: "product 1 desc",
		SKU:         "p123",
		CreatedAt:   time.Now().UTC().String(),
		UpdatedAt:   time.Now().UTC().String(),
		DeletedAt:   time.Now().UTC().String(),
	},
	&Product{
		ID:          234,
		Name:        "product2",
		Description: "product 2 desc",
		SKU:         "p234",
		CreatedAt:   time.Now().UTC().String(),
		UpdatedAt:   time.Now().UTC().String(),
		DeletedAt:   time.Now().UTC().String(),
	},
}
