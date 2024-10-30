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

func GetTasksHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filter := request.ToTaskFilter(r.URL.Query())
		filter.UserId = r.Context().Value(handler.UserIdKey).(int)
		result, err := service.FindTasks(filter, db)
		if err != nil {
			slog.Debug(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(resource.FromTaskPageResult(result))
	}
}

func UploadImageHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(r.PathValue("id"))

		r.ParseMultipartForm(2 << 20)
		file, _, err := r.FormFile("image")
		if err != nil {
			errStr := fmt.Sprintf("Error in reading the file %s\n", err)
			fmt.Println(errStr)
			fmt.Fprintf(w, "Error in reading the file")
			return
		}
		defer file.Close()
		filebytes, err := io.ReadAll(file)
		if err != nil {
			errStr := fmt.Sprintf("Error in reading the file buffer %s\n", err)
			fmt.Println(errStr)
			fmt.Fprintf(w, "Error in reading the file buffer")
			return
		}

		taskType := r.FormValue("type")
		if taskType == "book" {
			service.UpdateBookCover(id, filebytes, db)
		} else if taskType == "movie" {
			service.UpdateMoviePoster(id, filebytes, db)
		}

		fmt.Fprintf(w, "")
	}
}
