package config

import (
	"log"
	"os"
	"slices"
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

func LoadDatabaseOnly() {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found, reading from environment: %v", err)
	}

	C = &Config{
		App:      loadApp(),
		Database: loadDatabase(),
	}
}

func getEnv(key, fallback string) string {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
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

func getEnvPort(key string, fallback int) int {
	p := getEnvInt(key, fallback)
	if p < 1 || p > 65535 {
		log.Fatalf("Environment variable %s must be a port between 1 and 65535", key)
	}
	return p
}

func getEnvBool(key string, fallback bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		log.Fatalf("Environment variable %s must be true or false", key)
	}
	return b
}

func getEnvChoice(key, fallback string, choices ...string) string {
	v := getEnv(key, fallback)
	if !slices.Contains(choices, v) {
		log.Fatalf("Environment variable %s must be one of %v", key, choices)
	}
	return v
}
