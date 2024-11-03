package transport

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"matask/internal/service"
	"matask/internal/transport/handler"
	"matask/internal/transport/request"
	"matask/internal/transport/resource"
	"net/http"
	"strconv"
)

func GetTasksHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	filter := request.ToTaskFilter(r.URL.Query())
	filter.UserId = r.Context().Value(handler.UserIdKey).(int)
	result, err := service.FindTasks(filter, db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resource.FromTaskPageResult(result))
}

func UploadImageHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	userId := r.Context().Value(handler.UserIdKey).(int)

	r.ParseMultipartForm(2 << 20)
	file, _, err := r.FormFile("image")
	if err != nil {
		errStr := fmt.Sprintf("Error in reading the file %s\n", err)
		slog.Error(errStr)
		http.Error(w, errStr, http.StatusInternalServerError)
		return
	}
	defer file.Close()
	filebytes, err := io.ReadAll(file)
	if err != nil {
		errStr := fmt.Sprintf("Error in reading the file buffer %s\n", err)
		slog.Error(errStr)
		http.Error(w, errStr, http.StatusInternalServerError)
		return
	}

	taskType := r.FormValue("type")
	if taskType == "book" {
		err = service.UpdateBookCover(id, filebytes, userId, db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if taskType == "movie" {
		err = service.UpdateMoviePoster(id, filebytes, userId, db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	fmt.Fprintf(w, "")
}
