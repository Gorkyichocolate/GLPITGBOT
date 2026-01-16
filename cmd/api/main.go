package main

import (
	"log"

	"GLPITGBOT/internal/config"
	"GLPITGBOT/internal/http"
	"GLPITGBOT/internal/repository"

	"github.com/joho/godotenv"
)

func main() {
	// Загружаем .env (локально / dev)
	_ = godotenv.Load()

	// Загружаем конфигурацию
	cfg := config.Load()

	// подключаемся к Postgres
	if err := repository.ConnectPostgres(cfg); err != nil {
		log.Fatal("failed to connect postgres:", err)
	}
	defer repository.Pool.Close()

	// Запускаем HTTP API (Gin)
	if err := http.StartServer(cfg); err != nil {
		log.Fatal("http server error:", err)
	}
}
