package main

import (
	"github.com/duitocdai/internal/app"
	"github.com/duitocdai/internal/config"
)

func main() {
	///Initialize configuration
	cfg := config.NewConfig()

	///initialize application
	application := app.NewApplication(cfg)
	///start server
	if err := application.Run(); err != nil {
		panic(err)
	}
}
