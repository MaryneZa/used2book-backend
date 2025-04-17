package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
	"used2book-backend/internal/config"
	"used2book-backend/internal/models"
	"used2book-backend/internal/services"
	"used2book-backend/internal/utils"

	"github.com/dchest/uniuri"
	"golang.org/x/oauth2"
)

type AuthHandler struct {
	TokenService *services.TokenService
	AuthService  *services.AuthService
}

func (ah * AuthHandler) VerifyTokenHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	userID, err := utils.VerifyToken(req.Token, "access") // Your existing function
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	// Step 3: Success Response
	sendSuccessResponse(w, map[string]interface{}{
		"user_id" : userID,
	})
}


func (ah *AuthHandler) SignupEmailHandler(w http.ResponseWriter, r *http.Request) {

	var user models.AuthUser

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	user.Provider = "local"

	_, err := ah.AuthService.Signup(r.Context(), user)
	if err != nil {
		sendErrorResponse(w, http.StatusConflict, "Signup failed: "+err.Error()) // 409 Conflict if user exists
		return
	}

	sendSuccessResponse(w, map[string]interface{}{
		"success": true,
		"message": "Sign up successfully!",
		"user":    user,
	})
}

func (ah *AuthHandler) LoginEmailHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Parse JSON body for email, password, name, etc.

	var user models.AuthUser

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	user.Provider = "local"

	// Attempt to log the user in

	loginUser, err := ah.AuthService.Login(r.Context(), user)

	if err != nil {
		sendErrorResponse(w, http.StatusUnauthorized, "Login failed: "+err.Error())
		return
	}

	// Generate a JWT token
	accessToken, refreshToken, _ := ah.TokenService.GenerateTokens(r.Context(), loginUser.Email)

	// Set cookies for tokens
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(15 * time.Minute),
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  utils.RefreshTokenExpiration(),
	})

	// Successful login
	sendSuccessResponse(w, map[string]interface{}{
		"success": true,
		"message": "Login successful",
		"token":   accessToken,
	})
}

func (ah *AuthHandler) GoogleHandler(w http.ResponseWriter, r *http.Request) {
	oauthStateString := uniuri.New()
	log.Printf("Login oauthStateString: %s", oauthStateString)
	url := config.GoogleConfig.AuthCodeURL(oauthStateString)
	// log.Printf("Redirecting to: %s", url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (ah *AuthHandler) GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")         // ✅ Allow frontend
	w.Header().Set("Access-Control-Allow-Credentials", "true") // ✅ Allow cookies
	code := r.FormValue("code")

	googleUser, err := ah.processOAuthCallback(config.GoogleConfig, code)

	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	googleUser.Provider = "google"

	_, err = ah.AuthService.Login(r.Context(), *googleUser)
	if err != nil {
		// If "not found", we might do a signup
		_, err := ah.AuthService.Signup(r.Context(), *googleUser)
		if err != nil {
			http.Error(w, "Failed to create user: "+err.Error(), http.StatusConflict)
			return
		}
	}

	// Generate a JWT token
	accessToken, _, _ := ah.TokenService.GenerateTokens(r.Context(), googleUser.Email)

	// ✅ Redirect the user back to the frontend with the token in the URL
	redirectURL := fmt.Sprintf("http://localhost:3000/auth/callback?token=%s", accessToken)
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)

}

func (ah *AuthHandler) processOAuthCallback(config *oauth2.Config, code string) (*models.AuthUser, error) {


	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %v", err)
	}

	// Fetch user info from Google
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user info: %v", err)
	}
	defer response.Body.Close()

	// Read and parse user info
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read user info: %v", err)
	}

	var googleUser models.AuthUser
	if err := json.Unmarshal(contents, &googleUser); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %v", err)
	}

	// Check if email is verified
	if !googleUser.VerifiedEmail {
		return nil, fmt.Errorf("email is not verified")
	}

	googleUser.Provider = "google"

	name_split := strings.Split(googleUser.GoogleName, " ")
	googleUser.FirstName, googleUser.LastName = name_split[0], name_split[1]

	log.Printf("%v", googleUser)
	log.Printf("googleUser: %+v", googleUser)

	return &googleUser, nil
}

func (ah *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "Refresh token not provided", http.StatusUnauthorized)
		return
	}

	refreshToken := cookie.Value

	err = ah.TokenService.DeleteRefreshToken(r.Context(), refreshToken)
	if err != nil {
		http.Error(w, "Error logging out", http.StatusInternalServerError)
		return
	}

	// Clear cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Unix(0, 0), // Expire immediately
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Unix(0, 0), // Expire immediately
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Successfully logged out",
	})
}

func sendSuccessResponse(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(data)
}

func sendErrorResponse(w http.ResponseWriter, statusCode int, message string) {

	log.Println("error:", message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := map[string]interface{}{
		"success": false,
		"message": message,
	}
	json.NewEncoder(w).Encode(response)
}
