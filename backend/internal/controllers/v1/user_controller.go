package v1_controllers

import (
	"backend/internal/dto"
	"backend/internal/services"
	"backend/utils"
	"context"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserControllerV1 interface {
	LoginUser(c *gin.Context)
	CreateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
	UpdateUser(c *gin.Context)
}

type userControllerV1 struct {
	userService *services.UserService
}

func (u *userControllerV1) LoginUser(c *gin.Context) {
	ctx := context.Background()
	var loginDto dto.LoginUserDto
	if err := c.ShouldBindJSON(&loginDto); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	accessToken, refreshToken, err := u.userService.LoginUser(ctx, loginDto)
	if err != nil {
		c.JSON(401, utils.ErrorResponse("error", "Invalid credentials"))
		return
	}

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

func (u *userControllerV1) CreateUser(c *gin.Context) {
	ctx := context.Background()
	userDto := dto.CreateUserDto{}

	if err := c.ShouldBindJSON(&userDto); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	if err := u.userService.CreateUser(ctx, userDto); err != nil {
		c.Error(err).SetType(gin.ErrorTypeAny)
		return
	}

	c.JSON(200, utils.SuccessResponse("success", "User Created"))
}

func (u *userControllerV1) UpdateUser(c *gin.Context) {
	c.String(200, "User Profile Updated")
}

func (u *userControllerV1) DeleteUser(c *gin.Context) {
	c.String(200, "User Profile Deleted")
}

// 4. Factory function (optional)
func NewUserControllerV1() UserControllerV1 {
	return &userControllerV1{
		userService: services.NewUserService(),
	}
}

func getDomain(c *gin.Context) string {
	host := c.Request.Host
	if strings.Contains(host, ":") {
		parts := strings.Split(host, ":")
		return parts[0]
	}
	return host
}
