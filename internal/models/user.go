package models

import (
	"time"
	"database/sql"
)

// User represents a user in the system.
type User struct {
	ID             int       `json:"id" db:"id"`
	Email          string    `json:"email" db:"email"`
	Name           sql.NullString    `json:"username" db:"username"`
	Provider       string    `json:"provider" db:"provider"`
	HashedPassword string    `json:"hashed_password" db:"hashed_password"`
	ProfilePicture sql.NullString    `json:"picture" db:"picture"`
	Role           string    `json:"role,omitempty" db:"role"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type GetMe struct {
	ID             int       `json:"id" db:"id"`
	Email          string    `json:"email" db:"email"`
	Name           sql.NullString    `json:"username" db:"username"`
	ProfilePicture sql.NullString    `json:"picture" db:"picture"`
	Role           string    `json:"role,omitempty" db:"role"`
}

type AuthUser struct {
	Email          string `json:"email" db:"email"`
	Name           string `json:"username" db:"username"`
	Provider       string `json:"provider" db:"provider"`
	Password       string `json:"password" db:"password"`
	ProfilePicture string `json:"picture" db:"picture"`
	VerifiedEmail  bool   `json:"verified_email"`
}

type LoginUser struct {
	Email string `json:"email" db:"email"`
}

type OauthUser struct {
	Email string `json:"email" db:"email"`
}

// type GoogleUser struct {
//     ID            string `json:"id"`
//     Email         string `json:"email"`
//     VerifiedEmail bool   `json:"verified_email"`
//     Name          string `json:"name"`
//     ProfilePicture       string `json:"picture"`
//     GivenName     string `json:"given_name"`
//     FamilyName    string `json:"family_name"`
//     Link          string `json:"link"`
//     Gender        string `json:"gender"`
//     Locale        string `json:"locale"`
// }
