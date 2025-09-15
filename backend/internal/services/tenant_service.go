package services

import (
	"backend/internal/dto"
	"backend/internal/repositories"
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type TenantService struct {
	queries *repositories.Queries
}

func NewTenantService() *TenantService {
	queries := repositories.GetDB()
	return &TenantService{queries: queries}
}

func (s *TenantService) AddUserToTenantService(ctx context.Context, userId, tenantId string) error {
	parsedUserID, err := parseUUID(userId)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}
	parsedTenantID, err := parseUUID(tenantId)
	if err != nil {
		return fmt.Errorf("invalid tenant ID: %w", err)
	}

	newTenantUser := repositories.AddUserToTenantParams{
		TenantID: parsedTenantID,
		UserID:   parsedUserID,
	}
	err = s.queries.AddUserToTenant(ctx, newTenantUser)
	if err != nil {
		return fmt.Errorf("failed to add user to tenant: %w", err)
	}
	return nil

}

func (s *TenantService) CreateTenant(ctx context.Context, tenantDto dto.CreateTenantDto, userId string) (pgtype.UUID, error) {
	newTenant := repositories.CreateTenantParams{
		TenantName: tenantDto.TenantName,
		Domain:     tenantDto.Domain,
		Email:      tenantDto.Email,
	}

	TenantID, err := s.queries.CreateTenant(ctx, newTenant)
	if err != nil {
		log.Printf("TenantService - Failed to create tenant: %v", err)
		return pgtype.UUID{}, fmt.Errorf("failed to create tenant: %w", err)
	}

	err = s.AddUserToTenantService(ctx, userId, uuid.UUID(TenantID.Bytes).String())

	if err != nil {
		log.Printf("TenantService - Failed to add user to tenant: %v", err)
		return pgtype.UUID{}, fmt.Errorf("failed to add user to tenant: %w", err)
	}

	return TenantID, err
}

func (s *TenantService) UpdateTenant(ctx context.Context, tenantDto dto.UpdateTenantDto) error {
	updatedTenant := repositories.UpdateTenantParams{
		TenantName: tenantDto.TenantName,
		Domain:     tenantDto.Domain,
		Email:      tenantDto.Email,
	}
	_, err := s.queries.UpdateTenant(ctx, updatedTenant)

	if err != nil {
		log.Printf("TenantService - Failed to update tenant: %v", err)
		return fmt.Errorf("failed to update tenant: %w", err)
	}
	return nil
}

func (s *TenantService) DeleteTenant(ctx context.Context, tenantID pgtype.UUID) error {
	err := s.queries.DeleteTenant(ctx, tenantID)
	if err != nil {
		log.Printf("TenantService - Failed to delete tenant: %v", err)
		return fmt.Errorf("failed to delete tenant: %w", err)
	}
	return nil
}

func (s *TenantService) GetUserTenants(ctx context.Context, userID string) ([]repositories.TenantTenant, error) {
	parsedUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	pgUUID := pgtype.UUID{Bytes: parsedUUID, Valid: true}

	tenants, err := s.queries.GetTenantsByUserId(ctx, pgUUID)
	if err != nil {
		log.Printf("TenantService - Failed to get tenants for user %v: %v", userID, err)
		return nil, fmt.Errorf("failed to get tenants for user: %w", err)
	}
	return tenants, nil
}

func (s *TenantService) GetAllTenants(ctx context.Context, page int32, size int32) ([]repositories.TenantTenant, error) {
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * size
	params := repositories.ListTenantsParams{
		Limit:  size,
		Offset: offset,
	}
	tenants, err := s.queries.ListTenants(ctx, params)
	if err != nil {
		log.Printf("TenantService - Failed to get all tenants: %v", err)
		return nil, fmt.Errorf("failed to get all tenants: %w", err)
	}
	return tenants, nil
}
