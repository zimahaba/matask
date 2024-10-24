package service

import (
	"database/sql"
	"matask/internal/database"
	"matask/internal/model"
)

func FindBook(id int, db *sql.DB) model.Book {
	return database.FindBook(id, db)
}

func CreateBook(p model.Book, db *sql.DB) int {
	return database.SaveBook(p, db)
}
