package handlers

import (
	"log"
	"net/http"

	"github.com/ajaymahar/microservices/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	// handle the get request for the list of products
	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	}

	// catch all and return error for other request
	p.l.Println(r.Method, "method not allowed")
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {

	p.l.Println("Handle GET Products")
	//fetch the data from datastore
	pl := data.GetProducts()

	// serialize the list of products
	if err := pl.ToJSON(rw); err != nil {
		http.Error(rw, "unable to encode the json", http.StatusInternalServerError)
		return
	}
}
