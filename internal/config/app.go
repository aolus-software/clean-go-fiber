package config

type AppConfig struct {
	Name string
	Port string
	Env  string
	Key  string
	URL	 string
}

func loadApp() AppConfig {
	return AppConfig{
		Name: requireEnv("APP_NAME"),
		Port: getEnv("APP_PORT", "8080"),
		Env:  getEnv("APP_ENV", "development"),
		Key:  requireEnv("APP_KEY"),
		URL:  requireEnv("APP_URL"),
	}
}

