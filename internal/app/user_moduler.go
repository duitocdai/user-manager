package app

import (
	v1handler "github.com/duitocdai/internal/handler/v1"
	"github.com/duitocdai/internal/repository"
	"github.com/duitocdai/internal/routes"
	v1routes "github.com/duitocdai/internal/routes/v1"
	v1service "github.com/duitocdai/internal/service/v1"
)

type UserModule struct {
	routes routes.Route
}

func NewUserModule(ctx *ModuleContext) *UserModule {
	///Initialize reposiroty
	userRepo := repository.NewSqlUserRepository(ctx.DB)
	///Initialize service
	userService := v1service.NewUserService(userRepo, ctx.Redis)
	///Initialize hander
	userHandler := v1handler.NewUserHandler(userService)
	///Initialize routes
	userRoutes := v1routes.NewUserRoutes(userHandler)

	return &UserModule{
		routes: userRoutes,
	}

}

func (m *UserModule) Routes() routes.Route {
	return m.routes
}
