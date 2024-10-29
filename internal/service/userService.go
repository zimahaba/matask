package service

import (
	"database/sql"
	"matask/internal/database"
	"matask/internal/model"

	"golang.org/x/crypto/bcrypt"
)

func FindUser(id int, db *sql.DB) (model.MataskUser, error) {
	return database.FindUser(id, db)
}

func FindUserId(username string, db *sql.DB) (int, error) {
	return database.FindUserId(username, db)
}

func findPassword(username string, db *sql.DB) (string, error) {
	return database.FindPassword(username, db)
}

func CreateUser(user model.MataskUser, db *sql.DB) error {
	return database.SaveOrUpdateUser(user, db)
}

func VerifyCredentials(username string, password string, db *sql.DB) error {
	passwordHash, err := findPassword(username, db)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return err
	}
	return nil
}
