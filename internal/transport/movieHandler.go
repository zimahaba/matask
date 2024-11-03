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
	"strconv"
)

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

func CreateMovieHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var m request.MovieRequest
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userId := r.Context().Value(handler.UserIdKey).(int)
	movieId, err := service.SaveOrUpdateMovie(m.ToMovie(), userId, db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resource.IdResource{Id: movieId})
}

func UpdateMovieHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	var m request.MovieRequest
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	movie := m.ToMovie()
	movie.Id = id
	userId := r.Context().Value(handler.UserIdKey).(int)
	_, err = service.SaveOrUpdateMovie(movie, userId, db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "")
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
