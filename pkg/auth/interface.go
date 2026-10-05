package auth

import "github.com/duitocdai/internal/db/sqlc"

type TokenService interface {
	GenerateAccessToken(user sqlc.User) (string, error)
	GenerateTokenToken()
}
