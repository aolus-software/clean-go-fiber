package config

type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

func loadRedis() RedisConfig {
	return RedisConfig{
		Host: requireEnv("REDIS_HOST"),
		Port: requireEnv("REDIS_PORT"),
		Password: getEnv("REDIS_PASSWORD", ""),	
	}
}