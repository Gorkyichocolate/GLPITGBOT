package main

import (
	"GLPITGBOT/db"
	"GLPITGBOT/https"
	"GLPITGBOT/telegram"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	db.Connect()

	go telegram.Bot()
	r := gin.Default()
	if err := r.SetTrustedProxies([]string{"127.0.0.1", "::1"}); err != nil {
		log.Println("set trusted proxies error:", err)
	}

	api := r.Group("/api")
	{
		api.GET("/auth/check", https.CheckAuth)
		api.POST("/ticket/create", https.CreateTicket)
		api.POST("/ticket/update", https.NotifyTicketUpdate)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "7000"
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
