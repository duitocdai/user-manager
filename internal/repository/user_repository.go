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
func (ur *InMemoryUserRepository) Create() {

}
func (ur *InMemoryUserRepository) FindByUUID() {

}
func (ur *InMemoryUserRepository) Update() {

}
func (ur *InMemoryUserRepository) Delete() {

}
func (ur *InMemoryUserRepository) FindByEmail() {

}
