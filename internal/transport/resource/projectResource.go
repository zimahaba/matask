package resource

import (
	"matask/internal/model"
)

type ProjectResource struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	Started       Date   `json:"started"`
	Ended         Date   `json:"ended"`
	Description   string `json:"description"`
	Progress      int    `json:"progress"`
	DynamicFields map[string]interface{}
}

func FromProject(p model.Project) ProjectResource {
	return ProjectResource{
		Id:            p.Id,
		Name:          p.Task.Name,
		Started:       Date{p.Task.Started},
		Ended:         Date{p.Task.Ended},
		Description:   p.Description,
		Progress:      p.Progress,
		DynamicFields: p.DynamicFields,
	}
}

type ProjectPageResource struct {
	Projects      []ProjectResource `json:"projects"`
	Page          int               `json:"page"`
	Size          int               `json:"size"`
	TotalPages    int               `json:"totalPages"`
	TotalElements int               `json:"totalElements"`
}

func FromProjectPageResult(result model.ProjectPageResult) ProjectPageResource {
	projects := []ProjectResource{}
	for i := 0; i < len(result.Projects); i++ {
		b := result.Projects[i]
		resource := ProjectResource{Id: b.Id, Name: b.Name, Progress: b.Progress}
		projects = append(projects, resource)
	}
	return ProjectPageResource{
		Projects:      projects,
		Page:          result.Page,
		Size:          result.Size,
		TotalPages:    result.TotalPages,
		TotalElements: result.TotalElements,
	}
}
