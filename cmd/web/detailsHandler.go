package main

import (
	"bytes"
	"encoding/json"
	"io"
	"marketplace/pkg/models"
	"net/http"
	"strconv"
)

func (app *application) addInformation(w http.ResponseWriter, r *http.Request) {
	var newInformation models.Information

	body, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	err := json.NewDecoder(r.Body).Decode(&newInformation)

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	err = app.details.Insert(&newInformation)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated) // 201
}

func (app *application) getInformation(w http.ResponseWriter, r *http.Request) {
	informationIDString := r.URL.Query().Get("id")
	if informationIDString == "" {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	informationID, err := strconv.Atoi(informationIDString)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	information, err := app.details.GetInformation(informationID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	response, err := json.Marshal(information)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
