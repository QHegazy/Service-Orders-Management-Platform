package v1_routes

import (
	v1_controllers "backend/internal/controllers/v1"
	"backend/internal/controllers/v1/ws"
	"backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func ticketRoutes(r *gin.RouterGroup, hub *ws.Hub) {

	ticket := r.Group("/ticket")
	ticket.Use(middleware.ValidationErrorHandler())
	ticket.Use(middleware.DBErrorHandler())
	ticket.Use(middleware.AuthMiddleware())
	ticketController := v1_controllers.NewTicketControllerV1()
	ticket.POST("", middleware.RoleMiddleware("Customer"), ticketController.CreateTicket)
	ticket.PUT("/:id", middleware.RoleMiddleware("Admin"), ticketController.UpdateTicket)
	ticket.DELETE("/:id", middleware.RoleMiddleware("Admin"), ticketController.DeleteTicket)
	ticket.GET("/user", middleware.PaginationMiddleware(), middleware.RoleMiddleware("Admin", "Technician"), ticketController.ListTicketsByUserId)
	ticket.GET("/:id", ticketController.GetTicket)
	ticket.GET("/:id/comments", ticketController.ListCommentsByTicketID)
	ticket.GET("/customer", middleware.PaginationMiddleware(), middleware.RoleMiddleware("Customer"), ticketController.ListTicketsByCustomerId)

	r.GET("/ws/ticket/:ticketId", func(c *gin.Context) {
		ws.ServeWs(hub, c)
	})
}
