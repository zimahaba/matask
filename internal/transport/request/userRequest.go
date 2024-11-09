package request

import "matask/internal/model"

type CredentialsRequest struct {
	Username     string
	Password     string
	KeepLoggedIn bool
}

type UserRequest struct {
	Name     string
	Email    string
	Password string
}

func (request UserRequest) ToUser(password string) model.MataskUser {
	creds := model.UserCredentials{Username: request.Email, Password: password}
	return model.MataskUser{
		Name:        request.Name,
		Email:       request.Email,
		Credentials: creds,
	}
}
