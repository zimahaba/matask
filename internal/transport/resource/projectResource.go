package resource

import (
	"matask/internal/model"
)

type ProjectResource struct {
	Id          int
	Name        string
	Started     Date
	Ended       Date
	Description string
	Progress    int
}

func FromProject(p model.Project) ProjectResource {
	return ProjectResource{
		Id:          p.Id,
		Name:        p.Task.Name,
		Started:     Date{p.Task.Started},
		Ended:       Date{p.Task.Ended},
		Description: p.Description,
		Progress:    p.Progress,
	}
}
