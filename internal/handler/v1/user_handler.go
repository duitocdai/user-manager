package v1handler

import (
	"net/http"

	v1dto "github.com/duitocdai/internal/dto/v1"
	v1service "github.com/duitocdai/internal/service/v1"
	"github.com/duitocdai/internal/utils"
	"github.com/duitocdai/internal/validation"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	service v1service.UserService
}

func NewUserHandler(service v1service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (uh *UserHandler) GetAllUser(ctx *gin.Context) {
	var params v1dto.GetUserParams

	if err := ctx.ShouldBindQuery(&params); err != nil {
		utils.ReponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	users, total, err := uh.service.GetAllUser(ctx, params.Search, params.OrderBy, params.Sort, params.Page, params.Limit, false)
	if err != nil {
		utils.ReponseError(ctx, err)
		return
	}
	userDTO := v1dto.MapUsersToDTO(users)
	paginationResp := utils.NewPaginationResponse(userDTO, params.Page, params.Limit, total)
	utils.ReponseSuccess(ctx, http.StatusOK, "User list successfully", paginationResp)
}
func (uh *UserHandler) CreateUser(ctx *gin.Context) {
	var input v1dto.CreateUserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ReponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	user := input.MapCreateInputToModel()
	createdUser, err := uh.service.CreateUser(ctx, user)
	if err != nil {
		utils.ReponseError(ctx, err)
		return
	}
	userDTO := v1dto.MapUserToDTO(createdUser)
	utils.ReponseSuccess(ctx, http.StatusCreated, "user created successsfully", userDTO)
}
func (uh *UserHandler) GetUserByUUID(ctx *gin.Context) {
	var params v1dto.GetUserByUuidParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ReponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	userUuid, err := uuid.Parse(params.Uuid)
	if err != nil {
		utils.ReponseError(ctx, err)
		return
	}

	user, err := uh.service.GetUserByUUID(ctx, userUuid)
	if err != nil {
		utils.ReponseError(ctx, err)
		return
	}
	userDTO := v1dto.MapUserToDTO(user)
	utils.ReponseSuccess(ctx, http.StatusOK, "get An user successfully", userDTO)
}
func (uh *UserHandler) GetUserSoftDeleted(ctx *gin.Context) {
	var params v1dto.GetUserParams

	if err := ctx.ShouldBindQuery(&params); err != nil {
		utils.ReponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	users, total, err := uh.service.GetAllUser(ctx, params.Search, params.OrderBy, params.Sort, params.Page, params.Limit, true)
	if err != nil {
		utils.ReponseError(ctx, err)
		return
	}
	userDTO := v1dto.MapUsersToDTO(users)
	paginationResp := utils.NewPaginationResponse(userDTO, params.Page, params.Limit, total)
	utils.ReponseSuccess(ctx, http.StatusOK, "List user deleted successfully", paginationResp)
}
func (uh *UserHandler) UpdateUser(ctx *gin.Context) {
	var params v1dto.GetUserByUuidParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ReponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	userUuid, err := uuid.Parse(params.Uuid)
	if err != nil {
		utils.ReponseError(ctx, err)
		return
	}

	var input v1dto.UpdateUserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ReponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	user := input.MapUpdateInputToModel(userUuid)

	updatedUser, err := uh.service.UpdateUser(ctx, user)
	if err != nil {
		utils.ReponseError(ctx, err)
		return
	}
	userDTO := v1dto.MapUserToDTO(updatedUser)
	utils.ReponseSuccess(ctx, http.StatusOK, "user updated", userDTO)
}
func (uh *UserHandler) SoftDeleteUser(ctx *gin.Context) {
	var params v1dto.GetUserByUuidParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ReponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	userUuid, err := uuid.Parse(params.Uuid)
	if err != nil {
		utils.ReponseError(ctx, err)
		return
	}
	softDeleteUser, err := uh.service.SoftDeleteUser(ctx, userUuid)
	if err != nil {
		utils.ReponseError(ctx, err)
		return
	}
	userDTO := v1dto.MapUserToDTO(softDeleteUser)
	utils.ReponseSuccess(ctx, http.StatusOK, "user deleted", userDTO)
}

func (uh *UserHandler) RestoreUser(ctx *gin.Context) {
	var params v1dto.GetUserByUuidParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ReponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	userUuid, err := uuid.Parse(params.Uuid)
	if err != nil {
		utils.ReponseError(ctx, err)
		return
	}

	restoreUser, err := uh.service.RestoreUser(ctx, userUuid)
	if err != nil {
		utils.ReponseError(ctx, err)
		return
	}
	userDTO := v1dto.MapUserToDTO(restoreUser)
	utils.ReponseSuccess(ctx, http.StatusOK, "user restored", userDTO)
}

func (uh *UserHandler) DeleteUser(ctx *gin.Context) {
	var params v1dto.GetUserByUuidParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ReponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	userUuid, err := uuid.Parse(params.Uuid)
	if err != nil {
		utils.ReponseError(ctx, err)
		return
	}

	err = uh.service.DeleteUser(ctx, userUuid)
	if err != nil {
		utils.ReponseError(ctx, err)
		return
	}
	utils.ReponseSuccess(ctx, http.StatusOK, "user deleted")
}
