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

type BookFilter struct {
	Name          string
	Author        string
	Progress1     int
	Progress2     int
	Page          int
	Size          int
	SortField     string
	SortDirection string
	UserId        int
}

type BookPageResult struct {
	Books         []BookProjection
	Page          int
	Size          int
	TotalPages    int
	TotalElements int
}

type BookProjection struct {
	Id       int
	Name     string
	Author   string
	Progress int
}
