package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration values
type Config struct {
	MongoURI      string
	TelegramToken string
	PythonPath    string
	ScraperPath   string
	ServerPort    string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	if err := godotenv.Load("../../.env"); err != nil {
		return nil, err
	}

	return &Config{
		MongoURI:      os.Getenv("MONGO_ATLAS_URI"),
		TelegramToken: os.Getenv("ZENHA_TELEGRAM_TOKEN"),
		PythonPath:    os.Getenv("PYTHON_PATH"),
		ScraperPath:   os.Getenv("SCRAPER_PATH"),
		ServerPort:    os.Getenv("SERVER_PORT"),
	}, nil
}
