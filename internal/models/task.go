package models

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
	Name     string
	Type     string
	Started1 pq.NullTime
	Started2 pq.NullTime
	Ended1   pq.NullTime
	Ended2   pq.NullTime
	Page     int
	Size     int
}
