package service

import (
	"database/sql"
	"matask/internal/database"
	"matask/internal/model"
	"os"
)

func FindBook(id int, db *sql.DB) (model.Book, error) {
	return database.FindBook(id, db)
}

func SaveOrUpdateBook(p model.Book, userId int, db *sql.DB) (int, error) {
	return database.SaveOrUpdateBook(p, userId, db)
}

func UpdateBookCover(id int, filebytes []byte, db *sql.DB) {
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()

	coverPath := os.Getenv("COVER_PATH") + "cover"
	err = database.UpdateBookCover(id, coverPath, tx)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(coverPath, filebytes, 0666)
	if err != nil {
		panic(err)
	}

	if err = tx.Commit(); err != nil {
		panic(err)
	}
}

func DeleteBook(bookId int, db *sql.DB) error {
	return database.DeleteTaskCascade(bookId, "book", db)
}
