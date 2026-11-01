package config

type DatabaseConfig struct {
	URL string
}

func loadDatabase() DatabaseConfig {
	return DatabaseConfig{
		URL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/clean_go?sslmode=disable"),
	}
}
