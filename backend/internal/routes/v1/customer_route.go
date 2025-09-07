package v1_routes

import (
	v1_controllers "backend/internal/controllers/v1"
	"backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func customerRoutes(r *gin.RouterGroup) {

	customer := r.Group("/customer")
	customer.Use(middleware.ValidationErrorHandler())
	customer.Use(middleware.DBErrorHandler())
	customerController := v1_controllers.NewCustomerControllerV1()
	customer.POST("", customerController.CreateCustomer)
	customer.POST("/login",customerController.LoginCustomer)
	customer.PUT("", middleware.RoleMiddleware("Customer"), customerController.UpdateCustomer)
	customer.DELETE("", middleware.RoleMiddleware("Customer"), customerController.DeleteCustomer)
	// customer.GET("/:id", customerController.GetCustomerByID)
	// customer.GET("", customerController.ListAllCustomers)

}
