package v1_controllers

import (
	"backend/internal/dto"
	"backend/internal/services"
	"backend/utils"
	"context"
	"errors"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthControllerV1 interface {
	Login(c *gin.Context)
	Logout(c *gin.Context)
	Refresh(c *gin.Context)
}

type authControllerV1 struct {
	authService *services.AuthService
}

func (a *authControllerV1) Login(c *gin.Context) {
	ctx := context.Background()
	var loginDto dto.LoginDto
	if err := c.ShouldBindJSON(&loginDto); err != nil {
		log.Printf("Login - JSON binding error for request: %v", err)
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	log.Printf("Login attempt for email: %s", loginDto.Email)
	accessToken, refreshToken, err := a.authService.LoginService(ctx, loginDto)
	if err != nil {
		log.Printf("Login failed for email %s: %v", loginDto.Email, err)
		c.JSON(401, utils.ErrorResponse("error", "Invalid credentials"))
		return
	}

	log.Printf("Login successful for email: %s", loginDto.Email)

	c.SetCookie("token", refreshToken, 3600, "/", "localhost", false, true)

	if gin.Mode() == gin.ReleaseMode {
		domain := os.Getenv("DOMAIN")
		if domain == "" {
			domain = getDomain(c)
		}

		c.SetCookie(
			"token",
			refreshToken,
			7*24*3600,
			"/",
			domain,
			true,
			true,
		)

	}

	c.JSON(200, utils.SuccessResponse("success", gin.H{
		"access_token": accessToken,
	}))

}

func (a *authControllerV1) Logout(c *gin.Context) {
	ctx := context.Background()
	log.Printf("Logout attempt initiated")
	refreshToken, err := c.Cookie("token")
	if err != nil {
		log.Printf("Logout - refresh token not found in cookie: %v", err)
		c.JSON(400, utils.ErrorResponse("error", "Refresh token not found"))
		return
	}

	accessToken := c.GetHeader("Authorization")
	if accessToken == "" || len(accessToken) <= 7 || accessToken[:7] != "Bearer " {
		c.JSON(401, utils.ErrorResponse("error", "Authorization header missing or invalid"))
		return
	}
	accessToken = accessToken[7:]

	refreshTokenClaims, err := utils.DecodeJwtClaims(refreshToken)
	if err != nil {
		c.JSON(401, utils.ErrorResponse("error", "Invalid cookie token"))
		return
	}

	accessTokenClaims, err := utils.DecodeJwtClaims(accessToken)
	if err != nil {
		c.JSON(401, utils.ErrorResponse("error", "Invalid access token"))
		return
	}

	err = a.authService.LogoutService(ctx, refreshToken, accessToken, refreshTokenClaims.ExpiresAt.Unix(), accessTokenClaims.ExpiresAt.Unix())
	if err != nil {
		log.Printf("Logout service failed: %v", err)
		c.JSON(500, utils.ErrorResponse("error", "Failed to logout"))
		return
	}

	log.Printf("Logout successful")
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	c.JSON(200, utils.SuccessResponse("success", "Logged out successfully"))

}

func (a *authControllerV1) Refresh(c *gin.Context) {
	accessToken := c.GetHeader("Authorization")
	refreshToken, err := c.Cookie("token")
	if err != nil {
		c.JSON(400, utils.ErrorResponse("error", "invalid cookie token"))
		return
	}
	decodedRefreshToken, err := utils.DecodeToken(refreshToken)
	if err != nil {
		c.JSON(401, utils.ErrorResponse("error", "Invalid refresh token"))
		return
	}
	_, err = utils.ValidateToken(decodedRefreshToken)
	if err != nil {
		c.JSON(401, utils.ErrorResponse("error", "Invalid token"))
		return
	}
	if accessToken == "" || len(accessToken) <= 7 || accessToken[:7] != "Bearer " {
		c.JSON(401, utils.ErrorResponse("error", "Authorization header missing or invalid"))
		return
	}
	accessToken = accessToken[7:]
	_, err = utils.ValidateToken(accessToken)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			claims, err := utils.DecodeJwtClaims(accessToken)
			if err != nil {
				c.JSON(401, utils.ErrorResponse("error", "Invalid access token claims"))
				return
			}
			newAccessToken, err := utils.GenerateToken(utils.EntityData{
				ID:       claims.Data.ID,
				Username: claims.Data.Username,
				Belong:   claims.Data.Belong,
				Role:     claims.Data.Role,
			}, 15, "access")
			if err != nil {
				c.JSON(500, utils.ErrorResponse("error", "Failed to generate new access token"))
				return
			}
			c.JSON(200, utils.SuccessResponse("success", gin.H{"access_token": newAccessToken}))
			return
		}
		c.JSON(401, gin.H{"error": "Invalid access token"})
		return
	}

}

func NewAuthControllerV1() AuthControllerV1 {
	return &authControllerV1{
		authService: services.NewAuthService(),
	}
}
