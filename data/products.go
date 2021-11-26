package data

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

var (
	ErrProductNotFound = fmt.Errorf("product not found")
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

func (p *Product) FromJSON(r io.Reader) error {
	return json.NewDecoder(r).Decode(p)
}

// GetProducts returns the list of products
func GetProducts() Products {
	return productList
}

func AddProduct(p *Product) {
	p.ID = getNextID()
	productList = append(productList, p)
}

func getNextID() int {
	p := productList[len(productList)-1]
	return p.ID + 1
}

func UpdateProduct(id int, data *Product) error {
	_, pos, err := findProduct(id)
	if err != nil {
		return err
	}

	data.ID = id
	productList[pos] = data
	return nil
}

func findProduct(id int) (*Product, int, error) {
	for i, p := range productList {
		if p.ID == id {
			return p, i, nil
		}
	}
	return nil, -1, ErrProductNotFound
}

// Dummy data for the product
var productList = []*Product{
	&Product{
		ID:          1,
		Name:        "product1",
		Description: "product 1 desc",
		SKU:         "p123",
		CreatedAt:   time.Now().UTC().String(),
		UpdatedAt:   time.Now().UTC().String(),
		DeletedAt:   time.Now().UTC().String(),
	},
	&Product{
		ID:          2,
		Name:        "product2",
		Description: "product 2 desc",
		SKU:         "p234",
		CreatedAt:   time.Now().UTC().String(),
		UpdatedAt:   time.Now().UTC().String(),
		DeletedAt:   time.Now().UTC().String(),
	},
}
