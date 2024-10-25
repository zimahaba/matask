package request

import (
	"matask/internal/model"
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
