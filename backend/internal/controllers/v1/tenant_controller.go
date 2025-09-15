package v1_controllers

import (
	"backend/internal/dto"
	"backend/internal/services"
	"backend/utils"
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

type TenantControllerV1 interface {
	CreateTenant(c *gin.Context)
	UpdateTenant(c *gin.Context)
	DeleteTenant(c *gin.Context)
	GetUserTenants(c *gin.Context)
	AllTenants(c *gin.Context)
	AddUserToTenant(c *gin.Context)
}

type tenantControllerV1 struct {
	tenantService *services.TenantService
}

func (t *tenantControllerV1) CreateTenant(c *gin.Context) {
	ctx := context.Background()
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(400, utils.ErrorResponse("error", "User ID not found in token"))
		return
	}
	fmt.Println("asasdasd", userID)

	tenantDto := dto.CreateTenantDto{}
	if err := c.ShouldBindJSON(&tenantDto); err != nil {
		log.Printf("TenantController - JSON binding error for request: %v", err)
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	tenantID, err := t.tenantService.CreateTenant(ctx, tenantDto, userID)
	if err != nil {
		log.Printf("TenantController - Failed to create tenant: %v", err)
		c.Error(err).SetType(gin.ErrorTypeAny)
		return
	}

	c.JSON(200, utils.SuccessResponse("success", gin.H{"tenant_id": tenantID}))
}

func (t *tenantControllerV1) UpdateTenant(c *gin.Context) {
	ctx := context.Background()
	tenantDto := dto.UpdateTenantDto{}

	if err := c.ShouldBindJSON(&tenantDto); err != nil {
		log.Printf("TenantController - JSON binding error for request: %v", err)
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	if err := t.tenantService.UpdateTenant(ctx, tenantDto); err != nil {
		log.Printf("TenantController - Failed to update tenant: %v", err)
		c.Error(err).SetType(gin.ErrorTypeAny)
		return
	}

	c.JSON(200, utils.SuccessResponse("success", "Tenant Updated"))
}

func (t *tenantControllerV1) DeleteTenant(c *gin.Context) {

	c.JSON(200, utils.SuccessResponse("success", "Tenant Deleted"))
}

func (t *tenantControllerV1) GetUserTenants(c *gin.Context) {
	ctx := context.Background()
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(400, utils.ErrorResponse("error", "User ID not found in token"))
		return
	}

	parsedUUID, err := uuid.Parse(userID)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	tenants, err := t.tenantService.GetUserTenants(ctx, parsedUUID.String())
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeAny)
		return
	}

	c.JSON(200, utils.SuccessResponse("success", tenants))
}

func (t *tenantControllerV1) AllTenants(c *gin.Context) {
	ctx := context.Background()
	page, _ := c.Get("page")
	size, _ := c.Get("size")
	tenants, err := t.tenantService.GetAllTenants(ctx, page.(int32), size.(int32))
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeAny)
		return
	}

	c.JSON(200, utils.SuccessResponse("success", tenants))
}

func (t *tenantControllerV1) AddUserToTenant(c *gin.Context) {
	ctx := context.Background()
	addDto := dto.AddUserToTenantDto{}

	if err := c.ShouldBindJSON(&addDto); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	userUUID, err := uuid.Parse(addDto.UserID)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}
	tenantUUID, err := uuid.Parse(addDto.TenantID)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	err = t.tenantService.AddUserToTenantService(ctx, userUUID.String(), tenantUUID.String())
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeAny)
		return
	}

	c.JSON(200, utils.SuccessResponse("success", "User added to tenant successfully"))
}

func NewTenantControllerV1() TenantControllerV1 {
	return &tenantControllerV1{
		tenantService: services.NewTenantService(),
	}
}
