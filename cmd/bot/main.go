package main

import (
	"log"

	"GLPITGBOT/internal/config"
	"GLPITGBOT/internal/repository"
	"GLPITGBOT/internal/telegram"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	if err := repository.ConnectPostgres(cfg); err != nil {
		log.Fatal(err)
	}
	defer repository.Pool.Close()

	telegram.StartBot(cfg)
}
