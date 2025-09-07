package dto

type CreateUserDto struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email,min=6,max=100"`
	Password string `json:"password" binding:"required,min=6,max=100"`
	Role     string `json:"role" binding:"required,oneof=Admin Technician"`
}

type UpdateUserDto struct {
	Username *string `json:"username" binding:"omitempty,min=3,max=50"`
	Email    *string `json:"email" binding:"omitempty,email,min=6,max=100"`
	Password *string `json:"password" binding:"omitempty,min=6,max=100"`
}

type LoginDto struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=100"`
}
