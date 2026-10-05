package v1service

import (
	"github.com/duitocdai/internal/repository"
	"github.com/duitocdai/internal/utils"
	"github.com/duitocdai/pkg/auth"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepo     repository.UserRepository
	tokenService auth.TokenService
}

func NewAuthService(repo repository.UserRepository, tokenService auth.TokenService) *authService {
	return &authService{
		userRepo:     repo,
		tokenService: tokenService,
	}
}
func (as *authService) Login(ctx *gin.Context, email, password string) (string, int, error) {
	context := ctx.Request.Context()

	email = utils.NormalizeString(email)
	user, err := as.userRepo.GetByEmail(context, email)
	if err != nil {
		return "", 0, utils.NewError("invalid email or password", utils.ErrCodeUnauthorized)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.UserPassword), []byte(password)); err != nil {
		return "", 0, utils.NewError("invalid email or password", utils.ErrCodeUnauthorized)
	}
	accessToken, err := as.tokenService.GenerateAccessToken(user)
	if err != nil {
		return "", 0, utils.WrapError(err, "Unable to create accessToken", utils.ErrCodeInternal)
	}
	return accessToken, int(auth.AccessTokenTTL.Seconds()), nil
}
func (as *authService) Logout(ctx *gin.Context) error {

	return nil
}
