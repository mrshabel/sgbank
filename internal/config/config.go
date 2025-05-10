package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerAddr  string
	Env         ENV
	DatabaseURL string
}

type ENV string

const (
	DEV  ENV = "local"
	PROD ENV = "production"
)

func New() *Config {
	// load env config
	if err := godotenv.Load(); err != nil {
		log.Print("failed to load .env: ", err)
	}
	return &Config{
		ServerAddr:  getEnv("SERVER_ADDR", "localhost:8000"),
		Env:         ENV(getEnv("ENV", string(DEV))),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@postgres/sgbank?sslmode=disable"),
	}
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
