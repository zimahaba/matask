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

func GetMovieHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(r.PathValue("id"))
		m, err := service.FindMovie(id, db)
		if err != nil {
			// log
			// error response
		}
		json.NewEncoder(w).Encode(resource.FromMovie(m))
	}
}

func CreateMovieHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var m request.MovieRequest
		err := json.NewDecoder(r.Body).Decode(&m)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		userId := r.Context().Value(handler.UserIdKey).(int)
		movieId, err := service.SaveOrUpdateMovie(m.ToMovie(), userId, db)
		if err != nil {
			// log
			// error response
		}
		json.NewEncoder(w).Encode(resource.IdResource{Id: movieId})
	}
}

func UpdateMovieHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(r.PathValue("id"))
		var m request.MovieRequest
		err := json.NewDecoder(r.Body).Decode(&m)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		movie := m.ToMovie()
		movie.Id = id
		userId := r.Context().Value(handler.UserIdKey).(int)
		service.SaveOrUpdateMovie(movie, userId, db)
		fmt.Fprintf(w, "")
	}
}

func DeleteMovieHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(r.PathValue("id"))
		service.DeleteMovie(id, db)
		fmt.Fprintf(w, "")
	}
}
