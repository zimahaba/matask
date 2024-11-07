package transport

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"matask/internal/service"
	"matask/internal/transport/handler"
	"matask/internal/transport/request"
	"matask/internal/transport/resource"
	"net/http"
	"strconv"
)

func GetFilteredProjectsHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	filter := request.ToProjectFilter(r.URL.Query())
	filter.UserId = r.Context().Value(handler.UserIdKey).(int)
	result, err := service.FindFilteredProjects(filter, db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resource.FromProjectPageResult(result))
}

func GetProjectHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	userId := r.Context().Value(handler.UserIdKey).(int)
	p, err := service.FindProject(id, userId, db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resource.FromProject(p))
}

func SaveProjectHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var p request.ProjectRequest
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	project := p.ToProject()

	id, _ := strconv.Atoi(r.PathValue("id"))
	if id > 0 {
		project.Id = id
	}

	userId := r.Context().Value(handler.UserIdKey).(int)

	projectId, err := service.SaveProject(project, userId, db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resource.IdResource{Id: projectId})
}

func DeleteProjectHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	userId := r.Context().Value(handler.UserIdKey).(int)
	err := service.DeleteProject(id, userId, db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "")
}
