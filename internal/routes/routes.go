package routes

import (
	"github.com/duitocdai/internal/middleware"
	v1routes "github.com/duitocdai/internal/routes/v1"
	"github.com/duitocdai/internal/utils"
	"github.com/duitocdai/pkg/auth"
	"github.com/duitocdai/pkg/cache"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

type Route interface {
	Register(r *gin.RouterGroup)
}

func RegisterRoutes(r *gin.Engine, authService auth.TokenService, cacheService cache.RedisCacheService, routes ...Route) {
	httpLogger := utils.NewLoggerWithPath("http.log", "info")
	recoveryLogger := utils.NewLoggerWithPath("recovery.log", "warning")
	rateLimiterLoger := utils.NewLoggerWithPath("ratelimiter.log", "warning")
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(
		middleware.RateLimiterMiddleware(rateLimiterLoger),
		middleware.CORSMiddleware(),
		middleware.TraceMiddleware(),
		middleware.LoggerMiddleware(httpLogger),
		middleware.RecoveryMiddleware(recoveryLogger),
		middleware.ApiKeyMiddleware(),
	)

	v1api := r.Group("/api/v1")
	middleware.InitAuthMiddleware(authService, cacheService)
	protected := v1api.Group("")
	protected.Use(
		middleware.AuthMiddleware(),
	)

	for _, route := range routes {
		switch route.(type) {
		case *v1routes.AuthRoutes:
			route.Register(v1api)
		default:
			route.Register(protected)
		}
	}

	r.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(404, gin.H{
			"error": "Not Found",
			"path":  ctx.Request.URL.Path,
		})
	})
}
