package database

import (
	"database/sql"
	"fmt"
	"matask/internal/model"
	"time"

	"github.com/lib/pq"
)

var findTasksSelectSql = "SELECT t.id, t.name, t.type, t.started, t.ended"
var findTasksSelectCountSql = "SELECT COUNT(t.id)"
var findTasksFromWhereSql = `
		FROM task t
		WHERE ($1 = '' OR UPPER(t.name) like '%' || UPPER($1) || '%')
		AND ($2 = '' OR t.type = $2)
		AND (CAST($3 AS DATE) IS NULL OR t.started >= CAST($3 AS DATE))
		AND (CAST($4 AS DATE) IS NULL OR t.started <= CAST($4 AS DATE))
		AND (CAST($5 AS DATE) IS NULL OR (t.ended >= CAST($5 AS DATE) OR t.ended IS NULL))
		AND (CAST($6 AS DATE) IS NULL OR (t.ended <= CAST($6 AS DATE) OR t.ended IS NULL))
	`
var findTaskOrderSql = " ORDER BY %s %s "
var findTaskPageSql = " OFFSET $7 LIMIT $8 "

const (
	insertTaskSql = "INSERT INTO task (name, type, started, ended, created) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	updateTaskSql = "UPDATE task SET name = $2, started = $3, ended = $4 WHERE id = $1"
	deleteTaskSql = "DELETE FROM task WHERE id = $1"
)

var sortFieldMap = map[string]string{
	"id":   "id",
	"name": "name",
}

var sortDirectionMap = map[string]string{
	"ASC":  "ASC",
	"DESC": "DESC",
}

func FindTasks(f model.TaskFilter, db *sql.DB) model.TaskPageResult {
	sortField := sortFieldMap[f.SortField]
	if sortField == "" {
		sortField = "id"
	}
	sortDirection := sortDirectionMap[f.SortDirection]
	if sortDirection == "" {
		sortDirection = "ASC"
	}
	query := findTasksSelectSql + findTasksFromWhereSql + fmt.Sprintf(findTaskOrderSql, sortField, sortDirection) + findTaskPageSql
	countQuery := findTasksSelectCountSql + findTasksFromWhereSql

	var count int
	if err := db.QueryRow(countQuery, f.Name, f.Type, f.Started1, f.Started2, f.Ended1, f.Ended2).Scan(&count); err != nil {
		panic(err)
	}

	var offset int
	if f.Page <= 1 {
		offset = 0
	} else {
		offset = (f.Page - 1) * f.Size
	}

	rows, err := db.Query(query, f.Name, f.Type, f.Started1, f.Started2, f.Ended1, f.Ended2, offset, f.Size)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var tasks []model.TaskProjection

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
		taskProjection := model.TaskProjection{Task: t}
		tasks = append(tasks, taskProjection)
	}
	if err = rows.Err(); err != nil {
		panic(err)
	}

	totalPages := count / f.Size
	remainder := count % f.Size
	if totalPages == 0 || (totalPages > 0 && remainder > 0) {
		totalPages++
	}

	pageResult := model.TaskPageResult{Tasks: tasks, Page: f.Page, Size: f.Size, TotalPages: totalPages, TotalElements: count}
	return pageResult
}

func SaveOrUpdateTask(t model.Task, tx *sql.Tx) int {
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
	var err error
	if t.Id == 0 {
		err = tx.QueryRow(insertTaskSql, t.Name, t.Type, started, ended, now).Scan(&id)
	} else {
		_, err = tx.Exec(updateTaskSql, t.Id, t.Name, started, ended)
		id = t.Id
	}

	if err != nil {
		panic(err)
	}
	return id
}

func DeleteTaskCascade(childId int, childType string, db *sql.DB) {
	var findTaskIdSql string
	var deleteChildSql string
	if childType == "project" {
		findTaskIdSql = findProjectTaskIdSql
		deleteChildSql = deleteProjectSql
	} else if childType == "book" {
		findTaskIdSql = findBookTaskIdSql
		deleteChildSql = deleteBookSql
	} else if childType == "movie" {
		findTaskIdSql = findMovieTaskIdSql
		deleteChildSql = deleteMovieSql
	}

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()

	var taskId int
	if err := tx.QueryRow(findTaskIdSql, childId).Scan(&taskId); err != nil {
		panic(err)
	}
	_, err = tx.Exec(deleteChildSql, childId)
	if err != nil {
		panic(err)
	}
	_, err = tx.Exec(deleteTaskSql, taskId)
	if err != nil {
		panic(err)
	}

	if err = tx.Commit(); err != nil {
		panic(err)
	}
}
