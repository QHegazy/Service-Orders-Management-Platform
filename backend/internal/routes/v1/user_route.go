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
	user.POST("login", userController.LoginUser)
	user.POST("", userController.CreateUser)
	user.PUT("", userController.UpdateUser)
	user.DELETE("", userController.DeleteUser)

}
