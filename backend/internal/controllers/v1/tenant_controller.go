package v1_controllers

import (
	"backend/internal/dto"
	"backend/internal/services"
	"backend/utils"
	"context"

	"github.com/gin-gonic/gin"
)

type TenantControllerV1 interface {
	CreateTenant(c *gin.Context)
	UpdateTenant(c *gin.Context)
	DeleteTenant(c *gin.Context)
}

type tenantControllerV1 struct {
	tenantService *services.TenantService
}

func (t *tenantControllerV1) CreateTenant(c *gin.Context) {
	ctx := context.Background()
	tenantDto := dto.CreateTenantDto{}

	if err := c.ShouldBindJSON(&tenantDto); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	tenantID, err := t.tenantService.CreateTenant(ctx, tenantDto)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeAny)
		return
	}

	c.JSON(200, utils.SuccessResponse("success", gin.H{"tenant_id": tenantID}))
}

func (t *tenantControllerV1) UpdateTenant(c *gin.Context) {
	ctx := context.Background()
	tenantDto := dto.UpdateTenantDto{}

	if err := c.ShouldBindJSON(&tenantDto); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	if err := t.tenantService.UpdateTenant(ctx, tenantDto); err != nil {
		c.Error(err).SetType(gin.ErrorTypeAny)
		return
	}

	c.JSON(200, utils.SuccessResponse("success", "Tenant Updated"))
}

func (t *tenantControllerV1) DeleteTenant(c *gin.Context) {

	c.JSON(200, utils.SuccessResponse("success", "Tenant Deleted"))
}

func NewTenantControllerV1() TenantControllerV1 {
	return &tenantControllerV1{
		tenantService: services.NewTenantService(),
	}
}
