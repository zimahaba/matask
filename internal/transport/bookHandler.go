package transport

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"matask/internal/service"
	"matask/internal/transport/request"
	"matask/internal/transport/resource"
	"net/http"
	"strconv"
)

func GetBookHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(r.PathValue("id"))
		b := service.FindBook(id, db)
		json.NewEncoder(w).Encode(resource.FromBook(b))
	}
}

func CreateBookHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var b request.BookRequest
		err := json.NewDecoder(r.Body).Decode(&b)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		bookId := service.CreateBook(b.ToBook(), db)
		json.NewEncoder(w).Encode(resource.IdResource{Id: bookId})
	}
}

func UpdateBookHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "UpdateBookHandler")
	}
}
