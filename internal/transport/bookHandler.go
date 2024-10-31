package transport

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"matask/internal/service"
	"matask/internal/transport/handler"
	"matask/internal/transport/request"
	"matask/internal/transport/resource"
	"net/http"
	"os"
	"strconv"
)

func GetBookHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(r.PathValue("id"))
		userId := r.Context().Value(handler.UserIdKey).(int)
		b, err := service.FindBook(id, userId, db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(resource.FromBook(b))
	}
}

func GetBookCoverHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(r.PathValue("id"))
		userId := r.Context().Value(handler.UserIdKey).(int)
		coverPath, err := service.FindBookCoverPath(id, userId, db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fileBytes, err := os.ReadFile(coverPath)
		if err != nil {
			http.Error(w, "Image not found.", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		w.Write(fileBytes)
	}
}

func CreateBookHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var b request.BookRequest
		err := json.NewDecoder(r.Body).Decode(&b)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		userId := r.Context().Value(handler.UserIdKey).(int)
		bookId, err := service.SaveOrUpdateBook(b.ToBook(), userId, db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
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
			slog.Error(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		book := b.ToBook()
		book.Id = id
		userId := r.Context().Value(handler.UserIdKey).(int)
		_, err = service.SaveOrUpdateBook(book, userId, db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "")
	}
}

func DeleteBookHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(r.PathValue("id"))
		userId := r.Context().Value(handler.UserIdKey).(int)
		err := service.DeleteBook(id, userId, db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "")
	}
}
