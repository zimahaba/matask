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
	tx.QueryRow(insertMovieSql, synopsis, comments, year, rate, director, posterPath, taskId).Scan(&id)
	if err = tx.Commit(); err != nil {
		panic(err)
	}

	return id
}
