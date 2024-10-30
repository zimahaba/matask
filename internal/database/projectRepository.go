package database

import (
	"database/sql"
	"log/slog"
	"matask/internal/model"

	"github.com/lib/pq"
)

const (
	findProjectSql = `
		SELECT p.id, p.description, p.progress, p.dynamic_fields, t.id, t.name, t.type, t.started, t.ended 
		FROM project p
		INNER JOIN task t ON t.id = p.task_fk
		WHERE p.id = $1
		AND t.user_fk = $2
	`
	findProjectTaskIdSql = "SELECT t.id FROM project p INNER JOIN task t ON t.id = p.task_fk WHERE p.id = $1 AND t.user_fk = $2"
	insertProjectSql     = "INSERT INTO project (description, progress, dynamic_fields, task_fk) VALUES ($1, $2, $3, $4) RETURNING id"
	updateProjectSql     = `
		UPDATE project p 
		SET description = $3, progress = $4, dynamic_fields = $5 
		FROM task t 
		WHERE t.id = p.task_fk AND p.id = $1 AND t.user_fk = $2
	`
	deleteProjectSql = "DELETE FROM project p USING task t WHERE t.id = p.task_fk AND p.id = $1 AND t.user_fk = $2"
)

func FindProject(id int, userId int, db *sql.DB) (model.Project, error) {
	var p model.Project
	var description sql.NullString
	var started pq.NullTime
	var ended pq.NullTime
	if err := db.QueryRow(findProjectSql, id, userId).Scan(&p.Id, &description, &p.Progress, &p.DynamicFields, &p.Task.Id, &p.Task.Name, &p.Task.Type, &started, &ended); err != nil {
		slog.Error(err.Error())
		return model.Project{}, err
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
	return p, nil
}

func SaveOrUpdateProject(p model.Project, userId int, db *sql.DB) (int, error) {
	tx, err := db.Begin()
	if err != nil {
		slog.Error(err.Error())
		return -1, err
	}
	defer tx.Rollback()

	if p.Id != 0 {
		var taskId int
		if err := tx.QueryRow(findProjectTaskIdSql, p.Id, userId).Scan(&taskId); err != nil {
			slog.Error(err.Error())
			return -1, err
		}
		p.Task.Id = taskId
	}

	taskId, err := SaveOrUpdateTask(p.Task, userId, tx)
	if err != nil {
		slog.Error(err.Error())
		return -1, err
	}

	var id int
	var description sql.NullString
	if p.Description != "" {
		description = sql.NullString{String: p.Description, Valid: true}
	}

	if p.Id == 0 {
		err = tx.QueryRow(insertProjectSql, description, p.Progress, p.DynamicFields, taskId).Scan(&id)
	} else {
		_, err = tx.Exec(updateProjectSql, p.Id, userId, description, p.Progress, p.DynamicFields)
		id = p.Id
	}

	if err != nil {
		slog.Error(err.Error())
		return -1, err
	}

	if err = tx.Commit(); err != nil {
		slog.Error(err.Error())
		return -1, err
	}

	return id, nil
}
