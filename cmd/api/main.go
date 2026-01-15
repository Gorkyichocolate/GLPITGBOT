package main

import (
	"log"

	"GLPITGBOT/internal/config"
	"GLPITGBOT/internal/http"
	"GLPITGBOT/internal/repository"

	"github.com/joho/godotenv"
)

func main() {
	// 1️⃣ Загружаем .env (локально / dev)
	_ = godotenv.Load()

	// 2️⃣ Загружаем конфигурацию
	cfg := config.Load()

	// 3️⃣ Подключаемся к Postgres
	if err := repository.ConnectPostgres(cfg); err != nil {
		log.Fatal("failed to connect postgres:", err)
	}
	defer repository.Pool.Close()

	// 4️⃣ Запускаем HTTP API (Gin)
	if err := http.StartServer(cfg); err != nil {
		log.Fatal("http server error:", err)
	}
}
