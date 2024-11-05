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
	"strconv"
)

func GetFilteredMoviesHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	filter := request.ToMovieFilter(r.URL.Query())
	filter.UserId = r.Context().Value(handler.UserIdKey).(int)
	result, err := service.FindFilteredMovies(filter, db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resource.FromMoviePageResult(result))
}

func GetMovieHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	userId := r.Context().Value(handler.UserIdKey).(int)
	m, err := service.FindMovie(id, userId, db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resource.FromMovie(m))
}

func SaveMovieHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	r.ParseMultipartForm(2 << 20)
	var filebytes []byte
	file, _, err := r.FormFile("poster")
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

	m, err := request.ToMovie(r.Form)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	id, _ := strconv.Atoi(r.PathValue("id"))
	if id > 0 {
		m.Id = id
	}

	userId := r.Context().Value(handler.UserIdKey).(int)

	movieId, err := service.SaveOrUpdateMovie(m, filebytes, userId, db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resource.IdResource{Id: movieId})
}

func DeleteMovieHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	userId := r.Context().Value(handler.UserIdKey).(int)
	err := service.DeleteMovie(id, userId, db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "")
}
