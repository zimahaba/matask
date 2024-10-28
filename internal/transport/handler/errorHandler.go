package handler

import (
	"database/sql"
	"log"
	"net/http"
)

type ErrorHandler struct {
	*sql.DB
	H func(db *sql.DB, w http.ResponseWriter, r *http.Request) error
}

// ServeHTTP allows our Handler type to satisfy http.Handler.
func (h ErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.H(h.DB, w, r)
	if err != nil {
		log.Printf("ops ERROR - requestId=%v", r.Context().Value(requestIdKey))
		/*switch e := err.(type) {
		case Error:
			// We can retrieve the status here and write out a specific
			// HTTP status code.
			log.Printf("HTTP %d - %s", e.Status(), e)
			http.Error(w, e.Error(), e.Status())
		default:
			// Any error types we don't specifically look out for default
			// to serving a HTTP 500
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}*/
		http.Error(w, "an error occurred", 500)
	}
}
