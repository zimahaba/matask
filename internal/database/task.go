package database

import (
	"database/sql"
	"matask/internal/models"
	"time"
)

var insertTaskSql = "INSERT INTO task (name, type, started, ended, created, updated) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"

func SaveTask(t models.Task, tx *sql.Tx) int {
	now := time.Now()
	var id int
	err := tx.QueryRow(insertTaskSql, t.Name, t.Type, t.Started, t.Ended, now, now).Scan(&id)
	if err != nil {
		panic(err)
	}
	return id
}
