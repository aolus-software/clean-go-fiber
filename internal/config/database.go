package config

type DatabaseConfig struct {
	URL string
}

type DatabaseEnvConfig struct {
	HOST    string
	PORT    string
	USER    string
	PASSWORD string
	DBNAME  string
	SSLMODE string
	TIMEZONE string
}

func loadDatabaseUrl() string {
	host := requireEnv("DB_HOST")
	port := getEnv("DB_PORT", "5432")
	user := requireEnv("DB_USER")
	password := requireEnv("DB_PASSWORD")
	dbname := requireEnv("DB_NAME")
	sslmode := getEnv("DB_SSLMODE", "disable")
	timezone := getEnv("DB_TIMEZONE", "UTC")

	return "postgres://" + user + ":" + password + "@" + host + ":" + port + "/" + dbname + "?sslmode=" + sslmode + "&timezone=" + timezone
}

func loadDatabase() DatabaseConfig {
	return DatabaseConfig{
		URL: loadDatabaseUrl(),
	}
}