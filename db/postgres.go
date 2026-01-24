package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {
	// Загружаем .env из рабочей директории
	if err := godotenv.Load(); err != nil {
		log.Fatal("❌ .env file not found")
	}

	// Читаем переменные
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")

	// Защита от пустого порта (твоя текущая ошибка)
	if port == "" {
		log.Fatal("❌ DB_PORT is empty (check .env or working directory)")
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, name,
	)

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("❌ sql.Open error:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("❌ DB.Ping error:", err)
	}

	log.Println("✅ PostgreSQL connected")
}
