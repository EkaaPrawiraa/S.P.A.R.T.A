package configs

import (
	"os"
	"strconv"
)

type Config struct {
	DatabaseURL string
	OpenAIKey   string
	OpenAIModel string
	OpenAIBase  string
	RedisAddr   string
	RedisPass   string
	RedisDB     int
	JWTSecret   string
	Port        string
	AppEnv      string
}

func LoadConfig() *Config {
	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/gymdb?sslmode=disable"),
		OpenAIKey:   getEnv("OPENAI_API_KEY", ""),
		OpenAIModel: getEnv("OPENAI_MODEL", "gpt-4o-mini"),
		OpenAIBase:  getEnv("OPENAI_BASE_URL", "https://api.openai.com"),
		RedisAddr:   getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPass:   getEnv("REDIS_PASSWORD", ""),
		RedisDB:     getEnvInt("REDIS_DB", 0),
		JWTSecret:   getEnv("JWT_SECRET", "SUPER_SECRET"),
		Port:        getEnv("APP_PORT", "8080"),
		AppEnv:      getEnv("APP_ENV", "local"),
	}
}


func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	valStr, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return fallback
	}
	return val
}
