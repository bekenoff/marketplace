package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"marketplace/pkg/models"
	"net/http"
)

func (app *application) createEvent(w http.ResponseWriter, r *http.Request) {
	var newEvent models.Events

	body, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	err := json.NewDecoder(r.Body).Decode(&newEvent)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	err = app.event.Insert(&newEvent)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (app *application) getEvent(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get(":id")

	eventData, err := app.event.GetEventById(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.clientError(w, http.StatusNotFound)
		} else {
			app.serverError(w, err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(eventData)
}

func (app *application) GetAllEvents(w http.ResponseWriter, r *http.Request) {
	eventData, err := app.event.GetAllEvents()
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(eventData)
}
