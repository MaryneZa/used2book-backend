package config

import (
    "fmt"
    "log"
    "os"

    "github.com/joho/godotenv"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
)

var (
    GoogleLoginConfig  *oauth2.Config
    GoogleSignupConfig *oauth2.Config
)

func InitOAuth() {
    // Load environment variables from .env
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

	GoogleLoginConfig = &oauth2.Config{
        RedirectURL:  "http://localhost:3000/login/callback", // adjust as necessary
        ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
        ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
        Scopes: []string{
            "https://www.googleapis.com/auth/userinfo.profile",
            "https://www.googleapis.com/auth/userinfo.email",
        },
        Endpoint: google.Endpoint,
    }

	GoogleSignupConfig = &oauth2.Config{
        RedirectURL:  "http://localhost:3000/signup/callback", // adjust as necessary
        ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
        ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
        Scopes: []string{
            "https://www.googleapis.com/auth/userinfo.profile",
            "https://www.googleapis.com/auth/userinfo.email",
        },
        Endpoint: google.Endpoint,
    }
    fmt.Println("DEBUG - Client ID for Login:", GoogleLoginConfig.ClientID)
    fmt.Println("DEBUG - Client ID for Signup:", GoogleSignupConfig.ClientID)
}
