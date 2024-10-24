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

func GetMovieHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(r.PathValue("id"))
		m := service.FindMovie(id, db)
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

		movieId := service.CreateMovie(m.ToMovie(), db)
		json.NewEncoder(w).Encode(resource.IdResource{Id: movieId})
	}
}

func UpdateMovieHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "UpdateMovieHandler")
	}
}
