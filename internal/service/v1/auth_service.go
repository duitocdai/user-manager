package v1service

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/duitocdai/internal/db/sqlc"
	"github.com/duitocdai/internal/repository"
	"github.com/duitocdai/internal/utils"
	"github.com/duitocdai/pkg/auth"
	"github.com/duitocdai/pkg/cache"
	"github.com/duitocdai/pkg/logger"
	"github.com/duitocdai/pkg/mail"
	"github.com/duitocdai/pkg/rabbitmq"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/time/rate"
)

type authService struct {
	userRepo     repository.UserRepository
	tokenService auth.TokenService
	cacheService cache.RedisCacheService
	mailService  mail.EmailProviderService
	rabbitmq     rabbitmq.RabbitMQService
}
type LoginAttempt struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	mu              sync.Mutex
	clients         = make(map[string]*LoginAttempt)
	LoginAttemptTTL = 5 * time.Minute
	MaxLoginAttempt = 5
)

func NewAuthService(repo repository.UserRepository, tokenService auth.TokenService, cacheService cache.RedisCacheService, mailService mail.EmailProviderService, rabbitmqService rabbitmq.RabbitMQService) *authService {
	return &authService{
		userRepo:     repo,
		tokenService: tokenService,
		cacheService: cacheService,
		mailService:  mailService,
		rabbitmq:     rabbitmqService,
	}
}

func (as *authService) getClientIP(ctx *gin.Context) string {
	ip := ctx.ClientIP()
	if ip == "" {
		ip = ctx.Request.RemoteAddr
	}

	return ip
}

func (as *authService) getLoginAttempt(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	client, exists := clients[ip]
	if !exists {
		limiter := rate.NewLimiter(rate.Limit(float32(MaxLoginAttempt)/float32(LoginAttemptTTL.Seconds())), MaxLoginAttempt)
		newClient := &LoginAttempt{limiter, time.Now()}
		clients[ip] = newClient
		return limiter
	}

	client.lastSeen = time.Now()
	return client.limiter
}
func (as *authService) CheckLoginAttempt(ip string) error {
	limiter := as.getLoginAttempt(ip)

	if !limiter.Allow() {
		return utils.NewError(" Too many login attempt. Pls retry again later", utils.ErrrCodeTooManyRequest)
	}
	return nil
}
func (as *authService) CleanupClients(ip string) {
	mu.Lock()
	defer mu.Unlock()
	delete(clients, ip)
}
func (as *authService) Login(ctx *gin.Context, email, password string) (string, string, int, error) {
	context := ctx.Request.Context()
	ip := as.getClientIP(ctx)
	if err := as.CheckLoginAttempt(ip); err != nil {
		return "", "", 0, err
	}
	email = utils.NormalizeString(email)
	user, err := as.userRepo.GetByEmail(context, email)
	if err != nil {
		as.getLoginAttempt(ip)
		return "", "", 0, utils.NewError("invalid email or password", utils.ErrCodeUnauthorized)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.UserPassword), []byte(password)); err != nil {
		as.getLoginAttempt(ip)
		return "", "", 0, utils.NewError("invalid email or password", utils.ErrCodeUnauthorized)
	}
	accessToken, err := as.tokenService.GenerateAccessToken(user)
	if err != nil {
		return "", "", 0, utils.WrapError(err, "Unable to create accessToken", utils.ErrCodeInternal)
	}
	refreshToken, err := as.tokenService.GenerateRefreshToken(user)
	if err != nil {
		return "", "", 0, utils.WrapError(err, "Unable to create refreshToken", utils.ErrCodeInternal)
	}
	if err := as.tokenService.StoreRefreshToken(refreshToken); err != nil {
		return "", "", 0, utils.WrapError(err, "Can not save refresh Token", utils.ErrCodeInternal)
	}
	as.CleanupClients(ip)

	return accessToken, refreshToken.Token, int(auth.AccessTokenTTL.Seconds()), nil
}
func (as *authService) Logout(ctx *gin.Context, refreshToken string) error {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return utils.NewError("Missing Authorization header", utils.ErrCodeUnauthorized)
	}

	accessToken := strings.TrimPrefix(authHeader, "Bearer ")

	_, claims, err := as.tokenService.ParseToken(accessToken)
	if err != nil {
		return utils.NewError("Invalid access Token", utils.ErrCodeUnauthorized)
	}

	if jti, ok := claims["jti"].(string); ok {
		expUnix, _ := claims["exp"].(float64)
		exp := time.Unix(int64(expUnix), 0)
		key := "blacklist:" + jti
		ttl := time.Until(exp)
		as.cacheService.Set(key, "revoked", ttl)
	}
	_, err = as.tokenService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return utils.NewError("Refresh Token is invalid or revoked", utils.ErrCodeUnauthorized)
	}
	if err := as.tokenService.RevokedRefreshToken(refreshToken); err != nil {
		return utils.WrapError(err, "Unabled to revoke token", utils.ErrCodeUnauthorized)
	}
	return nil
}

func (as *authService) RefreshToken(ctx *gin.Context, refreshTokenString string) (string, string, int, error) {
	context := ctx.Request.Context()

	token, err := as.tokenService.ValidateRefreshToken(refreshTokenString)
	if err != nil {
		return "", "", 0, utils.NewError("Refresh Token invalid or revoked", utils.ErrCodeInternal)
	}
	userUuid, _ := uuid.Parse(token.UserUUID)
	user, err := as.userRepo.GetByUUID(context, userUuid)
	if err != nil {
		return "", "", 0, utils.NewError("user not found", utils.ErrCodeInternal)
	}
	accessToken, err := as.tokenService.GenerateAccessToken(user)
	if err != nil {
		return "", "", 0, utils.WrapError(err, "Unable to create accessToken", utils.ErrCodeInternal)
	}
	refreshToken, err := as.tokenService.GenerateRefreshToken(user)
	if err != nil {
		return "", "", 0, utils.WrapError(err, "Unable to create refreshToken", utils.ErrCodeInternal)
	}
	if err := as.tokenService.RevokedRefreshToken(refreshTokenString); err != nil {
		return "", "", 0, utils.WrapError(err, "Unable to revoked refreshToken", utils.ErrCodeInternal)
	}

	if err := as.tokenService.StoreRefreshToken(refreshToken); err != nil {
		return "", "", 0, utils.WrapError(err, "Can not save refresh Token", utils.ErrCodeInternal)
	}

	return accessToken, refreshToken.Token, int(auth.AccessTokenTTL.Seconds()), nil
}
func (as *authService) RequestForgotPassword(ctx *gin.Context, email string) error {
	context := ctx.Request.Context()
	traceID := logger.GetTraceId(context)
	rateLimitKey := fmt.Sprintf("reset:ratelimit:%s", email)
	if exists, err := as.cacheService.Exists(rateLimitKey); err == nil && exists {
		return utils.NewError("Please wait before requesting another password", utils.ErrrCodeTooManyRequest)
	}
	user, err := as.userRepo.GetByEmail(context, email)
	if err != nil {
		return utils.NewError("Email Not Found", utils.ErrCodeNotFound)
	}
	token, err := utils.GenerateRandomString(16)
	if err != nil {
		return utils.NewError("fail to generate", utils.ErrCodeInternal)
	}
	if err := as.cacheService.Set("reset:"+token, user.UserUuid, time.Hour); err != nil {
		return utils.NewError("Fail to store reset token", utils.ErrCodeInternal)
	}

	if err := as.cacheService.Set(rateLimitKey, "1", 5*time.Minute); err != nil {
		return utils.NewError("Fail to store rate limit reset password", utils.ErrCodeInternal)
	}

	resetLink := fmt.Sprintf("https://app.example.com/reset-password?token=%s", token)

	mailContent := &mail.Email{
		To: []mail.Address{
			{Email: email},
		},
		Subject: "Password Reset Request",
		Text:    fmt.Sprintf("Hi %s, \n\n You requested to reset your password please click the link below to reset it:\n%s\n\nThe link will expire in 1 hour. \n\n Best Regard.", user.UserEmail, resetLink),
	}

	///cần update thành publish vào rabbitMQ
	message := mail.QueueMessage{
		Email:   *mailContent,
		TraceID: traceID,
	}
	if err := as.rabbitmq.Publish(context, "auth_email_queue", message); err != nil {
		return utils.NewError("Failed to send password reset email", utils.ErrCodeInternal)
	}

	// if err := as.mailService.SendMail(context, mailContent); err != nil {
	// 	return utils.NewError("Failed to send email reset password", utils.ErrCodeInternal)
	// }
	return nil
}
func (as *authService) ResetPassword(ctx *gin.Context, token, password string) error {
	context := ctx.Request.Context()
	var UserUUIDStr string
	err := as.cacheService.Get("reset:"+token, &UserUUIDStr)
	if err == redis.Nil || UserUUIDStr == "" {
		return utils.NewError("INVALID OR EXP", utils.ErrCodeNotFound)
	}
	if err != nil {
		return utils.NewError("fail to get reset Token", utils.ErrCodeInternal)
	}

	userUuid, err := uuid.Parse(UserUUIDStr)
	if err != nil {
		return utils.WrapError(err, "Uuid is invalid", utils.ErrCodeInternal)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return utils.WrapError(err, "fail to hash password", utils.ErrCodeInternal)
	}
	var input = sqlc.UpdatePasswordParams{
		UserPassword: string(hashedPassword),
		UserUuid:     userUuid,
	}

	_, err = as.userRepo.UpdatePassword(context, input)

	if err != nil {
		return utils.WrapError(err, "fail to update new password", utils.ErrCodeInternal)
	}

	err = as.cacheService.Clear("reset:" + token)
	return nil
}
