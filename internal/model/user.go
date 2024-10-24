package model

type MataskUser struct {
	Id    int
	Name  string
	Email string
}

type UserCredentials struct {
	Id       int
	Username string
	Password string
}
