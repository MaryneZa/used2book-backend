package services

import (
	"context"
	"used2book-backend/internal/models"
	"used2book-backend/internal/repository/mysql"
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


func (us *UserService) SetUserPreferredGenres(ctx context.Context, userID int, genreIDs []int) error {
	return us.userRepo.AddUserPreferredGenres(ctx, userID, genreIDs)
}

func (us *UserService) GetUserPreferences(ctx context.Context, userID int) ([]models.Genre, error) {
	return us.userRepo.GetUserPreferredGenres(ctx, userID)
}

func (us *UserService) CreateBankAccount(ctx context.Context, bank *models.BankAccount) (int, error) {
	return us.userRepo.CreateBankAccount(ctx, bank)
}

func (us *UserService) GetMe(ctx context.Context, userID int) (*models.GetMe, error) {

	user, err := us.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user, err

}

func (us *UserService) EditAccountInfo(ctx context.Context, userID int, firstName string, lastName string, phoneNumber string)  error {
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

func (us *UserService) EditProfile(ctx context.Context, userID int, first_name string, last_name string, address string, quote string, bio string, phone_number string) error {
	return us.userRepo.EditProfile(ctx, userID, first_name, last_name, address, quote, bio, phone_number)
}


func (us *UserService) AddBookToLibrary(ctx context.Context, userID int, bookID int, reading_status int)  (bool, error) {
	return us.userRepo.AddBookToLibrary(ctx , userID, bookID, reading_status)
}

func (us *UserService) AddBookToWishlist(ctx context.Context, userID int, bookID int)  (bool, error) {
	return us.userRepo.AddBookToWishlist(ctx , userID, bookID)
}

func (us *UserService) AddBookToListing(ctx context.Context, userID int, bookID int, price float32, allow_offer bool, imageURLs []string, seller_note string, phone_number string)  (bool, error) {
	return us.userRepo.AddBookToListing(ctx , userID, bookID, price, allow_offer, imageURLs, seller_note, phone_number)
}

func (us *UserService) CountUsers() (int, error) {
	return us.userRepo.CountUsers()
}

func (us *UserService) GetUserLibrary(ctx context.Context, userID int) ([]models.UserLibrary, error){
	return us.userRepo.GetUserLibrary(ctx, userID)
}

func (us *UserService) GetAllListings(ctx context.Context) ([]models.UserListing, error){
	return us.userRepo.GetAllListings(ctx)
}

func (us *UserService) GetPurchasedListingsByUserID(ctx context.Context, userID int) ([]models.MyPurchase, error) {
	return us.userRepo.GetPurchasedListingsByUserID(ctx, userID)
}

func (us *UserService) GetMyOrders(ctx context.Context, sellerID int) ([]models.MyOrder, error) {
	return us.userRepo.GetMyOrders(ctx, sellerID)
}

func (us *UserService) GetUsersByBookInWishlist(ctx context.Context, bookID int) ([]models.WishlistUser, error) {
	return us.userRepo.GetUsersByBookInWishlist(ctx, bookID)
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

// func (us *UserService) GetListingWithBookByID(ctx context.Context, listingID int) (*models.ListingDetails, error) {
// 	return us.userRepo.GetListingWithBookByID(ctx, listingID)
// }

// // UpdateStripeAccountID sets the stripe_account_id for a user
// func (us *UserService) UpdateStripeAccountID(ctx context.Context, userID int, accountID string) error {
//     return us.userRepo.UpdateStripeAccountID(ctx, userID, accountID)
// }

// GetListingByID looks up a single listing
func (us *UserService) GetListingByID(ctx context.Context, listingID int) (*models.ListingDetails, error) {
    return us.userRepo.GetListingByID(ctx, listingID)
}

func (us *UserService) RemoveListing(ctx context.Context, userID int, listingID int) error {
    return us.userRepo.RemoveListing(ctx, userID, listingID)
}

func (us *UserService) MarkListingAsSold(ctx context.Context, listingID int, buyerID int) error {
	return us.userRepo.MarkListingAsSold(ctx, listingID, buyerID)
}

func (us *UserService) UpdateGender(ctx context.Context, userID int, gender string) error {
    return us.userRepo.UpdateGender(ctx, userID, gender)
}

func (us *UserService) GetGender(ctx context.Context, userID int) (string, error) {
    return us.userRepo.GetGender(ctx, userID)
}

func (us *UserService) AddToCart(ctx context.Context, userID int, listingID int) (int, error) {
	return us.userRepo.AddToCart(ctx, userID, listingID)
}

func (us *UserService) RemoveFromCart(ctx context.Context, userID int, listingID int) error {
	return us.userRepo.RemoveFromCart(ctx,userID,listingID)

}
func (us *UserService) GetCart(ctx context.Context, userID int) ([]models.CartItem, error) {
	return us.userRepo.GetCart(ctx, userID)
}

func (us *UserService) AddToOffers(ctx context.Context, buyerID int, listingID int, offeredPrice float64) (int, error) {
    return us.userRepo.AddToOffers(ctx, buyerID, listingID, offeredPrice)
}

func (us *UserService) GetBuyerOffers(ctx context.Context, buyerID int) ([]models.OfferItem, error) {
    return us.userRepo.GetBuyerOffers(ctx, buyerID)
}

func (us *UserService) GetSellerOffers(ctx context.Context, sellerID int) ([]models.OfferItem, error) {
    return us.userRepo.GetSellerOffers(ctx, sellerID)
}

func (us *UserService) RemoveFromOffers(ctx context.Context, buyerID int, listingID int) error {
    return us.userRepo.RemoveFromOffers(ctx, buyerID, listingID)
}

func (us *UserService) AcceptOffer(ctx context.Context, sellerID int, offerID int) (int, error) {
    return us.userRepo.AcceptOffer(ctx, sellerID, offerID)
}

func (us *UserService) RejectOffer(ctx context.Context, sellerID int, offerID int) (int, error) {
    return us.userRepo.RejectOffer(ctx, sellerID, offerID)
}

func (us *UserService) ReserveListing(ctx context.Context, listingID int, buyerID int) (bool, error) {
	return us.userRepo.ReserveListing(ctx, listingID, buyerID)
}

// service/user_service.go
func (us *UserService) GetAcceptedOffer(ctx context.Context, offerID int) (*models.OfferItem, error) {
    return us.userRepo.GetAcceptedOffer(ctx, offerID)
}

func (us *UserService) GetOfferByID(ctx context.Context, offerID int) (*models.OfferItem, error) {
    return us.userRepo.GetOfferByID(ctx, offerID)
}

func (us *UserService) ReserveListingForOffer(ctx context.Context, listingID int, buyerID int) (bool, error) {
    return us.userRepo.ReserveListingForOffer(ctx, listingID, buyerID)
}

func (us *UserService) ExpireReservedListing(ctx context.Context, listingID int) error {
	return us.userRepo.ExpireReservedListing(ctx, listingID)
}

// service/user_service.go
func (us *UserService) RevertOfferReservation(ctx context.Context, listingID int, offerID int) error {
    return us.userRepo.RevertOfferReservation(ctx, listingID, offerID)
}

func (us *UserService) CreateTransaction(ctx context.Context, stripe_session_id string, buyerID int, listingID int, offer_id *int, amount float64, status string) error {
	return us.userRepo.CreateTransaction(ctx, stripe_session_id, buyerID, listingID, offer_id, amount, status)
}

func (us *UserService) UpdateTransactionStatus(ctx context.Context, listingID int, status string) error {
	return us.userRepo.UpdateTransactionStatus(ctx, listingID, status)
}

func (us *UserService) IsListingReserved(ctx context.Context, listingID int) (bool, bool, error) {
	return us.userRepo.IsListingReserved(ctx, listingID)
}

func (us *UserService) GetAllUserReview(ctx context.Context) ([]models.UserReview, error) {
	return us.userRepo.GetAllUserReview(ctx)
}

func (us *UserService) GetAllUserPreferred(ctx context.Context) ([]models.UserPreferred, error) {
	return us.userRepo.GetAllUserPreferred(ctx)
}

// service/user_service.go
func (us *UserService) CreatePost(ctx context.Context, userID int, content string, imageURLs []string, genreID *int, bookID *int) (models.Post, error) {
    return us.userRepo.CreatePost(ctx, userID, content, imageURLs, genreID, bookID)
}

// GetAllPosts retrieves all posts
func (us *UserService) GetAllPosts(ctx context.Context) ([]models.Post, error) {
    return us.userRepo.GetAllPosts(ctx)
}

// GetPostsByUserID retrieves posts by user ID
func (us *UserService) GetPostsByUserID(ctx context.Context, userID int) ([]models.Post, error) {
    return us.userRepo.GetPostsByUserID(ctx, userID)
}

// GetPostByPostID retrieves a post by its ID
func (us *UserService) GetPostByPostID(ctx context.Context, postID int) (models.Post, error) {
    return us.userRepo.GetPostByPostID(ctx, postID)
}

// CreateComment creates a new comment
func (us *UserService) CreateComment(ctx context.Context, postID, userID int, content string) (models.Comment, error) {
    return us.userRepo.CreateComment(ctx, postID, userID, content)
}

// GetCommentsByPostID retrieves comments for a post
func (us *UserService) GetCommentsByPostID(ctx context.Context, postID int) ([]models.Comment, error) {
    return us.userRepo.GetCommentsByPostID(ctx, postID)
}

// CreateLike adds a like to a post
func (us *UserService) CreateLike(ctx context.Context, postID, userID int) (models.Like, error) {
    return us.userRepo.CreateLike(ctx, postID, userID)
}

// RemoveLike removes a like from a post
func (us *UserService) RemoveLike(ctx context.Context, postID, userID int) error {
    return us.userRepo.RemoveLike(ctx, postID, userID)
}

// GetLikeCountByPostID gets the like count for a post
func (us *UserService) GetLikeCountByPostID(ctx context.Context, postID int) (int, error) {
    return us.userRepo.GetLikeCountByPostID(ctx, postID)
}

// IsPostLikedByUser checks if a user has liked a post
func (us *UserService) IsPostLikedByUser(ctx context.Context, postID, userID int) (bool, error) {
    return us.userRepo.IsPostLikedByUser(ctx, postID, userID)
}

func (us *UserService) DeleteUserLibraryByID(ctx context.Context, bookID int) (bool, error) {
    return us.userRepo.DeleteUserLibraryByID(ctx, bookID)
}

func (us *UserService) CreateBookRequest(ctx context.Context, req *models.BookRequest) (int, error) {
    return us.userRepo.CreateBookRequest(ctx, req)
}

func (us *UserService) GetBookRequests(ctx context.Context) ([]*models.BookRequest, error) {
	return us.userRepo.GetBookRequests(ctx)
}


