package service

import (
	"database/sql"
	"matask/internal/database"
	"matask/internal/model"
)

func FindProject(id int, db *sql.DB) model.Project {
	return database.FindProject(id, db)
}

func SaveOrUpdateProject(p model.Project, userId int, db *sql.DB) int {
	return database.SaveOrUpdateProject(p, userId, db)
}

func DeleteProject(projectId int, db *sql.DB) {
	database.DeleteTaskCascade(projectId, "project", db)
}
