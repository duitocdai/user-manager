package service

import (
	"log"

	repository "github.com/duitocdai/internal/repository"
)

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}
func (us *userService) GetAllUser() {
	us.repo.FindAll()
	log.Println("getALLuSER INTO userService")
}
func (us *userService) CreateUser() {

}
func (us *userService) GetUserByUUID() {

}
func (us *userService) UpdateUser() {

}
func (us *userService) DeleteUser() {

}
