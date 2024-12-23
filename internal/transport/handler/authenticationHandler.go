package handler

import (
	"context"
	"log/slog"
	"matask/internal/service"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func Auth(next MataskHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := r.Cookie(service.TOKEN_COOKIE_NAME)
		if err != nil {
			slog.Error(err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims := &service.Claims{}

		token, err := jwt.ParseWithClaims(tokenCookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
			return service.JwtKey, nil
		})

		if err != nil || !token.Valid {
			slog.Error(err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userId, err := service.FindUserId(claims.Username, next.DB)
		if err != nil {
			slog.Error(err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		newCtx := context.WithValue(r.Context(), UserIdKey, userId)
		rWithId := r.WithContext(newCtx)
		next.ServeHTTP(w, rWithId)
	})
}
