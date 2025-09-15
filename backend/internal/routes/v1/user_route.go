package v1_routes

import (
	v1_controllers "backend/internal/controllers/v1"
	"backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func userRoutes(r *gin.RouterGroup) {

	user := r.Group("/user")
	user.Use(middleware.ValidationErrorHandler())
	user.Use(middleware.DBErrorHandler())
	userController := v1_controllers.NewUserControllerV1()

	// Public routes
	user.POST("", userController.CreateUser)

	authUser := user.Group("")
	authUser.Use(middleware.AuthMiddleware())
	authUser.Use(middleware.RoleMiddleware("Admin"))

	// Admin only routes
	user.POST("/technician", userController.CreateTechnicianUser)
	user.GET("/technicians", userController.GetTechnicians)
	user.PUT("", userController.UpdateUser)
	user.DELETE("/:id", userController.DeleteUser)

}
