package services

import (
	"backend/internal/dto"
	"backend/internal/repositories"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type TicketService struct {
	queries *repositories.Queries
}

func NewTicketService() *TicketService {
	return &TicketService{queries: repositories.GetDB()}
}

// Helper functions
func parseUUID(s string) (pgtype.UUID, error) {
	if s == "" {
		return pgtype.UUID{}, nil
	}
	id, err := uuid.Parse(s)
	if err != nil {
		return pgtype.UUID{}, err
	}
	return pgtype.UUID{Bytes: id, Valid: true}, nil
}

func makeText(s string) pgtype.Text {
	if s == "" {
		return pgtype.Text{}
	}
	return pgtype.Text{String: s, Valid: true}
}

func makeStatus(s string) repositories.NullTicketTicketStatus {
	if s == "" {
		return repositories.NullTicketTicketStatus{}
	}
	return repositories.NullTicketTicketStatus{
		TicketTicketStatus: repositories.TicketTicketStatus(s),
		Valid:              true,
	}
}

func makePriority(s string) repositories.NullTicketTicketPriority {
	if s == "" {
		return repositories.NullTicketTicketPriority{}
	}
	return repositories.NullTicketTicketPriority{
		TicketTicketPriority: repositories.TicketTicketPriority(s),
		Valid:                true,
	}
}

func (s *TicketService) CreateTicket(ctx context.Context, dto dto.CreateTicketDto) (repositories.TicketTicket, error) {
	tenantUUID, err := parseUUID(dto.TenantID)
	if err != nil {
		return repositories.TicketTicket{}, fmt.Errorf("invalid tenant_id: %w", err)
	}
	customerUUID, err := parseUUID(dto.CustomerID)
	if err != nil {
		return repositories.TicketTicket{}, fmt.Errorf("invalid customer_id: %w", err)
	}
	assignedToUUID, err := parseUUID(dto.AssignedTo)
	if err != nil {
		return repositories.TicketTicket{}, fmt.Errorf("invalid assigned_to: %w", err)
	}

	params := repositories.CreateTicketParams{
		TenantID:    tenantUUID,
		CustomerID:  customerUUID,
		AssignedTo:  assignedToUUID,
		Title:       dto.Title,
		Description: makeText(dto.Description),
		Status:      makeStatus("OPEN"),
		Priority:    makePriority(dto.Priority),
	}

	ticket, err := s.queries.CreateTicket(ctx, params)
	if err != nil {
		return repositories.TicketTicket{}, fmt.Errorf("failed to create ticket: %w", err)
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
		Status:      makeStatus(dto.Status),
		Priority:    makePriority(dto.Priority),
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

func (s *TicketService) ListTicketsByTenantID(ctx context.Context, tenantID pgtype.UUID, limit, offset int32) ([]repositories.TicketTicket, error) {
	params := repositories.ListTicketsByTenantIDParams{
		TenantID: tenantID,
		Limit:    limit,
		Offset:   offset,
	}
	tickets, err := s.queries.ListTicketsByTenantID(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list tickets by tenant: %w", err)
	}
	return tickets, nil
}

func (s *TicketService) ListTicketsByCustomerID(ctx context.Context, customerID pgtype.UUID, limit, offset int32) ([]repositories.TicketTicket, error) {
	params := repositories.ListTicketsByCustomerIDParams{
		CustomerID: customerID,
		Limit:      limit,
		Offset:     offset,
	}
	tickets, err := s.queries.ListTicketsByCustomerID(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list tickets by customer: %w", err)
	}
	return tickets, nil
}

func (s *TicketService) ListTicketsByAssignedTo(ctx context.Context, assignedTo pgtype.UUID, limit, offset int32) ([]repositories.TicketTicket, error) {
	params := repositories.ListTicketsByAssignedToParams{
		AssignedTo: assignedTo,
		Limit:      limit,
		Offset:     offset,
	}
	tickets, err := s.queries.ListTicketsByAssignedTo(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list tickets by assigned user: %w", err)
	}
	return tickets, nil
}
