package services

import (
	"database/sql"
	"matask/internal/database"
	"matask/internal/models"
)

func FindTasks(filter models.TaskFilter, db *sql.DB) []models.Task {
	return database.FindTasks(filter, db)
}
