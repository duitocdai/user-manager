package repository

import (
	"log"

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
func (ur *InMemoryUserRepository) FindAll() {
	log.Println("find all repository ")
}
func (ur *InMemoryUserRepository) Create(user model.User) error {
	ur.users = append(ur.users, user)
	return nil
}
func (ur *InMemoryUserRepository) FindByUUID() {

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
