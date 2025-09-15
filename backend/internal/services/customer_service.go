package services

import (
	"backend/internal/dto"
	"backend/internal/repositories"
	"backend/utils"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type CustomerService struct {
	queries *repositories.Queries
}

func NewCustomerService() *CustomerService {
	queries := repositories.GetDB()
	return &CustomerService{queries: queries}
}

func (s *CustomerService) CreateCustomer(ctx context.Context, customerDto dto.CreateCustomerDto) error {
	fmt.Println(customerDto)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(customerDto.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	newCustomer := repositories.CreateCustomerParams{
		FirstName: customerDto.FirstName,
		LastName:  customerDto.LastName,
		Username:  customerDto.Username,
		Password:  string(hashedPassword),
		Email:     pgtype.Text{String: customerDto.Email, Valid: true},
	}
	_, err = s.queries.CreateCustomer(ctx, newCustomer)
	return err
}

func (s *CustomerService) LoginCustomer(ctx context.Context, loginCustomerDto dto.LoginDto) (string, string, error) {

	customer, err := s.queries.GetCustomerByEmail(ctx, pgtype.Text{String: loginCustomerDto.Email, Valid: true})
	if err != nil {
		log.Printf("CustomerService - Customer not found for email %s: %v", loginCustomerDto.Email, err)
		return "", "", fmt.Errorf("customer not found: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(loginCustomerDto.Password))
	if err != nil {
		log.Printf("CustomerService - Password verification failed for %s: %v", loginCustomerDto.Email, err)
		return "", "", fmt.Errorf("invalid password: %w", err)
	}
	claims := utils.EntityData{
		ID:       customer.ID.String(),
		Username: customer.Username,
		Role:     "Customer",
	}

	accessToken, err := utils.GenerateToken(
		claims,
		time.Minute*60,
		"access",
	)
	if err != nil {
		log.Printf("CustomerService - Failed to generate access token for %s: %v", loginCustomerDto.Email, err)
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := utils.GenerateToken(
		utils.EntityData{ID: customer.ID.String()},
		time.Hour*24*7,
		"refresh",
	)
	encodedRefreshToken, _ := utils.EncodeToken(refreshToken)

	if err != nil {
		log.Printf("CustomerService - Failed to generate refresh token for %s: %v", loginCustomerDto.Email, err)
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return accessToken, encodedRefreshToken, nil
}
