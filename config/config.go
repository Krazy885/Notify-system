package config

import (
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	DiscordToken  string
	TelegramToken string
	DBUser        string
	DBPass        string
	DBHost        string
	DBName        string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	return &Config{
		DiscordToken:  os.Getenv("DISCORD_TOKEN"),
		TelegramToken: os.Getenv("TELEGRAM_TOKEN"),
		DBUser:        os.Getenv("DB_USER"),
		DBPass:        os.Getenv("DB_PASS"),
		DBHost:        os.Getenv("DB_HOST"),
		DBName:        os.Getenv("DB_NAME"),
	}, nil
}