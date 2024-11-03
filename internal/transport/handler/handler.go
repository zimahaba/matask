package handler

import (
	"database/sql"
	"net/http"
)

type MataskHandler struct {
	DB *sql.DB
	F  func(w http.ResponseWriter, r *http.Request, db *sql.DB)
}

func (h MataskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.F(w, r, h.DB)
}
