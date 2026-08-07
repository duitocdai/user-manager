package repository

import model "github.com/duitocdai/internal/models"

type UserRepository interface {
	FindAll() ([]model.User, error)
	Create(user model.User) error
	FindByUUID(uuid string) (model.User, bool)
	Update(uuid string, user model.User) error
	Delete(uuid string) error
	FindByEmail(email string) (model.User, bool)
}
