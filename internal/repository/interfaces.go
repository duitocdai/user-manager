package repository

import (
	"context"

	"github.com/duitocdai/internal/db/sqlc"
)

type UserRepository interface {
	FindAll()
	Create(ctx context.Context, userParams sqlc.CreateUserParams)
	FindByUUID(uuid string)
	Update(uuid string)
	Delete(uuid string)
	FindByEmail(email string)
}
