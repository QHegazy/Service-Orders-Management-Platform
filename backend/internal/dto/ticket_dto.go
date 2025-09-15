package dto

type CreateTicketDto struct {
	TenantID    string `json:"tenant_id" binding:"required,uuid"`
	CustomerID  string `json:"customer_id" binding:"required,uuid"`
	AssignedTo  string `json:"assigned_to" binding:"omitempty,uuid"`
	Title       string `json:"title" binding:"required,min=1,max=200"`
	Description string `json:"description" binding:"omitempty,max=2000"`
	Priority    string `json:"priority" binding:"omitempty,oneof=LOW MEDIUM HIGH URGENT CRITICAL"`
}

type UpdateTicketDto struct {
	AssignedTo  string `json:"assigned_to" binding:"omitempty,uuid"`
	Title       string `json:"title" binding:"omitempty,min=1,max=200"`
	Description string `json:"description" binding:"omitempty,max=2000"`
	Status      string `json:"status" binding:"omitempty,oneof=OPEN IN_PROGRESS RESOLVED CLOSED REOPENED"`
	Priority    string `json:"priority" binding:"omitempty,oneof=LOW MEDIUM HIGH URGENT CRITICAL"`
}

type TicketCommentDto struct {
	TicketID   string `json:"ticket_id" binding:"required,uuid"`
	AuthorID   string `json:"user_id" binding:"required,uuid"`
	AuthorType string `json:"author_type" binding:"required,oneof=User Customer"`
	Comment    string `json:"comment" binding:"required,min=1,max=1000"`
}
