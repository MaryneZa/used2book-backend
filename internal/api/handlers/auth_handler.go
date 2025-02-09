package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
	"used2book-backend/internal/config"
	"used2book-backend/internal/models"
	"used2book-backend/internal/services"
	"used2book-backend/internal/utils"

	"github.com/dchest/uniuri"
	"golang.org/x/oauth2"
)

type AuthHandler struct {
	UserService  *services.UserService
	TokenService *services.TokenService
}

// POST /auth/google/verify
// The frontend obtains the Google ID token (JWT) from the Google JS client
// and sends it here in the request body. We verify it and then create/login the user.

// Post
func (ah *AuthHandler) SignupEmailHandler(w http.ResponseWriter, r *http.Request) {

	// 1. Parse JSON body for email, password, name, etc.
	var user models.AuthUser

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	user.Provider = "local"

	// 2. Check if user with same email already exists
	_, err := ah.UserService.Signup(r.Context(), user)
	if err != nil {
		sendErrorResponse(w, http.StatusConflict, "Signup failed: "+err.Error()) // 409 Conflict if user exists
		return
	}

	// Step 3: Success Response
	sendSuccessResponse(w, map[string]interface{}{
		"success": true,
		"message": "Sign up successfully!",
	})
}

func (ah *AuthHandler) LoginEmailHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Parse JSON body for email, password, name, etc.

	var user models.AuthUser

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	user.Provider = "local"

	// Attempt to log the user in

	loginUser, err := ah.UserService.Login(r.Context(), user)

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

func (ah *AuthHandler) GetMeHandler(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from request context (set by middleware)
	userID, ok := r.Context().Value("user_id").(int)

	log.Println("userID: ", userID)
	if !ok || userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := ah.UserService.GetMe(r.Context(), userID)
	log.Println("user: ", user)

	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Successful login
	sendSuccessResponse(w, map[string]interface{}{
		"user": user,
	})
}

func (ah *AuthHandler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<a href='/user/google-provider'>Sign up with Google</a>")

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

	user, err := ah.UserService.Login(r.Context(), *googleUser)
	if err != nil {
		// If "not found", we might do a signup
		user, err = ah.UserService.Signup(r.Context(), *googleUser)
		if err != nil {
			http.Error(w, "Failed to create user: "+err.Error(), http.StatusConflict)
			return
		}
	}

	// Generate a JWT token
	accessToken, _ , _ := ah.TokenService.GenerateTokens(r.Context(), user.Email)

	// ✅ Redirect the user back to the frontend with the token in the URL
	redirectURL := fmt.Sprintf("http://localhost:3000/auth/callback?token=%s", accessToken)
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)

}

func (ah *AuthHandler) processOAuthCallback(config *oauth2.Config, code string) (*models.AuthUser, error) {

	// Exchange the code for a token
	// log.Printf("Authorization code: %s", code)
	// log.Printf("Redirect URI: %s", config.RedirectURL)

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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := map[string]interface{}{
		"success": false,
		"message": message,
	}
	json.NewEncoder(w).Encode(response)
}

// func (ah *AuthHandler) VerifyGoogleTokenHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	// 1. Parse JSON body, expecting { "token": "<google_id_token>" }
// 	var body struct {
// 		Token string `json:"token"`
// 	}
// 	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
// 		http.Error(w, "Invalid JSON", http.StatusBadRequest)
// 		return
// 	}

// 	// 2. Verify the token with Google or your chosen method
// 	googleUser, err := ah.verifyGoogleIDToken(body.Token)
// 	if err != nil {
// 		http.Error(w, "Invalid Google token: "+err.Error(), http.StatusUnauthorized)
// 		return
// 	}
// 	log.Printf("Google user verified: %#v", googleUser)

// 	// 3. Login or Signup the user in your DB
// 	// We'll treat googleUser as an AuthUser with Provider="google"
// 	googleUser.Provider = "google"
// 	user, err := ah.UserService.Login(r.Context(), *googleUser)
// 	if err != nil {
// 		// If "not found", we might do a signup
// 		user, err = ah.UserService.Signup(r.Context(), *googleUser)
// 		if err != nil {
// 			http.Error(w, "Failed to create user: "+err.Error(), http.StatusConflict)
// 			return
// 		}
// 	}

// 	// 4. Generate tokens, set cookies
// 	accessToken, refreshToken, _ := ah.TokenService.GenerateTokens(r.Context(), user.Email)

// 	http.SetCookie(w, &http.Cookie{
// 		Name:     "access_token",
// 		Value:    accessToken,
// 		HttpOnly: true,
// 		Secure:   true,
// 		SameSite: http.SameSiteStrictMode,
// 		Expires:  time.Now().Add(15 * time.Minute),
// 	})

// 	http.SetCookie(w, &http.Cookie{
// 		Name:     "refresh_token",
// 		Value:    refreshToken,
// 		HttpOnly: true,
// 		Secure:   true,
// 		SameSite: http.SameSiteStrictMode,
// 		Expires:  utils.RefreshTokenExpiration(),
// 	})

// 	// Respond with success
// 	sendSuccessResponse(w, map[string]interface{}{
// 		"success": true,
// 		"message": "Google login successful",
// 		"user":    user,
// 	})
// }

// // This function actually verifies the Google ID token and returns an AuthUser
// func (ah *AuthHandler) verifyGoogleIDToken(idToken string) (*models.AuthUser, error) {
// 	// Option A: Use a direct HTTP request to Google "tokeninfo"
// 	// Option B: Use "google.golang.org/api/idtoken" package
// 	// Below is a simplified HTTP-based approach:

// 	resp, err := http.Get("https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=" + idToken)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to contact Google tokeninfo: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		body, _ := io.ReadAll(resp.Body)
// 		return nil, fmt.Errorf("bad status code %d from tokeninfo: %s", resp.StatusCode, body)
// 	}

// 	var data map[string]interface{}
// 	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
// 		return nil, fmt.Errorf("failed to decode tokeninfo: %v", err)
// 	}

// 	// data should have keys like "email", "name", "sub", etc.
// 	if data["email"] == nil {
// 		return nil, fmt.Errorf("token has no email")
// 	}

// 	// You should also check "aud" matches your Google Client ID for security,
// 	// e.g. data["aud"] == <your_client_id>

// 	verifiedEmail := data["verified_email"] == "true" ||
// 		data["verified_email"] == true

// 	if !verifiedEmail {
// 		return nil, fmt.Errorf("google email not verified")
// 	}

// 	user := &models.AuthUser{
// 		Email:         data["email"].(string),
// 		Name:          data["name"].(string), // might be missing
// 		VerifiedEmail: verifiedEmail,
// 	}
// 	return user, nil
// func (ah *AuthHandler) LoginGoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {

// 	code := r.FormValue("code")

// 	googleUser, err := ah.processOAuthCallback(config.GoogleLoginConfig, code)

// 	if err != nil {
// 		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	// Attempt to log the user in
// 	user, err := ah.UserService.Login(r.Context(), *googleUser)

// 	if err != nil {
// 		sendErrorResponse(w, http.StatusUnauthorized, "Login failed: "+err.Error())
// 		return
// 	}

// 	// Generate a JWT token
// 	accessToken, refreshToken, _ := ah.TokenService.GenerateTokens(r.Context(), user.Email)

// 	// Set cookies for tokens
// 	http.SetCookie(w, &http.Cookie{
// 		Name:     "access_token",
// 		Value:    accessToken,
// 		HttpOnly: true,
// 		Secure:   true,
// 		SameSite: http.SameSiteStrictMode,
// 		Expires:  time.Now().Add(15 * time.Minute),
// 	})

// 	http.SetCookie(w, &http.Cookie{
// 		Name:     "refresh_token",
// 		Value:    refreshToken,
// 		HttpOnly: true,
// 		Secure:   true,
// 		SameSite: http.SameSiteStrictMode,
// 		Expires:  utils.RefreshTokenExpiration(),
// 	})

// 	// Successful login
// 	sendSuccessResponse(w, map[string]interface{}{
// 		"success":      true,
// 		"message":      "Login successful",
// 		"googleUser":   googleUser, // Include the user object if needed
// 		"user":         user,
// 		"token": accessToken,
// 	})
// }

// func (ah *AuthHandler) SignupGoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
// 	code := r.FormValue("code")
// 	if code == "" {
// 		sendErrorResponse(w, http.StatusBadRequest, "Missing OAuth code")
// 		return
// 	}

// 	// Step 1: Process OAuth Callback
// 	googleUser, err := ah.processOAuthCallback(config.GoogleSignupConfig, code)
// 	if err != nil {
// 		sendErrorResponse(w, http.StatusInternalServerError, "OAuth processing failed: "+err.Error())
// 		return
// 	}

// 	// Step 2: Sign Up the User
// 	_, err = ah.UserService.Signup(r.Context(), *googleUser)
// 	if err != nil {
// 		sendErrorResponse(w, http.StatusConflict, "Signup failed: "+err.Error()) // 409 Conflict if user exists
// 		return
// 	}

// 	// Step 3: Success Response
// 	sendSuccessResponse(w, map[string]interface{}{
// 		"success": true,
// 		"message": "Sign up successfully!",
// 	})
// }
// }
// func (ah *AuthHandler) LoginGoogleHandler(w http.ResponseWriter, r *http.Request) {
// 	oauthStateString := uniuri.New()
// 	log.Printf("Login oauthStateString: %s", oauthStateString)
// 	url := config.GoogleLoginConfig.AuthCodeURL(oauthStateString)
// 	// log.Printf("Redirecting to: %s", url)
// 	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
// }

// func (ah *AuthHandler) SignupGoogleHandler(w http.ResponseWriter, r *http.Request) {
// 	oauthStateString := uniuri.New()
// 	log.Printf("Signup oauthStateString: %s", oauthStateString)
// 	url := config.GoogleSignupConfig.AuthCodeURL(oauthStateString)
// 	// log.Printf("Redirecting to: %s", url)
// 	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
// }
