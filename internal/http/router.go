package http

import (
	"GLPITGBOT/internal/http/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/health", handlers.Health)

	r.POST("/glpi/ticket", handlers.CreateTicket)

	r.GET("/glpi/ticket/:id?view=status", handlers.TicketStatus)

	r.GET("/glpi/ticket/:id?view=history", handlers.GetLastTickets)

	return r
}
