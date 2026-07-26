// Package config
package config

import "os"

type Config struct {
	DatabaseURL string
	JWTSecret   string
}

var defaults = map[string]string{
	"DATABASE_URL": "postgres://postgres:maiphuong@localhost:5432/wodo_app_db?sslmode=disable",
	"JWT_SECRET":   "the-secret-jwt-key",
}

func getEnv(key string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return defaults[key]
}

func Load() *Config {
	return &Config{
		DatabaseURL: getEnv("DATABASE_URL"),
		JWTSecret:   getEnv("JWT_SECRET"),
	}
}
