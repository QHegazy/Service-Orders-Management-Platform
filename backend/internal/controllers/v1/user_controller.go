package v1_controllers

import (
	"backend/internal/dto"
	"backend/internal/services"
	"backend/utils"
	"context"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserControllerV1 interface {
	CreateUser(c *gin.Context)
	CreateTechnicianUser(c *gin.Context)
	GetTechnicians(c *gin.Context)
	DeleteUser(c *gin.Context)
	UpdateUser(c *gin.Context)
}

type userControllerV1 struct {
	userService *services.UserService
}

func (u *userControllerV1) CreateUser(c *gin.Context) {
	ctx := context.Background()
	userDto := dto.CreateUserDto{}

	if err := c.ShouldBindJSON(&userDto); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	if err := u.userService.CreateUser(ctx, userDto, "Admin"); err != nil {
		c.Error(err).SetType(gin.ErrorTypeAny)
		return
	}

	c.JSON(200, utils.SuccessResponse("success", "User Created"))
}

func (u *userControllerV1) CreateTechnicianUser(c *gin.Context) {
	ctx := context.Background()
	userDto := dto.CreateUserDto{}

	if err := c.ShouldBindJSON(&userDto); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	if err := u.userService.CreateUser(ctx, userDto, "Technician"); err != nil {
		c.Error(err).SetType(gin.ErrorTypeAny)
		return
	}

	c.JSON(200, utils.SuccessResponse("success", "Tech User Created"))
}

func (u *userControllerV1) GetTechnicians(c *gin.Context) {
	ctx := context.Background()

	technicians, err := u.userService.GetTechnicians(ctx)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeAny)
		return
	}

	c.JSON(200, utils.SuccessResponse("success", technicians))
}

func (u *userControllerV1) UpdateUser(c *gin.Context) {
	ctx := context.Background()
	userID := c.GetString("userID")

	var updateDto dto.UpdateUserDto
	if err := c.ShouldBindJSON(&updateDto); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	if err := u.userService.UpdateUserProfile(ctx, userID, updateDto); err != nil {
		c.Error(err).SetType(gin.ErrorTypeAny)
		return
	}

	c.JSON(200, utils.SuccessResponse("success", "Profile updated successfully"))
}

func (u *userControllerV1) DeleteUser(c *gin.Context) {
	ctx := context.Background()
	userID := c.Param("id")

	if userID == "" {
		c.JSON(400, utils.ErrorResponse("error", "User ID is required"))
		return
	}

	if err := u.userService.DeleteUserByID(ctx, userID); err != nil {
		c.Error(err).SetType(gin.ErrorTypeAny)
		return
	}

	c.JSON(200, utils.SuccessResponse("success", "User deleted successfully"))
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
