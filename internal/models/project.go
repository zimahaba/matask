package models

type Project struct {
	Id          int
	Description string
	Progress    int
	Task        Task
}
