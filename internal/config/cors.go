package config

import "strings"

type CorsConfig struct {
	AllowedOrigins string
}

func loadCors() CorsConfig {
	raw := getEnv("ALLOWED_HOST", "")
	if strings.TrimSpace(raw) == "" {
		return CorsConfig{AllowedOrigins: "*"}
	}
	origins := strings.Split(raw, ",")
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}
	return CorsConfig{AllowedOrigins: strings.Join(origins, ",")}
}
