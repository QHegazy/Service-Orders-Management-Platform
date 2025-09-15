package v1_controllers

import (
	"backend/internal/dto"
	"backend/internal/services"
	"backend/utils"
	"context"

	"github.com/gin-gonic/gin"
)

type CustomerControllerV1 interface {
	// LoginCustomer(c *gin.Context)
	CreateCustomer(c *gin.Context)
	DeleteCustomer(c *gin.Context)
	UpdateCustomer(c *gin.Context)
}

type customerControllerV1 struct {
	customerService *services.CustomerService
}

// func (u *customerControllerV1) LoginCustomer(c *gin.Context) {
// 	ctx := context.Background()
// 	var loginDto dto.LoginDto
// 	if err := c.ShouldBindJSON(&loginDto); err != nil {
// 		c.Error(err).SetType(gin.ErrorTypeBind)
// 		return
// 	}

// 	accessToken, refreshToken, err := u.customerService.LoginCustomer(ctx, loginDto)
// 	if err != nil {
// 		c.JSON(401, utils.ErrorResponse("error", "Invalid credentials"))
// 		return
// 	}

// 	c.SetCookie("token", refreshToken, 3600, "/", "localhost", false, true)

// 	if gin.Mode() == gin.ReleaseMode {
// 		domain := os.Getenv("DOMAIN")
// 		if domain == "" {
// 			domain = getDomain(c)
// 		}

// 		c.SetCookie(
// 			"token",
// 			refreshToken,
// 			7*24*3600,
// 			"/",
// 			domain,
// 			true,
// 			true,
// 		)

// 	}

// 	c.JSON(200, utils.SuccessResponse("success", gin.H{
// 		"access_token": accessToken,
// 	}))
// }

func (u *customerControllerV1) CreateCustomer(c *gin.Context) {
	ctx := context.Background()
	customerDto := dto.CreateCustomerDto{}

	if err := c.ShouldBindJSON(&customerDto); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	if err := u.customerService.CreateCustomer(ctx, customerDto); err != nil {
		c.Error(err).SetType(gin.ErrorTypeAny)
		return
	}

	c.JSON(200, utils.SuccessResponse("success", "Customer Created"))
}

func (u *customerControllerV1) UpdateCustomer(c *gin.Context) {
	c.String(200, "Customer Profile Updated")
}

func (u *customerControllerV1) DeleteCustomer(c *gin.Context) {
	c.String(200, "Customer Profile Deleted")
}

// 4. Factory function (optional)
func NewCustomerControllerV1() CustomerControllerV1 {
	return &customerControllerV1{
		customerService: services.NewCustomerService(),
	}
}
