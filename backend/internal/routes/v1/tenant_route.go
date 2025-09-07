package v1_routes

import (
	v1_controllers "backend/internal/controllers/v1"
	"backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func tenantRoutes(r *gin.RouterGroup) {

	tenant := r.Group("/tenant")
	tenant.Use(middleware.ValidationErrorHandler())
	tenant.Use(middleware.DBErrorHandler())
	tenant.Use(middleware.AuthMiddleware())
	tenant.Use(middleware.RoleMiddleware("Admin"))
	tenantController := v1_controllers.NewTenantControllerV1()
	tenant.POST("", tenantController.CreateTenant)
	tenant.PUT("", tenantController.UpdateTenant)
	tenant.DELETE("", tenantController.DeleteTenant)

}
