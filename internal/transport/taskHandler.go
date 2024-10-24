package transport

import (
	"database/sql"
	"encoding/json"
	"matask/internal/service"
	"matask/internal/transport/request"
	"net/http"
)

func GetTasksHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filter := request.ToTaskFilter(r.URL.Query())
		result := service.FindTasks(filter, db)
		json.NewEncoder(w).Encode(result)
	}
}
