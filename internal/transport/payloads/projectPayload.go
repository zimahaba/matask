package payload

import (
	"fmt"
	"matask/internal/models"
)

type ProjectPayload struct {
	Name        string
	Type        string
	Started     MataskTime
	Ended       MataskTime
	Description string
	Progress    int
}

func ToProject(payload ProjectPayload) models.Project {
	fmt.Printf("payload: %v", payload)
	task := models.Task{
		Name: payload.Name,
		Type: payload.Type,
	}
	if !payload.Started.Time.IsZero() {
		task.Started = payload.Started.Time
	}
	if !payload.Ended.Time.IsZero() {
		task.Ended = payload.Ended.Time
	}
	return models.Project{
		Description: payload.Description,
		Progress:    payload.Progress,
		Task:        task,
	}
}
