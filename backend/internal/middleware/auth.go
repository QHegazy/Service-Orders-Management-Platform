package middleware

import (
	"backend/utils"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("AuthMiddleware - Processing request for path: %s", c.Request.URL.Path)
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Printf("AuthMiddleware - Missing Authorization header for path: %s", c.Request.URL.Path)
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse("error", "Authorization header is required"))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse("error", "Invalid Authorization header format"))
			log.Printf("AuthMiddleware - Invalid Authorization header format for path %s: %s", c.Request.URL.Path, authHeader)
			return
		}

		tokenString := parts[1]
		log.Printf("AuthMiddleware - Validating token for path: %s", c.Request.URL.Path)
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse("error", fmt.Sprintf("Invalid token: %v", err)))
			log.Printf("AuthMiddleware - Token validation error for path %s: %v", c.Request.URL.Path, err)
			return
		}

		log.Printf("AuthMiddleware - Token validated successfully for user: %s, role: %s", claims.Data.Username, claims.Data.Role)
		fmt.Println("claims")
		c.Set("userID", claims.Data.ID)
		c.Set("username", claims.Data.Username)
		c.Set("role", claims.Data.Role)
		c.Set("tenantID", claims.Data.Belong)

		c.Next()
	}
}
