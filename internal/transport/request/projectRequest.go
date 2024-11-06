package request

import (
	"matask/internal/model"
	"strconv"
)

type ProjectRequest struct {
	Name          string
	Started       Date
	Ended         Date
	Description   string
	Progress      int
	DynamicFields map[string]interface{}
}

func (request ProjectRequest) ToProject() model.Project {
	task := model.Task{
		Name:    request.Name,
		Type:    "project",
		Started: request.Started.Time,
		Ended:   request.Ended.Time,
	}
	return model.Project{
		Description:   request.Description,
		Progress:      request.Progress,
		DynamicFields: request.DynamicFields,
		Task:          task,
	}
}

func ToProjectFilter(query map[string][]string) model.ProjectFilter {
	var filter model.ProjectFilter
	if len(query["name"]) > 0 {
		filter.Name = query["name"][0]
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
