package model

import (
	"time"

	"github.com/lib/pq"
)

type Task struct {
	Id      int
	Name    string
	Type    string
	Started time.Time
	Ended   time.Time
	Created time.Time
	User    MataskUser
}

type TaskFilter struct {
	Name          string
	Type          string
	Started1      pq.NullTime
	Started2      pq.NullTime
	Ended1        pq.NullTime
	Ended2        pq.NullTime
	Page          int
	Size          int
	SortField     string
	SortDirection string
	UserId        int
}

type TaskProjection struct {
	Task    Task
	ChildId int
}

type TaskPageResult struct {
	Tasks         []TaskProjection
	Page          int
	Size          int
	TotalPages    int
	TotalElements int
}
