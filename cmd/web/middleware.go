package main

import (
	"fmt"
	"log"
	"marketplace/pkg/models"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")
		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *application) jwtAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		// Удаляем префикс "Bearer "
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Проверяем токен
		claims := &models.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Здесь можно проверить алгоритм
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Println("Unexpected signing method:", token.Header["alg"])
				return nil, http.ErrNoLocation
			}
			// Вернуть секретный ключ
			return []byte("your_secret_key"), nil // Замените на ваш секретный ключ
		})

		if err != nil {
			log.Printf("Error parsing token: %v", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Логируем успешную аутентификацию
		log.Printf("User authenticated: %v", claims.UserID)

		// Передаем управление следующему обработчику
		next.ServeHTTP(w, r)
	})
}
