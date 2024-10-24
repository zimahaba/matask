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

func GetProjectHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(r.PathValue("id"))
		p := service.FindProject(id, db)
		json.NewEncoder(w).Encode(resource.FromProject(p))
	}
}

func CreateProjectHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var p request.ProjectRequest
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		projectId := service.CreateProject(p.ToProject(), db)
		json.NewEncoder(w).Encode(resource.IdResource{Id: projectId})
	}
}

func UpdateProjectHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "UpdateProjectHandler")
	}
}
