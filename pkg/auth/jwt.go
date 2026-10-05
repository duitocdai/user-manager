package auth

import (
	"encoding/json"
	"time"

	"github.com/duitocdai/internal/db/sqlc"
	"github.com/duitocdai/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTService struct {
}
type Claims struct {
}
type EncryptPayload struct {
	UserUUID string `json:"user_uuid"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

var (
	jwtSecret     = []byte(utils.GetEnv("JWT_SECRET", "8f2c9a7d4e1b6c0f3a5e9d8b2c7f1a6e5d0b4c9f7a3e8d1b6c2f5a9e4d7b0c1f"))
	jwtEncryptKey = []byte(utils.GetEnv("JWT_ENCRYPT_KEY", "K8v2Qp7mX4nR9tL3cW6zY1aF5hJ0sD"))
)

const (
	AccessTokenTTL = 15 * time.Minute
)

func NewJWTService() *JWTService {
	return &JWTService{}
}

// / Access Token
func (js *JWTService) GenerateAccessToken(user sqlc.User) (string, error) {
	payload := &EncryptPayload{
		UserUUID: user.UserUuid.String(),
		Email:    user.UserEmail,
		Role:     string(rune(user.UserLevel)),
	}

	rawData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	encrypt, err := utils.EncryptAES(rawData, jwtEncryptKey)
	if err != nil {
		return "", err
	}
	claims := jwt.MapClaims{
		"data": encrypt,
		"jti":  uuid.NewString(),
		"exp":  time.Now().Add(AccessTokenTTL).Unix(),
		"iat":  time.Now().Unix(),
		"iss":  "khoa hoc lap tring golang",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// / Refresh Token
func (js *JWTService) GenerateToken() {

}

func (js *JWTService) ParseToken(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, nil, utils.NewError("invalid token", utils.ErrCodeUnauthorized)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, nil, utils.NewError("invalid token", utils.ErrCodeUnauthorized)
	}
	return token, claims, nil
}
