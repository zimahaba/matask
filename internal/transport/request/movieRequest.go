package request

import (
	"log/slog"
	"matask/internal/model"
	"strconv"
	"time"
)

type MovieRequest struct {
	Name     string
	Started  Date
	Ended    Date
	Synopsis string
	Comments string
	Year     string
	Rate     int
	Director string
	Actors   []string
}

func (request MovieRequest) ToMovie() model.Movie {
	task := model.Task{
		Name:    request.Name,
		Type:    "movie",
		Started: request.Started.Time,
		Ended:   request.Ended.Time,
	}
	return model.Movie{
		Synopsis: request.Synopsis,
		Comments: request.Comments,
		Year:     request.Year,
		Rate:     request.Rate,
		Director: request.Director,
		Actors:   model.Actors{Actors: request.Actors},
		Task:     task,
	}
}

func ToMovie(query map[string][]string) (model.Movie, error) {
	var err error
	var name string
	if len(query["name"]) > 0 {
		name = query["name"][0]
	}
	var director string
	if len(query["director"]) > 0 {
		director = query["director"][0]
	}
	var year string
	if len(query["year"]) > 0 {
		year = query["year"][0]
	}
	var genre string
	if len(query["genre"]) > 0 {
		genre = query["genre"][0]
	}
	/*var actors []string
	if len(query["actors"]) > 0 {

	}*/
	var started time.Time
	if len(query["started"]) > 0 {
		startedValue := query["started"][0]
		if startedValue != "" {
			started, err = time.Parse(time.DateOnly, startedValue)
			if err != nil {
				slog.Error(err.Error())
				return model.Movie{}, err
			}
		}
	}
	var ended time.Time
	if len(query["ended"]) > 0 {
		endedValue := query["ended"][0]
		if endedValue != "" {
			ended, err = time.Parse(time.DateOnly, endedValue)
			if err != nil {
				slog.Error(err.Error())
				return model.Movie{}, err
			}
		}
	}
	var synopsis string
	if len(query["synopsis"]) > 0 {
		synopsis = query["synopsis"][0]
	}
	var comments string
	if len(query["comments"]) > 0 {
		comments = query["comments"][0]
	}
	var rate int
	if len(query["rate"]) > 0 {
		rate, err = strconv.Atoi(query["rate"][0])
		if err != nil {
			slog.Error(err.Error())
			return model.Movie{}, err
		}
	}
	task := model.Task{
		Name:    name,
		Type:    "movie",
		Started: started,
		Ended:   ended,
	}
	return model.Movie{
		Director: director,
		Synopsis: synopsis,
		Comments: comments,
		Year:     year,
		Rate:     rate,
		Genre:    genre,
		Task:     task,
	}, nil
}

func ToMovieFilter(query map[string][]string) model.MovieFilter {
	var filter model.MovieFilter
	if len(query["name"]) > 0 {
		filter.Name = query["name"][0]
	}
	if len(query["director"]) > 0 {
		filter.Director = query["director"][0]
	}
	if len(query["actor"]) > 0 {
		filter.Actor = query["actor"][0]
	}
	if len(query["year"]) > 0 {
		filter.Year = query["year"][0]
	}
	if len(query["page"]) > 0 {
		page, err := strconv.Atoi(query["page"][0])
		if err != nil {
			filter.Page = 1
		} else {
			filter.Page = page
		}
	} else {
		filter.Page = 1
	}
	if len(query["size"]) > 0 {
		size, err := strconv.Atoi(query["size"][0])
		if err != nil {
			filter.Size = 10
		} else {
			filter.Size = size
		}
	} else {
		filter.Size = 10
	}

	if len(query["sortField"]) > 0 {
		filter.SortField = query["sortField"][0]
	} else {
		filter.SortField = "id"
	}

	if len(query["sortDirection"]) > 0 {
		filter.SortDirection = query["sortDirection"][0]
	} else {
		filter.SortDirection = "ASC"
	}
	return filter
}
