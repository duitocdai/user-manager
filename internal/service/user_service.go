package service

import (
	"log"

	model "github.com/duitocdai/internal/models"
	repository "github.com/duitocdai/internal/repository"
	"github.com/duitocdai/internal/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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
func (us *userService) CreateUser(user model.User) (model.User, error) {
	user.Email = utils.NormalizeString(user.Email)

	if _, exist := us.repo.FindByEmail(user.Email); exist {
		return model.User{}, utils.NewError("email already exists", utils.ErrCodeConfict)
	}

	user.UUID = uuid.New().String()

	hassedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, utils.WrapError(err, "failed to hash password", utils.ErrCodeInternal)
	}
	user.Password = string(hassedPassword)

	if err := us.repo.Create(user); err != nil {
		return model.User{}, utils.WrapError(err, "failed tp create user", utils.ErrCodeInternal)
	}
	return user, nil
}
func (us *userService) GetUserByUUID() {

}
func (us *userService) UpdateUser() {

}
func (us *userService) DeleteUser() {

}
