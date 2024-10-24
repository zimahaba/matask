package database

import (
	"database/sql"
	"fmt"
	"matask/internal/model"
	"time"

	"github.com/lib/pq"
)

var findTasksSql = `
	SELECT t.id, t.name, t.type, t.started, t.ended
	FROM task t
	WHERE ($1 = '' OR UPPER(t.name) like '%' || UPPER($1) || '%')
	AND ($2 = '' OR t.type = $2)
	AND (CAST($3 AS DATE) IS NULL OR t.started >= CAST($3 AS DATE))
	AND (CAST($4 AS DATE) IS NULL OR t.started <= CAST($4 AS DATE))
	AND (CAST($5 AS DATE) IS NULL OR (t.ended >= CAST($5 AS DATE) OR t.ended IS NULL))
	AND (CAST($6 AS DATE) IS NULL OR (t.ended <= CAST($6 AS DATE) OR t.ended IS NULL))
`

var insertTaskSql = "INSERT INTO task (name, type, started, ended, created) VALUES ($1, $2, $3, $4, $5) RETURNING id"

func FindTasks(f model.TaskFilter, db *sql.DB) []model.Task {
	fmt.Printf("filter: %v.\n", f)
	rows, err := db.Query(findTasksSql, f.Name, f.Type, f.Started1, f.Started2, f.Ended1, f.Ended2)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var tasks []model.Task

	for rows.Next() {
		var t model.Task
		var started pq.NullTime
		var ended pq.NullTime
		if err := rows.Scan(&t.Id, &t.Name, &t.Type, &started, &ended); err != nil {
			panic(err)
		}
		if started.Valid {
			t.Started = started.Time
		}
		if ended.Valid {
			t.Ended = ended.Time
		}
		tasks = append(tasks, t)
	}
	if err = rows.Err(); err != nil {
		panic(err)
	}
	return tasks
}

func SaveTask(t model.Task, tx *sql.Tx) int {
	var started *time.Time
	if !t.Started.IsZero() {
		started = &t.Started
	}
	var ended *time.Time
	if !t.Ended.IsZero() {
		ended = &t.Ended
	}

	now := time.Now()
	var id int
	err := tx.QueryRow(insertTaskSql, t.Name, t.Type, started, ended, now).Scan(&id)
	if err != nil {
		panic(err)
	}
	return id
}
