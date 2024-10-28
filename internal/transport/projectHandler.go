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
		userId := r.Context().Value(handler.UserIdKey).(int)
		projectId := service.SaveOrUpdateProject(p.ToProject(), userId, db)
		json.NewEncoder(w).Encode(resource.IdResource{Id: projectId})
	}
}

func UpdateProjectHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(r.PathValue("id"))
		var p request.ProjectRequest
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		project := p.ToProject()
		project.Id = id
		userId := r.Context().Value(handler.UserIdKey).(int)
		service.SaveOrUpdateProject(project, userId, db)
		fmt.Fprintf(w, "")
	}
}

func DeleteProjectHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(r.PathValue("id"))
		service.DeleteProject(id, db)
		fmt.Fprintf(w, "")
	}
}
