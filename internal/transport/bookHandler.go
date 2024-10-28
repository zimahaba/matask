package transport

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"matask/internal/service"
	"matask/internal/transport/handler"
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

		userId := r.Context().Value(handler.UserIdKey).(int)
		bookId := service.SaveOrUpdateBook(b.ToBook(), userId, db)
		json.NewEncoder(w).Encode(resource.IdResource{Id: bookId})
	}
}

func UpdateBookHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(r.PathValue("id"))
		var b request.BookRequest
		err := json.NewDecoder(r.Body).Decode(&b)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		book := b.ToBook()
		book.Id = id
		userId := r.Context().Value(handler.UserIdKey).(int)
		service.SaveOrUpdateBook(book, userId, db)
		fmt.Fprintf(w, "")
	}
}

func DeleteBookHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(r.PathValue("id"))
		service.DeleteBook(id, db)
		fmt.Fprintf(w, "")
	}
}
