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
		userId := r.Context().Value(handler.UserIdKey).(int)
		p, err := service.FindProject(id, userId, db)
		if err != nil {
			// log
			// error response
		}
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
		projectId, err := service.SaveOrUpdateProject(p.ToProject(), userId, db)
		if err != nil {
			// log
			// error response
		}
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
		projectId, err := service.SaveOrUpdateProject(project, userId, db)
		if err != nil {
			fmt.Print(projectId)
			// log
			// error response
		}
		fmt.Fprintf(w, "")
	}
}

func DeleteProjectHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(r.PathValue("id"))
		err := service.DeleteProject(id, db)
		if err != nil {
			// log
			// error response
		}
		fmt.Fprintf(w, "")
	}
}
