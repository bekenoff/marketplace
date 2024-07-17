package main

import (
	"bytes"
	"encoding/json"
	"io"
	"marketplace/pkg/models"
	"net/http"
	"strconv"
)

func (app *application) addFavorite(w http.ResponseWriter, r *http.Request) {
	var favorite models.Favorites

	body, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	err := json.NewDecoder(r.Body).Decode(&favorite)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	err = app.fav.Insert(&favorite)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (app *application) getFavorites(w http.ResponseWriter, r *http.Request) {
	clientIDStr := r.URL.Query().Get("client_id")
	if clientIDStr == "" {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	clientID, err := strconv.Atoi(clientIDStr)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	favorites, err := app.fav.GetByClientID(clientID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	response, err := json.Marshal(favorites)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
