package service

import (
	"database/sql"
	"matask/internal/database"
	"matask/internal/model"
)

func FindBook(id int, db *sql.DB) model.Book {
	return database.FindBook(id, db)
}

func SaveOrUpdateBook(p model.Book, db *sql.DB) int {
	return database.SaveOrUpdateBook(p, db)
}

func DeleteBook(bookId int, db *sql.DB) {
	database.DeleteTaskCascade(bookId, "book", db)
}
