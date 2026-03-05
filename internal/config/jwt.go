package config

type JWTConfig struct {
	SecretKey string
	RefreshSecretKey string
	ExpiresIn int
	RefreshExpiresIn int
}

func loadJWT() JWTConfig {
	return JWTConfig{
		SecretKey: requireEnv("JWT_SECRET_KEY"),
		RefreshSecretKey: requireEnv("JWT_REFRESH_SECRET_KEY"),
		ExpiresIn: getEnvInt("JWT_EXPIRES_IN", 7200),
		RefreshExpiresIn: getEnvInt("JWT_REFRESH_EXPIRES_IN", 604800),
	}
}