package transport

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"matask/internal/service"
	"matask/internal/transport/request"
	"net/http"
	"os"
)

func GetTasksHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filter := request.ToTaskFilter(r.URL.Query())
		result := service.FindTasks(filter, db)
		json.NewEncoder(w).Encode(result)
	}
}

func UploadCoverHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		r.ParseMultipartForm(2 << 20)
		file, _, err := r.FormFile("cover")
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

		os.WriteFile(os.Getenv("COVER_PATH"), filebytes, 0666)

		fmt.Fprintf(w, "UploadCoverHandler")
	}
}
