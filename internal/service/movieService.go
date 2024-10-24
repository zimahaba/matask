package service

import (
	"database/sql"
	"matask/internal/database"
	"matask/internal/model"
)

func FindMovie(id int, db *sql.DB) model.Movie {
	return database.FindMovie(id, db)
}

func CreateMovie(p model.Movie, db *sql.DB) int {
	return database.SaveMovie(p, db)
}
