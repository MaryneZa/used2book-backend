package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID              int            `json:"id" db:"id"`
	Email           string         `json:"email" db:"email"`
	Provider        string         `json:"provider" db:"provider"`
	HashedPassword  string         `json:"hashed_password" db:"hashed_password"`
	ProfilePicture  string         `json:"picture_profile" db:"picture_profile"`
	BackgroundPicture string       `json:"picture_background" db:"picture_background"`
	FirstName       string         `json:"first_name,omitempty" db:"first_name"`
	LastName        string         `json:"last_name,omitempty" db:"last_name"`
	PhoneNumber     sql.NullString `json:"phone_number" db:"phone_number"`
	OmiseAccountID  sql.NullString `json:"omise_account_id" db:"omise_account_id"` // ✅ Added Omise account ID
	Quote           string         `json:"quote" db:"quote"`
	Bio             string         `json:"bio" db:"bio"`
	Gender 		string `json:"gender" db:"gender"`
	Role            string         `json:"role,omitempty" db:"role"`
	CreatedAt       time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at" db:"updated_at"`
}

// User represents a user in the system.

type GetMe struct {
	ID             int            `json:"id" db:"id"`
	Email          string         `json:"email" db:"email"`
	FirstName           string    `json:"first_name,omitempty" db:"first_name"`
	LastName           string    `json:"last_name,omitempty" db:"last_name"`
	ProfilePicture string `json:"picture_profile" db:"picture_profile"`
	BackgroundPicture string `json:"picture_background" db:"picture_background"`
	PhoneNumber        sql.NullString   `json:"phone_number" db:"phone_number"`
	Gender 		string `json:"gender" db:"gender"`
	Quote        string    `json:"quote" db:"quote"`
	Bio       string    `json:"bio" db:"bio"`
	Role           string         `json:"role,omitempty" db:"role"`
	OmiseAccountID  sql.NullString `json:"omise_account_id" db:"omise_account_id"` // ✅ Added Omise account ID

}

type GetAllUsers struct {
	ID             int            `json:"id" db:"id"`
	Email          string         `json:"email" db:"email"`
	FirstName           string    `json:"first_name,omitempty" db:"first_name"`
	LastName           string    `json:"last_name,omitempty" db:"last_name"`
	ProfilePicture string `json:"picture_profile" db:"picture_profile"`
	BackgroundPicture string `json:"picture_background" db:"picture_background"`
	PhoneNumber        sql.NullString   `json:"phone_number" db:"phone_number"`
	Gender 		string `json:"gender" db:"gender"`
	Quote        string    `json:"quote" db:"quote"`
	Bio       string    `json:"bio" db:"bio"`
	Role           string         `json:"role,omitempty" db:"role"`
}

type GetUserProfile struct {
	ID             int            `json:"id" db:"id"`
	FirstName           string    `json:"first_name,omitempty" db:"first_name"`
	LastName           string    `json:"last_name,omitempty" db:"last_name"`
	ProfilePicture string `json:"picture_profile" db:"picture_profile"`
	BackgroundPicture string `json:"picture_background" db:"picture_background"`
	Quote        string    `json:"quote" db:"quote"`
	Bio       string    `json:"bio" db:"bio"`
}

type AuthUser struct {
	Email          string `json:"email" db:"email"`
	GoogleName     string `json:"name" db:"name"`
	FirstName           string    `json:"first_name,omitempty" db:"first_name"`
	LastName           string    `json:"last_name,omitempty" db:"last_name"`
	Provider       string `json:"provider" db:"provider"`
	Password       string `json:"password" db:"password"`
	ProfilePicture string `json:"picture_profile" db:"picture_profile"`
	VerifiedEmail  bool   `json:"verified_email"`
	Role           string         `json:"role,omitempty" db:"role"`
}

type LoginUser struct {
	Email string `json:"email" db:"email"`
}

type OauthUser struct {
	Email string `json:"email" db:"email"`
}

type UserAddLibraryForm struct {
	BookID             int            `json:"book_id" db:"book_id"`
	Status          string          `json:"status" db:"status"`
	Price		float32          `json:"price" db:"price"`
	AllowOffer     bool          `json:"allow_offers" db:"allow_offers"`
	SellerNote          string          `json:"seller_note" db:"seller_note"`
}


type UserLibrary struct {
	ID        int    `json:"id" db:"id"` // ✅ Primary key (necessary)
	UserID    int    `json:"user_id" db:"user_id"` // ✅ Necessary (foreign key)
	BookID    int    `json:"book_id" db:"book_id"` // ✅ Necessary (foreign key)
	Status    string `json:"status" db:"status"` // ✅ Necessary ('owned', 'not_own', 'wishlist')
}



type UserReview struct {
	ID        int    `json:"id" db:"id"` // ✅ Primary key (necessary)
	UserID    int    `json:"user_id" db:"user_id"` // ✅ Necessary (foreign key)
	BookID    int    `json:"book_id" db:"book_id"` // ✅ Necessary (foreign key)
	Rating    float32    `json:"rating" db:"rating"`
}


// UserListing represents a book listing for sale.
type UserListing struct {
	ID            int             `json:"id" db:"id"`
	SellerID      int             `json:"seller_id" db:"seller_id"`
	BookID        int             `json:"book_id" db:"book_id"`
	Status        string          `json:"status" db:"status"`
	Price         float32         `json:"price" db:"price"`
	AllowOffer    bool            `json:"allow_offers" db:"allow_offers"`
	SellerOmiseID sql.NullString  `json:"seller_omise_id" db:"seller_omise_id"` // ✅ Added Omise ID
}

// ListingDetails combines Listing + Book Data.
type ListingDetails struct {
    ListingID    int     `json:"listing_id"`
    SellerID     int     `json:"seller_id"`
    SellerOmiseID sql.NullString `json:"seller_omise_id" db:"seller_omise_id"` // ✅ Seller's Omise ID
    BookID       int     `json:"book_id"`
    Price        float32 `json:"price"`
    Status       string  `json:"status"`
    AllowOffers  bool    `json:"allow_offers"`
	SellerNote          string          `json:"seller_note" db:"seller_note"`

    // Book details
    Title         string    `json:"title"`
    Author        string    `json:"author"`
    Description   string    `json:"description,omitempty"`
    Language      string    `json:"language,omitempty"`
    ISBN          string    `json:"isbn,omitempty"`
    Publisher     string    `json:"publisher,omitempty"`
    PublishDate   time.Time `json:"publish_date,omitempty"`
    CoverImageURL string    `json:"cover_image_url,omitempty"`
    AverageRating string    `json:"average_rating,omitempty"`
    NumRatings    string    `json:"num_ratings,omitempty"`
	ImageURLs     []string  `json:"image_urls"`
}


type CartItem struct {
    ID            int     `json:"id"`
    UserID        int     `json:"user_id"`
    ListingID     int     `json:"listing_id"`
    BookID        int     `json:"book_id"`
    Price         float32 `json:"price"`
    AllowOffers   bool    `json:"allow_offers"`
    SellerID      int     `json:"seller_id"`
    BookTitle     string  `json:"book_title"`
    BookAuthor    string  `json:"book_author"`
    CoverImageURL string  `json:"cover_image_url"`
    ImageURL      string  `json:"image_url,omitempty"` // Added for first listing image
}

type UserPreferred struct {
    UserID       int     `json:"user_id"`
    GenreID 	 int     `json:"genre_id"`
}

// type GoogleUser struct {
//     ID            string `json:"id"`
//     Email         string `json:"email"`
//     VerifiedEmail bool   `json:"verified_email"`
//     Name          string `json:"name"`
//     ProfilePicture       string `json:"picture_profile"`
//     GivenName     string `json:"given_name"`
//     FamilyName    string `json:"family_name"`
//     Link          string `json:"link"`
//     Gender        string `json:"gender"`
//     Locale        string `json:"locale"`
// }
