package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func Logging(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId := uuid.New()
		newCtx := context.WithValue(r.Context(), requestIdKey, requestId)
		rWithId := r.WithContext(newCtx)
		log.Printf("Received request: %s %s - requestId=%v ", r.Method, r.URL.Path, requestId)
		next.ServeHTTP(w, rWithId)
	})
}
