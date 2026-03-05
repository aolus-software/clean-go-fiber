package config

import "strings"

type CorsConfig struct {
	Host []string
	Method []string
}

func loadCors() CorsConfig {
	host := strings.Split(getEnv("CORS_ALLOWED_HOST", "*"), `,`)
	method := strings.Split(getEnv("CORS_ALLOWED_METHOD", "GET,POST,PUT,PATCH,DELETE,HEAD"), `,`)

	return CorsConfig{
		Host: host,
		Method: method,
	}
}