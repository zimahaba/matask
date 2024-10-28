package database

import (
	"database/sql"
	"fmt"
	"matask/internal/model"
)

const (
	insertUserSql        = "INSERT INTO matask_user (name, email, user_credentials_fk) VALUES ($1, $2, $3) RETURNING id"
	insertCredentialsSql = "INSERT INTO user_credentials (username, password) VALUES ($1, $2) RETURNING id"
)

func SaveOrUpdateUser(u model.MataskUser, db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()

	var credentialsId int
	if u.Id == 0 {
		credentialsId, err = SaveUserCredentials(u.Credentials, tx)
		if err != nil {
			return err
		}
	}

	var id int
	if u.Id == 0 {
		err = tx.QueryRow(insertUserSql, u.Name, u.Email, credentialsId).Scan(&id)
		if err != nil {
			panic(err)
		}
	} else {
		_, err = tx.Exec("", u.Id, u.Name)
		if err != nil {
			panic(err)
		}
		id = u.Id
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func SaveUserCredentials(c model.UserCredentials, tx *sql.Tx) (int, error) {
	var id int
	fmt.Println(c.Username)
	err := tx.QueryRow(insertCredentialsSql, c.Username, c.Password).Scan(&id)
	if err != nil {
		return id, err
	}
	return id, nil
}
