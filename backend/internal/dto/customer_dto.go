package dto

type CreateCustomerDto struct {
	FirstName string `json:"first_name" binding:"required,min=1,max=100"`
	LastName  string `json:"last_name" binding:"required,min=1,max=100"`
	Username  string `json:"username" binding:"required,min=3,max=50"`
	Email     string `json:"email" binding:"omitempty,email,max=100"`
	Password  string `json:"password" binding:"required,min=6,max=100"`
}

type UpdateCustomerRequest struct {
	FirstName string `json:"first_name" binding:"omitempty,min=1,max=100"`
	LastName  string `json:"last_name" binding:"omitempty,min=1,max=100"`
	Username  string `json:"username" binding:"omitempty,min=3,max=50"`
	Email     string `json:"email" binding:"omitempty,email,max=100"`
	Password  string `json:"password" binding:"omitempty,min=6,max=100"`
}
