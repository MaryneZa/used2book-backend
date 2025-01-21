package services

import (
	"context"
	"time"
	"used2book-backend/internal/models"
	"used2book-backend/internal/repository/mysql"
	"fmt"
)

type UserService struct {
	userRepo *mysql.UserRepository
}

func NewUserService(repo *mysql.UserRepository) *UserService {
	return &UserService{userRepo: repo}
}
func (us * UserService) Signup(ctx context.Context, reqUser models.SignupUser) (*models.User, error){
	existing, err := us.userRepo.FindByEmail(ctx, reqUser.Email)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		return nil, fmt.Errorf("user already exists")
	}

	newUser := &models.User{
		Email: reqUser.Email,
		VerifiedEmail: reqUser.VerifiedEmail,
		Name: reqUser.Name,
		ProfilePictureURL: reqUser.ProfilePictureURL,
		Role: "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := us.userRepo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

func (us * UserService) Login(ctx context.Context, reqUser models.LoginUser) (string, error){
	existing, err := us.userRepo.FindByEmail(ctx, reqUser.Email)
	if err != nil {
		return "", err
	}

	if existing == nil {
		return "", fmt.Errorf("There is no this account: ", reqUser.Email)
	}	
	return reqUser.Email, nil
	
}