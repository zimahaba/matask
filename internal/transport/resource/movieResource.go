package resource

import (
	"matask/internal/model"
)

type MovieResource struct {
	Id       int      `json:"id"`
	Name     string   `json:"name"`
	Started  Date     `json:"started"`
	Ended    Date     `json:"ended"`
	Synopsis string   `json:"synopsis"`
	Comments string   `json:"comments"`
	Year     string   `json:"year"`
	Rate     int      `json:"rate"`
	Director string   `json:"director"`
	Actors   []string `json:"actors"`
	Genre    string   `json:"genre"`
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
		Actors:   m.Actors.Actors,
		Genre:    m.Genre,
	}
}

type MoviePageResource struct {
	Movies        []MovieResource `json:"movies"`
	Page          int             `json:"page"`
	Size          int             `json:"size"`
	TotalPages    int             `json:"totalPages"`
	TotalElements int             `json:"totalElements"`
}

func FromMoviePageResult(result model.MoviePageResult) MoviePageResource {
	movies := []MovieResource{}
	for i := 0; i < len(result.Movies); i++ {
		m := result.Movies[i]
		resource := MovieResource{Id: m.Id, Name: m.Name, Director: m.Director, Actors: m.Actors, Year: m.Year}
		movies = append(movies, resource)
	}
	return MoviePageResource{
		Movies:        movies,
		Page:          result.Page,
		Size:          result.Size,
		TotalPages:    result.TotalPages,
		TotalElements: result.TotalElements,
	}
}
