package db

import (
	"context"
	"go-rest-template/pkg/config"
	"go-rest-template/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectToDB() *pgxpool.Pool {
	url := config.APP().DB_URL
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		logger.Panic("Invalid DB config: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		logger.Panic("Unable to connect to database pool: %v", err)
	}

	return pool
}
