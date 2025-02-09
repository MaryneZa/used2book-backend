package services

import (
	"context"
	"log"
	"used2book-backend/internal/repository/mysql"
	"used2book-backend/internal/utils"
	"used2book-backend/internal/models"

)

type TokenService struct {
	tokenRepo *mysql.TokenRepository
	userRepo *mysql.UserRepository
}

func NewTokenService(token_repo *mysql.TokenRepository, user_repo *mysql.UserRepository) *TokenService {
	return &TokenService{tokenRepo: token_repo, userRepo: user_repo}
}

func (ts *TokenService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	if ts.userRepo == nil {
        log.Fatal("userRepo is not initialized")
    }
	user, err := ts.userRepo.FindByEmail(ctx, email)
	if err != nil {
		log.Printf("Error finding user by email: %v", err)
		return nil, err
	}
	if user == nil {
		log.Println("User not found")
		return nil, nil
	}
	return user, nil
}

func (ts *TokenService) GenerateTokens(ctx context.Context, email string) (string, string, error) {
	log.Println("Starting GenerateTokens")

	user, err := ts.GetUserByEmail(ctx, email)
	if err != nil {
		return "", "", err
	}
	log.Printf("User ID: %v", user.ID)

	err = ts.tokenRepo.DeleteTokensById(ctx, user.ID)
	if err != nil {
		return "", "", err
	}
	// Generate access token
	log.Println("Generating access token")
	accessToken, err := utils.GenerateAccessToken(user.ID)
	if err != nil {
		log.Printf("Error generating access token: %v", err)
		return "", "", err
	}

	// Generate refresh token
	log.Println("Generating refresh token")
	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		log.Printf("Error generating refresh token: %v", err)
		return "", "", err
	}

	// Store refresh token in the database
	log.Println("Storing refresh token in the database")
	err = ts.tokenRepo.StoreRefreshToken(ctx, user.ID, refreshToken, utils.RefreshTokenExpiration())
	if err != nil {
		log.Printf("Error storing refresh token: %v", err)
		return "", "", err
	}

	log.Printf("accesss token : %s", accessToken)
	log.Printf("refresh token : %s", refreshToken)

	return accessToken, refreshToken, nil
}

func (ts *TokenService) ValidateRefreshToken(ctx context.Context, token string) (int, error) {
	return ts.tokenRepo.ValidateRefreshToken(ctx, token)
}

func (ts *TokenService) UpdateRefreshToken(ctx context.Context, userID int, refresh_token string) error {
	return ts.tokenRepo.UpdateRefreshToken(ctx, userID, refresh_token)
}

func (ts *TokenService) DeleteRefreshToken(ctx context.Context, token string) error {
	return ts.tokenRepo.DeleteRefreshToken(ctx, token)
}

func (ts *TokenService) CleanupExpiredTokens(ctx context.Context) error {
	return ts.tokenRepo.DeleteExpiredTokens(ctx)
}
