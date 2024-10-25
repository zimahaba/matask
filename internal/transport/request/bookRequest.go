package request

import (
	"matask/internal/model"
)

type BookRequest struct {
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

func (request BookRequest) ToBook() model.Book {
	task := model.Task{
		Name:    request.Name,
		Type:    "book",
		Started: request.Started.Time,
		Ended:   request.Ended.Time,
	}
	return model.Book{
		Progress: request.Progress,
		Author:   request.Author,
		Synopsis: request.Synopsis,
		Comments: request.Comments,
		Year:     request.Year,
		Rate:     request.Rate,
		//CoverPath: ,
		Task: task,
	}
}
