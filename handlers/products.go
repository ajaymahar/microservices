package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/ajaymahar/microservices/data"
	"github.com/gorilla/mux"
)

// Products is product handlers with different methods
type Products struct {
	l *log.Logger
}

// ProductKey is key for data.Product
type ProductKey struct{}

// Factory function to create new Product
func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// GetProducts method to accepet Product GET request
//
func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {

	p.l.Println("Handle GET Products")
	//fetch the data from datastore
	pl := data.GetProducts()

	// Encode the list of products into JSON format
	if err := pl.ToJSON(rw); err != nil {
		http.Error(rw, "unable to encode the json", http.StatusInternalServerError)
		return
	}
}

// CreateProduct method will create new data.Product with POST request.
func (p *Products) CreateProduct(rw http.ResponseWriter, r *http.Request) {

	p.l.Println("Method Post to create product")

	// fetch the product from the request context, added by MiddlewareProductValidator
	prod, ok := r.Context().Value(ProductKey{}).(data.Product)
	if !ok {

		http.Error(rw, "type assersion failed for the product", http.StatusInternalServerError)
		return
	}

	// Add product to the datastore
	data.AddProduct(&prod)

}

// UpdateProduct method will update the Product with PUT request
func (p *Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {

	p.l.Println("Method PUT to update the product")

	// Fetch the id from request URL
	vars := mux.Vars(r)

	// convert string id to int type, to pass it to data.UpdateProduct method
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		p.l.Printf("bad request id should be number got <%v>", vars["id"])
		http.Error(rw, "bad request id must be a number", http.StatusBadRequest)
		return
	}

	// fetch the product from the request context, added by MiddlewareProductValidator
	prod, ok := r.Context().Value(ProductKey{}).(data.Product)
	if !ok {

		http.Error(rw, "type assersion failed for the product", http.StatusInternalServerError)
		return
	}

	// update the product from the datastore
	if err := data.UpdateProduct(id, &prod); err != nil {
		p.l.Println(err.Error())
		http.Error(rw, "not found", http.StatusNotFound)
		return
	}
}

// MiddlewareProductValidator is middleware which run before running the actual handlerFunc
func (p *Products) MiddlewareProductValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		var prod data.Product
		if err := prod.FromJSON(r.Body); err != nil {
			defer r.Body.Close()
			p.l.Println("unable to parse the json")
			http.Error(rw, "unable to parse the json", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), ProductKey{}, prod)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)
	})
}
