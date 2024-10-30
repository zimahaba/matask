package service

import (
	"database/sql"
	"matask/internal/database"
	"matask/internal/model"
	"os"
)

func FindBook(id int, userId int, db *sql.DB) (model.Book, error) {
	return database.FindBook(id, userId, db)
}

func SaveOrUpdateBook(p model.Book, userId int, db *sql.DB) (int, error) {
	return database.SaveOrUpdateBook(p, userId, db)
}

func UpdateBookCover(id int, filebytes []byte, userId int, db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	coverPath := os.Getenv("COVER_PATH") + "cover"
	err = database.UpdateBookCover(id, coverPath, userId, tx)
	if err != nil {
		return err
	}

	err = os.WriteFile(coverPath, filebytes, 0666)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func DeleteBook(bookId int, userId int, db *sql.DB) error {
	return database.DeleteTaskCascade(bookId, "book", userId, db)
}
