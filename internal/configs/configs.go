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

	return config
}
