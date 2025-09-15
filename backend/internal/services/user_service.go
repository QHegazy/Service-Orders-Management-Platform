package services

import (
	"backend/internal/dto"
	"backend/internal/repositories"
	"backend/utils"
	"context"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	queries *repositories.Queries
}

func NewUserService() *UserService {
	queries := repositories.GetDB()
	return &UserService{queries: queries}
}

func (s *UserService) CreateUser(ctx context.Context, userDto dto.CreateUserDto, role string) error {
	if !utils.IsBusinessEmail(userDto.Email) {
		log.Printf("UserService - Creating user with business email: %s", userDto.Email)
		return fmt.Errorf("business emails are not allowed for user creation: %s", userDto.Email)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDto.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	newUser := repositories.CreateUserParams{
		Username: userDto.Username,
		Password: string(hashedPassword),
		Email:    userDto.Email,
		Role:     repositories.UserRole(role),
	}

	_, err = s.queries.CreateUser(ctx, newUser)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) LoginUser(ctx context.Context, loginUserDto dto.LoginDto) (string, string, error) {

	user, err := s.queries.GetTenantIdsByUserMail(ctx, loginUserDto.Email)
	if err != nil {
		log.Printf("UserService - User not found for email %s: %v", loginUserDto.Email, err)
		return "", "", fmt.Errorf("user not found: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginUserDto.Password))
	if err != nil {
		log.Printf("UserService - Password verification failed for %s: %v", loginUserDto.Email, err)
		return "", "", fmt.Errorf("invalid password: %w", err)
	}
	tenantIds, ok := user.TenantIds.([]string)
	claims := utils.EntityData{
		ID:       user.ID.String(),
		Username: user.Username,
		Role:     string(user.Role),
	}
	if ok {
		claims.Belong = tenantIds
	}

	accessToken, err := utils.GenerateToken(
		claims,
		time.Minute*60,
		"access",
	)
	if err != nil {
		log.Printf("UserService - Failed to generate access token for %s: %v", loginUserDto.Email, err)
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := utils.GenerateToken(
		utils.EntityData{ID: user.ID.String()},
		time.Hour*24*7,
		"refresh",
	)
	encodedRefreshToken, _ := utils.EncodeToken(refreshToken)

	if err != nil {
		log.Printf("UserService - Failed to generate refresh token for %s: %v", loginUserDto.Email, err)
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return accessToken, encodedRefreshToken, nil
}

func (s *UserService) GetTechnicians(ctx context.Context) ([]repositories.GetUsersByRoleRow, error) {
	technicians, err := s.queries.GetUsersByRole(ctx, repositories.UserRoleTechnician)
	if err != nil {
		return nil, fmt.Errorf("failed to get technicians: %w", err)
	}
	return technicians, nil
}

func (s *UserService) UpdateUserProfile(ctx context.Context, userID string, userDto dto.UpdateUserDto) error {
	// Parse UUID
	id, err := parseUUID(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	var hashedPassword string
	if userDto.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(userDto.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		hashedPassword = string(hash)
	}

	params := repositories.UpdateUserParams{
		ID:       id,
		Username: userDto.Username,
		Email:    userDto.Email,
		Password: hashedPassword,
	}

	_, err = s.queries.UpdateUser(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (s *UserService) DeleteUserByID(ctx context.Context, userID string) error {
	id, err := parseUUID(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	err = s.queries.DeleteUser(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
