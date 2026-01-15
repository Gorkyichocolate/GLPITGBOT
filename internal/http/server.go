package http

import (
	"log"

	"GLPITGBOT/internal/config"
)

func StartServer(cfg *config.Config) error {
	r := SetupRouter()

	log.Println("HTTP server started on :8080")
	return r.Run(":8080")
}
