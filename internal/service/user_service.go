package service

import (
	"strings"

	"github.com/duitocdai/internal/models"
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
func (us *userService) GetAllUser(search string, page, limit int) ([]models.User, error) {
	users, err := us.repo.FindAll()
	if err != nil {
		return nil, utils.WrapError(err, "failed to fetch users", utils.ErrCodeInternal)
	}
	var filterUsers []model.User
	if search != "" {
		search = strings.ToLower(search)
		for _, user := range users {
			name := strings.ToLower(user.Name)
			email := strings.ToLower(user.Email)
			if strings.Contains(name, search) || strings.Contains(email, search) {
				filterUsers = append(filterUsers, user)
			}
		}
	} else {
		filterUsers = users
	}
	start := (page - 1) * limit
	if start >= len(filterUsers) {
		return []models.User{}, nil
	}

	end := start + limit
	if end > len(filterUsers) {
		end = len(filterUsers)
	}
	return filterUsers[start:end], nil
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
func (us *userService) GetUserByUUID(uuid string) (model.User, error) {
	user, found := us.repo.FindByUUID(uuid)
	if !found {
		return model.User{}, utils.NewError("user not found", utils.ErrCodeNotFound)
	}
	return user, nil
}
func (us *userService) UpdateUser() {

}
func (us *userService) DeleteUser() {

}
