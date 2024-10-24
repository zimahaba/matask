package database

import (
	"database/sql"
	"matask/internal/model"

	"github.com/lib/pq"
)

var findProjectSql = `
	SELECT p.id, p.description, p.progress, t.name, t.started, t.ended 
	FROM project p
	INNER JOIN task t ON t.id = p.task_fk
	WHERE p.id = $1
`

var insertProjectSql = "INSERT INTO project (description, progress, task_fk) VALUES ($1, $2, $3) RETURNING id"

func FindProject(id int, db *sql.DB) model.Project {
	var p model.Project
	var started pq.NullTime
	var ended pq.NullTime
	if err := db.QueryRow(findProjectSql, id).Scan(&p.Id, &p.Description, &p.Progress, &p.Task.Name, &started, &ended); err != nil {
		panic(err)
	}
	if started.Valid {
		p.Task.Started = started.Time
	}
	if ended.Valid {
		p.Task.Ended = ended.Time
	}
	return p
}

func SaveProject(p model.Project, db *sql.DB) int {
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()

	taskId := SaveTask(p.Task, tx)

	var id int
	tx.QueryRow(insertProjectSql, p.Description, p.Progress, taskId).Scan(&id)
	if err = tx.Commit(); err != nil {
		panic(err)
	}

	return id
}
