package routes

import (
	"github.com/duitocdai/internal/middleware"
	"github.com/duitocdai/internal/utils"
	"github.com/duitocdai/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Route interface {
	Register(r *gin.RouterGroup)
}

func RegisterRoutes(r *gin.Engine, routes ...Route) {
	httpLogger := newLoggerWithPath("../../internal/logs/http.log", "info")
	recoveryLogger := newLoggerWithPath("../../internal/logs/recovery.log", "warning")
	rateLimiterLoger := newLoggerWithPath("../../internal/logs/ratelimiter.log", "warning")
	r.Use(
		middleware.RateLimiterMiddleware(rateLimiterLoger),
		middleware.LoggerMiddleware(httpLogger),
		middleware.RecoveryMiddlware(recoveryLogger),
		middleware.ApiKeyMiddleware(),
		middleware.AuthMiddleware(),
	)
	v1api := r.Group("/api/v1")

	for _, route := range routes {
		route.Register(v1api)
	}
}

func newLoggerWithPath(path string, level string) *zerolog.Logger {
	config := logger.LoggerConfig{
		Level:      level,
		Filename:   path,
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     5,
		Compress:   true,
		IsDev:      utils.GetEnv("APP_ENV", "development"),
	}
	return logger.NewLogger(config)
}
