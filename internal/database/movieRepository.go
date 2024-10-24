package database

import (
	"database/sql"
	"fmt"
	"matask/internal/model"

	"github.com/lib/pq"
)

var findMovieSql = `
	SELECT m.id, m.synopsis, m.comments, m.year, m.rate, m.director, m.poster_path, t.name, t.started, t.ended 
	FROM movie m
	INNER JOIN task t ON t.id = m.task_fk
	WHERE m.id = $1
`

var insertMovieSql = "INSERT INTO movie (synopsis, comments, year, rate, director, poster_path, task_fk) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"

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

func SaveMovie(m model.Movie, db *sql.DB) int {
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()

	taskId := SaveTask(m.Task, tx)
	fmt.Printf("movie: %v.\n", m)

	var id int
	tx.QueryRow(insertMovieSql, m.Synopsis, m.Comments, m.Year, m.Rate, m.Director, m.PosterPath, taskId).Scan(&id)
	if err = tx.Commit(); err != nil {
		panic(err)
	}

	return id
}
