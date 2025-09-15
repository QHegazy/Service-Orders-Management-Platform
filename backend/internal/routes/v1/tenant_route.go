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

	tenantController := v1_controllers.NewTenantControllerV1()

	// Admin only routes
	adminTenant := tenant.Group("")
	adminTenant.Use(middleware.RoleMiddleware("Admin"))
	adminTenant.POST("", tenantController.CreateTenant)
	adminTenant.PUT("", tenantController.UpdateTenant)
	adminTenant.DELETE("", tenantController.DeleteTenant)
	adminTenant.POST("add-to", tenantController.AddUserToTenant)
	tenant.GET("all", middleware.PaginationMiddleware(), middleware.RoleMiddleware("Customer"), tenantController.AllTenants)

	// User routes (authenticated users can see their own tenants)
	userTenant := tenant.Group("")
	userTenant.GET("user", tenantController.GetUserTenants)

}
