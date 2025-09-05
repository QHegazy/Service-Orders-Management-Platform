package v1_routes

import (
	"github.com/gin-gonic/gin"
)

func V1RoutesRegister(r *gin.Engine) {
	v1 := r.Group("/v1")

	v1.GET("", func(c *gin.Context) {
		c.String(200, "OK")
	})
	userRoutes(v1)
}
