package config

import (
	"log"
	"time"
)

const (
	AccessTokenLifetime  = 24 * time.Hour
	RefreshTokenLifetime = 7 * 24 * time.Hour
)

type JWTConfig struct {
	Secret        string
	RefreshSecret string
}

func loadJWT() JWTConfig {
	secret := requireEnv("JWT_SECRET")
	refresh := requireEnv("JWT_REFRESH_SECRET")
	if secret == refresh {
		log.Fatal("JWT_REFRESH_SECRET must differ from JWT_SECRET")
	}
	return JWTConfig{Secret: secret, RefreshSecret: refresh}
}
