package transport

import (
	"encoding/json"
	"fmt"
	"matask/internal/models"
	"net/http"
)

func GetProjectHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "GetProjectHandler")
	}
}

func GetProjectsPaginatedHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintf(w, "GetProjectsPaginatedHandler")
	}
}

func CreateProjectHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var project models.Project

		err := json.NewDecoder(r.Body).Decode(&project)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Printf("project: %v", project)

		fmt.Fprintf(w, "CreateProjectHandler")
	}
}

func UpdateProjectHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "UpdateProjectHandler")
	}
}
