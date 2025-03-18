package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	// GoogleLoginConfig  *oauth2.Config
	// GoogleSignupConfig *oauth2.Config
	GoogleConfig *oauth2.Config
)

func InitOAuth() {
	// Load environment variables from .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// log.Println("ENV - main" ,os.Getenv("ENV"))

	// if os.Getenv("ENV") != "production" {
    //     if err := godotenv.Load(); err != nil {
    //         log.Println("Warning: .env file not found, using system environment variables - oauth")
    //     }
    // }

	GoogleConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:6951/auth/callback", // adjust as necessary
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}
	// GoogleLoginConfig = &oauth2.Config{
	//     RedirectURL:  "http://localhost:6951/login/callback", // adjust as necessary
	//     ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	//     ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	//     Scopes: []string{
	//         "https://www.googleapis.com/auth/userinfo.profile",
	//         "https://www.googleapis.com/auth/userinfo.email",
	//     },
	//     Endpoint: google.Endpoint,
	// }

	// GoogleSignupConfig = &oauth2.Config{
	//     RedirectURL:  "http://localhost:6951/signup/callback", // adjust as necessary
	//     ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	//     ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	//     Scopes: []string{
	//         "https://www.googleapis.com/auth/userinfo.profile",
	//         "https://www.googleapis.com/auth/userinfo.email",
	//     },
	//     Endpoint: google.Endpoint,
	// }
	// fmt.Println("DEBUG - Client ID for Login:", GoogleLoginConfig.ClientID)
	// fmt.Println("DEBUG - Client ID for Signup:", GoogleSignupConfig.ClientID)
}
