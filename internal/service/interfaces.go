package service

type UserService interface {
	GetAllUser()
	CreateUser()
	GetUserByUUID()
	UpdateUser()
	DeleteUser()
}
