package service

import (
	"database/sql"
	"matask/internal/database"
	"matask/internal/model"
)

func FindTasks(filter model.TaskFilter, db *sql.DB) (model.TaskPageResult, error) {
	return database.FindTasks(filter, db)
}
