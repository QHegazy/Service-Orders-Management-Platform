package services

import (
	"backend/internal/dto"
	"backend/internal/repositories"
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type TicketService struct {
	queries *repositories.Queries
}

func NewTicketService() *TicketService {
	return &TicketService{queries: repositories.GetDB()}
}

func makeText(s string) pgtype.Text {
	if s == "" {
		return pgtype.Text{}
	}
	return pgtype.Text{String: s, Valid: true}
}

func (s *TicketService) CreateTicket(ctx context.Context, dto dto.CreateTicketDto) (pgtype.UUID, error) {
	tenantUUID, err := parseUUID(dto.TenantID)
	if err != nil {
		return pgtype.UUID{}, fmt.Errorf("invalid tenant_id: %w", err)
	}

	customerUUID, err := parseUUID(dto.CustomerID)
	if err != nil {
		return pgtype.UUID{}, fmt.Errorf("invalid customer_id: %w", err)
	}

	if dto.AssignedTo == "" {
		technicianIDs, err := s.queries.GetTechnicianIDsFromTenantID(ctx, tenantUUID)
		if err != nil {
			return pgtype.UUID{}, fmt.Errorf("failed to get technicians for tenant: %w", err)
		}
		if len(technicianIDs) > 0 {
			rand.Seed(time.Now().UnixNano())
			randomIndex := rand.Intn(len(technicianIDs))
			dto.AssignedTo = technicianIDs[randomIndex].String()
		}
	}

	assignedToUUID, err := parseUUID(dto.AssignedTo)
	if err != nil {
		return pgtype.UUID{}, fmt.Errorf("invalid assigned_to: %w", err)
	}

	params := repositories.CreateTicketParams{
		TenantID:    tenantUUID,
		CustomerID:  customerUUID,
		Title:       dto.Title,
		Description: makeText(dto.Description),
		Status:      "OPEN",
		AssignedTo:  assignedToUUID,
		Priority:    dto.Priority,
	}

	ticket, err := s.queries.CreateTicket(ctx, params)
	if err != nil {
		return pgtype.UUID{}, fmt.Errorf("failed to create ticket: %w", err)
	}
	return ticket, nil
}
func (s *TicketService) UpdateTicket(ctx context.Context, ticketID pgtype.UUID, dto dto.UpdateTicketDto) (repositories.TicketTicket, error) {
	assignedToUUID, err := parseUUID(dto.AssignedTo)
	if err != nil {
		return repositories.TicketTicket{}, fmt.Errorf("invalid assigned_to: %w", err)
	}

	params := repositories.UpdateTicketParams{
		ID:          ticketID,
		AssignedTo:  assignedToUUID,
		Title:       dto.Title,
		Description: makeText(dto.Description),
		Status:      dto.Status,
		Priority:    dto.Priority,
	}

	ticket, err := s.queries.UpdateTicket(ctx, params)
	if err != nil {
		return repositories.TicketTicket{}, fmt.Errorf("failed to update ticket: %w", err)
	}
	return ticket, nil
}

func (s *TicketService) DeleteTicket(ctx context.Context, ticketID pgtype.UUID) error {
	if err := s.queries.DeleteTicket(ctx, ticketID); err != nil {
		return fmt.Errorf("failed to delete ticket: %w", err)
	}
	return nil
}

func (s *TicketService) GetTicketByID(ctx context.Context, ticketID pgtype.UUID) (repositories.TicketTicket, error) {
	ticket, err := s.queries.GetTicketByID(ctx, ticketID)
	if err != nil {
		return repositories.TicketTicket{}, fmt.Errorf("failed to get ticket: %w", err)
	}
	return ticket, nil
}

func (s *TicketService) ListTicketsByTenantID(ctx context.Context, tenantID pgtype.UUID, page, size int32) ([]repositories.TicketTicket, error) {
	params := repositories.ListTicketsByTenantIDParams{
		TenantID: tenantID,
		Limit:    size,
		Offset:   page,
	}
	tickets, err := s.queries.ListTicketsByTenantID(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list tickets by tenant: %w", err)
	}
	return tickets, nil
}

func (s *TicketService) ListTicketsByCustomerID(ctx context.Context, customerID string, page, size int32) ([]repositories.ListTicketsByCustomerIDRow, error) {
	customerUUID, err := parseUUID(customerID)
	if err != nil {
		return nil, fmt.Errorf("invalid customer_id: %w", err)
	}
	fmt.Println(page, size)
	params := repositories.ListTicketsByCustomerIDParams{
		CustomerID: customerUUID,
		Limit:      size,
		Offset:     page,
	}
	tickets, err := s.queries.ListTicketsByCustomerID(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list tickets by customer: %w", err)
	}
	return tickets, nil
}

func (s *TicketService) ListTicketsByAssignedTo(ctx context.Context, assignedTo pgtype.UUID, page, size int32) ([]repositories.TicketTicket, error) {
	params := repositories.ListTicketsByAssignedToParams{
		AssignedTo: assignedTo,
		Limit:      size,
		Offset:     page,
	}
	tickets, err := s.queries.ListTicketsByAssignedTo(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list tickets by assigned user: %w", err)
	}
	return tickets, nil
}

func (s *TicketService) ListTicketsByUserID(ctx context.Context, userID pgtype.UUID, page, size int32) ([]repositories.ListTicketsByUserIdRow, error) {
	params := repositories.ListTicketsByUserIdParams{
		UserID: userID,
		Limit:  size,
		Offset: page,
	}
	tickets, err := s.queries.ListTicketsByUserId(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list tickets by user: %w", err)
	}
	return tickets, nil
}

func (s *TicketService) ListCommentsByTicketID(ctx context.Context, ticketID pgtype.UUID, page, size int32) ([]repositories.ListCommentsByTicketIDRow, error) {
	params := repositories.ListCommentsByTicketIDParams{
		TicketID: ticketID,
		Limit:    size,
		Offset:   page,
	}
	comments, err := s.queries.ListCommentsByTicketID(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list comments by ticket ID: %w", err)
	}
	return comments, nil
}

func (s *TicketService) CreateComment(ctx context.Context, ticketIDStr, content, authorIDStr, authorType string) (pgtype.UUID, error) {
	ticketUUID, err := parseUUID(ticketIDStr)
	if err != nil {
		return pgtype.UUID{}, fmt.Errorf("invalid ticket_id: %w", err)
	}

	authorUUID, err := parseUUID(authorIDStr)
	if err != nil {
		return pgtype.UUID{}, fmt.Errorf("invalid author_id: %w", err)
	}

	params := repositories.CreateCommentParams{
		TicketID:   ticketUUID,
		Comment:    content,
		AuthorID:   authorUUID,
		AuthorType: authorType,
	}

	commentID, err := s.queries.CreateComment(ctx, params)
	if err != nil {
		return pgtype.UUID{}, fmt.Errorf("failed to create comment: %w", err)
	}
	return commentID, nil
}
