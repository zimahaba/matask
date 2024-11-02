package request

import (
	"log/slog"
	"matask/internal/model"
	"strconv"
	"time"
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

func ToBook(query map[string][]string) (model.Book, error) {
	var err error
	var name string
	if len(query["name"]) > 0 {
		name = query["name"][0]
	}
	var author string
	if len(query["author"]) > 0 {
		author = query["author"][0]
	}
	var year string
	if len(query["year"]) > 0 {
		year = query["year"][0]
	}
	/*var genre string
	if len(query["genre"]) > 0 {
		genre = query["genre"][0]
	}*/
	var started time.Time
	if len(query["started"]) > 0 {
		started, err = time.Parse(time.DateOnly, query["started"][0])
		if err != nil {
			slog.Error(err.Error())
			return model.Book{}, err
		}
	}
	var ended time.Time
	if len(query["ended"]) > 0 {
		ended, err = time.Parse(time.DateOnly, query["ended"][0])
		if err != nil {
			slog.Error(err.Error())
			return model.Book{}, err
		}
	}
	task := model.Task{
		Name:    name,
		Type:    "book",
		Started: started,
		Ended:   ended,
	}
	return model.Book{
		//Progress: request.Progress,
		Author: author,
		//Synopsis: request.Synopsis,
		//Comments: request.Comments,
		Year: year,
		//Rate:     request.Rate,
		//CoverPath: ,
		Task: task,
	}, nil
}

func ToBookFilter(query map[string][]string) model.BookFilter {
	var filter model.BookFilter
	if len(query["name"]) > 0 {
		filter.Name = query["name"][0]
	}
	if len(query["author"]) > 0 {
		filter.Author = query["author"][0]
	}
	if len(query["progress1"]) > 0 {
		p1, err := strconv.Atoi(query["progress1"][0])
		if err != nil {
			filter.Progress1 = -1
		} else {
			filter.Progress1 = p1
		}
	} else {
		filter.Progress1 = -1
	}
	if len(query["progress2"]) > 0 {
		p2, err := strconv.Atoi(query["progress2"][0])
		if err != nil {
			filter.Progress2 = -1
		} else {
			filter.Progress2 = p2
		}
	} else {
		filter.Progress2 = -1
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
