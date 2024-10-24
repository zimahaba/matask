package model

type Book struct {
	Id        int
	Progress  int
	Author    string
	Synopsis  string
	Comments  string
	Year      string
	Rate      int
	CoverPath string
	Task      Task
}
