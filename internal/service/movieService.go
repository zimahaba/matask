package service

import (
	"database/sql"
	"matask/internal/database"
	"matask/internal/model"
	"os"
)

func FindMovie(id int, db *sql.DB) model.Movie {
	return database.FindMovie(id, db)
}

func SaveOrUpdateMovie(p model.Movie, db *sql.DB) int {
	return database.SaveOrUpdateMovie(p, db)
}

func UpdateMoviePoster(id int, filebytes []byte, db *sql.DB) {
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()

	posterPath := os.Getenv("POSTER_PATH") + "poster"
	err = database.UpdateMoviePoster(id, posterPath, tx)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(posterPath, filebytes, 0666)
	if err != nil {
		panic(err)
	}

	if err = tx.Commit(); err != nil {
		panic(err)
	}
}

func DeleteMovie(movieId int, db *sql.DB) {
	database.DeleteTaskCascade(movieId, "movie", db)
}
