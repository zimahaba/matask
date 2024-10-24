package request

import (
	"matask/internal/model"
)

type ProjectRequest struct {
	Name        string
	Type        string
	Started     Date
	Ended       Date
	Description string
	Progress    int
}

func (request ProjectRequest) ToProject() model.Project {
	task := model.Task{
		Name: request.Name,
		Type: request.Type,
	}
	return model.Project{
		Description: request.Description,
		Progress:    request.Progress,
		Task:        task,
	}
}
