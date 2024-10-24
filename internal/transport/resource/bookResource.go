package resource

import (
	"matask/internal/model"
)

type BookResource struct {
	Id       int
	Name     string
	Started  Date
	Ended    Date
	Progress int
	Author   string
	Synopsis string
	Comments string
	Year     string
	Rate     int
	//CoverImage
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
		//CoverImage
	}
}
