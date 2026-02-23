package main

import (
	"GLPITGBOT/db"
	"GLPITGBOT/https"
	"GLPITGBOT/telegram"

	"github.com/gin-gonic/gin"
)

func main() {
	db.Connect()

	go telegram.Bot()
	r := gin.Default()

	api := r.Group("/api")
	{
		api.GET("/auth/check", https.CheckAuth)
		api.POST("/ticket/create", https.CreateTicket)
		api.POST("/ticket/update", https.NotifyTicketUpdate)
	}

	r.Run(":8080")
}
