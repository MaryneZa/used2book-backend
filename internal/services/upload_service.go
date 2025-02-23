package services

import (
	"fmt"
	"io"
	"used2book-backend/internal/utils"
	"used2book-backend/internal/repository/mysql"
)

// ✅ UploadService handles image upload logic
type UploadService struct {
	userRepo *mysql.UserRepository
}

// ✅ NewUploadService creates a new service instance
func NewUploadService(userRepo *mysql.UserRepository) *UploadService {
	return &UploadService{userRepo: userRepo}
}

// ✅ UploadProfileImage uploads an image & saves the URL
func (s *UploadService) UploadProfileImage(userID int, file io.Reader, fileName string) (string, error) {
	// ✅ Upload to ImageKit.io
	uploadURL, err := utils.UploadToImageKit(file, fileName)
	if err != nil {
		return "", fmt.Errorf("upload failed: %v", err)
	}

	// ✅ Save Image URL to Database
	err = s.userRepo.SaveProfileImage(userID, uploadURL)
	if err != nil {
		return "", fmt.Errorf("failed to update user record: %v", err)
	}

	return uploadURL, nil
}

// ✅ UploadProfileImage uploads an image & saves the URL
func (s *UploadService) UploadBackgroundImage(userID int, file io.Reader, fileName string) (string, error) {
	// ✅ Upload to ImageKit.io
	uploadURL, err := utils.UploadToImageKit(file, fileName)
	if err != nil {
		return "", fmt.Errorf("upload failed: %v", err)
	}

	// ✅ Save Image URL to Database
	err = s.userRepo.SaveBackgroundImage(userID, uploadURL)
	if err != nil {
		return "", fmt.Errorf("failed to update user record: %v", err)
	}

	return uploadURL, nil
}