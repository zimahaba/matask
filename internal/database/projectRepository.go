package database

import (
	"database/sql"
	"matask/internal/model"

	"github.com/lib/pq"
)

var findProjectSql = `
	SELECT p.id, p.description, p.progress, t.id, t.name, t.type, t.started, t.ended 
	FROM project p
	INNER JOIN task t ON t.id = p.task_fk
	WHERE p.id = $1
`

var findProjectTaskId = "SELECT t.id FROM project p INNER JOIN task t ON t.id = p.task_fk WHERE p.id = $1"

var insertProjectSql = "INSERT INTO project (description, progress, task_fk) VALUES ($1, $2, $3) RETURNING id"

var updateProjectSql = "UPDATE project SET description = $2, progress = $3 WHERE id = $1"

func FindProject(id int, db *sql.DB) model.Project {
	var p model.Project
	var description sql.NullString
	var started pq.NullTime
	var ended pq.NullTime
	if err := db.QueryRow(findProjectSql, id).Scan(&p.Id, &description, &p.Progress, &p.Task.Id, &p.Task.Name, &p.Task.Type, &started, &ended); err != nil {
		panic(err)
	}
	if description.Valid {
		p.Description = description.String
	}
	if started.Valid {
		p.Task.Started = started.Time
	}
	if ended.Valid {
		p.Task.Ended = ended.Time
	}
	return p
}

func SaveOrUpdateProject(p model.Project, db *sql.DB) int {
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()

	if p.Id != 0 {
		var taskId int
		if err := db.QueryRow(findProjectTaskId, p.Id).Scan(&taskId); err != nil {
			panic(err)
		}
		p.Task.Id = taskId
	}

	taskId := SaveOrUpdateTask(p.Task, tx)

	var id int
	var description sql.NullString
	if p.Description != "" {
		description = sql.NullString{String: p.Description, Valid: true}
	}

	if p.Id == 0 {
		tx.QueryRow(insertProjectSql, description, p.Progress, taskId).Scan(&id)
	} else {
		_, err = tx.Exec(updateProjectSql, p.Id, description, p.Progress)
		if err != nil {
			panic(err)
		}
		id = p.Id
	}

	if err = tx.Commit(); err != nil {
		panic(err)
	}

	return id
}
