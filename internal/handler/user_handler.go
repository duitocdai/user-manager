package handler

import (
	"log"

	"github.com/duitocdai/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (uh *UserHandler) GetAllUser(ctx *gin.Context) {
	log.Println("get all user")
	uh.service.GetAllUser()
}
func (uh *UserHandler) CreateUser(ctx *gin.Context) {

}
func (uh *UserHandler) GetUserByUUID(ctx *gin.Context) {

}
func (uh *UserHandler) UpdateUser(ctx *gin.Context) {

}
func (uh *UserHandler) DeleteUser(ctx *gin.Context) {

}
