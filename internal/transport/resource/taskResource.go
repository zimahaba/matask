package resource

import "matask/internal/model"

type TaskResource struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Started Date   `json:"started"`
	Ended   Date   `json:"ended"`
	ChildId int    `json:"childId"`
}

type TaskPageResource struct {
	Tasks         []TaskResource `json:"tasks"`
	Page          int            `json:"page"`
	Size          int            `json:"size"`
	TotalPages    int            `json:"totalPages"`
	TotalElements int            `json:"totalElements"`
}

func FromTaskPageResult(result model.TaskPageResult) TaskPageResource {
	tasks := []TaskResource{}
	for i := 0; i < len(result.Tasks); i++ {
		t := result.Tasks[i]
		resource := TaskResource{Id: t.Task.Id, Name: t.Task.Name, Type: t.Task.Type, Started: Date{t.Task.Started}, Ended: Date{t.Task.Ended}, ChildId: t.ChildId}
		tasks = append(tasks, resource)
	}
	return TaskPageResource{
		Tasks:         tasks,
		Page:          result.Page,
		Size:          result.Size,
		TotalPages:    result.TotalPages,
		TotalElements: result.TotalElements,
	}
}
