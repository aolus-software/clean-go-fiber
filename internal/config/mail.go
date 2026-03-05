package config

type MailConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	MailFrom string
}

func loadMail() MailConfig {
	return MailConfig{
		Host:     requireEnv("MAIL_HOST"),
		Port:     requireEnv("MAIL_PORT"),
		User:     requireEnv("MAIL_USER"),
		Password: requireEnv("MAIL_PASSWORD"),
		MailFrom: requireEnv("MAIL_FROM"),
	}
}