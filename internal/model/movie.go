package model

type Movie struct {
	Id         int
	Synopsis   string
	Comments   string
	Year       string
	Rate       int
	Actors     string
	Director   string
	PosterPath string
	Task       Task
}
