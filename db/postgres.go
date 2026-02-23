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
	if err := godotenv.Load(); err != nil {
		log.Fatal("❌ .env file not found")
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")

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

	if _, err = DB.Exec(`
		ALTER TABLE users
		ADD COLUMN IF NOT EXISTS session_token TEXT NOT NULL DEFAULT ''
	`); err != nil {
		log.Fatal("❌ schema sync error (users.session_token):", err)
	}

	log.Println("✅ PostgreSQL connected")
}
