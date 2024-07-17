package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"marketplace/pkg/models"
	"net/http"
)

func (app *application) createRec(w http.ResponseWriter, r *http.Request) {
	var newRec models.Recommendation

	body, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	err := json.NewDecoder(r.Body).Decode(&newRec)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	err = app.rec.Insert(&newRec)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (app *application) getRec(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get(":id")

	recData, err := app.rec.GetRecommendationById(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.clientError(w, http.StatusNotFound)
		} else {
			app.serverError(w, err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(recData)
}

func (app *application) GetAllRec(w http.ResponseWriter, r *http.Request) {
	recData, err := app.rec.GetAllRecommendations()
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(recData)
}
