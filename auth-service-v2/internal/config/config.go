package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port       string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPass     string
	DBName     string
	JWTSecret  string
	JWTTtlMins int
}

func NewConfigFromEnv() *Config {
	ttl := 60
	if v := os.Getenv("JWT_TTL_MINUTES"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			ttl = n
		}
	}

	return &Config{
		Port:       getEnv("PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPass:     getEnv("DB_PASS", "postgres"),
		DBName:     getEnv("DB_NAME", "authdb"),
		JWTSecret:  getEnv("JWT_SECRET", "secret"),
		JWTTtlMins: ttl,
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
