package config

import (
	"log"
	"time"
)

type AppConfig struct {
	Name          string
	Port          int
	Env           string
	Timezone      *time.Location
	LogLevel      string
	EnableAPIDocs bool
	ClientURL     string
}

func loadApp() AppConfig {
	tzName := getEnv("APP_TIMEZONE", "UTC")
	tz, err := time.LoadLocation(tzName)
	if err != nil {
		log.Fatalf("Environment variable APP_TIMEZONE is not a valid IANA timezone: %s", tzName)
	}

	return AppConfig{
		Name:          getEnv("APP_NAME", "Clean Go Fiber"),
		Port:          getEnvPort("APP_PORT", 3000),
		Env:           getEnvChoice("APP_ENV", "development", "development", "staging", "production"),
		Timezone:      tz,
		LogLevel:      getEnvChoice("LOG_LEVEL", "info", "info", "warn", "debug"),
		EnableAPIDocs: getEnvBool("ENABLE_API_DOCS", false),
		ClientURL:     getEnv("CLIENT_URL", "http://localhost:3000"),
	}
}
