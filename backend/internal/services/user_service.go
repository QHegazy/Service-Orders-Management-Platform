package services

import (
	"backend/internal/dto"
	"backend/internal/repositories"
	"backend/utils"
	"context"
	"fmt"
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

func (s *UserService) CreateUser(ctx context.Context, userDto dto.CreateUserDto) error {
	fmt.Println(userDto)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDto.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	newUser := repositories.CreateUserParams{
		Username: userDto.Username,
		Password: string(hashedPassword),
		Email:    userDto.Email,
		Role:     repositories.UserRole(userDto.Role),
	}
	_, err = s.queries.CreateUser(ctx, newUser)
	return err
}

func (s *UserService) LoginUser(ctx context.Context, loginUserDto dto.LoginUserDto) (string, string, error) {
	user, err := s.queries.GetUserByUsername(ctx, loginUserDto.Username)
	if err != nil {
		return "", "", fmt.Errorf("user not found: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginUserDto.Password))
	if err != nil {
		return "", "", fmt.Errorf("invalid password: %w", err)
	}
	claims := utils.EntityData{
		ID:       user.ID.String(),
		Username: user.Username,
		Belong:   user.TenantID.String(),
		Role:     string(user.Role),
	}
	accessToken, err := utils.GenerateToken(
		claims,
		time.Minute*15,
		"access",
	)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := utils.GenerateToken(
		utils.EntityData{ID: user.ID.String()},
		time.Hour*24*7,
		"refresh",
	)
	encodedRefreshToken, _ := utils.EncodeToken(refreshToken)

	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return accessToken, encodedRefreshToken, nil
}
