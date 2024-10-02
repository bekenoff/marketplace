package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
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
		cookie, err := r.Cookie("token")
		if err != nil {
			fmt.Println("Error getting cookie:", err) // Добавьте отладочный вывод
			http.Error(w, "Unauthorized1", http.StatusUnauthorized)
			return
		}

		tokenStr := cookie.Value
		if tokenStr == "" {
			http.Error(w, "Unauthorized2", http.StatusUnauthorized)
			return
		}

		// Удаляем префикс "Bearer " из токена, если он есть
		if strings.HasPrefix(tokenStr, "Bearer ") {
			tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
		}

		claims := &Claims{}

		// Парсим токен и извлекаем claims
		tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtKey, nil
		})

		if err != nil || !tkn.Valid {
			http.Error(w, "Unauthorized3", http.StatusUnauthorized)
			return
		}

		// Устанавливаем claims в контекст запроса для последующего использования
		ctx := context.WithValue(r.Context(), "claims", claims)
		r = r.WithContext(ctx)

		// Продолжаем обработку запроса
		next.ServeHTTP(w, r)
	})
}
