package service

import (
	"github.com/duitocdai/internal/models"
	model "github.com/duitocdai/internal/models"
)

type UserService interface {
	GetAllUser(search string, page, limit int) ([]models.User, error)
	CreateUser(user model.User) (model.User, error)
	GetUserByUUID(uuid string) (model.User, error)
	UpdateUser(uuid string, user model.User) (model.User, error)
	DeleteUser(uuid string) error
}
