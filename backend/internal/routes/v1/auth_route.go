package v1_routes

import (
	v1_controllers "backend/internal/controllers/v1"
	"backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func authRoutes(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	auth.Use(middleware.ValidationErrorHandler())
	auth.Use(middleware.DBErrorHandler())
	authController := v1_controllers.NewAuthControllerV1()
	auth.POST("/login", authController.Login)
	auth.POST("/logout", middleware.AuthMiddleware(), authController.Logout)
	auth.POST("/refresh", authController.Refresh)
}
