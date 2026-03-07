package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Mail     MailConfig
	Cors     CorsConfig
	Jwt      JWTConfig
}

var C *Config

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found, reading from environment: %v", err)
	}

	C = &Config{
		App:      loadApp(),
		Database: loadDatabase(),
		Redis:    loadRedis(),
		Mail:     loadMail(),
		Cors:     loadCors(),
		Jwt:      loadJWT(),
	}
}

func getEnv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}

func requireEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("Environment variable %s is required but not set", key)
	}
	return v
}

func getEnvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		log.Fatalf("Environment variable %s is not a valid integer", key)
	}
	return i
}
