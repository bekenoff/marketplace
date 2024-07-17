package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) getProducts(w http.ResponseWriter, r *http.Request) {
	products, err := app.product.GetAllProducts()
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(products)
	if err != nil {
		app.serverError(w, err)
	}
}
