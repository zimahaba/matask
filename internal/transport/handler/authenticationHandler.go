package handler

import (
	"context"
	"fmt"
	"log"
	"matask/internal/transport"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func Auth(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("Authorization")
		fmt.Printf("token: %v.\v", token)

		claims := &transport.Claims{}

		tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return transport.JwtKey, nil
		})

		if err != nil || !tkn.Valid {
			log.Printf("error: %v.\n", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		newCtx := context.WithValue(r.Context(), userIdKey, 1)
		rWithId := r.WithContext(newCtx)
		next.ServeHTTP(w, rWithId)
	})
}
