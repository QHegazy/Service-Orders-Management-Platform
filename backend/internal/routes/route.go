package routes

import (
	v1_routes "backend/internal/routes/v1"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes() *gin.Engine {
	gin_handler := gin.Default()
	v1_routes.V1RoutesRegister(gin_handler)

	return gin_handler

}
