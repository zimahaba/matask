package transport

import (
	"fmt"
	"net/http"
)

func GetBookHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "GetBookHandler")
	}
}

func GetBooksPaginatedHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "GetBooksPaginatedHandler")
	}
}

func CreateBookHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "CreateBookHandler")
	}
}

func UpdateBookHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "UpdateBookHandler")
	}
}
