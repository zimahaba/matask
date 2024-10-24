package transport

import (
	"database/sql"
	"encoding/json"
	"matask/internal/services"
	payload "matask/internal/transport/payloads"
	"net/http"
)

func GetTasksHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filter := payload.ToTaskFilter(r.URL.Query())
		result := services.FindTasks(filter, db)
		json.NewEncoder(w).Encode(result)
	}
}
