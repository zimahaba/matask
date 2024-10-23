package transport

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"matask/internal/services"
	payload "matask/internal/transport/payloads"
	resource "matask/internal/transport/resources"
	"net/http"
)

func GetProjectHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		fmt.Printf("project id: %v", id)
		fmt.Fprintf(w, "GetProjectHandler")
	}
}

func GetProjectsPaginatedHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Print(r.URL.Query().Get("page"))
		fmt.Fprintf(w, "GetProjectsPaginatedHandler")
	}
}

func CreateProjectHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var p payload.ProjectPayload
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		project := payload.ToProject(p)
		projectId := services.CreateProject(project, db)
		json.NewEncoder(w).Encode(resource.IdResource{Id: projectId})
	}
}

func UpdateProjectHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "UpdateProjectHandler")
	}
}
