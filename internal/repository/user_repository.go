package repository

import (
	model "github.com/duitocdai/internal/models"
)

type InMemoryUserRepository struct {
	users []model.User
}

func NewInMemoryUserRepository() UserRepository {
	return &InMemoryUserRepository{
		users: make([]model.User, 0),
	}
}
func (ur *InMemoryUserRepository) FindAll() ([]model.User, error) {
	return ur.users, nil
}
func (ur *InMemoryUserRepository) Create(user model.User) error {
	ur.users = append(ur.users, user)
	return nil
}
func (ur *InMemoryUserRepository) FindByUUID(uuid string) (model.User, bool) {
	for _, user := range ur.users {
		if user.UUID == uuid {
			return user, true
		}
	}
	return model.User{}, false
}
func (ur *InMemoryUserRepository) Update() {

}
func (ur *InMemoryUserRepository) Delete() {

}
func (ur *InMemoryUserRepository) FindByEmail(email string) (model.User, bool) {
	for _, user := range ur.users {
		if user.Email == email {
			return user, true
		}
	}

	return model.User{}, false
}
