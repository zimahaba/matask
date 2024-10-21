package transport

import (
	"fmt"
	"net/http"
)

func GetMovieHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "GetMovieHandler")
	}
}

func GetMoviesPaginatedHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "GetMoviesPaginatedHandler")
	}
}

func CreateMovieHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "CreateMovieHandler")
	}
}

func UpdateMovieHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "UpdateMovieHandler")
	}
}
