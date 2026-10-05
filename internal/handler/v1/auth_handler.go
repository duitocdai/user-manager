package v1handler

import (
	"net/http"

	v1dto "github.com/duitocdai/internal/dto/v1"
	v1service "github.com/duitocdai/internal/service/v1"
	"github.com/duitocdai/internal/utils"
	"github.com/duitocdai/internal/validation"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service v1service.AuthService
}

func NewAuthHandler(service v1service.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (ah *AuthHandler) Login(ctx *gin.Context) {
	var input v1dto.LoginInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ReponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	accessToken, exprireAt, err := ah.service.Login(ctx, input.Email, input.Password)
	if err != nil {
		utils.ReponseError(ctx, err)
		return
	}
	response := v1dto.LoginResponse{
		AccessToken: accessToken,
		ExpiresIn:   exprireAt,
	}
	utils.ReponseSuccess(ctx, http.StatusOK, "Login Successfully", response)
}
func (ah *AuthHandler) Logout(ctx *gin.Context) {

}
