package database

import (
	"database/sql"
	"fmt"
	"matask/internal/model"
)

const (
	findUserSql          = "SELECT u.id, u.name, u.email FROM matask_user u WHERE u.id = $1"
	findUserIdSql        = "SELECT u.id FROM matask_user u WHERE u.email = $1"
	findPasswordSql      = "SELECT uc.password FROM user_credentials uc WHERE uc.username = $1"
	insertUserSql        = "INSERT INTO matask_user (name, email, user_credentials_fk) VALUES ($1, $2, $3) RETURNING id"
	insertCredentialsSql = "INSERT INTO user_credentials (username, password) VALUES ($1, $2) RETURNING id"
)

func FindUser(id int, db *sql.DB) (model.MataskUser, error) {
	var user model.MataskUser
	if err := db.QueryRow(findUserSql, id).Scan(&user.Id, &user.Name, &user.Email); err != nil {
		return user, err
	}
	return user, nil
}

func FindUserId(email string, db *sql.DB) (int, error) {
	var userId int
	if err := db.QueryRow(findUserIdSql, email).Scan(&userId); err != nil {
		return userId, err
	}
	return userId, nil
}

func FindPassword(username string, db *sql.DB) (string, error) {
	var password string
	if err := db.QueryRow(findPasswordSql, username).Scan(&password); err != nil {
		return "", err
	}
	return password, nil
}

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
