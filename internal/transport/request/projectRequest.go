package request

import (
	"matask/internal/models"
)

type ProjectRequest struct {
	Name        string
	Type        string
	Started     Date
	Ended       Date
	Description string
	Progress    int
}

func (request ProjectRequest) ToProject() models.Project {
	task := models.Task{
		Name: request.Name,
		Type: request.Type,
	}
	return models.Project{
		Description: request.Description,
		Progress:    request.Progress,
		Task:        task,
	}
}
