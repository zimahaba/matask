package service

import (
	"database/sql"
	"matask/internal/database"
	"matask/internal/model"
)

func FindProject(id int, db *sql.DB) model.Project {
	return database.FindProject(id, db)
}

func CreateProject(p model.Project, db *sql.DB) int {
	return database.SaveProject(p, db)
}
