package dto

type CreateTenantRequest struct {
	TenantName string `json:"tenant_name" binding:"required,min=1,max=100"`
	Domain     string `json:"domain" binding:"required,min=3,max=100"`
	LogoUrl    string `json:"logo_url" binding:"omitempty,url,max=500"`
}

type UpdateTenantRequest struct {
	TenantName string `json:"tenant_name" binding:"omitempty,min=1,max=100"`
	Domain     string `json:"domain" binding:"omitempty,min=3,max=100"`
	LogoUrl    string `json:"logo_url" binding:"omitempty,url,max=500"`
	IsActive   *bool  `json:"is_active" binding:"omitempty"`
}

