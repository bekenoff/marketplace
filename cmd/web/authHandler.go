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
	"strings"
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
	// Получение токена из заголовков запроса
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Удаляем префикс "Bearer " из токена, если он есть
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

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
	json.NewEncoder(w).Encode(user)
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

func (app *application) loginClient(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Telephone int    `json:"telephone"`
		Password  string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Ошибка при декодировании запроса", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Получение хэшированного пароля из базы данных
	storedPassword, err := app.client.GetPasswordByTelephone(credentials.Telephone)
	if err != nil {
		http.Error(w, "Неверные учетные данные", http.StatusUnauthorized)
		return
	}

	// Сравнение паролей
	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(credentials.Password)); err != nil {
		http.Error(w, "Неверные учетные данные", http.StatusUnauthorized)
		return
	}

	accessToken, refreshToken, err := app.createTokens(credentials.Telephone)
	if err != nil {
		http.Error(w, "Ошибка при создании токенов", http.StatusInternalServerError)
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

// Создание токенов
func (app *application) createTokens(telephone int) (string, string, error) {
	// Создание access токена
	accessClaims := &Claims{
		Telephone: telephone,
		UserID:    telephone, // Предположим, что ID клиента равен телефону
		IsRefresh: false,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)), // Действителен 15 минут
		},
	}
	accessToken, err := app.signToken(accessClaims)
	if err != nil {
		return "", "", err
	}

	// Создание refresh токена
	refreshClaims := &Claims{
		Telephone: telephone,
		UserID:    telephone,
		IsRefresh: true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // Действителен 7 дней
		},
	}
	refreshToken, err := app.signToken(refreshClaims)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// Подпись токена
func (app *application) signToken(claims *Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// Обновление токенов с использованием refresh токена
func (app *application) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var refreshToken struct {
		Token string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&refreshToken); err != nil {
		http.Error(w, "Ошибка при декодировании запроса", http.StatusBadRequest)
		return
	}

	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(refreshToken.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !tkn.Valid || !claims.IsRefresh {
		http.Error(w, "Неверный или истекший refresh токен", http.StatusUnauthorized)
		return
	}

	// Создание новых токенов
	newAccessToken, newRefreshToken, err := app.createTokens(claims.Telephone)
	if err != nil {
		http.Error(w, "Ошибка при создании новых токенов", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
