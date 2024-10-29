package handler

import (
	"context"
	"database/sql"
	"fmt"
	"matask/internal/service"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func Auth(next http.Handler, db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := r.Cookie("token")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims := &service.Claims{}

		token, err := jwt.ParseWithClaims(tokenCookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
			return service.JwtKey, nil
		})

		if err != nil || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userId, err := service.FindUserId(claims.Username, db)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		fmt.Printf("going through: %v.\n", userId)
		newCtx := context.WithValue(r.Context(), UserIdKey, userId)
		rWithId := r.WithContext(newCtx)
		next.ServeHTTP(w, rWithId)
	})
}
