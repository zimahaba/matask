package database

import (
	"database/sql"
	"log/slog"
	"matask/internal/model"
)

func FindUser(id int, db *sql.DB) (model.MataskUser, error) {
	const query = "SELECT u.id, u.name, u.email FROM matask_user u WHERE u.id = $1"

	var user model.MataskUser
	if err := db.QueryRow(query, id).Scan(&user.Id, &user.Name, &user.Email); err != nil {
		slog.Error(err.Error())
		return user, err
	}
	return user, nil
}

func FindUserId(email string, db *sql.DB) (int, error) {
	const query = "SELECT u.id FROM matask_user u WHERE u.email = $1"

	var userId int
	if err := db.QueryRow(query, email).Scan(&userId); err != nil {
		slog.Error(err.Error())
		return userId, err
	}
	return userId, nil
}

func FindUsernameByRefreshToken(refreshToken string, db *sql.DB) (string, error) {
	const query = "SELECT uc.username FROM user_credentials uc WHERE uc.refresh_token = $1"

	var username string
	if err := db.QueryRow(query, refreshToken).Scan(&username); err != nil {
		slog.Error(err.Error())
		return "", err
	}
	return username, nil
}

func FindPassword(username string, db *sql.DB) (string, error) {
	const query = "SELECT uc.password FROM user_credentials uc WHERE uc.username = $1"

	var password string
	if err := db.QueryRow(query, username).Scan(&password); err != nil {
		slog.Error(err.Error())
		return "", err
	}
	return password, nil
}

func SaveOrUpdateUser(u model.MataskUser, db *sql.DB) error {
	const query = "INSERT INTO matask_user (name, email, user_credentials_fk) VALUES ($1, $2, $3) RETURNING id"

	tx, err := db.Begin()
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	defer tx.Rollback()

	var credentialsId int
	if u.Id == 0 {
		credentialsId, err = SaveUserCredentials(u.Credentials, tx)
		if err != nil {
			slog.Error(err.Error())
			return err
		}
	}

	var id int
	if u.Id == 0 {
		err = tx.QueryRow(query, u.Name, u.Email, credentialsId).Scan(&id)
		if err != nil {
			slog.Error(err.Error())
			return err
		}
	} else {
		_, err = tx.Exec("", u.Id, u.Name)
		if err != nil {
			slog.Error(err.Error())
			return err
		}
		id = u.Id
	}

	if err = tx.Commit(); err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
}

func SaveUserCredentials(c model.UserCredentials, tx *sql.Tx) (int, error) {
	const query = "INSERT INTO user_credentials (username, password) VALUES ($1, $2) RETURNING id"

	var id int
	err := tx.QueryRow(query, c.Username, c.Password).Scan(&id)
	if err != nil {
		slog.Error(err.Error())
		return id, err
	}
	return id, nil
}

func UpsertRefreshToken(refreshToken string, username string, db *sql.DB) error {
	const query = "UPDATE user_credentials SET refresh_token = $2 WHERE username = $1"

	_, err := db.Exec(query, username, refreshToken)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	return nil
}
