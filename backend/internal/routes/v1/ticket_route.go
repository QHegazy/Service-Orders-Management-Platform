package v1_routes

import (
	v1_controllers "backend/internal/controllers/v1"
	"backend/internal/controllers/v1/ws"
	"backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func ticketRoutes(r *gin.RouterGroup) {

	ticket := r.Group("/ticket")
	ticket.Use(middleware.ValidationErrorHandler())
	ticket.Use(middleware.DBErrorHandler())
	ticketController := v1_controllers.NewTicketControllerV1()
	ticket.POST("", ticketController.CreateTicket)
	ticket.PUT("/:id", ticketController.UpdateTicket)
	ticket.DELETE("/:id", ticketController.DeleteTicket)
	ticket.GET("/connect", ws.HandleWS)
	// ticket.GET("", ticketController.ListTickets)
	// ticket.GET("/comments/:id", ticketController.ListCommentsByTicketID)
}
