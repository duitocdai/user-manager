package service

import model "github.com/duitocdai/internal/models"

type UserService interface {
	GetAllUser()
	CreateUser(user model.User) (model.User, error)
	GetUserByUUID()
	UpdateUser()
	DeleteUser()
}
