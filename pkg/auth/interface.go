package auth

import (
	"github.com/duitocdai/internal/db/sqlc"
	"github.com/golang-jwt/jwt/v5"
)

type TokenService interface {
	GenerateAccessToken(user sqlc.User) (string, error)
	GenerateRefreshToken(user sqlc.User) (RefreshToken, error)
	ParseToken(tokenString string) (*jwt.Token, jwt.MapClaims, error)
	DecryptAccessTokenPayload(tokenString string) (*EncryptPayload, error)
	StoreRefreshToken(token RefreshToken) error
	RevokedRefreshToken(token string) error
	ValidateRefreshToken(token string) (RefreshToken, error)
}
