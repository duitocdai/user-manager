package v1service

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/duitocdai/internal/db/sqlc"
	"github.com/duitocdai/internal/repository"
	"github.com/duitocdai/internal/utils"
	"github.com/duitocdai/pkg/cache"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo  repository.UserRepository
	cache *cache.RedisCacheService
}

func NewUserService(repo repository.UserRepository, redisClient *redis.Client) UserService {
	return &userService{
		repo:  repo,
		cache: cache.NewRedisCacheService(redisClient),
	}
}
func (us *userService) GetAllUser(ctx *gin.Context, search, orderBy, sort string, page, limit int32, deleted bool) ([]sqlc.User, int32, error) {
	context := ctx.Request.Context()

	cacheKey := us.generateCacheKey(search, orderBy, sort, page, limit, deleted)
	/**Get Cache redis **/
	var cacheData struct {
		Users []sqlc.User `json:"users"`
		Total int32       `json:"total"`
	}

	if err := us.cache.Get(cacheKey, &cacheData); err == nil && cacheData.Users != nil {
		return cacheData.Users, cacheData.Total, nil
	}

	if sort == "" {
		sort = "desc"
	}

	if orderBy == "" {
		orderBy = "user_created_at"
	}

	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		limitInt := utils.GetIntEnv("LIMIT_ITEM_ON_PER_PAGE", 10)

		limit = int32(limitInt)
	}
	offset := (page - 1) * limit

	users, err := us.repo.GetAllV2(context, search, orderBy, sort, limit, offset, deleted)
	if err != nil {
		return []sqlc.User{}, 0, utils.WrapError(err, "failed to fetch users", utils.ErrCodeInternal)
	}

	total, err := us.repo.CountUsers(context, search, deleted)
	if err != nil {
		return []sqlc.User{}, 0, utils.WrapError(err, "failed to count users", utils.ErrCodeInternal)
	}
	/// create cache dataa
	cacheData = struct {
		Users []sqlc.User `json:"users"`
		Total int32       `json:"total"`
	}{
		Users: users,
		Total: int32(total),
	}
	us.cache.Set("getAllUsers", cacheData, 5*time.Minute)
	return users, int32(total), nil
}
func (us *userService) CreateUser(ctx *gin.Context, input sqlc.CreateUserParams) (sqlc.User, error) {
	context := ctx.Request.Context()
	input.UserEmail = utils.NormalizeString(input.UserEmail)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.UserPassword), bcrypt.DefaultCost)
	if err != nil {
		return sqlc.User{}, utils.WrapError(err, "failed to hashpassword", utils.ErrCodeInternal)
	}
	input.UserPassword = string(hashedPassword)
	user, err := us.repo.Create(context, input)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return sqlc.User{}, utils.NewError("email already exists", utils.ErrCodeConfict)
		}
		return sqlc.User{}, utils.WrapError(err, "failed to created", utils.ErrCodeInternal)
	}
	///-- Clear cache redis
	if err := us.cache.Clear("user:*"); err != nil {
		log.Printf("fail to clear cache: %v", err)
	}
	return user, nil
}
func (us *userService) GetUserByUUID(ctx *gin.Context, userUuid uuid.UUID) (sqlc.User, error) {
	context := ctx.Request.Context()
	user, err := us.repo.GetByUUID(context, userUuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.User{}, utils.NewError("user not found", utils.ErrCodeNotFound)
		}
		return sqlc.User{}, utils.WrapError(err, "failed to get an user", utils.ErrCodeInternal)
	}
	return user, nil
}
func (us *userService) UpdateUser(ctx *gin.Context, input sqlc.UpdateUserParams) (sqlc.User, error) {
	context := ctx.Request.Context()
	if input.UserPassword != nil && *input.UserPassword != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*input.UserPassword), bcrypt.DefaultCost)
		if err != nil {
			return sqlc.User{}, utils.WrapError(err, "failed to hashpassword", utils.ErrCodeInternal)
		}

		hashed := string(hashedPassword)
		input.UserPassword = &hashed
	}

	updatedUser, err := us.repo.Update(context, input)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.User{}, utils.NewError("user not found", utils.ErrCodeNotFound)
		}
		return sqlc.User{}, utils.WrapError(err, "failed to updated", utils.ErrCodeInternal)
	}
	return updatedUser, nil
}
func (us *userService) SoftDeleteUser(ctx *gin.Context, uuid uuid.UUID) (sqlc.User, error) {
	context := ctx.Request.Context()
	softDeleteUser, err := us.repo.SoftDelete(context, uuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.User{}, utils.NewError("user not found", utils.ErrCodeNotFound)
		}
		return sqlc.User{}, utils.WrapError(err, "failed to delete user", utils.ErrCodeInternal)
	}
	return softDeleteUser, nil
}
func (us *userService) RestoreUser(ctx *gin.Context, uuid uuid.UUID) (sqlc.User, error) {
	context := ctx.Request.Context()
	restoreUser, err := us.repo.RestoreUser(context, uuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.User{}, utils.NewError("user not found", utils.ErrCodeNotFound)
		}
		return sqlc.User{}, utils.WrapError(err, "failed to restore user", utils.ErrCodeInternal)
	}
	return restoreUser, nil
}
func (us *userService) DeleteUser(ctx *gin.Context, uuid uuid.UUID) error {
	context := ctx.Request.Context()
	_, err := us.repo.Delete(context, uuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NewError("user not found", utils.ErrCodeNotFound)
		}
		return utils.WrapError(err, "failed to delete user", utils.ErrCodeInternal)
	}
	return nil
}

func (us *userService) generateCacheKey(search, orderBy, sort string, page, limit int32, deleted bool) string {
	search = strings.TrimSpace(search)
	if search == "" {
		search = "none"
	}
	orderBy = strings.TrimSpace(orderBy)
	if orderBy == "" {
		orderBy = "user_created_at"
	}

	sort = strings.ToLower(strings.TrimSpace(sort))
	if sort == "" {
		sort = "desc"
	}

	return fmt.Sprintf("users:%s:%s:%s:%d:%d:%t", search, orderBy, sort, page, limit, deleted)
}
