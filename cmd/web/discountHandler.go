package main

import (
	"bytes"
	"encoding/json"
	"io"
	"marketplace/pkg/models"
	"net/http"
)

func (app *application) addDiscount(w http.ResponseWriter, r *http.Request) {
	var discount models.Discount

	body, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	err := json.NewDecoder(r.Body).Decode(&discount)

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	err = app.discount.Insert(&discount)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated) // 201
}
