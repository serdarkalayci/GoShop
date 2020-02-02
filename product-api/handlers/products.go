package handlers

import (
	"context"
	"goshop/product-api/data"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	err := lp.ToJson(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
		return
	}
}

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	product := r.Context().Value(KeyProduct{}).(data.Product)
	err := product.FromJson(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
		return
	}
	p.l.Printf("Product: %#v", product)
	data.AddProduct(&product)
}

func (p *Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Invalid id for product", http.StatusBadGateway)
		return
	}
	product := r.Context().Value(KeyProduct{}).(data.Product)
	err = data.UpdateProduct(id, &product)
	if err != nil {
		http.Error(rw, "Unable to update product json", http.StatusInternalServerError)
		return
	}
	p.l.Printf("Product: %#v", product)
}

type KeyProduct struct{}

func (p Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// validate the product json
		product := &data.Product{}
		err := product.FromJson(r.Body)
		if err != nil {
			http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
			return
		}
		// add the product to the context
		ctx := context.WithValue(r.Context(), KeyProduct{}, product)
		req := r.WithContext(ctx)
		// call the next handler in chain
		next.ServeHTTP(rw, req)
	})
}
