package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"marketplace/pkg/models"
	"net/http"
)

func (app *application) createSight(w http.ResponseWriter, r *http.Request) {
	var newSight models.Sight

	body, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	err := json.NewDecoder(r.Body).Decode(&newSight)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	err = app.sight.Insert(&newSight)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (app *application) getSight(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get(":id")

	sightData, err := app.sight.GetSightById(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.clientError(w, http.StatusNotFound)
		} else {
			app.serverError(w, err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(sightData)
}

func (app *application) getAllSights(w http.ResponseWriter, r *http.Request) {
	sightData, err := app.sight.GetAllSights()
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(sightData)
}
