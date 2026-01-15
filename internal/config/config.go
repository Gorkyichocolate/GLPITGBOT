package config

import (
	"log"
	"os"
)

type Config struct {
	// Telegram
	TGBotToken string

	//GLPI
	GLPI string

	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

// Load загружает конфигурацию из environment variables
func Load() *Config {
	cfg := &Config{
		TGBotToken: getEnv("TG_BOT_TOKEN", ""),
		GLPI:       getEnv("GLPI", ""),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", ""),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", ""),
	}

	validate(cfg)
	return cfg
}

func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func validate(cfg *Config) {
	if cfg.TGBotToken == "" {
		log.Fatal("TG_BOT_TOKEN is required")
	}

	if cfg.DBUser == "" || cfg.DBPassword == "" || cfg.DBName == "" {
		log.Fatal("Database configuration is incomplete")
	}
}
