package models

import "time"

// User represents a user in the system.
type User struct {
    ID                string    `json:"id" db:"id"`
    Email             string    `json:"email" db:"email"`
    VerifiedEmail     bool      `json:"verified_email" db:"verified_email`
    Name              string    `json:"name" db:"name"`
    ProfilePictureURL string    `json:"picture" db:"picture"`
    Role              string    `json:"role,omitempty" db:"role"`
    CreatedAt         time.Time `json:"created_at" db:"created_at"`
    UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

type SignupUser struct {
    Email             string    `json:"email" db:"email"`
    VerifiedEmail     bool      `json:"verified_email" db:"verified_email`
    Name              string    `json:"name" db:"name"`
    ProfilePictureURL string    `json:"picture" db:"picture"`
}

type LoginUser struct {
    Email             string    `json:"email" db:"email"`
    VerifiedEmail     bool      `json:"verified_email" db:"verified_email`
}

type OauthUser struct {
    Email             string    `json:"email" db:"email"`
    VerifiedEmail     bool      `json:"verified_email" db:"verified_email`
}
// type GoogleUser struct {
//     ID            string `json:"id"`
//     Email         string `json:"email"`
//     VerifiedEmail bool   `json:"verified_email"`
//     Name          string `json:"name"`
//     ProfilePictureURL       string `json:"picture"`
//     GivenName     string `json:"given_name"`
//     FamilyName    string `json:"family_name"`
//     Link          string `json:"link"`
//     Gender        string `json:"gender"`
//     Locale        string `json:"locale"`
// }


