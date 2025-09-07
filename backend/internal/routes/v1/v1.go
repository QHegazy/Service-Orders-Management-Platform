package v1_routes

import (
	"backend/utils"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func V1RoutesRegister(r *gin.Engine) {
	v1 := r.Group("/v1")
	v1.POST("/refresh", func(c *gin.Context) {
		// 1️⃣ Get tokens
		RefreshToken, _ := c.Cookie("token")
		accessToken := c.GetHeader("Authorization")
		if accessToken == "" || len(accessToken) <= 7 || accessToken[:7] != "Bearer " {
			c.JSON(401, gin.H{"error": "Authorization header missing or invalid"})
			return
		}
		accessToken = accessToken[7:]

		decodedRefreshToken, err := utils.DecodeToken(RefreshToken)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid refresh token"})
			return
		}

		_, err = utils.ValidateToken(decodedRefreshToken)
		if err != nil {
			fmt.Println(err)
			c.JSON(401, gin.H{"error": "Invalid token"})
			return
		}

		_, err = utils.ValidateToken(accessToken)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				claims, err := utils.DecodeJwtClaims(accessToken)
				if err != nil {
					c.JSON(401, gin.H{"error": "Invalid access token claims"})
					return
				}
				newToken, err := utils.GenerateToken(utils.EntityData{
					ID:       claims.Data.ID,
					Username: claims.Data.Username,
					Belong:   claims.Data.Belong,
					Role:     claims.Data.Role,
				}, 15, "access")
				if err != nil {
					c.JSON(500, gin.H{"error": "Failed to generate new access token"})
					return
				}
				c.JSON(200, utils.SuccessResponse("success", gin.H{"access_token": newToken}))
				return
			}
			c.JSON(401, gin.H{"error": "Invalid access token"})
			return
		}
	})
	v1.GET("", func(c *gin.Context) {
		c.String(200, "OK")
	})
	userRoutes(v1)

	customerRoutes(v1)
	ticketRoutes(v1)

}
