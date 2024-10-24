package resource

import (
	"matask/internal/models"
)

type ProjectResource struct {
	Id          int
	Name        string
	Description string
	Progress    int
	Started     Date
	Ended       Date
}

func FromProject(p models.Project) ProjectResource {
	return ProjectResource{
		Id:          p.Id,
		Name:        p.Task.Name,
		Description: p.Description,
		Progress:    p.Progress,
		Started:     Date{p.Task.Started},
		Ended:       Date{p.Task.Ended},
	}
}
