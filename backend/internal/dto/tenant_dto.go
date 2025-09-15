package dto

type CreateTenantDto struct {
	TenantName string `json:"tenant_name" binding:"required,min=1,max=100"`
	Domain     string `json:"domain" binding:"required,min=3,max=100"`
	Email      string `json:"email" binding:"required,email"`
}

type UpdateTenantDto struct {
	TenantName string `json:"tenant_name" binding:"omitempty,min=1,max=100"`
	Domain     string `json:"domain" binding:"omitempty,min=3,max=100"`
	Email      string `json:"email" binding:"omitempty,email"`
	Status     string `json:"status" binding:"omitempty,oneof=ACTIVE INACTIVE SUSPENDED"`
}

type AddUserToTenantDto struct {
	UserID   string `json:"user_id" binding:"required,uuid"`
	TenantID string `json:"tenant_id" binding:"required,uuid"`
}
