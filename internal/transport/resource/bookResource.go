package resource

import (
	"matask/internal/model"
)

type BookResource struct {
	Id       int    `json:"id"omm`
	Name     string `json:"name"`
	Started  Date   `json:"started,omitempty"`
	Ended    Date   `json:"ended,omitempty"`
	Progress int    `json:"progress"`
	Author   string `json:"author"`
	Synopsis string `json:"synopsis,omitempty"`
	Comments string `json:"comments,omitempty"`
	Year     string `json:"year,omitempty"`
	Rate     int    `json:"rate,omitempty"`
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

type BookPageResource struct {
	Books         []BookResource
	Page          int
	Size          int
	TotalPages    int
	TotalElements int
}

func FromBookPageResult(result model.BookPageResult) BookPageResource {
	books := []BookResource{}
	for i := 0; i < len(result.Books); i++ {
		b := result.Books[i]
		resource := BookResource{Id: b.Id, Name: b.Name, Author: b.Author, Progress: b.Progress}
		books = append(books, resource)
	}
	return BookPageResource{
		Books:         books,
		Page:          result.Page,
		Size:          result.Size,
		TotalPages:    result.TotalPages,
		TotalElements: result.TotalElements,
	}
}
