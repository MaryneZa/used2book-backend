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

type UserService struct {
	userRepo *mysql.UserRepository
}

func NewUserService(repo *mysql.UserRepository) *UserService {
	return &UserService{userRepo: repo}
}
func (us *UserService) Signup(ctx context.Context, reqUser models.AuthUser) (*models.User, error) {
	existing, err := us.userRepo.FindByEmail(ctx, reqUser.Email)
	if err != nil {
		return nil, err
	}
	log.Println("user provider %s", reqUser.Provider)
	if existing != nil {
		return nil, fmt.Errorf("user already exists")
	}

	// 3. Hash password (bcrypt, argon2, etc.)

	// Before hashing, just declare hashedPassword as an empty string:
	var hashedPassword string

	if reqUser.Provider == "local" {
		hashed, err := utils.HashedPassword(reqUser.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		hashedPassword = hashed
	}

	newUser := &models.User{
		Email:          reqUser.Email,
		HashedPassword: hashedPassword,
		Provider:       reqUser.Provider,
		Role:           "user",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := us.userRepo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

func (us *UserService) Login(ctx context.Context, reqUser models.AuthUser) (*models.User, error) {
	user, err := us.userRepo.FindByEmail(ctx, reqUser.Email)
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

// ✅ Fetch user by ID
func (us *UserService) GetUserByID(ctx context.Context, userID int) (*models.User, error) {
	return us.userRepo.FindByID(ctx, userID)
}

func (us *UserService) GetMe(ctx context.Context, userID int) (*models.GetMe, error) {

	user, err := us.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	me := &models.GetMe{
		ID: 			user.ID,
		Email:          user.Email,
		Name:			user.Name,
		Role:           "user",
	}
	
	return me, err

}
