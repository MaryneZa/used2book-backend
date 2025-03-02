package services

import (
	"context"
	"used2book-backend/internal/models"
	"used2book-backend/internal/repository/mysql"
	"database/sql"
)

type UserService struct {
	userRepo *mysql.UserRepository
}

func NewUserService(repo *mysql.UserRepository) *UserService {
	return &UserService{userRepo: repo}
}

func (us *UserService) GetAllUsers(ctx context.Context) ([]models.GetAllUsers, error) {
	return us.userRepo.GetAllUsers(ctx)
}

// ✅ Fetch user by ID
func (us *UserService) GetUserByID(ctx context.Context, userID int) (*models.GetMe, error) {
	return us.userRepo.FindByID(ctx, userID)
}

func (us * UserService) UpdateOmiseAccountID(ctx context.Context, userID int, omiseAccountID string) error {
	return us.userRepo.UpdateOmiseAccountID(ctx, userID, omiseAccountID)
}

func (us *UserService) SetUserPreferredGenres(ctx context.Context, userID int, genreIDs []int) error {
	return us.userRepo.AddUserPreferredGenres(ctx, userID, genreIDs)
}


func (us *UserService) GetMe(ctx context.Context, userID int) (*models.GetMe, error) {

	user, err := us.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user, err

}

func (us *UserService) EditAccountInfo(ctx context.Context, userID int, firstName string, lastName string, phoneNumber sql.NullString)  error {
	return us.userRepo.EditAccountInfo(ctx, userID, firstName, lastName, phoneNumber)
}

func (us *UserService) EditPhoneNumber(ctx context.Context, userID int, phoneNumber string)  error {
	return us.userRepo.EditPhoneNumber(ctx, userID, phoneNumber)
}

func (us *UserService) IsPhoneNumberTaken(ctx context.Context, phoneNumber string) (bool, error) {
	return us.userRepo.IsPhoneNumberTaken(ctx, phoneNumber)
}

func (us *UserService) EditName(ctx context.Context, userID int, firstName string, lastName string)  error {
	return us.userRepo.EditName(ctx, userID, firstName, lastName)
}

func (us *UserService) EditPreferrence(ctx context.Context, userID int, quote string, bio string)  error {
	return us.userRepo.EditPreferrence(ctx, userID, quote, bio)
}


func (us *UserService) AddBookToLibrary(ctx context.Context, userID int, bookID int, own_status string, price float32, allow_offer bool)  (bool, error) {
	return us.userRepo.AddBookToLibrary(ctx , userID, bookID, own_status, price, allow_offer)
}

func (us *UserService) CountUsers() (int, error) {
	return us.userRepo.CountUsers()
}

func (us *UserService) GetUserLibrary(ctx context.Context, userID int) ([]models.UserLibrary, error){
	return us.userRepo.GetUserLibrary(ctx, userID)
}

func (us *UserService) GetAllListings(ctx context.Context, userID int) ([]models.UserListing, error){
	return us.userRepo.GetAllListings(ctx, userID)
}

func (us *UserService) GetAllListingsByBookID(ctx context.Context, userID int, bookID int) ([]models.UserListing, error){
	return us.userRepo.GetAllListingsByBookID(ctx, userID, bookID)
}

func (us *UserService) GetMyListings(ctx context.Context, userID int) ([]models.UserListing, error){
	return us.userRepo.GetMyListings(ctx, userID)
}

func (us *UserService) GetWishlistByUserID(ctx context.Context, userID int) ([]models.Book, error) {
	return us.userRepo.GetWishlistByUserID(ctx, userID)
}

func (us *UserService) IsBookInWishlist(ctx context.Context, userID int, bookID int) (bool, error) {
	return us.userRepo.IsBookInWishlist(ctx, userID, bookID)
}

func (us *UserService) GetListingWithBookByID(ctx context.Context, listingID int) (*models.ListingDetails, error) {
	return us.userRepo.GetListingWithBookByID(ctx, listingID)
}

// // UpdateStripeAccountID sets the stripe_account_id for a user
// func (us *UserService) UpdateStripeAccountID(ctx context.Context, userID int, accountID string) error {
//     return us.userRepo.UpdateStripeAccountID(ctx, userID, accountID)
// }

// GetListingByID looks up a single listing
func (us *UserService) GetListingByID(ctx context.Context, listingID int) (*models.ListingDetails, error) {
    return us.userRepo.GetListingByID(ctx, listingID)
}

func (us *UserService) MarkListingAsSold(ctx context.Context, listingID int, buyerID int, transactionAmount float32) error {
	return us.userRepo.MarkListingAsSold(ctx, listingID, buyerID, transactionAmount)
}


