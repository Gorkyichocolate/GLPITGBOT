package main

import (
	"GLPITGBOT/db"
	"GLPITGBOT/https"
	"GLPITGBOT/telegram"

	"github.com/gin-gonic/gin"
)

func main() {
	db.Connect()

	telegram.Bot()
	r := gin.Default()

	api := r.Group("/api")
	{
		api.GET("/auth/check", https.CheckAuth)
	}

	r.Run(":8080")
}
