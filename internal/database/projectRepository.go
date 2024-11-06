package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"matask/internal/model"

	"github.com/lib/pq"
)

var findFilteredProjectsSelectSql = "SELECT p.id, t.name, p.progress"
var findFilteredProjectsSelectCountSql = "SELECT COUNT(p.id)"
var findFilteredProjectsFromWhereSql = `
		FROM project p
		INNER JOIN task t ON t.id = p.task_fk
		WHERE ($1 = '' OR UPPER(t.name) like '%' || UPPER($1) || '%')
		AND (p.progress >= COALESCE(CAST($2 AS INTEGER), 0))
		AND (p.progress <= COALESCE(CAST($3 AS INTEGER), 100))
		AND t.user_fk = $4
	`
var findFilteredProjectsOrderSql = " ORDER BY %s %s "
var findFilteredProjectsPageSql = " OFFSET $5 LIMIT $6 "

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

var projectSortFieldMap = map[string]string{
	"id":       "p.id",
	"name":     "t.name",
	"progress": "p.progress",
}

func FindFilteredProjects(f model.ProjectFilter, db *sql.DB) (model.ProjectPageResult, error) {
	sortField := projectSortFieldMap[f.SortField]
	if sortField == "" {
		sortField = "id"
	}
	sortDirection := sortDirectionMap[f.SortDirection]
	if sortDirection == "" {
		sortDirection = "ASC"
	}
	query := findFilteredProjectsSelectSql + findFilteredProjectsFromWhereSql + fmt.Sprintf(findFilteredProjectsOrderSql, sortField, sortDirection) + findFilteredProjectsPageSql
	countQuery := findFilteredProjectsSelectCountSql + findFilteredProjectsFromWhereSql

	var count int
	progress1 := sql.NullInt64{}
	if f.Progress1 >= 0 {
		progress1 = sql.NullInt64{Int64: int64(f.Progress1), Valid: true}
	}
	progress2 := sql.NullInt64{}
	if f.Progress2 >= 0 {
		progress2 = sql.NullInt64{Int64: int64(f.Progress2), Valid: true}
	}
	if err := db.QueryRow(countQuery, f.Name, progress1, progress2, f.UserId).Scan(&count); err != nil {
		slog.Error(err.Error())
		return model.ProjectPageResult{}, err
	}

	var offset int
	if f.Page <= 1 {
		offset = 0
	} else {
		offset = (f.Page - 1) * f.Size
	}

	rows, err := db.Query(query, f.Name, progress1, progress2, f.UserId, offset, f.Size)
	if err != nil {
		slog.Error(err.Error())
		return model.ProjectPageResult{}, err
	}
	defer rows.Close()

	var projects []model.ProjectProjection

	for rows.Next() {
		var projectProjection model.ProjectProjection
		if err := rows.Scan(&projectProjection.Id, &projectProjection.Name, &projectProjection.Progress); err != nil {
			slog.Error(err.Error())
			return model.ProjectPageResult{}, err
		}
		projects = append(projects, projectProjection)
	}
	if err = rows.Err(); err != nil {
		slog.Error(err.Error())
		return model.ProjectPageResult{}, err
	}

	totalPages := count / f.Size
	remainder := count % f.Size
	if totalPages == 0 || (totalPages > 0 && remainder > 0) {
		totalPages++
	}

	pageResult := model.ProjectPageResult{Projects: projects, Page: f.Page, Size: f.Size, TotalPages: totalPages, TotalElements: count}
	return pageResult, nil
}

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
