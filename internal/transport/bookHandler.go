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
		b, err := service.FindBook(id, db)
		if err != nil {
			// log
			// error response
		}
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
		bookId, err := service.SaveOrUpdateBook(b.ToBook(), userId, db)
		if err != nil {
			// log
			// error response
		}
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
		bookId, err := service.SaveOrUpdateBook(book, userId, db)
		if err != nil {
			fmt.Print(bookId)
			// log
			// error response
		}
		fmt.Fprintf(w, "")
	}
}

func DeleteBookHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(r.PathValue("id"))
		err := service.DeleteBook(id, db)
		if err != nil {
			// log
			// error response
		}
		fmt.Fprintf(w, "")
	}
}
