package database

import (
	"database/sql"
	"log/slog"
	"matask/internal/model"
	"os"
	"strings"

	"github.com/lib/pq"
)

const findMovieTaskIdSql = "SELECT t.id FROM movie m INNER JOIN task t ON t.id = m.task_fk WHERE m.id = $1 AND t.user_fk = $2"
const deleteMovieSql = "DELETE FROM movie m USING task t WHERE t.id = m.task_fk AND m.id = $1 AND t.user_fk = $2"

var movieSortFieldMap = map[string]string{
	"id":       "m.id",
	"name":     "t.name",
	"director": "m.director",
	"year":     "m.year",
}

func FindFilteredMovies(f model.MovieFilter, db *sql.DB) (model.MoviePageResult, error) {
	query, countQuery := buildMovieFilteredQueries(f)

	var count int
	if err := db.QueryRow(countQuery, f.Name, f.Director, f.Year, f.UserId).Scan(&count); err != nil { // TODO: add actors
		slog.Error(err.Error())
		return model.MoviePageResult{}, err
	}

	offset := getOffset(f.Page, f.Size)

	rows, err := db.Query(query, f.Name, f.Director, f.Year, f.UserId, offset, f.Size) // TODO: add actors
	if err != nil {
		slog.Error(err.Error())
		return model.MoviePageResult{}, err
	}
	defer rows.Close()

	var movies []model.MovieProjection

	for rows.Next() {
		var movieProjection model.MovieProjection
		var director sql.NullString
		var year sql.NullString
		// TODO: add actors
		if err := rows.Scan(&movieProjection.Id, &movieProjection.Name, &director, &year); err != nil {
			slog.Error(err.Error())
			return model.MoviePageResult{}, err
		}
		if director.Valid {
			movieProjection.Director = director.String
		}
		if year.Valid {
			movieProjection.Year = year.String
		}
		movies = append(movies, movieProjection)
	}
	if err = rows.Err(); err != nil {
		slog.Error(err.Error())
		return model.MoviePageResult{}, err
	}

	totalPages := calculateTotalPages(count, f.Size)

	pageResult := model.MoviePageResult{Movies: movies, Page: f.Page, Size: f.Size, TotalPages: totalPages, TotalElements: count}
	return pageResult, nil
}

func buildMovieFilteredQueries(f model.MovieFilter) (string, string) {
	selectQuery := "SELECT m.id, t.name, m.director, m.year" // TODO: add actors
	selectCount := "SELECT COUNT(m.id)"
	from := `  		FROM movie m
					INNER JOIN task t ON t.id = m.task_fk
					WHERE ($1 = '' OR UPPER(t.name) like '%' || UPPER($1) || '%')
					AND ($2 = '' OR UPPER(m.director) like '%' || UPPER($2) || '%' )
					AND ($3 = '' OR UPPER(m.year) like '%' || UPPER($3) || '%' )
					AND t.user_fk = $4`
	order := getOrderQuery(f.SortField, f.SortDirection, movieSortFieldMap)
	offsetLimit := "OFFSET $5 LIMIT $6 "

	query := selectQuery + from + order + offsetLimit
	countQuery := selectCount + from

	return query, countQuery
}

func FindMovie(id int, userId int, db *sql.DB) (model.Movie, error) {
	query := `SELECT m.id, m.synopsis, m.comments, m.year, m.rate, m.director, m.actors, m.genre, m.poster_path, t.name, t.started, t.ended 
			FROM movie m
			INNER JOIN task t ON t.id = m.task_fk
			WHERE m.id = $1
			AND t.user_fk = $2`
	var m model.Movie
	var synopsis sql.NullString
	var comments sql.NullString
	var year sql.NullString
	var rate sql.NullInt32
	var director sql.NullString
	var genre sql.NullString
	var posterPath sql.NullString
	var started pq.NullTime
	var ended pq.NullTime
	if err := db.QueryRow(query, id, userId).Scan(&m.Id, &synopsis, &comments, &year, &rate, &director, &m.Actors, &genre, &posterPath, &m.Task.Name, &started, &ended); err != nil {
		slog.Error(err.Error())
		return model.Movie{}, err
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
	if genre.Valid {
		m.Genre = genre.String
	}
	if posterPath.Valid {
		m.PosterPath = posterPath.String
	}
	if started.Valid {
		m.Task.Started = started.Time
	}
	if ended.Valid {
		m.Task.Ended = ended.Time
	}
	return m, nil
}

func FindMoviePosterPath(id int, userId int, db *sql.DB) (string, error) {
	query := "SELECT m.poster_path FROM movie m INNER JOIN task t ON t.id = m.task_fk WHERE m.id = $1 AND t.user_fk = $2"
	var posterPath string
	if err := db.QueryRow(query, id, userId).Scan(&posterPath); err != nil {
		slog.Error(err.Error())
		return "", err
	}
	return posterPath, nil
}

func SaveOrUpdateMovie(m model.Movie, filebytes []byte, userId int, db *sql.DB) (int, error) {
	query := `INSERT INTO movie (synopsis, comments, year, rate, director, genre, actors, poster_path, task_fk) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`

	tx, err := db.Begin()
	if err != nil {
		slog.Error(err.Error())
		return -1, err
	}
	defer tx.Rollback()

	if m.Id != 0 {
		query = `UPDATE movie m 
				 SET synopsis = $3, comments = $4, year = $5, rate = $6, director = $7, genre = $8, actors = $9 
				 FROM task t 
				 WHERE t.id = m.task_fk AND m.id = $1 AND t.user_fk = $2`
		var taskId int
		if err := tx.QueryRow(findMovieTaskIdSql, m.Id, userId).Scan(&taskId); err != nil {
			slog.Error(err.Error())
			return -1, err
		}
		m.Task.Id = taskId
	}

	taskId, err := SaveOrUpdateTask(m.Task, userId, tx)
	if err != nil {
		slog.Error(err.Error())
		return -1, err
	}

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
	var genre sql.NullString
	if m.Genre != "" {
		genre = sql.NullString{String: m.Genre, Valid: true}
	}

	var posterPath sql.NullString
	basePath := os.Getenv("POSTER_PATH")
	if len(filebytes) > 0 {
		fullPath := basePath + strings.ReplaceAll(m.Task.Name, " ", "_")
		posterPath = sql.NullString{String: fullPath, Valid: true}
	}

	if m.Id == 0 {
		err = tx.QueryRow(query, synopsis, comments, year, rate, director, genre, m.Actors, posterPath, taskId).Scan(&id)
	} else {
		_, err = tx.Exec(query, m.Id, userId, synopsis, comments, year, rate, director, genre, m.Actors)
		id = m.Id
	}
	if err != nil {
		slog.Error(err.Error())
		return -1, err
	}

	// check if user dir exists, if not create dir
	if len(filebytes) > 0 {
		err = os.WriteFile(posterPath.String, filebytes, 0666)
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

func UpdateMoviePoster(id int, posterPath string, userId int, tx *sql.Tx) error {
	updateMoviePosterSql := "UPDATE movie SET poster_path = $3 FROM task t WHERE WHERE t.id = m.task_fk AND m.id = $1 AND t.user_fk = $2"
	_, err := tx.Exec(updateMoviePosterSql, id, posterPath, userId)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	return nil
}
