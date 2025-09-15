package services

import (
	"backend/internal/dto"
	"backend/internal/repositories"
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type InvoiceService struct {
	queries *repositories.Queries
}

func NewInvoiceService() *InvoiceService {
	queries := repositories.GetDB()
	return &InvoiceService{queries: queries}
}

func (s *InvoiceService) CreateInvoice(ctx context.Context, invoiceDto dto.CreateInvoiceRequest) error {
	ticketUUID, err := parseUUID(invoiceDto.TicketID)
	if err != nil {
		return fmt.Errorf("invalid ticket ID: %w", err)
	}

	// Convert float64 amount to pgtype.Numeric
	amountNumeric := pgtype.Numeric{}
	if err := amountNumeric.Scan(strconv.FormatFloat(invoiceDto.Amount, 'f', -1, 64)); err != nil {
		return fmt.Errorf("failed to convert amount to pgtype.Numeric: %w", err)
	}

	dueDate, err := time.Parse("2006-01-02", invoiceDto.DueDate)
	if err != nil {
		return fmt.Errorf("invalid due date format: %w", err)
	}

	newInvoice := repositories.CreateInvoiceParams{
		TicketID: ticketUUID,
		Amount:   amountNumeric,
		Currency: invoiceDto.Currency,
		DueDate:  pgtype.Date{Time: dueDate, Valid: true},
	}

	_, err = s.queries.CreateInvoice(ctx, newInvoice)
	if err != nil {
		log.Printf("InvoiceService - Failed to create invoice: %v", err)
		return fmt.Errorf("failed to create invoice: %w", err)
	}
	return nil
}
