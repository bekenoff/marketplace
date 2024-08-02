package main

import (
	"bytes"
	"encoding/json"
	"io"
	"marketplace/pkg/models"
	"net/http"
)

func (app *application) addImage(w http.ResponseWriter, r *http.Request) {

	var image models.Image

	body, _ := io.ReadAll(r.Body)

	r.Body = io.NopCloser(bytes.NewBuffer(body))

	err := json.NewDecoder(r.Body).Decode(&image)

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	err = app.image.Insert(&image)

	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)

}
