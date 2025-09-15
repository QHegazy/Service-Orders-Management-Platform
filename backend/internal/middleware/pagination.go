package middleware

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

func PaginationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		limitStr := c.DefaultQuery("size", "15")
		offsetStr := c.DefaultQuery("page", "0")

		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			log.Printf("PaginationMiddleware - Invalid limit parameter: %v", err)
			limit = 10
		}
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			log.Printf("PaginationMiddleware - Invalid offset parameter: %v", err)
			offset = 0
		}
		c.Set("size", int32(limit))
		c.Set("page", int32(offset))
		c.Next()
	}
}
