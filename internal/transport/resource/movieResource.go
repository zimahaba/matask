package resource

import (
	"matask/internal/model"
)

type MovieResource struct {
	Id       int
	Name     string
	Started  Date
	Ended    Date
	Synopsis string
	Comments string
	Year     string
	Rate     int
	Director string
	//CoverImage
}

func FromMovie(m model.Movie) MovieResource {
	return MovieResource{
		Id:       m.Id,
		Name:     m.Task.Name,
		Started:  Date{m.Task.Started},
		Ended:    Date{m.Task.Ended},
		Synopsis: m.Synopsis,
		Comments: m.Comments,
		Year:     m.Year,
		Rate:     m.Rate,
		Director: m.Director,
		//CoverImage
	}
}
