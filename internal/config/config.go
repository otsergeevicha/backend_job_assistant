package config

import (
	"os"
	"time"
)

type Config struct {
	BotToken       string
	SessionSecret  string
	SessionTTL     time.Duration
	TelegramMaxAge time.Duration
}

func Load() Config {
	return Config{
		BotToken:       mustEnv("TELEGRAM_BOT_TOKEN"),
		SessionSecret:  mustEnv("SESSION_SECRET"),
		SessionTTL:     7 * 24 * time.Hour,
		TelegramMaxAge: 24 * time.Hour,
	}
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(key + " is required")
	}
	return v
}
