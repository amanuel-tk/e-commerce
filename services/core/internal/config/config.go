package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	RedisAddr   string
	LogLevel    string
	LogFormat   string
}

func LoadConfig() (*Config, error) {

	_ = godotenv.Load()

	cfg := &Config{
		Port:        getEnv("CORE_PORT", "3000"),
		DatabaseURL: os.Getenv("CORE_DATABASE_URL"),
		RedisAddr:   getEnv("CORE_REDIS_ADDR", "localhost:6379"),
		LogLevel:    getEnv("CORE_LOG_LEVEL", "info"),
		LogFormat:   getEnv("CORE_LOG_FORMAT", "text"),
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	return cfg, nil

}

func getEnv(key string, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}
