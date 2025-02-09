package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DiscordKey string
	OpenAIKey  string
	NewsKey    string
	SportsKey  string
}

var AppConfig *Config

func LoadConfig() *Config {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cfg := &Config{
		DiscordKey: os.Getenv("DISCORD_KEY"),
		OpenAIKey:  os.Getenv("OPENAI_KEY"),
		NewsKey:    os.Getenv("NEWS_KEY"),
		SportsKey:  os.Getenv("SPORTS_RADAR_KEY"),
	}

	if cfg.OpenAIKey == "" || cfg.NewsKey == "" || cfg.SportsKey == "" || cfg.DiscordKey == "" {
		log.Fatal("Missing one or more API keys in environment variables.")
	}

	AppConfig = cfg
	return cfg
}
