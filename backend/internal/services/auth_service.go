package services

import (
	"backend/internal/dto"
	"backend/internal/repositories"
	"backend/utils"
	"context"
	"log"
)

type AuthService struct {
	queries *repositories.Queries
}

func NewAuthService() *AuthService {
	queries := repositories.GetDB()
	return &AuthService{queries: queries}
}

func (a *AuthService) LoginService(ctx context.Context, loginUserDto dto.LoginDto) (string, string, error) {
	if utils.IsBusinessEmail(loginUserDto.Email) {
		refreshToken, access_token, err := NewUserService().LoginUser(ctx, loginUserDto)
		if err != nil {
			log.Printf("AuthService - UserService login failed for %s: %v", loginUserDto.Email, err)
		} else {
			log.Printf("AuthService - UserService login successful for %s", loginUserDto.Email)
		}
		return refreshToken, access_token, err
	} else {
		log.Printf("AuthService - Routing to CustomerService for customer email: %s", loginUserDto.Email)
		refreshToken, access_token, err := NewCustomerService().LoginCustomer(ctx, loginUserDto)
		if err != nil {
			log.Printf("AuthService - CustomerService login failed for %s: %v", loginUserDto.Email, err)
		} else {
			log.Printf("AuthService - CustomerService login successful for %s", loginUserDto.Email)
		}
		return refreshToken, access_token, err
	}
}

func (a *AuthService) LogoutService(ctx context.Context, refreshToken, accessToken string, refreshTokenExpiration, accessTokenExpiration int64) error {

	err := utils.AddBlackListToken(refreshToken, refreshTokenExpiration)
	if err != nil {
		log.Printf("AuthService - Failed to blacklist refresh token: %v", err)
		return err
	}

	err = utils.AddBlackListToken(accessToken, accessTokenExpiration)
	if err != nil {
		log.Printf("AuthService - Failed to blacklist access token: %v", err)
		return err
	}

	return nil
}
