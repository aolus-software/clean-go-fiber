package config

type MailConfig struct {
	Host   string
	Port   int
	Secure bool
	From   string
	User   string
	Pass   string
}

func loadMail() MailConfig {
	return MailConfig{
		Host:   getEnv("MAIL_HOST", "smtp.example.com"),
		Port:   getEnvPort("MAIL_PORT", 587),
		Secure: getEnvBool("MAIL_SECURE", false),
		From:   getEnv("MAIL_FROM", "Clean Go Fiber <your_email@example.com>"),
		User:   getEnv("MAIL_USER", "your_email@example.com"),
		Pass:   getEnv("MAIL_PASS", "your_email_password"),
	}
}
