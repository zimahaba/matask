package request

import (
	"fmt"
	"matask/internal/model"
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
	fmt.Printf("dyn: %v.\n", request.DynamicFields)
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
