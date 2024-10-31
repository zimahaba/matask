package resource

import (
	"matask/internal/model"
)

type BookResource struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Started  Date   `json:"started"`
	Ended    Date   `json:"ended"`
	Progress int    `json:"progress"`
	Author   string `json:"author"`
	Synopsis string `json:"synopsis"`
	Comments string `json:"comments"`
	Year     string `json:"year"`
	Rate     int    `json:"rate"`
}

func FromBook(b model.Book) BookResource {
	return BookResource{
		Id:       b.Id,
		Name:     b.Task.Name,
		Started:  Date{b.Task.Started},
		Ended:    Date{b.Task.Ended},
		Progress: b.Progress,
		Author:   b.Author,
		Synopsis: b.Synopsis,
		Comments: b.Comments,
		Year:     b.Year,
		Rate:     b.Rate,
	}
}
