package dto

type CreateTicketRequest struct {
	TenantID    string `json:"tenant_id" binding:"required,uuid"`
	CustomerID  string `json:"customer_id" binding:"required,uuid"`
	AssignedTo  string `json:"assigned_to" binding:"omitempty,uuid"`
	Title       string `json:"title" binding:"required,min=1,max=200"`
	Description string `json:"description" binding:"omitempty,max=2000"`
	Priority    string `json:"priority" binding:"omitempty,oneof=LOW MEDIUM HIGH URGENT CRITICAL"`
}

type UpdateTicketRequest struct {
	AssignedTo  string `json:"assigned_to" binding:"omitempty,uuid"`
	Title       string `json:"title" binding:"omitempty,min=1,max=200"`
	Description string `json:"description" binding:"omitempty,max=2000"`
	Status      string `json:"status" binding:"omitempty,oneof=OPEN IN_PROGRESS RESOLVED CLOSED REOPENED"`
	Priority    string `json:"priority" binding:"omitempty,oneof=LOW MEDIUM HIGH URGENT CRITICAL"`
}

type CreateTicketCommentRequest struct {
	Comment string `json:"comment" binding:"required,min=1,max=1000"`
}

type TicketCommentResponse struct {
	ID        string `json:"id"`
	TicketID  string `json:"ticket_id"`
	Comment   string `json:"comment"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
