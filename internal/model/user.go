package model

import "time"

type MataskUser struct {
	Id          int
	Name        string
	Email       string
	Birthday    time.Time
	Credentials UserCredentials
}

type UserCredentials struct {
	Id       int
	Username string
	Password string
}
