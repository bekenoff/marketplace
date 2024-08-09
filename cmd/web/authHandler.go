package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"marketplace/pkg/models"
	"net/http"
	"strconv"

	"github.com/sfreiberg/gotwilio"
)

func (app *application) signupClient(w http.ResponseWriter, r *http.Request) {
	var newClient models.Client

	body, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	err := json.NewDecoder(r.Body).Decode(&newClient)

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	err = app.client.Insert(newClient.Telephone, newClient.Password)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated) // 201
}

func (app *application) getUserById(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("id")
	if userIDStr == "" {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	user, err := app.client.GetUserById(userIDStr)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.clientError(w, http.StatusNotFound)
		} else {
			app.serverError(w, err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(user)
}

func (app *application) signupClientLaw(w http.ResponseWriter, r *http.Request) {
	var newClient models.ClientLaw

	body, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	err := json.NewDecoder(r.Body).Decode(&newClient)

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	err = app.client.InsertLaw(&newClient)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated) // 201
}

func (app *application) loginClient(w http.ResponseWriter, r *http.Request) {
	var client models.Client

	body, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	err := json.NewDecoder(r.Body).Decode(&client)

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	clientId, err := app.client.Authenticate(client.Password, client.Telephone)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			app.clientError(w, http.StatusBadRequest)
			return
		} else {
			app.serverError(w, err)

			return
		}
	}

	responseUser, err := app.client.GetUserById(strconv.Itoa(clientId))

	_, err = w.Write(responseUser)
	if err != nil {
		return
	}
}

func (app *application) loginAdmin(w http.ResponseWriter, r *http.Request) {
	var client models.Client

	body, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	err := json.NewDecoder(r.Body).Decode(&client)

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	clientId, err := app.client.AuthenticateAdmin(client.Telephone, client.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			app.clientError(w, http.StatusBadRequest)
			return
		} else {
			app.serverError(w, err)

			return
		}
	}

	responseUser, err := app.client.GetUserByIdAdmin(strconv.Itoa(clientId))

	_, err = w.Write(responseUser)
	if err != nil {
		return
	}
}

// RECOVERY

// Recovery

func (app *application) Recoverybysms(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("idclient")
	if err != nil {
		http.Error(w, "Ошибка получения cookie", http.StatusBadRequest)
		return
	}

	idclient := cookie.Value

	clientphone, err := app.client.GetClientPhoneById(idclient)

	fmt.Print(clientphone)

	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}

	link := fmt.Sprintf("http://localhost:4000/client-password-recovery?id=%s", idclient)

	sendSMS(clientphone, link)

}

func sendSMS(recipient string, message string) error {
	accountSid := "AC17c5b66f4964850573f2ea5a06a4aa9e"
	authToken := "2084ef8187bf3aebb4d5ad92f7a80708"
	twilio := gotwilio.NewTwilioClient(accountSid, authToken)

	from := "+14692405277" // Номер Twilio, с которого будет отправлено SMS
	to := recipient
	body := message

	_, _, err := twilio.SendSMS(from, to, body, "", "")

	return err
}

func (app *application) updatePassword(w http.ResponseWriter, r *http.Request) {

	id_string := r.URL.Query().Get("id")

	id, err := strconv.Atoi(id_string)

	if err != nil {
		app.clientError(w, http.StatusInternalServerError)
		return
	}

	type pass struct {
		OldPassword string `json:"oldpassword"`
		NewPassword string `json:"newpassword"`
	}
	var clientpass pass

	body, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	err = json.NewDecoder(r.Body).Decode(&clientpass)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	err = app.client.ChangePassword(id, clientpass.OldPassword, clientpass.NewPassword)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK) // 200
}
