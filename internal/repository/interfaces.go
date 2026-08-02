package repository

import model "github.com/duitocdai/internal/models"

type UserRepository interface {
	FindAll()
	Create(user model.User) error
	FindByUUID()
	Update()
	Delete()
	FindByEmail(email string) (model.User, bool)
}
