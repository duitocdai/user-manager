package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/duitocdai/internal/config"
	"github.com/duitocdai/internal/db"
	"github.com/duitocdai/internal/db/sqlc"
	"github.com/duitocdai/internal/routes"
	"github.com/duitocdai/internal/utils"
	"github.com/duitocdai/internal/validation"
	"github.com/duitocdai/pkg/auth"
	"github.com/duitocdai/pkg/cache"
	"github.com/duitocdai/pkg/logger"
	"github.com/duitocdai/pkg/mail"
	"github.com/duitocdai/pkg/rabbitmq"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Module interface {
	Routes() routes.Route
}

type Application struct {
	config  *config.Config
	router  *gin.Engine
	modules []Module
}
type ModuleContext struct {
	DB    sqlc.Querier
	Redis *redis.Client
}

func NewApplication(cfg *config.Config) (*Application, error) {
	if err := validation.InitValidator(); err != nil {
		logger.Log.Fatal().Err(err).Msg("Validator init failed ")
		return nil, err
	}
	r := gin.Default()
	if err := db.InitDB(); err != nil {
		logger.Log.Fatal().Err(err).Msg("DB init false ")
	}
	redisClient := config.NewRedisClient()
	cacheRedisService := cache.NewRedisCacheService(redisClient)
	tokenService := auth.NewJWTService(cacheRedisService)
	mailLogger := utils.NewLoggerWithPath("mail.log", "info")
	factory, err := mail.NewProviderFactory(mail.ProviderMailtrap)
	if err != nil {
		mailLogger.Error().Err(err).Msg("fail to create mail provider factory")
		return nil, err
	}
	mailService, err := mail.NewMailService(cfg, mailLogger, factory)
	if err != nil {
		mailLogger.Error().Err(err).Msg("fail to initilize mail service")
		return nil, err
	}
	rabbitmqLogger := utils.NewLoggerWithPath("worker.log", "info")
	rabbitmqService, _ := rabbitmq.NewRabbitMQService(
		utils.GetEnv("RABBITMQ_URL", "amqp://guest:guest@rabbitmq:5672/"), rabbitmqLogger,
	)
	ctx := &ModuleContext{
		DB:    db.DB,
		Redis: redisClient,
	}
	modules := []Module{
		NewUserModule(ctx),
		NewAuthModule(ctx, tokenService, cacheRedisService, mailService, rabbitmqService),
	}
	routes.RegisterRoutes(r, tokenService, cacheRedisService, getModulRoutes(modules)...)
	return &Application{
		config:  cfg,
		router:  r,
		modules: modules,
	}, nil
}

func (a *Application) Run() error {
	srv := &http.Server{
		Addr:    a.config.ServerAddress,
		Handler: a.router,
	}

	quit := make(chan os.Signal, 1)
	///syscall.SIGINT -> ctrol c
	///syscall.SIGTERM -> kill service
	///syscall.SIGHUP -> reload service
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		logger.Log.Info().Msgf("server is running %s ", a.config.ServerAddress)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			logger.Log.Fatal().Err(err).Msg("fail to start server")
		}
	}()
	<-quit
	logger.Log.Warn().Msg("shut down signal received ...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Fatal().Err(err).Msg("server forced to shutdown:")
	}
	logger.Log.Info().Msg("Sever exited gracefully")
	return nil
}
func getModulRoutes(modules []Module) []routes.Route {
	routeList := make([]routes.Route, len(modules))
	for i, module := range modules {
		routeList[i] = module.Routes()
	}
	return routeList
}
