package services

import (
	"context"
	"fmt"
	"log"
	"time"
	"used2book-backend/internal/models"
	"used2book-backend/internal/repository/mysql"
	"used2book-backend/internal/utils"
)

type AuthService struct {
	userRepo *mysql.UserRepository
}

func NewAuthService(repo *mysql.UserRepository) *AuthService {
	return &AuthService{userRepo: repo}
}
func (as *AuthService) Signup(ctx context.Context, reqUser models.AuthUser) (*models.GetMe, error) {
	existing, err := as.userRepo.FindByEmail(ctx, reqUser.Email)

	if err != nil {
		return nil, err
	}
	log.Println("user provider %s", reqUser.Provider)
	if existing != nil {
		return nil, fmt.Errorf("user already exists")
	}

	var hashedPassword string

	if reqUser.Provider == "local" {
		hashed, err := utils.HashedPassword(reqUser.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		hashedPassword = hashed
	}

	if reqUser.Role == "" {
		reqUser.Role = "user"
	}

	newUser := &models.User{
		Email:          reqUser.Email,
		HashedPassword: hashedPassword,
		Provider:       reqUser.Provider,
		FirstName:      reqUser.FirstName,
		LastName:       reqUser.LastName,
		Role:           reqUser.Role,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	id, err := as.userRepo.Create(ctx, newUser)
	if err != nil {
		return nil, err
	}

	user, err := as.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (as *AuthService) Login(ctx context.Context, reqUser models.AuthUser) (*models.User, error) {
	user, err := as.userRepo.FindByEmail(ctx, reqUser.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, fmt.Errorf("there is no this account: ", reqUser.Email)
	}

	// Verify the password
	if reqUser.Provider == "local" {
		if !utils.CheckPasswordHash(user.HashedPassword, reqUser.Password) {
			return nil, fmt.Errorf("invalid email or password")
		}
	}


	return user, nil

}



