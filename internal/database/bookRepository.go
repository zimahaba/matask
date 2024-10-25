package database

import (
	"database/sql"
	"matask/internal/model"

	"github.com/lib/pq"
)

var findBookSql = `
	SELECT b.id, b.progress, b.author, b.synopsis, b.comments, b.year, b.rate, b.cover_path, t.name, t.started, t.ended 
	FROM book b
	INNER JOIN task t ON t.id = b.task_fk
	WHERE b.id = $1
`

var insertBookSql = "INSERT INTO book (progress, author, synopsis, comments, year, rate, cover_path, task_fk) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id"

func FindBook(id int, db *sql.DB) model.Book {
	var b model.Book
	var author sql.NullString
	var synopsis sql.NullString
	var comments sql.NullString
	var year sql.NullString
	var rate sql.NullInt32
	var coverPath sql.NullString
	var started pq.NullTime
	var ended pq.NullTime
	if err := db.QueryRow(findBookSql, id).Scan(&b.Id, &b.Progress, &author, &synopsis, &comments, &year, &rate, &coverPath, &b.Task.Name, &started, &ended); err != nil {
		panic(err)
	}
	if author.Valid {
		b.Author = author.String
	}
	if synopsis.Valid {
		b.Synopsis = synopsis.String
	}
	if comments.Valid {
		b.Comments = comments.String
	}
	if year.Valid {
		b.Year = year.String
	}
	if rate.Valid {
		b.Rate = int(rate.Int32)
	}
	if coverPath.Valid {
		b.CoverPath = coverPath.String
	}
	if started.Valid {
		b.Task.Started = started.Time
	}
	if ended.Valid {
		b.Task.Ended = ended.Time
	}
	return b
}

func SaveBook(b model.Book, db *sql.DB) int {
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()

	taskId := SaveOrUpdateTask(b.Task, tx)

	var id int
	var author sql.NullString
	if b.Author != "" {
		author = sql.NullString{String: b.Author, Valid: true}
	}
	var synopsis sql.NullString
	if b.Synopsis != "" {
		synopsis = sql.NullString{String: b.Synopsis, Valid: true}
	}
	var comments sql.NullString
	if b.Comments != "" {
		comments = sql.NullString{String: b.Comments, Valid: true}
	}
	var year sql.NullString
	if b.Year != "" {
		year = sql.NullString{String: b.Year, Valid: true}
	}
	var rate sql.NullInt32
	if b.Rate != 0 {
		rate = sql.NullInt32{Int32: int32(b.Rate), Valid: true}
	}
	var coverPath sql.NullString
	if b.CoverPath != "" {
		coverPath = sql.NullString{String: b.CoverPath, Valid: true}
	}

	tx.QueryRow(insertBookSql, b.Progress, author, synopsis, comments, year, rate, coverPath, taskId).Scan(&id)
	if err = tx.Commit(); err != nil {
		panic(err)
	}

	return id
}
