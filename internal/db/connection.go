package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/duitocdai/internal/config"
	"github.com/duitocdai/internal/db/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

var DB sqlc.Querier

func InitDB() error {
	connStr := config.NewConfig().DNS()

	conf, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return fmt.Errorf("error passing db config: %v", err)
	}

	conf.MaxConns = 50
	conf.MinConns = 5
	conf.MaxConnLifetime = 30 * time.Minute
	conf.MaxConnIdleTime = 5 * time.Minute
	conf.HealthCheckPeriod = 1 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	DBPool, err := pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		return fmt.Errorf("error creating db pool: %v", err)
	}
	DB = sqlc.New(DBPool)
	if err := DBPool.Ping(ctx); err != nil {
		return fmt.Errorf("ping error : %v", err)
	}
	log.Println("connected")
	return nil
}
