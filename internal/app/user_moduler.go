package app

import (
	"github.com/duitocdai/internal/handler"
	"github.com/duitocdai/internal/repository"
	"github.com/duitocdai/internal/routes"
	"github.com/duitocdai/internal/service"
)

type UserModule struct {
	routes routes.Route
}

func NewUserModule() *UserModule {
	///Initialize reposiroty
	userRepo := repository.NewInMemoryUserRepository()
	///Initialize service
	userService := service.NewUserService(userRepo)
	///Initialize hander
	userHandler := handler.NewUserHandler(userService)
	///Initialize routes
	userRoutes := routes.NewUserRoutes(userHandler)

	return &UserModule{
		routes: userRoutes,
	}

}

func (m *UserModule) Routes() routes.Route {
	return m.routes
}
