package main

import (
	"StudyHub/backend/internal/config"
	"StudyHub/backend/pgk/postgres"
	"context"
	"fmt"
)

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbConnString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName)

	_ = postgres.New(ctx, dbConnString)

}
