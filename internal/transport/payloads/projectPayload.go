package payload

import (
	"matask/internal/models"
)

type ProjectPayload struct {
	Name        string
	Type        string
	Started     Date
	Ended       Date
	Description string
	Progress    int
}

func ToProject(payload ProjectPayload) models.Project {
	task := models.Task{
		Name: payload.Name,
		Type: payload.Type,
	}
	return models.Project{
		Description: payload.Description,
		Progress:    payload.Progress,
		Task:        task,
	}
}
