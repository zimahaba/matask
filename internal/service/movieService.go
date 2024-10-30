package service

import (
	"database/sql"
	"matask/internal/database"
	"matask/internal/model"
	"os"
)

func FindMovie(id int, userId int, db *sql.DB) (model.Movie, error) {
	return database.FindMovie(id, userId, db)
}

func SaveOrUpdateMovie(p model.Movie, userId int, db *sql.DB) (int, error) {
	return database.SaveOrUpdateMovie(p, userId, db)
}

func UpdateMoviePoster(id int, filebytes []byte, userId int, db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	posterPath := os.Getenv("POSTER_PATH") + "poster"
	err = database.UpdateMoviePoster(id, posterPath, userId, tx)
	if err != nil {
		return err
	}

	err = os.WriteFile(posterPath, filebytes, 0666)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func DeleteMovie(movieId int, userId int, db *sql.DB) error {
	return database.DeleteTaskCascade(movieId, "movie", userId, db)
}
