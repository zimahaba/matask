package database

import (
	"database/sql"
	"log/slog"
	"matask/internal/model"
	"os"
	"strings"

	"github.com/lib/pq"
)

const findBookTaskIdSql = "SELECT t.id FROM book b INNER JOIN task t ON t.id = b.task_fk WHERE b.id = $1 AND t.user_fk = $2"
const deleteBookSql = "DELETE FROM book b USING task t WHERE t.id = b.task_fk AND b.id = $1 AND t.user_fk = $2"

var bookSortFieldMap = map[string]string{
	"id":       "b.id",
	"name":     "t.name",
	"author":   "b.author",
	"progress": "b.progress",
}

func FindFilteredBooks(f model.BookFilter, db *sql.DB) (model.BookPageResult, error) {
	query, countQuery := buildBookFilteredQueries(f)

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

	offset := getOffset(f.Page, f.Size)

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

	totalPages := calculateTotalPages(count, f.Size)

	pageResult := model.BookPageResult{Books: books, Page: f.Page, Size: f.Size, TotalPages: totalPages, TotalElements: count}
	return pageResult, nil
}

func buildBookFilteredQueries(f model.BookFilter) (string, string) {
	selectQuery := "SELECT b.id, t.name, b.author, b.progress"
	selectCount := "SELECT COUNT(b.id)"
	from := `		FROM book b
					INNER JOIN task t ON t.id = b.task_fk
					WHERE ($1 = '' OR UPPER(t.name) like '%' || UPPER($1) || '%')
					AND ($2 = '' OR UPPER(b.author) like '%' || UPPER($2) || '%' )
					AND (b.progress >= COALESCE(CAST($3 AS INTEGER), 0))
					AND (b.progress <= COALESCE(CAST($4 AS INTEGER), 100))
					AND t.user_fk = $5`
	order := getOrderQuery(f.SortField, f.SortDirection, bookSortFieldMap)
	offsetLimit := "OFFSET $6 LIMIT $7 "

	query := selectQuery + from + order + offsetLimit
	countQuery := selectCount + from

	return query, countQuery
}

func FindBook(id int, userId int, db *sql.DB) (model.Book, error) {
	query := `SELECT b.id, b.progress, b.author, b.synopsis, b.comments, b.year, b.rate, b.genre, b.cover_path, t.name, t.started, t.ended
			  FROM book b
			  INNER JOIN task t ON t.id = b.task_fk
			  WHERE b.id = $1 
			  AND t.user_fk = $2`

	var b model.Book
	var synopsis sql.NullString
	var comments sql.NullString
	var year sql.NullString
	var rate sql.NullInt32
	var genre sql.NullString
	var coverPath sql.NullString
	var started pq.NullTime
	var ended pq.NullTime
	if err := db.QueryRow(query, id, userId).Scan(&b.Id, &b.Progress, &b.Author, &synopsis, &comments, &year, &rate, &genre, &coverPath, &b.Task.Name, &started, &ended); err != nil {
		slog.Error(err.Error())
		return model.Book{}, err
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
	if genre.Valid {
		b.Genre = genre.String
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
	query := "SELECT b.cover_path FROM book b INNER JOIN task t ON t.id = b.task_fk WHERE b.id = $1 AND t.user_fk = $2"

	var coverPath string
	if err := db.QueryRow(query, id, userId).Scan(&coverPath); err != nil {
		slog.Error(err.Error())
		return "", err
	}
	return coverPath, nil
}

func SaveOrUpdateBook(b model.Book, filebytes []byte, userId int, db *sql.DB) (int, error) {
	query := `INSERT INTO book (progress, author, synopsis, comments, year, genre, rate, cover_path, task_fk) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`

	tx, err := db.Begin()
	if err != nil {
		slog.Error(err.Error())
		return -1, err
	}
	defer tx.Rollback()

	if b.Id != 0 {
		query = `UPDATE book b
				 SET progress = $3, author = $4, synopsis = $5, comments = $6, year = $7, genre = $8, rate = $9
				 FROM task t 
				 WHERE t.id = b.task_fk AND b.id = $1 AND t.user_fk = $2`
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
	var genre sql.NullString
	if b.Genre != "" {
		genre = sql.NullString{String: b.Genre, Valid: true}
	}
	var rate sql.NullInt32
	if b.Rate != 0 {
		rate = sql.NullInt32{Int32: int32(b.Rate), Valid: true}
	}

	var coverPath sql.NullString
	basePath := os.Getenv("COVER_PATH")
	if len(filebytes) > 0 {
		fullPath := basePath + strings.ReplaceAll(b.Task.Name, " ", "_")
		coverPath = sql.NullString{String: fullPath, Valid: true}
	}
	if b.Id == 0 {
		err = tx.QueryRow(query, b.Progress, author, synopsis, comments, year, genre, rate, coverPath, taskId).Scan(&id)
	} else {
		_, err = tx.Exec(query, b.Id, userId, b.Progress, author, synopsis, comments, year, genre, rate)
		id = b.Id
	}
	if err != nil {
		slog.Error(err.Error())
		return -1, err
	}

	// check if user dir exists, if not create dir
	if len(filebytes) > 0 {
		err = os.WriteFile(coverPath.String, filebytes, 0666)
		if err != nil {
			slog.Error(err.Error())
			return -1, err
		}
	}

	if err = tx.Commit(); err != nil {
		slog.Error(err.Error())
		return -1, err
	}

	return id, nil
}

func UpdateBookCover(id int, coverPath string, userId int, tx *sql.Tx) error {
	updateBookCoverSql := "UPDATE book b SET cover_path = $3 FROM task t WHERE t.id = b.task_fk AND b.id = $1 AND t.user_fk = $2"
	_, err := tx.Exec(updateBookCoverSql, id, userId, coverPath)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	return nil
}
