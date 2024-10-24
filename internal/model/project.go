package model

type Project struct {
	Id          int
	Description string
	Progress    int
	Task        Task
}
