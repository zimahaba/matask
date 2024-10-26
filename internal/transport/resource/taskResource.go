package resource

import "matask/internal/model"

type TaskResource struct {
	Id      int
	Name    string
	Type    string
	Started Date
	Ended   Date
	ChildId int
}

type TaskPageResource struct {
	Tasks         []TaskResource
	Page          int
	Size          int
	TotalPages    int
	TotalElements int
}

func FromTaskPageResult(result model.TaskPageResult) TaskPageResource {
	tasks := []TaskResource{}
	for i := 0; i < len(result.Tasks); i++ {
		t := result.Tasks[i]
		resource := TaskResource{Id: t.Task.Id, Name: t.Task.Name, Type: t.Task.Type, Started: Date{t.Task.Started}, Ended: Date{t.Task.Ended}}
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
