package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/duitocdai/internal/config"
	"github.com/duitocdai/internal/db"
	"github.com/duitocdai/internal/db/sqlc"
	"github.com/duitocdai/internal/routes"
	"github.com/duitocdai/internal/validation"
	"github.com/duitocdai/pkg/auth"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

func NewApplication(cfg *config.Config) *Application {
	if err := validation.InitValidator(); err != nil {
		log.Fatalf("Validator init failed %v", err)
	}
	LoadEnv()
	r := gin.Default()
	if err := db.InitDB(); err != nil {
		log.Fatalf("database init false: %v", err)
	}
	redisClient := config.NewRedisClient()
	tokenService := auth.NewJWTService()
	ctx := &ModuleContext{
		DB:    db.DB,
		Redis: redisClient,
	}
	modules := []Module{
		NewUserModule(ctx),
		NewAuthModule(ctx, tokenService),
	}
	routes.RegisterRoutes(r, getModulRoutes(modules)...)
	return &Application{
		config:  cfg,
		router:  r,
		modules: modules,
	}
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
		log.Printf("server is running %s \n", a.config.ServerAddress)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("fail to start server: %v", err)
		}
	}()
	<-quit
	log.Println("shut down signal received ...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}
	log.Println("server exited gracefully")
	return nil
}
func getModulRoutes(modules []Module) []routes.Route {
	routeList := make([]routes.Route, len(modules))
	for i, module := range modules {
		routeList[i] = module.Routes()
	}
	return routeList
}

func LoadEnv() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Println("No .env file found")
	}
}
