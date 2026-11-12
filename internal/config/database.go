package config

import (
	"log"
	"time"
)

type DatabaseConfig struct {
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
	SlowQuery       time.Duration
}

func loadDatabase() DatabaseConfig {
	maxOpen := getEnvInt("DB_POOL_MAX_OPEN", 25)
	if maxOpen < 1 {
		log.Fatal("Environment variable DB_POOL_MAX_OPEN must be at least 1")
	}
	maxIdle := getEnvInt("DB_POOL_MAX_IDLE", 10)
	if maxIdle < 0 {
		log.Fatal("Environment variable DB_POOL_MAX_IDLE must not be negative")
	}
	slowMS := getEnvInt("DB_SLOW_QUERY_MS", 200)
	if slowMS < 0 {
		log.Fatal("Environment variable DB_SLOW_QUERY_MS must not be negative")
	}

	return DatabaseConfig{
		URL:             getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/clean_go?sslmode=disable"),
		MaxOpenConns:    maxOpen,
		MaxIdleConns:    maxIdle,
		ConnMaxLifetime: getEnvDuration("DB_POOL_CONN_MAX_LIFETIME", 5*time.Minute),
		ConnMaxIdleTime: getEnvDuration("DB_POOL_CONN_MAX_IDLE_TIME", 0),
		SlowQuery:       time.Duration(slowMS) * time.Millisecond,
	}
}
