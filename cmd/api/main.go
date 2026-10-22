package main

import (
	"path/filepath"

	"github.com/duitocdai/internal/app"
	"github.com/duitocdai/internal/config"
	"github.com/duitocdai/internal/utils"
	"github.com/duitocdai/pkg/logger"
	"github.com/joho/godotenv"
)

func main() {
	rootDir := utils.MustGetWorkingDir()

	logFile := filepath.Join(rootDir, "internal/logs/app.log")
	logger.InitLogger(logger.LoggerConfig{
		Level:      "info",
		Filename:   logFile,
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     5,
		Compress:   true,
		IsDev:      utils.GetEnv("APP_ENV", "development"),
	})
	if err := godotenv.Load(filepath.Join(rootDir, ".env")); err != nil {
		logger.Log.Warn().Msg("Unable to load .env")
	} else {
		logger.Log.Info().Msg("Loaded .env api proccess successfully")
	}
	///Initialize configuration
	cfg := config.NewConfig()

	///initialize application
	application, err := app.NewApplication(cfg)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to initialize application")
	}
	///start server
	if err := application.Run(); err != nil {
		logger.Log.Fatal().Err(err).Msg("Application run fail")

	}
}
