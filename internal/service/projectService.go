package service

import (
	"database/sql"
	"matask/internal/database"
	"matask/internal/model"
)

func FindProject(id int, userId int, db *sql.DB) (model.Project, error) {
	return database.FindProject(id, userId, db)
}

func SaveOrUpdateProject(p model.Project, userId int, db *sql.DB) (int, error) {
	return database.SaveOrUpdateProject(p, userId, db)
}

func DeleteProject(projectId int, userId int, db *sql.DB) error {
	return database.DeleteTaskCascade(projectId, "project", userId, db)
}
