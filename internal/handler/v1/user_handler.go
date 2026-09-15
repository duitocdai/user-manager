package v1handler

import (
	"log"
	"net/http"

	v1dto "github.com/duitocdai/internal/dto/v1"
	v1service "github.com/duitocdai/internal/service/v1"
	"github.com/duitocdai/internal/utils"
	"github.com/duitocdai/internal/validation"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service v1service.UserService
}

func NewUserHandler(service v1service.UserService) *UserHandler {
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

	utils.ReponseSuccess(ctx, http.StatusOK, "usersDTO")
}
func (uh *UserHandler) CreateUser(ctx *gin.Context) {
	var input v1dto.CreateUserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ReponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	log.Println("vl")
	utils.ReponseSuccess(ctx, http.StatusCreated, "&userDTO")
}
func (uh *UserHandler) GetUserByUUID(ctx *gin.Context) {
	var params GetUserByUuidParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ReponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	utils.ReponseSuccess(ctx, http.StatusCreated, "")
}
func (uh *UserHandler) UpdateUser(ctx *gin.Context) {
	var params GetUserByUuidParam

	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ReponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	var input v1dto.UpdateUserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ReponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	utils.ReponseSuccess(ctx, http.StatusOK, " &userDTO")
}
func (uh *UserHandler) DeleteUser(ctx *gin.Context) {
	var params GetUserByUuidParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ReponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	utils.ReponseStatusCode(ctx, http.StatusNoContent)
}
