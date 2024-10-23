package models

import (
	"time"
)

type Task struct {
	Id      int
	Name    string
	Type    string
	Started time.Time
	Ended   time.Time
	Created time.Time
	Updated time.Time
	User    MataskUser
}
