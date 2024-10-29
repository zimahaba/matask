package resource

import "matask/internal/model"

type UserResource struct {
	Username string `json:"username"`
}

func FromUser(user model.MataskUser) UserResource {
	return UserResource{Username: user.Email}
}
