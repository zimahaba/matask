package transport

import (
	"database/sql"
	"encoding/json"
	"matask/internal/services"
	"matask/internal/transport/request"
	"net/http"
)

func GetTasksHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filter := request.ToTaskFilter(r.URL.Query())
		result := services.FindTasks(filter, db)
		json.NewEncoder(w).Encode(result)
	}
}
