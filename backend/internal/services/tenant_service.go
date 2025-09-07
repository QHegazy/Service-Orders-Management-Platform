package services

import (
	"backend/internal/dto"
	"backend/internal/repositories"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type TenantService struct {
	queries *repositories.Queries
}

func NewTenantService() *TenantService {
	queries := repositories.GetDB()
	return &TenantService{queries: queries}
}

func (s *TenantService) AddUserToTenantService(ctx context.Context, userId, tenantId pgtype.UUID) error {

	newTenantUser := repositories.AddUserToTenantParams{
		TenantID: tenantId,
		UserID:   userId,
	}
	err := s.queries.AddUserToTenant(ctx, newTenantUser)
	if err != nil {
		return fmt.Errorf("failed to add user to tenant: %w", err)
	}
	return nil

}

func (s *TenantService) CreateTenant(ctx context.Context, tenantDto dto.CreateTenantDto) (pgtype.UUID, error) {
	newTenant := repositories.CreateTenantParams{
		TenantName: tenantDto.TenantName,
		Domain:     tenantDto.Domain,
		Email:      tenantDto.Email,
	}
	TenantID, err := s.queries.CreateTenant(ctx, newTenant)
	if err != nil {
		return pgtype.UUID{}, fmt.Errorf("failed to create tenant: %w", err)
	}
	fmt.Println(TenantID)

	return TenantID, err
}

func (s *TenantService) UpdateTenant(ctx context.Context, tenantDto dto.UpdateTenantDto) error {
	updatedTenant := repositories.UpdateTenantParams{
		TenantName: tenantDto.TenantName,
		Domain:     tenantDto.Domain,
	}
	_, err := s.queries.UpdateTenant(ctx, updatedTenant)

	if err != nil {
		return fmt.Errorf("failed to update tenant: %w", err)
	}
	return nil
}

func (s *TenantService) DeleteTenant(ctx context.Context, tenantID pgtype.UUID) error {
	err := s.queries.DeleteTenant(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("failed to delete tenant: %w", err)
	}
	return nil
}
