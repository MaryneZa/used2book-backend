package services

import (
	"context"
	"log"
	"used2book-backend/internal/models"
	"used2book-backend/internal/repository/mysql"
	"github.com/joho/godotenv"
	"os"
	"errors"
	"fmt"
)

// BookService handles book-related operations
type BookService struct {
	bookRepo *mysql.BookRepository
}

// NewBookService initializes BookService
func NewBookService(bookRepo *mysql.BookRepository) *BookService {
	return &BookService{bookRepo: bookRepo}
}

func (bs *BookService) GetAllBooks(ctx context.Context) ([]models.Book, error) {
	return bs.bookRepo.GetAllBooks(ctx)
}


func (bs *BookService) GetBookByID(ctx context.Context, bookID int) (*models.Book, error) {
	return bs.bookRepo.GetBookByID(ctx, bookID)
}

func (bs *BookService) GetAllBookGenres(ctx context.Context) ([]models.BookGenre, error) {
	return bs.bookRepo.GetAllBookGenres(ctx)
}

// // GetBookWithRatings retrieves a book with its ratings
// func (bs *BookService) GetBookWithRatings(ctx context.Context, bookID int) (*models.BookWithRatings, error) {
// 	return bs.bookRepo.GetBookWithRatings(ctx, bookID)
// }
// // AddOrUpdateUserRating allows a user to add or update their rating for a book
// func (bs *BookService) AddOrUpdateUserRating(ctx context.Context, userID int, bookID int, rating float64) error {
// 	return bs.bookRepo.UpdateUserBookRating(ctx, userID, bookID, rating)
// }

// // DeleteUserRating removes a user's rating for a book
// func (bs *BookService) DeleteUserRating(ctx context.Context, userID int, bookID int) error {
// 	return bs.bookRepo.DeleteUserRating(ctx, userID, bookID)
// }

// SyncBooksFromGoogleSheets loads books from Google Sheets to MySQL
func (bs *BookService) SyncBooksFromGoogleSheets(sheetID, apiKey string) error {
	return bs.bookRepo.SyncBooksFromGoogleSheets(sheetID, apiKey)
}

func (bs *BookService) CountBooks() (int, error) {
	return bs.bookRepo.CountBooks()
}

func (bs *BookService) GetGenresByBookID(ctx context.Context, bookID int) ([]string, error) {
	return bs.bookRepo.GetGenresByBookID(ctx, bookID)
}

func (bs *BookService) AddBookReview(ctx context.Context, userID int, bookID int, rating float32, comment string) error {
	return bs.bookRepo.AddBookReview(ctx, userID, bookID, rating, comment)
}

func (bs *BookService) GetReviewsByBookID(ctx context.Context, bookID int) ([]models.BookReview, error) {
	return bs.bookRepo.GetReviewsByBookID(ctx, bookID)
}

func (bs *BookService) GetAllGenres(ctx context.Context) ([]models.Genre, error) {
	return bs.bookRepo.GetAllGenres(ctx)
}






// SyncBooksIfNeeded checks if books exist, then syncs from Google Sheets if empty
func (bs *BookService) SyncBooksIfNeeded() {
	bookCount, err := bs.bookRepo.CountBooks()
	if err != nil {
		log.Fatalf("Failed to check book count: %v", err)
	}

	if bookCount == 0 {
		log.Println("📚 No books found in database, syncing from Google Sheets...")
		sheetID, err := getGoogleSheetID()
		if err != nil {
			fmt.Errorf("failed to get GOOGLE_SHEET_ID: %v", err)
		}
		apiKey, err := getGoogleSheetAPIKey()
		if err != nil {
			fmt.Errorf("failed to get GOOGLE_SHEET_API_KEY: %v", err)
		}

		if sheetID != "" && apiKey != "" {
			if err := bs.SyncBooksFromGoogleSheets(sheetID, apiKey); err != nil {
				log.Printf("❌ Failed to sync books from Google Sheets: %v", err)
			} else {
				log.Println("✅ Books synced successfully from Google Sheets!")
			}
		} else {
			log.Println("⚠️ Missing GOOGLE_SHEET_ID or GOOGLE_SHEET_API_KEY, skipping book sync.")
		}
	} else {
		log.Println("✅ Books already exist in database, skipping sync.")
	}
}
func getGoogleSheetID() (string, error) {
    if err := godotenv.Load(); err != nil {
        return "", errors.New("failed to load .env file")
    }

	// log.Println("ENV - gg_id" ,os.Getenv("ENV"))

	// if os.Getenv("ENV") != "production" {
    //     if err := godotenv.Load(); err != nil {
    //         log.Println("Warning: .env file not found, using system environment variables - gg_id")
    //     }
    // }
    secret := os.Getenv("GOOGLE_SHEET_ID")
    
    if secret == "" {
        return "", errors.New("GOOGLE_SHEET_ID is not set in .env file")
    }
    return secret, nil
}
func getGoogleSheetAPIKey() (string, error) {

    if err := godotenv.Load(); err != nil {
        return "", errors.New("failed to load .env file")
    }

	// log.Println("ENV - gg_key" ,os.Getenv("ENV"))

	// if os.Getenv("ENV") != "production" {
    //     if err := godotenv.Load(); err != nil {
    //         log.Println("Warning: .env file not found, using system environment variables - key")
    //     }
    // }
    secret := os.Getenv("GOOGLE_SHEET_API_KEY")
    
    if secret == "" {
        return "", errors.New("GOOGLE_SHEET_API_KEY is not set in .env file")
    }
    return secret, nil
}


