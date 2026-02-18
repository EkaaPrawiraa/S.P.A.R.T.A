package configs

import (
	"os"
)

type Config struct {
	DatabaseURL string
	OpenAIKey   string
	Port        string
}

func LoadConfig() *Config {
	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/gymdb?sslmode=disable"),
		OpenAIKey:   getEnv("OPENAI_API_KEY", ""),
		Port:        getEnv("APP_PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
