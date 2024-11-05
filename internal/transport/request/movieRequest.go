package request

import (
	"matask/internal/model"
	"strconv"
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
