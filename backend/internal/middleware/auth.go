package middleware

import (
	"backend/utils"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse("error", "Authorization header is required"))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse("error", "Invalid Authorization header format"))
			return
		}

		tokenString := parts[1]
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse("error", fmt.Sprintf("Invalid token: %v", err)))
			return
		}
		c.Set("userID", claims.Data.ID)
		c.Set("username", claims.Data.Username)
		c.Set("role", claims.Data.Role)
		c.Set("tenantID", claims.Data.Belong)

		c.Next()
	}
}
