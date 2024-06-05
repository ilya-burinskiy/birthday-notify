package configs

import (
	"os"
	"time"
)

const AuthTokenExp = 24 * time.Hour
const SecretKey = "secret"

type Config struct {
	RunAddr string
	DSN     string

	SMTPAuthIdentity string
	SMTPAuthUsername string
	SMTPAuthPassword string
	SMTPHost         string
	SMTPPort         string
}

func Parse() Config {
	config := Config{
		RunAddr: "localhost:8000",
		DSN:     "postgres://birthday-notify:password@localhost:5432/birthday-notify",
	}

	if envRunAdd := os.Getenv("RUN_ADDRESS"); envRunAdd != "" {
		config.RunAddr = envRunAdd
	}
	if envDSN := os.Getenv("DATABASE_URI"); envDSN != "" {
		config.DSN = envDSN
	}

	if envSMTPAuthIdentity := os.Getenv("SMTP_AUTH_IDENTITY"); envSMTPAuthIdentity != "" {
		config.SMTPAuthIdentity = envSMTPAuthIdentity
	}
	if envSMTPAuthUsername := os.Getenv("SMTP_AUTH_USERNAME"); envSMTPAuthUsername != "" {
		config.SMTPAuthUsername = envSMTPAuthUsername
	}
	if envSMTPAuthPassword := os.Getenv("SMTP_AUTH_PASSWORD"); envSMTPAuthPassword != "" {
		config.SMTPAuthPassword = envSMTPAuthPassword
	}
	if envSMTPHost := os.Getenv("SMTP_HOST"); envSMTPHost != "" {
		config.SMTPHost = envSMTPHost
	}
	if envSMTPPort := os.Getenv("SMTP_PORT"); envSMTPPort != "" {
		config.SMTPPort = envSMTPPort
	}

	return config
}
