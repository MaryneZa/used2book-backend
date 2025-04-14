package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID                int            `json:"id" db:"id"`
	Email             string         `json:"email" db:"email"`
	Provider          string         `json:"provider" db:"provider"`
	HashedPassword    string         `json:"hashed_password" db:"hashed_password"`
	ProfilePicture    string         `json:"picture_profile" db:"picture_profile"`
	BackgroundPicture string         `json:"picture_background" db:"picture_background"`
	FirstName         string         `json:"first_name,omitempty" db:"first_name"`
	LastName          string         `json:"last_name,omitempty" db:"last_name"`
	Address          string         `json:"address,omitempty" db:"address"`
	PhoneNumber       string `json:"phone_number" db:"phone_number"`
	Quote             string         `json:"quote" db:"quote"`
	Bio               string         `json:"bio" db:"bio"`
	Gender            string         `json:"gender" db:"gender"`
	Role              string         `json:"role,omitempty" db:"role"`
	CreatedAt         time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at" db:"updated_at"`
}

type BankAccount struct {
	ID                int
	UserID            int
	BankName          string
	AccountNumber     string
	AccountHolderName string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// User represents a user in the system.

type GetMe struct {
	ID                int            `json:"id" db:"id"`
	Email             string         `json:"email" db:"email"`
	FirstName         string         `json:"first_name,omitempty" db:"first_name"`
	LastName          string         `json:"last_name,omitempty" db:"last_name"`
	ProfilePicture    string         `json:"picture_profile" db:"picture_profile"`
	BackgroundPicture string         `json:"picture_background" db:"picture_background"`
	PhoneNumber       string `json:"phone_number" db:"phone_number"`
	Gender            string         `json:"gender" db:"gender"`
	Quote             string         `json:"quote" db:"quote"`
	Bio               string         `json:"bio" db:"bio"`
	Role              string         `json:"role,omitempty" db:"role"`
	HasBankAccount    bool    `json:"has_bank_account"` // ✅ just a boolean
	Address          string         `json:"address,omitempty" db:"address"`
}

type WishlistUser struct {
	UserID    int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	ProfilePicture  string `json:"picture_profile" db:"picture_profile"`
}

type BookRequest struct {
	ID              int    `json:"id"`
	UserID          int    `json:"user_id"`
	Title           string `json:"title"`
	ISBN            string `json:"isbn"`
	Note            string `json:"note"`
	UserFirstName   string `json:"user_first_name"`
	UserLastName    string `json:"user_last_name"`
	UserEmail       string `json:"user_email"`
	UserPictureProfile string `json:"user_picture_profile"`
    CreatedAt     time.Time `json:"created_at,omitempty" db:"created_at"`
}


type GetAllUsers struct {
	ID                int            `json:"id" db:"id"`
	Email             string         `json:"email" db:"email"`
	FirstName         string         `json:"first_name,omitempty" db:"first_name"`
	LastName          string         `json:"last_name,omitempty" db:"last_name"`
	ProfilePicture    string         `json:"picture_profile" db:"picture_profile"`
	BackgroundPicture string         `json:"picture_background" db:"picture_background"`
	PhoneNumber       string `json:"phone_number" db:"phone_number"`
	Gender            string         `json:"gender" db:"gender"`
	Quote             string         `json:"quote" db:"quote"`
	Bio               string         `json:"bio" db:"bio"`
	Role              string         `json:"role,omitempty" db:"role"`
}

type GetUserProfile struct {
	ID                int    `json:"id" db:"id"`
	FirstName         string `json:"first_name,omitempty" db:"first_name"`
	LastName          string `json:"last_name,omitempty" db:"last_name"`
	ProfilePicture    string `json:"picture_profile" db:"picture_profile"`
	BackgroundPicture string `json:"picture_background" db:"picture_background"`
	Quote             string `json:"quote" db:"quote"`
	Bio               string `json:"bio" db:"bio"`
}

type AuthUser struct {
	Email          string `json:"email" db:"email"`
	GoogleName     string `json:"name" db:"name"`
	FirstName      string `json:"first_name,omitempty" db:"first_name"`
	LastName       string `json:"last_name,omitempty" db:"last_name"`
	Provider       string `json:"provider" db:"provider"`
	Password       string `json:"password" db:"password"`
	ProfilePicture string `json:"picture_profile" db:"picture_profile"`
	VerifiedEmail  bool   `json:"verified_email"`
	Role           string `json:"role,omitempty" db:"role"`
}

type LoginUser struct {
	Email string `json:"email" db:"email"`
}

type OauthUser struct {
	Email string `json:"email" db:"email"`
}

type UserAddListingForm struct {
	BookID     int     `json:"book_id" db:"book_id"`
	Status     string  `json:"status" db:"status"`
	Price      float32 `json:"price" db:"price"`
	AllowOffer bool    `json:"allow_offers" db:"allow_offers"`
	SellerNote string  `json:"seller_note" db:"seller_note"`
	PhoneNumber  string `json:"phone_number" db:"phone_number"`

}

type UserLibrary struct {
	ID     int    `json:"id" db:"id"`           // ✅ Primary key (necessary)
	UserID int    `json:"user_id" db:"user_id"` // ✅ Necessary (foreign key)
	BookID int    `json:"book_id" db:"book_id"` // ✅ Necessary (foreign key)
	ReadingStatus int `json:"reading_status" db:"reading_status"`   // ✅ Necessary (0,1)
}

type UserReview struct {
	ID     int     `json:"id" db:"id"`           // ✅ Primary key (necessary)
	UserID int     `json:"user_id" db:"user_id"` // ✅ Necessary (foreign key)
	BookID int     `json:"book_id" db:"book_id"` // ✅ Necessary (foreign key)
	Rating float32 `json:"rating" db:"rating"`
}

// UserListing represents a book listing for sale.
type UserListing struct {
	ID            int            `json:"id" db:"id"`
	SellerID      int            `json:"seller_id" db:"seller_id"`
	BookID        int            `json:"book_id" db:"book_id"`
	Status        string         `json:"status" db:"status"`
	Price         float32        `json:"price" db:"price"`
	AllowOffer    bool           `json:"allow_offers" db:"allow_offers"`
	ImageURLs     []string  `json:"image_urls"`
}



type UserPreferred struct {
	UserID  int `json:"user_id"`
	GenreID int `json:"genre_id"`
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
