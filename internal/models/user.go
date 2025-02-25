package models

import (
	"database/sql"
	"time"
)

// User represents a user in the system.
type User struct {
	ID             int            `json:"id" db:"id"`
	Email          string         `json:"email" db:"email"`
	Provider       string         `json:"provider" db:"provider"`
	HashedPassword string         `json:"hashed_password" db:"hashed_password"`
	ProfilePicture string `json:"picture_profile" db:"picture_profile"`
	BackgroundPicture string `json:"picture_background" db:"picture_background"`
	FirstName           string    `json:"first_name,omitempty" db:"first_name"`
	LastName           string    `json:"last_name,omitempty" db:"last_name"`
	PhoneNumber        sql.NullString    `json:"phone_number" db:"phone_number"`
	Quote        string    `json:"quote" db:"quote"`
	Bio       string    `json:"bio" db:"bio"`
	Role           string         `json:"role,omitempty" db:"role"`
	CreatedAt      time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at" db:"updated_at"`
}

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
	// PersonalNotes          string          `json:"personal_notes" db:"personal_notes"`
}



type UserListing struct {
	ID         int     `json:"id" db:"id"` // Primary key
	SellerID   int     `json:"seller_id" db:"seller_id"` // Instead of UserID (clearer meaning)
	BookID     int     `json:"book_id" db:"book_id"` // Foreign key
	Status     string  `json:"status" db:"status"` // Listing status
	Price      float32 `json:"price" db:"price"` // Sale price
	AllowOffer bool    `json:"allow_offers" db:"allow_offers"` // Boolean for offers
}


type UserLibrary struct {
	ID        int    `json:"id" db:"id"` // ✅ Primary key (necessary)
	UserID    int    `json:"user_id" db:"user_id"` // ✅ Necessary (foreign key)
	BookID    int    `json:"book_id" db:"book_id"` // ✅ Necessary (foreign key)
	Status    string `json:"status" db:"status"` // ✅ Necessary ('owned', 'not_own', 'wishlist')
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
