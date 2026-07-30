// Package config
package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL            string
	JWTSecret              string
	AccessTokenExpiration  time.Duration
	RefreshTokenExpiration time.Duration
}

var defaults = map[string]string{
	"DATABASE_URL":             "postgres://postgres:maiphuong@localhost:5432/wodo_app_db?sslmode=disable",
	"JWT_SECRET":               "the-secret-jwt-key",
	"ACCESS_TOKEN_EXPIRATION":  "15m",
	"REFRESH_TOKEN_EXPIRATION": "720h",
}

func getEnv(key string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return defaults[key]
}

func parseDurationEnv(key string) time.Duration {
	valStr := getEnv(key)
	duration, err := time.ParseDuration(valStr)
	if err != nil {
		log.Printf("Warning: invalid duration for %s, using fallback", key)
		duration, _ = time.ParseDuration(defaults[key])
	}
	return duration
}

func Load() *Config {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("Note: .env file not found, using environment variables or defaults")
	}
	return &Config{
		DatabaseURL:            getEnv("DATABASE_URL"),
		JWTSecret:              getEnv("JWT_SECRET"),
		AccessTokenExpiration:  parseDurationEnv("ACCESS_TOKEN_EXPIRATION"),
		RefreshTokenExpiration: parseDurationEnv("REFRESH_TOKEN_EXPIRATION"),
	}
}
