package database

import (
	"database/sql"
	"fmt"
	"matask/internal/models"
)

var insertProjectSql = "INSERT INTO project (description, progress, task_fk) VALUES ($1, $2, $3) RETURNING id"

func SaveProject(p models.Project, db *sql.DB) int {
	fmt.Printf("project: %v\n\n", p)

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()

	taskId := SaveTask(p.Task, tx)
	fmt.Printf("taskId: %v\n\n", taskId)

	var id int
	tx.QueryRow(insertProjectSql, p.Description, p.Progress, taskId).Scan(&id)
	if err = tx.Commit(); err != nil {
		panic(err)
	}

	return id
}
