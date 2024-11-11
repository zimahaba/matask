package service

import (
	"database/sql"
	"matask/internal/database"
	"matask/internal/model"
)

func FindFilteredTasks(filter model.TaskFilter, db *sql.DB) (model.TaskPageResult, error) {
	return database.FindFilteredTasks(filter, db)
}
