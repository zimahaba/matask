package transport

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"matask/internal/service"
	"matask/internal/transport/handler"
	"matask/internal/transport/request"
	"matask/internal/transport/resource"
	"net/http"
	"os"
	"strconv"
)

func GetFilteredBooksHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	filter := request.ToBookFilter(r.URL.Query())
	filter.UserId = r.Context().Value(handler.UserIdKey).(int)
	result, err := service.FindFilteredBooks(filter, db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resource.FromBookPageResult(result))
}

func GetBookHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	userId := r.Context().Value(handler.UserIdKey).(int)
	b, err := service.FindBook(id, userId, db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resource.FromBook(b))
}

func GetBookCoverHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	userId := r.Context().Value(handler.UserIdKey).(int)
	coverPath, err := service.FindBookCoverPath(id, userId, db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fileBytes, err := os.ReadFile(coverPath)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Book cover not found.", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)
	w.Write(fileBytes)
}

func SaveBookHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	r.ParseMultipartForm(2 << 20)
	var filebytes []byte
	file, _, err := r.FormFile("cover")
	if err == nil {
		defer file.Close()
		filebytes, err = io.ReadAll(file)
		if err != nil {
			errStr := fmt.Sprintf("Error in reading the file buffer %s\n", err)
			slog.Error(errStr)
			http.Error(w, errStr, http.StatusInternalServerError)
			return
		}
	}
	fmt.Printf("form: %v.\n", r.Form)
	b, err := request.ToBook(r.Form)
	id, _ := strconv.Atoi(r.PathValue("id"))
	fmt.Printf("id: %v.\n", id)
	if id > 0 {
		b.Id = id
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	userId := r.Context().Value(handler.UserIdKey).(int)

	bookId, err := service.SaveOrUpdateBook(b, filebytes, userId, db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resource.IdResource{Id: bookId})
}

func DeleteBookHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	userId := r.Context().Value(handler.UserIdKey).(int)
	err := service.DeleteBook(id, userId, db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "")
}
