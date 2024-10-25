package service

import (
	"database/sql"
	"matask/internal/database"
	"matask/internal/model"
)

func FindMovie(id int, db *sql.DB) model.Movie {
	return database.FindMovie(id, db)
}

func SaveOrUpdateMovie(p model.Movie, db *sql.DB) int {
	return database.SaveOrUpdateMovie(p, db)
}
