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
	Genre    string `json:"genre"`
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
		Genre:    b.Genre,
	}
}

type BookPageResource struct {
	Books         []BookResource `json:"books"`
	Page          int            `json:"page"`
	Size          int            `json:"size"`
	TotalPages    int            `json:"totalPages"`
	TotalElements int            `json:"totalElements"`
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
