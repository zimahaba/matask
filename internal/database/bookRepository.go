package database

import (
	"database/sql"
	"matask/internal/model"

	"github.com/lib/pq"
)

const (
	findBookSql = `
		SELECT b.id, b.progress, b.author, b.synopsis, b.comments, b.year, b.rate, b.cover_path, t.name, t.started, t.ended 
		FROM book b
		INNER JOIN task t ON t.id = b.task_fk
		WHERE b.id = $1
	`
	findBookTaskIdSql  = "SELECT t.id FROM book b INNER JOIN task t ON t.id = b.task_fk WHERE b.id = $1"
	insertBookSql      = "INSERT INTO book (progress, author, synopsis, comments, year, rate, task_fk) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
	updateBookSql      = "UPDATE book SET progress = $2, author = $3, synopsis = $4, comments = $5, year = $6, rate = $7 WHERE id = $1"
	updateBookCoverSql = "UPDATE book SET cover_path = $2 WHERE id = $1"
	deleteBookSql      = "DELETE FROM book WHERE id = $1"
)

func FindBook(id int, db *sql.DB) (model.Book, error) {
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
		return model.Book{}, err
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
	return b, nil
}

func SaveOrUpdateBook(b model.Book, userId int, db *sql.DB) (int, error) {
	tx, err := db.Begin()
	if err != nil {
		return -1, err
	}
	defer tx.Rollback()

	if b.Id != 0 {
		var taskId int
		if err := tx.QueryRow(findBookTaskIdSql, b.Id).Scan(&taskId); err != nil {
			return -1, err
		}
		b.Task.Id = taskId
	}

	taskId, err := SaveOrUpdateTask(b.Task, userId, tx)
	if err != nil {
		return -1, err
	}

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

	if b.Id == 0 {
		err = tx.QueryRow(insertBookSql, b.Progress, author, synopsis, comments, year, rate, taskId).Scan(&id)
		if err != nil {
			return -1, err
		}
	} else {
		_, err = tx.Exec(updateBookSql, b.Id, b.Progress, author, synopsis, comments, year, rate)
		if err != nil {
			return -1, err
		}
		id = b.Id
	}

	if err = tx.Commit(); err != nil {
		return -1, err
	}

	return id, nil
}

func UpdateBookCover(id int, coverPath string, tx *sql.Tx) error {
	_, err := tx.Exec(updateBookCoverSql, id, coverPath)
	return err
}
