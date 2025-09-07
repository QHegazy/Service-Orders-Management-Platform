package dto

type CreateTenantDto struct {
	TenantName string `form:"tenant_name" binding:"required,min=1,max=100"`
	Domain     string `form:"domain" binding:"required,min=3,max=100"`
	Email      string `json:"email" binding:"required,email"`
}

type UpdateTenantDto struct {
	TenantName string `form:"tenant_name" binding:"omitempty,min=1,max=100"`
	Domain     string `form:"domain" binding:"omitempty,min=3,max=100"`
	Email      string `json:"email" binding:"omitempty,email"`
}
