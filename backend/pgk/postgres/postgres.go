package postgres

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

//shared lib for code to connect for the Postgres

// should accept Config.Postgres and based on this connect to the database

// connection string-postgresql://postgres:postgres@localhost:5431/campusfit?sslmode=disable
func New(ctx context.Context, connString string) *pgxpool.Pool {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Fatal("failure to connect to the postgres", err)
	}

	//ping to check for errors
	if err := pool.Ping(ctx); err != nil {
		log.Fatal("failure to ping the postgres", err)
	}

	return pool
}
