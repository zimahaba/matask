package service

import (
	"database/sql"
	"matask/internal/database"
	"matask/internal/model"
)

func CreateUser(user model.MataskUser, db *sql.DB) error {
	return database.SaveOrUpdateUser(user, db)
}
