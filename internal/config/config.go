package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken       string
	SessionSecret  string
	ServerPort     string
	SessionTTL     time.Duration
	TelegramMaxAge time.Duration
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		BotToken:       mustEnv("TELEGRAM_TOKEN"),
		SessionSecret:  mustEnv("SECRET_KEY"),
		ServerPort:     mustEnv("SERVER_PORT"),
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
