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
	accessToken, refreshToken, exprireAt, err := ah.service.Login(ctx, input.Email, input.Password)
	if err != nil {
		utils.ReponseError(ctx, err)
		return
	}
	response := v1dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    exprireAt,
	}
	utils.ReponseSuccess(ctx, http.StatusOK, "Login Successfully", response)
}
func (ah *AuthHandler) Logout(ctx *gin.Context) {
	var input v1dto.RefreshTokenInput

	if err := ctx.ShouldBindBodyWithJSON(&input); err != nil {
		utils.ReponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	if err := ah.service.Logout(ctx, input.RefreshToken); err != nil {
		utils.ReponseError(ctx, err)
		return
	}
	utils.ReponseSuccess(ctx, http.StatusOK, "Logout Successfully")
}
func (ah *AuthHandler) RefreshToken(ctx *gin.Context) {
	var input v1dto.RefreshTokenInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ReponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	accessToken, refreshToken, exprireAt, err := ah.service.RefreshToken(ctx, input.RefreshToken)
	if err != nil {
		utils.ReponseError(ctx, err)
		return
	}
	response := v1dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    exprireAt,
	}
	utils.ReponseSuccess(ctx, http.StatusOK, "Login Successfully", response)
}
func (ah *AuthHandler) RequestForgotPassword(ctx *gin.Context) {
	var input v1dto.RequestPasswordInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ReponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	err := ah.service.RequestForgotPassword(ctx, input.Email)
	if err != nil {
		utils.ReponseError(ctx, err)
	}
	utils.ReponseSuccess(ctx, http.StatusOK, "Reset link sent to email")
}
func (ah *AuthHandler) ResetPassword(ctx *gin.Context) {
	var input v1dto.ResetPasswordInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ReponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	err := ah.service.ResetPassword(ctx, input.Token, input.NewPassword)
	if err != nil {
		utils.ReponseError(ctx, err)
		return
	}
	utils.ReponseSuccess(ctx, http.StatusOK, "Password reset successfully")
}
