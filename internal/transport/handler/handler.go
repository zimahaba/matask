package handler

import (
	"database/sql"
	"net/http"

	"github.com/redis/go-redis/v9"
)

type MataskHandler struct {
	DB *sql.DB
	F  func(w http.ResponseWriter, r *http.Request, db *sql.DB)
}

type MataskTTLHandler struct {
	DB    *sql.DB
	Redis *redis.Client
	F     func(w http.ResponseWriter, r *http.Request, db *sql.DB, redis *redis.Client)
}

func (h MataskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.F(w, r, h.DB)
}

func (h MataskTTLHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.F(w, r, h.DB, h.Redis)
}
