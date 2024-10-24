package service

import (
	"database/sql"
	"matask/internal/database"
	"matask/internal/model"
)

func FindTasks(filter model.TaskFilter, db *sql.DB) []model.Task {
	return database.FindTasks(filter, db)
}
