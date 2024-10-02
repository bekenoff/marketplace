package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"marketplace/pkg/models"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sfreiberg/gotwilio"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("your_secret_key")

// Claims struct for storing token data
type Claims struct {
	Telephone int  `json:"telephone"`
	UserID    int  `json:"user_id"`
	IsRefresh bool `json:"is_refresh"`
	jwt.RegisteredClaims
}

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
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// Если токен отсутствует, возвращаем ошибку 401
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		app.serverError(w, err)
		return
	}

	tokenStr := cookie.Value

	claims := &Claims{}

	// Проверка токена
	tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		app.serverError(w, err)
		return
	}

	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Используем userID из токена
	user, err := app.client.GetUserById(strconv.Itoa(claims.UserID))
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
	var credentials struct {
		Telephone int    `json:"telephone"`
		Password  string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		app.clientError(w, http.StatusBadRequest) // 400 Bad Request
		return
	}
	defer r.Body.Close()

	storedPassword, err := app.client.GetPasswordByTelephone(credentials.Telephone)
	if err != nil {
		app.clientError(w, http.StatusUnauthorized) // 401 Unauthorized
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(credentials.Password)); err != nil {
		app.clientError(w, http.StatusUnauthorized) // 401 Unauthorized
		return
	}

	accessToken, refreshToken, err := app.createTokens(credentials.Telephone)
	if err != nil {
		app.serverError(w, err)
		return
	}

	response := map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"telephone":     credentials.Telephone,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (app *application) createTokens(telephone int) (string, string, error) {
	accessClaims := &Claims{
		Telephone: telephone,
		UserID:    telephone,
		IsRefresh: false,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}
	accessToken, err := app.signToken(accessClaims)
	if err != nil {
		return "", "", err
	}

	refreshClaims := &Claims{
		Telephone: telephone,
		UserID:    telephone,
		IsRefresh: true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		},
	}
	refreshToken, err := app.signToken(refreshClaims)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil

}

func NewRefreshToken() (string, error) {
	b := make([]byte, 32)

	// Используем rand.Read для генерации случайных байтов
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}

func (app *application) CreateSession(ctx context.Context, user models.Client, accessToken string) (models.Tokens, error) {
	var (
		res models.Tokens
		err error
	)

	userIDStr := strconv.Itoa(user.Id)

	res.AccessToken = accessToken

	// Генерируем только RefreshToken
	res.RefreshToken, err = NewRefreshToken()
	if err != nil {
		return res, err
	}

	// Создание и сохранение сессии с RefreshToken
	session := models.Session{
		RefreshToken: res.RefreshToken,
		ExpiresAt:    time.Now().Add(1 * time.Hour),
	}

	err = app.client.SetSession(ctx, userIDStr, session)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (app *application) signToken(claims *Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
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

func (app *application) refreshToken(w http.ResponseWriter, r *http.Request) {
	var refreshToken struct {
		Token string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&refreshToken); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(refreshToken.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !tkn.Valid || !claims.IsRefresh {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	newAccessToken, err := app.createAccessToken(claims)
	if err != nil {
		app.serverError(w, err)
		return
	}

	response := map[string]interface{}{
		"access_token": newAccessToken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (app *application) createAccessToken(claims *Claims) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute)
	newClaims := &Claims{
		Telephone: claims.Telephone,
		UserID:    claims.UserID,
		IsRefresh: false,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	newAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	return newAccessToken.SignedString(jwtKey)
}
