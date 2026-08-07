package handler

import (
	"net/http"

	"github.com/duitocdai/internal/dto"
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

type GetUserByUuidParam struct {
	Uuid string `uri:"uuid" binding:"uuid"`
}
type GetUserParams struct {
	Search string `form:"search" binding:"omitempty,min=3,max=50,search"`
	Page   int    `form:"page" binding:"omitempty,gte=1,lte=100"`
	Limit  int    `form:"limit" binding:"omitempty,gte=1,lte=100"`
}

func (uh *UserHandler) GetAllUser(ctx *gin.Context) {
	var params GetUserParams

	if err := ctx.ShouldBindQuery(&params); err != nil {
		utils.ReponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	if params.Page == 0 {
		params.Page = 1
	}
	if params.Limit == 0 {
		params.Limit = 10
	}
	users, err := uh.service.GetAllUser(params.Search, params.Page, params.Limit)
	if err != nil {
		utils.ReponseError(ctx, err)
		return
	}
	usersDTO := dto.MapUsersToDTO(users)
	utils.ReponseSuccess(ctx, http.StatusOK, usersDTO)
}
func (uh *UserHandler) CreateUser(ctx *gin.Context) {
	var input dto.CreateUserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ReponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	user := input.MapCreateInputToModel()
	createdUser, err := uh.service.CreateUser(user)
	if err != nil {
		utils.ReponseError(ctx, err)
		return
	}

	userDTO := dto.MapUserToDTO(createdUser)
	utils.ReponseSuccess(ctx, http.StatusCreated, &userDTO)
}
func (uh *UserHandler) GetUserByUUID(ctx *gin.Context) {
	var params GetUserByUuidParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ReponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	user, err := uh.service.GetUserByUUID(params.Uuid)
	if err != nil {
		utils.ReponseError(ctx, err)
		return
	}
	userDTO := dto.MapUserToDTO(user)
	utils.ReponseSuccess(ctx, http.StatusCreated, &userDTO)
}
func (uh *UserHandler) UpdateUser(ctx *gin.Context) {
	var params GetUserByUuidParam

	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ReponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	var input dto.UpdateUserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ReponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	users := input.MapUpdateInputToModel()
	updatedUser, err := uh.service.UpdateUser(params.Uuid, users)
	if err != nil {
		utils.ReponseError(ctx, err)
		return
	}

	userDTO := dto.MapUserToDTO(updatedUser)

	utils.ReponseSuccess(ctx, http.StatusOK, &userDTO)
}
func (uh *UserHandler) DeleteUser(ctx *gin.Context) {
	var params GetUserByUuidParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ReponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	if err := uh.service.DeleteUser(params.Uuid); err != nil {
		utils.ReponseError(ctx, err)
		return
	}
	utils.ReponseStatusCode(ctx, http.StatusNoContent)
}
