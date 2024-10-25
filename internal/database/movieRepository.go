package database

import (
	"database/sql"
	"matask/internal/model"

	"github.com/lib/pq"
)

const (
	findMovieSql = `
		SELECT m.id, m.synopsis, m.comments, m.year, m.rate, m.director, m.poster_path, t.name, t.started, t.ended 
		FROM movie m
		INNER JOIN task t ON t.id = m.task_fk
		WHERE m.id = $1
	`
	findMovieTaskIdSql = "SELECT t.id FROM movie m INNER JOIN task t ON t.id = m.task_fk WHERE m.id = $1"
	insertMovieSql     = "INSERT INTO movie (synopsis, comments, year, rate, director, poster_path, task_fk) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
	updateMovieSql     = "UPDATE movie SET synopsis = $2, comments = $3, year = $4, rate = $5, director = $6 WHERE id = $1"
	deleteMovieSql     = "DELETE FROM movie WHERE id = $1"
)

func FindMovie(id int, db *sql.DB) model.Movie {
	var m model.Movie
	var synopsis sql.NullString
	var comments sql.NullString
	var year sql.NullString
	var rate sql.NullInt32
	var director sql.NullString
	var coverPath sql.NullString
	var started pq.NullTime
	var ended pq.NullTime
	if err := db.QueryRow(findMovieSql, id).Scan(&m.Id, &synopsis, &comments, &year, &rate, &director, &coverPath, &m.Task.Name, &started, &ended); err != nil {
		panic(err)
	}
	if synopsis.Valid {
		m.Synopsis = synopsis.String
	}
	if comments.Valid {
		m.Comments = comments.String
	}
	if year.Valid {
		m.Year = year.String
	}
	if rate.Valid {
		m.Rate = int(rate.Int32)
	}
	if director.Valid {
		m.Director = director.String
	}
	if coverPath.Valid {
		m.PosterPath = coverPath.String
	}
	if started.Valid {
		m.Task.Started = started.Time
	}
	if ended.Valid {
		m.Task.Ended = ended.Time
	}
	return m
}

func SaveOrUpdateMovie(m model.Movie, db *sql.DB) int {
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()

	if m.Id != 0 {
		var taskId int
		if err := tx.QueryRow(findMovieTaskIdSql, m.Id).Scan(&taskId); err != nil {
			panic(err)
		}
		m.Task.Id = taskId
	}

	taskId := SaveOrUpdateTask(m.Task, tx)

	var id int
	var synopsis sql.NullString
	if m.Synopsis != "" {
		synopsis = sql.NullString{String: m.Synopsis, Valid: true}
	}
	var comments sql.NullString
	if m.Comments != "" {
		comments = sql.NullString{String: m.Comments, Valid: true}
	}
	var year sql.NullString
	if m.Year != "" {
		year = sql.NullString{String: m.Year, Valid: true}
	}
	var rate sql.NullInt32
	if m.Rate != 0 {
		rate = sql.NullInt32{Int32: int32(m.Rate), Valid: true}
	}
	var director sql.NullString
	if m.Director != "" {
		director = sql.NullString{String: m.Director, Valid: true}
	}
	var posterPath sql.NullString
	if m.PosterPath != "" {
		posterPath = sql.NullString{String: m.PosterPath, Valid: true}
	}

	if m.Id == 0 {
		tx.QueryRow(insertMovieSql, synopsis, comments, year, rate, director, posterPath, taskId).Scan(&id)
		if err != nil {
			panic(err)
		}
	} else {
		_, err = tx.Exec(updateMovieSql, m.Id, synopsis, comments, year, rate, director)
		if err != nil {
			panic(err)
		}
		id = m.Id
	}

	if err = tx.Commit(); err != nil {
		panic(err)
	}

	return id
}
