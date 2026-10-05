package app

import (
	v1handler "github.com/duitocdai/internal/handler/v1"
	"github.com/duitocdai/internal/repository"
	"github.com/duitocdai/internal/routes"
	v1routes "github.com/duitocdai/internal/routes/v1"
	v1service "github.com/duitocdai/internal/service/v1"
	"github.com/duitocdai/pkg/auth"
)

type AuthModule struct {
	routes routes.Route
}

func NewAuthModule(ctx *ModuleContext, tokenService auth.TokenService) *AuthModule {
	userRepo := repository.NewSqlUserRepository(ctx.DB)
	authService := v1service.NewAuthService(userRepo, tokenService)
	authHandler := v1handler.NewAuthHandler(authService)
	authRoutes := v1routes.NewAuthRoutes(authHandler)

	return &AuthModule{
		routes: authRoutes,
	}

}

func (m *AuthModule) Routes() routes.Route {
	return m.routes
}
