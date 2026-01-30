package config

import (
	"log"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string `env:"DB_HOST"`
	DBPort     string `env:"DB_PORT"`
	DBUser     string `env:"DB_USER"`
	DBPass     string `env:"DB_PASS"`
	DBName     string `env:"DB_NAME"`
	JwtKey     string `env:"JWT_KEY"`
	BucketName string `env:"AWS_S3_BUCKET"`
	AWS_S3_URL string `env:"AWS_S3_URL"`
}

func Load() Config {
	var cnf Config
	// Load .env file if it exists (for local dev), but don't fail if it doesn't (for Docker)
	_ = godotenv.Load(".env")

	err := env.Parse(&cnf)
	if err != nil {
		log.Fatal("cannot parse environment variables", err)
	}
	return cnf
}
