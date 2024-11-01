package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"matask/internal/model"

	"github.com/lib/pq"
)

var findFilteredBooksSelectSql = "SELECT b.id, t.name, b.author, b.progress"
var findFilteredBooksSelectCountSql = "SELECT COUNT(b.id)"
var findFilteredBooksFromWhereSql = `
		FROM book b
		INNER JOIN task t ON t.id = b.task_fk
		WHERE ($1 = '' OR UPPER(t.name) like '%' || UPPER($1) || '%')
		AND ($2 = '' OR UPPER(b.author) like '%' || UPPER($2) || '%' )
		AND (b.progress >= COALESCE(CAST($3 AS INTEGER), 0))
		AND (b.progress <= COALESCE(CAST($4 AS INTEGER), 100))
		AND t.user_fk = $5
	`
var findFilteredBooksOrderSql = " ORDER BY %s %s "
var findFilteredBooksPageSql = " OFFSET $6 LIMIT $7 "

const (
	findBookSql = `
		SELECT b.id, b.progress, b.author, b.synopsis, b.comments, b.year, b.rate, b.cover_path, t.name, t.started, t.ended 
		FROM book b
		INNER JOIN task t ON t.id = b.task_fk
		WHERE b.id = $1 
		AND t.user_fk = $2
	`
	findBookTaskIdSql    = "SELECT t.id FROM book b INNER JOIN task t ON t.id = b.task_fk WHERE b.id = $1 AND t.user_fk = $2"
	findBookCoverPathSql = "SELECT b.cover_path FROM book b INNER JOIN task t ON t.id = b.task_fk WHERE b.id = $1 AND t.user_fk = $2"
	insertBookSql        = "INSERT INTO book (progress, author, synopsis, comments, year, rate, task_fk) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
	updateBookSql        = `
		UPDATE book b
		SET progress = $3, author = $4, synopsis = $5, comments = $6, year = $7, rate = $8
		FROM task t 
		WHERE t.id = b.task_fk AND b.id = $1 AND t.user_fk = $2
		`
	updateBookCoverSql = "UPDATE book b SET cover_path = $3 FROM task t WHERE t.id = b.task_fk AND b.id = $1 AND t.user_fk = $2"
	deleteBookSql      = "DELETE FROM book b USING task t WHERE t.id = b.task_fk AND b.id = $1 AND t.user_fk = $2"
)

func FindFilteredBooks(f model.BookFilter, db *sql.DB) (model.BookPageResult, error) {
	sortField := sortFieldMap[f.SortField]
	if sortField == "" {
		sortField = "id"
	}
	sortDirection := sortDirectionMap[f.SortDirection]
	if sortDirection == "" {
		sortDirection = "ASC"
	}
	query := findFilteredBooksSelectSql + findFilteredBooksFromWhereSql + fmt.Sprintf(findFilteredBooksOrderSql, sortField, sortDirection) + findFilteredBooksPageSql
	countQuery := findFilteredBooksSelectCountSql + findFilteredBooksFromWhereSql

	var count int
	progress1 := sql.NullInt64{}
	if f.Progress1 >= 0 {
		progress1 = sql.NullInt64{Int64: int64(f.Progress1), Valid: true}
	}
	progress2 := sql.NullInt64{}
	if f.Progress2 >= 0 {
		progress2 = sql.NullInt64{Int64: int64(f.Progress2), Valid: true}
	}
	if err := db.QueryRow(countQuery, f.Name, f.Author, progress1, progress2, f.UserId).Scan(&count); err != nil {
		slog.Error(err.Error())
		return model.BookPageResult{}, err
	}

	var offset int
	if f.Page <= 1 {
		offset = 0
	} else {
		offset = (f.Page - 1) * f.Size
	}

	rows, err := db.Query(query, f.Name, f.Author, progress1, progress2, f.UserId, offset, f.Size)
	if err != nil {
		slog.Error(err.Error())
		return model.BookPageResult{}, err
	}
	defer rows.Close()

	var books []model.BookProjection

	for rows.Next() {
		var bookProjection model.BookProjection
		var author sql.NullString
		if err := rows.Scan(&bookProjection.Id, &bookProjection.Name, &author, &bookProjection.Progress); err != nil {
			slog.Error(err.Error())
			return model.BookPageResult{}, err
		}
		if author.Valid {
			bookProjection.Author = author.String
		}
		books = append(books, bookProjection)
	}
	if err = rows.Err(); err != nil {
		slog.Error(err.Error())
		return model.BookPageResult{}, err
	}

	totalPages := count / f.Size
	remainder := count % f.Size
	if totalPages == 0 || (totalPages > 0 && remainder > 0) {
		totalPages++
	}

	pageResult := model.BookPageResult{Books: books, Page: f.Page, Size: f.Size, TotalPages: totalPages, TotalElements: count}
	return pageResult, nil
}

func FindBook(id int, userId int, db *sql.DB) (model.Book, error) {
	var b model.Book
	var author sql.NullString
	var synopsis sql.NullString
	var comments sql.NullString
	var year sql.NullString
	var rate sql.NullInt32
	var coverPath sql.NullString
	var started pq.NullTime
	var ended pq.NullTime
	if err := db.QueryRow(findBookSql, id, userId).Scan(&b.Id, &b.Progress, &author, &synopsis, &comments, &year, &rate, &coverPath, &b.Task.Name, &started, &ended); err != nil {
		slog.Error(err.Error())
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

func FindBookCoverPath(id int, userId int, db *sql.DB) (string, error) {
	var coverPath string
	if err := db.QueryRow(findBookCoverPathSql, id, userId).Scan(&coverPath); err != nil {
		slog.Error(err.Error())
		return "", err
	}
	return coverPath, nil
}

func SaveOrUpdateBook(b model.Book, userId int, db *sql.DB) (int, error) {
	tx, err := db.Begin()
	if err != nil {
		slog.Error(err.Error())
		return -1, err
	}
	defer tx.Rollback()

	if b.Id != 0 {
		var taskId int
		if err := tx.QueryRow(findBookTaskIdSql, b.Id, userId).Scan(&taskId); err != nil {
			slog.Error(err.Error())
			return -1, err
		}
		b.Task.Id = taskId
	}

	taskId, err := SaveOrUpdateTask(b.Task, userId, tx)
	if err != nil {
		slog.Error(err.Error())
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
	} else {
		_, err = tx.Exec(updateBookSql, b.Id, userId, b.Progress, author, synopsis, comments, year, rate)
		id = b.Id
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

func UpdateBookCover(id int, coverPath string, userId int, tx *sql.Tx) error {
	_, err := tx.Exec(updateBookCoverSql, id, userId, coverPath)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	return nil
}
