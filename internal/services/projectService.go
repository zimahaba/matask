package services

import (
	"database/sql"
	"matask/internal/database"
	"matask/internal/models"
)

func CreateProject(p models.Project, db *sql.DB) int {
	return database.SaveProject(p, db)
}
