package services

import (
	"database/sql"
	"matask/internal/database"
	"matask/internal/models"
)

func FindProject(id int, db *sql.DB) models.Project {
	return database.FindProject(id, db)
}

func CreateProject(p models.Project, db *sql.DB) int {
	return database.SaveProject(p, db)
}
