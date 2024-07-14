package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

type Config struct {
	HTTP_PORT     string
	GRPC_PORT     string
	DB_HOST       string
	DB_PORT       int
	DB_USER       string
	DB_NAME       string
	DB_PASSWORD   string
	ACCESS_TOKEN  string
	REFRESH_TOKEN string
}

func coalesce(env string, defaultValue interface{}) interface{} {
	value, exists := os.LookupEnv(env)
	if !exists {
		return defaultValue
	}
	return value
}

func Load() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := Config{}

	cfg.HTTP_PORT = cast.ToString(coalesce("HTTP_PORT", ":8081"))
	cfg.GRPC_PORT = cast.ToString(coalesce("GRPC_PORT", ":50050"))
	cfg.DB_HOST = cast.ToString(coalesce("DB_HOST", "localhost"))
	cfg.DB_PORT = cast.ToInt(coalesce("DB_PORT", 5432))
	cfg.DB_USER = cast.ToString(coalesce("DB_USER", "postgres"))
	cfg.DB_NAME = cast.ToString(coalesce("DB_NAME", "postgres"))
	cfg.DB_PASSWORD = cast.ToString(coalesce("DB_PASSWORD", "password"))
	cfg.ACCESS_TOKEN = cast.ToString(coalesce("ACCESS_TOKEN", "my_secret_key"))
	cfg.REFRESH_TOKEN = cast.ToString(coalesce("REFRESH_TOKEN", "my_secret_key"))

	return cfg
}
