package database

import (
	"database/sql"
	"log/slog"
	"matask/internal/model"
	"time"

	"github.com/lib/pq"
)

const (
	insertTaskSql = "INSERT INTO task (name, type, started, ended, created, user_fk) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	updateTaskSql = "UPDATE task SET name = $3, started = $4, ended = $5 WHERE id = $1 AND user_fk = $2"
	deleteTaskSql = "DELETE FROM task WHERE id = $1 AND user_fk = $2"
)

var taskSortFieldMap = map[string]string{
	"id":   "b.id",
	"name": "t.name",
}

func FindTasks(f model.TaskFilter, db *sql.DB) (model.TaskPageResult, error) {
	query, countQuery := buildTaskFilteredQueries(f)

	var count int
	if err := db.QueryRow(countQuery, f.Name, f.Type, f.Started1, f.Started2, f.Ended1, f.Ended2, f.UserId).Scan(&count); err != nil {
		slog.Error(err.Error())
		return model.TaskPageResult{}, err
	}

	offset := getOffset(f.Page, f.Size)

	rows, err := db.Query(query, f.Name, f.Type, f.Started1, f.Started2, f.Ended1, f.Ended2, f.UserId, offset, f.Size)
	if err != nil {
		slog.Error(err.Error())
		return model.TaskPageResult{}, err
	}
	defer rows.Close()

	var tasks []model.TaskProjection

	for rows.Next() {
		var taskProjection model.TaskProjection
		var started pq.NullTime
		var ended pq.NullTime
		if err := rows.Scan(&taskProjection.Task.Id, &taskProjection.Task.Name, &taskProjection.Task.Type, &started, &ended, &taskProjection.ChildId); err != nil {
			slog.Error(err.Error())
			return model.TaskPageResult{}, err
		}
		if started.Valid {
			taskProjection.Task.Started = started.Time
		}
		if ended.Valid {
			taskProjection.Task.Ended = ended.Time
		}
		tasks = append(tasks, taskProjection)
	}
	if err = rows.Err(); err != nil {
		slog.Error(err.Error())
		return model.TaskPageResult{}, err
	}

	totalPages := calculateTotalPages(count, f.Size)

	pageResult := model.TaskPageResult{Tasks: tasks, Page: f.Page, Size: f.Size, TotalPages: totalPages, TotalElements: count}
	return pageResult, nil
}

func buildTaskFilteredQueries(f model.TaskFilter) (string, string) {
	selectQuery := `SELECT t.id, t.name, t.type, t.started, t.ended,
					CASE
						WHEN t.type = 'project' THEN p.id
						WHEN t.type = 'book' then b.id
						WHEN t.type = 'movie' then m.id
					END`
	selectCount := "SELECT COUNT(t.id)"
	from := `		FROM task t
					LEFT OUTER JOIN project p ON p.task_fk = t.id
					LEFT OUTER JOIN book b ON b.task_fk = t.id
					LEFT OUTER JOIN movie m ON m.task_fk = t.id
					WHERE ($1 = '' OR UPPER(t.name) like '%' || UPPER($1) || '%')
					AND ($2 = '' OR t.type = $2)
					AND (CAST($3 AS DATE) IS NULL OR t.started >= CAST($3 AS DATE))
					AND (CAST($4 AS DATE) IS NULL OR t.started <= CAST($4 AS DATE))
					AND (CAST($5 AS DATE) IS NULL OR (t.ended >= CAST($5 AS DATE) OR t.ended IS NULL))
					AND (CAST($6 AS DATE) IS NULL OR (t.ended <= CAST($6 AS DATE) OR t.ended IS NULL))
					AND t.user_fk = $7`
	order := getOrderQuery(f.SortField, f.SortDirection, taskSortFieldMap)
	offsetLimit := "OFFSET $8 LIMIT $9 "

	query := selectQuery + from + order + offsetLimit
	countQuery := selectCount + from

	return query, countQuery
}

func SaveOrUpdateTask(t model.Task, userId int, tx *sql.Tx) (int, error) {
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
		err = tx.QueryRow(insertTaskSql, t.Name, t.Type, started, ended, now, userId).Scan(&id)
	} else {
		_, err = tx.Exec(updateTaskSql, t.Id, userId, t.Name, started, ended)
		id = t.Id
	}

	if err != nil {
		slog.Error(err.Error())
		return -1, err
	}
	return id, nil
}

func DeleteTaskCascade(childId int, childType string, userId int, db *sql.DB) error {
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
		slog.Error(err.Error())
		return err
	}
	defer tx.Rollback()

	var taskId int
	if err := tx.QueryRow(findTaskIdSql, childId, userId).Scan(&taskId); err != nil {
		slog.Error(err.Error())
		return err
	}
	_, err = tx.Exec(deleteChildSql, childId, userId)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	_, err = tx.Exec(deleteTaskSql, taskId, userId)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	if err = tx.Commit(); err != nil {
		slog.Error(err.Error())
		return err
	}
	return nil
}
