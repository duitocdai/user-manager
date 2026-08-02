package handler

import (
	"net/http"

	"github.com/duitocdai/internal/dto"
	model "github.com/duitocdai/internal/models"
	"github.com/duitocdai/internal/service"
	"github.com/duitocdai/internal/utils"
	"github.com/duitocdai/internal/validation"
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

}
func (uh *UserHandler) CreateUser(ctx *gin.Context) {
	// var input dto.CreateUserInput
	// if err := ctx.ShouldBindJSON(&input); err != nil {
	// 	utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
	// 	return
	// }

	// user := input.MapCreateInputToModel()

	// createdUser, err := uh.service.CreateUser(user)
	// if err != nil {
	// 	utils.ResponseError(ctx, err)
	// 	return
	// }

	// userDTO := dto.MapUserToDTO(createdUser)

	// utils.ResponseSuccess(ctx, http.StatusCreated, &userDTO)
	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, validation.HandleValidationErrors(err))
		return
	}
	createdUser, err := uh.service.CreateUser(user)
	if err != nil {
		utils.ReponseError(ctx, err)
		return
	}

	userDTO := dto.MapUserToDTO(createdUser)
	utils.ReponseSuccess(ctx, http.StatusCreated, &userDTO)
}
func (uh *UserHandler) GetUserByUUID(ctx *gin.Context) {

}
func (uh *UserHandler) UpdateUser(ctx *gin.Context) {

}
func (uh *UserHandler) DeleteUser(ctx *gin.Context) {

}
