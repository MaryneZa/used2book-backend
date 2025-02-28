package handlers

import (
	"net/http"
	"encoding/json"
	"log"
	"used2book-backend/internal/services"
	"used2book-backend/internal/models"
	"io"
	"strconv"
	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	UserService  *services.UserService
	UploadService *services.UploadService
}

func (uh *UserHandler) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {

	// Call the BookService method to get the total book count
	users, err := uh.UserService.GetAllUsers(r.Context())
	if err != nil {
		// Handle the error, e.g., return a 500 Internal Server Error
		http.Error(w, "Failed to get all users", http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(w, map[string]interface{}{
		"users": users,
	})
	
}

func (uh *UserHandler) GetMeHandler(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from request context (set by middleware)
	userID, ok := r.Context().Value("user_id").(int)

	log.Println("userID: ", userID)
	if !ok || userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := uh.UserService.GetMe(r.Context(), userID)
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

// ✅ UploadProfileImage HTTP Handler
func (uh *UserHandler) UploadProfileImageHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10MB max file size
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "File too large")
		return
	}

	// ✅ Get file from form-data
	file, handler, err := r.FormFile("image")
	log.Println("UploadProfileImageHandler: error from UploadService:", err)

	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Error retrieving file")
		return
	}
	defer file.Close()

	// ✅ Get userID from request context (assuming authentication middleware sets it)
	userID := r.Context().Value("user_id").(int)

	// ✅ Upload image using Service Layer
	uploadURL, err := uh.UploadService.UploadProfileImage(userID, file, handler.Filename)
	log.Println("UploadProfileImageHandler: error from UploadService:", err)

	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Image upload failed")
		return
	}

	// ✅ Send response
	sendSuccessResponse(w, map[string]interface{}{
		"success": true,
		"image_profile_url": uploadURL,
	})
}

// ✅ UploadProfileImage HTTP Handler
func (uh *UserHandler) UploadBackgroundImageHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10MB max file size
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "File too large")
		return
	}

	// ✅ Get file from form-data
	file, handler, err := r.FormFile("image")
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Error retrieving file")
		return
	}
	defer file.Close()

	// ✅ Get userID from request context (assuming authentication middleware sets it)
	userID := r.Context().Value("user_id").(int)

	// ✅ Upload image using Service Layer
	uploadURL, err := uh.UploadService.UploadBackgroundImage(userID, file, handler.Filename)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Image upload failed")
		return
	}

	// ✅ Send response
	sendSuccessResponse(w, map[string]interface{}{
		"success": true,
		"image_background_url": uploadURL,
	})
}

func (uh *UserHandler) EditAccountInfoHandler(w http.ResponseWriter, r *http.Request) {

	// 1. Parse JSON body for email, password, name, etc.
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	userID := r.Context().Value("user_id").(int)

	user.Provider = "local"

	// 2. Check if user with same email already exists
	err := uh.UserService.EditAccountInfo(r.Context(), userID, user.FirstName, user.LastName, user.PhoneNumber)
	if err != nil {
		sendErrorResponse(w, http.StatusConflict, "Edit Account Info "+err.Error()) // 409 Conflict if user exists
		return
	}

	// Step 3: Success Response
	sendSuccessResponse(w, map[string]interface{}{
		"success": true,
		"message": "Edited info successfully!",
	})
}

func (uh *UserHandler) EditUserNameHandler(w http.ResponseWriter, r *http.Request) {

	// 1. Parse JSON body for email, password, name, etc.
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	userID := r.Context().Value("user_id").(int)
	user.Provider = "local"

	// 2. Check if user with same email already exists
	err := uh.UserService.EditName(r.Context(), userID, user.FirstName, user.LastName)
	if err != nil {
		sendErrorResponse(w, http.StatusConflict, "Edit Username "+err.Error()) // 409 Conflict if user exists
		return
	}

	// Step 3: Success Response
	sendSuccessResponse(w, map[string]interface{}{
		"success": true,
		"message": "Edited username successfully!",
	})
}

func (uh *UserHandler) EditPreferrenceHandler(w http.ResponseWriter, r *http.Request) {

	// 1. Parse JSON body for email, password, name, etc.
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	userID := r.Context().Value("user_id").(int)

	user.Provider = "local"

	// 2. Check if user with same email already exists
	err := uh.UserService.EditPreferrence(r.Context(), userID, user.Quote, user.Bio)
	if err != nil {
		sendErrorResponse(w, http.StatusConflict, "Edit Preferrence "+err.Error()) // 409 Conflict if user exists
		return
	}

	// Step 3: Success Response
	sendSuccessResponse(w, map[string]interface{}{
		"success": true,
		"message": "Edited preferrence successfully!",
	})
}

func (uh *UserHandler) AddBookToLibraryHandler(w http.ResponseWriter, r *http.Request) {
	// Log the incoming request body
	body, _ := io.ReadAll(r.Body)
	log.Println("Request Body:", string(body)) // ✅ Print raw JSON body for debugging

	var user models.UserAddLibraryForm

	if err := json.Unmarshal(body, &user); err != nil {
		log.Println("JSON Decode Error:", err) // ✅ Log the error
		sendErrorResponse(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		sendErrorResponse(w, http.StatusUnauthorized, "User ID missing")
		return
	}

	_, err := uh.UserService.AddBookToLibrary(r.Context(), userID, user.BookID, user.Status, user.Price, user.AllowOffer)
	if err != nil {
		sendErrorResponse(w, http.StatusConflict, "Edit Preference "+err.Error()) 
		return
	}

	sendSuccessResponse(w, map[string]interface{}{
		"success": true,
		"message": "Book added to library successfully!",
	})
}

func (uh *UserHandler) AddBookToWishListHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("📡 Incoming request: AddBookToWishListHandler")

	bookIDStr := chi.URLParam(r, "bookID")
	log.Println("📖 Extracted Book ID String:", bookIDStr)

	if bookIDStr == "" {
		sendErrorResponse(w, http.StatusBadRequest, "Book ID is required")
		return
	}

	bookID, err := strconv.Atoi(bookIDStr)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid book ID: "+err.Error())
		return
	}

	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		sendErrorResponse(w, http.StatusUnauthorized, "User ID missing")
		return
	}

	// ✅ Call service to toggle wishlist and get status
	inWishlist, err := uh.UserService.AddBookToLibrary(r.Context(), userID, bookID, "wishlist", 0, false)
	if err != nil {
		log.Println("❌ Wishlist Error:", err)
		sendErrorResponse(w, http.StatusConflict, "Wishlist error: "+err.Error()) 
		return
	}

	// ✅ Return updated wishlist status
	log.Println("✅ Wishlist updated successfully")
	sendSuccessResponse(w, map[string]interface{}{
		"success": true,
		"message": "Wishlist updated successfully!",
		"in_wishlist": inWishlist, // ✅ Return wishlist status for button toggle
	})
}



func (uh *UserHandler) GetUserCount(w http.ResponseWriter, r *http.Request) {

	// Call the BookService method to get the total book count
	count, err := uh.UserService.CountUsers()
	if err != nil {
		// Handle the error, e.g., return a 500 Internal Server Error
		http.Error(w, "Failed to count books", http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(w, map[string]interface{}{
		"count": count,
	})
	
}

func (uh *UserHandler) GetAllListingsHandler(w http.ResponseWriter, r *http.Request) {

	// Call the BookService method to get the total book count
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
        sendErrorResponse(w, http.StatusUnauthorized, "User ID missing")
        return
    }

	listing, err := uh.UserService.GetAllListings(r.Context(), userID)
	if err != nil {
		// Handle the error, e.g., return a 500 Internal Server Error
		http.Error(w, "Failed to get listing", http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(w, map[string]interface{}{
		"listing": listing,
	})
	
}

func (uh *UserHandler) GetMyListingsHandler(w http.ResponseWriter, r *http.Request) {

	// Call the BookService method to get the total book count
	userID := r.Context().Value("user_id").(int)


	listing, err := uh.UserService.GetMyListings(r.Context(), userID)
	if err != nil {
		// Handle the error, e.g., return a 500 Internal Server Error
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to get listing: "+err.Error()) // 409 Conflict if user exists
		return
	}

	sendSuccessResponse(w, map[string]interface{}{
		"listing": listing,
	})
	
}



func (uh *UserHandler) GetMyLibraryHandler(w http.ResponseWriter, r *http.Request) {

	// Call the BookService method to get the total book count
	userID := r.Context().Value("user_id").(int)


	library, err := uh.UserService.GetUserLibrary(r.Context(), userID)
	if err != nil {
		// Handle the error, e.g., return a 500 Internal Server Error
		http.Error(w, "Failed to get library", http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(w, map[string]interface{}{
		"library": library,
	})
	
}

func (uh *UserHandler) GetUserListingsHandler(w http.ResponseWriter, r *http.Request) {

	
	userIDStr := chi.URLParam(r, "userID")

    userID, err := strconv.Atoi(userIDStr) // Convert to int
	log.Println("userID: ", userID)
    if err != nil {
        sendErrorResponse(w, http.StatusBadRequest, "Invalid book ID")
        return
    }
	listing, err := uh.UserService.GetMyListings(r.Context(), userID)
	if err != nil {
		// Handle the error, e.g., return a 500 Internal Server Error
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to get listing: "+err.Error()) // 409 Conflict if user exists
		return
	}

	sendSuccessResponse(w, map[string]interface{}{
		"listing": listing,
	})
	
}



func (uh *UserHandler) GetUserLibraryHandler(w http.ResponseWriter, r *http.Request) {
	
	userIDStr := chi.URLParam(r, "userID")

    userID, err := strconv.Atoi(userIDStr) // Convert to int
	log.Println("userID: ", userID)
    if err != nil {
        sendErrorResponse(w, http.StatusBadRequest, "Invalid book ID")
        return
    }

	library, err := uh.UserService.GetUserLibrary(r.Context(), userID)
	if err != nil {
		// Handle the error, e.g., return a 500 Internal Server Error
		http.Error(w, "Failed to get library", http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(w, map[string]interface{}{
		"library": library,
	})
	
}

func (uh *UserHandler) GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from request context (set by middleware)

	userIDStr := chi.URLParam(r, "userID")

    userID, err := strconv.Atoi(userIDStr) // Convert to int
	log.Println("userID: ", userID)
    if err != nil {
        sendErrorResponse(w, http.StatusBadRequest, "Invalid book ID")
        return
    }

	user, err := uh.UserService.GetMe(r.Context(), userID)
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

func (uh *UserHandler) IsBookInWishlistHandler(w http.ResponseWriter, r *http.Request) {
	// Parse user ID and book ID from request (e.g., query parameters)
	// userIDStr := r.URL.Query().Get("user_id")
	userID := r.Context().Value("user_id").(int)

	bookIDStr := chi.URLParam(r, "bookID")

	// Validate that both user ID and book ID are provided
	if bookIDStr == "" {
		http.Error(w, "Book ID are required", http.StatusBadRequest)
		return
	}

	bookID, err := strconv.Atoi(bookIDStr)
	if err != nil {
		http.Error(w, "Invalid Book ID", http.StatusBadRequest)
		return
	}

	// Call the UserService to check if the book is in the wishlist
	isInWishlist, err := uh.UserService.IsBookInWishlist(r.Context(), userID, bookID)
	if err != nil {
		// Handle the error, e.g., return a 500 Internal Server Error
		http.Error(w, "Failed to check book wishlist status", http.StatusInternalServerError)
		return
	}

	// Send the response indicating whether the book is in the wishlist
	sendSuccessResponse(w, map[string]interface{}{
		"in_wishlist": isInWishlist,
	})
}

func (uh *UserHandler) GetMyWishlist(w http.ResponseWriter, r *http.Request) {
	
	userID := r.Context().Value("user_id").(int)

	// Call the UserService to get the user's wishlist
	wishlist, err := uh.UserService.GetWishlistByUserID(r.Context(), userID)
	if err != nil {
		// Handle the error, e.g., return a 500 Internal Server Error
		http.Error(w, "Failed to fetch user wishlist", http.StatusInternalServerError)
		return
	}

	// Send the wishlist as a JSON response
	sendSuccessResponse(w, map[string]interface{}{
		"wishlist": wishlist,
	})
}

func (uh *UserHandler) GetUserWishlist(w http.ResponseWriter, r *http.Request) {
	
	userIDStr := chi.URLParam(r, "userID")

    userID, err := strconv.Atoi(userIDStr) // Convert to int
	log.Println("userID: ", userID)
    if err != nil {
        sendErrorResponse(w, http.StatusBadRequest, "Invalid book ID")
        return
    }

	// Call the UserService to get the user's wishlist
	wishlist, err := uh.UserService.GetWishlistByUserID(r.Context(), userID)
	if err != nil {
		// Handle the error, e.g., return a 500 Internal Server Error
		http.Error(w, "Failed to fetch user wishlist", http.StatusInternalServerError)
		return
	}

	// Send the wishlist as a JSON response
	sendSuccessResponse(w, map[string]interface{}{
		"wishlist": wishlist,
	})
}

func (uh *UserHandler) GetListingWithBookByIDHandler(w http.ResponseWriter, r *http.Request) {
	
	listingIDStr := chi.URLParam(r, "listingID")

    listingID, err := strconv.Atoi(listingIDStr) // Convert to int
	log.Println("listingID: ", listingID)
    if err != nil {
        sendErrorResponse(w, http.StatusBadRequest, "Invalid listing ID")
        return
    }

	// Call the UserService to get the user's wishlist
	listing, err := uh.UserService.GetListingWithBookByID(r.Context(), listingID)
	if err != nil {
		// Handle the error, e.g., return a 500 Internal Server Error
		http.Error(w, "Failed to fetch listing", http.StatusInternalServerError)
		return
	}

	// Send the wishlist as a JSON response
	sendSuccessResponse(w, map[string]interface{}{
		"listing": listing,
	})
}

// func (uh *UserHandler) EditPhoneNumberHandler(w http.ResponseWriter, r *http.Request) {

// 	// 1. Parse JSON body for email, password, name, etc.
// 	var user models.User

// 	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
// 		sendErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
// 		return
// 	}

// 	user.Provider = "local"

// 	// 2. Check if user with same email already exists
// 	err := uh.UserService.EditPhoneNumber(r.Context(), user.ID, user.PhoneNumber)
// 	if err != nil {
// 		sendErrorResponse(w, http.StatusConflict, "Edit Phone Number: "+err.Error()) // 409 Conflict if user exists
// 		return
// 	}

// 	// Step 3: Success Response
// 	sendSuccessResponse(w, map[string]interface{}{
// 		"success": true,
// 		"message": "Edited info successfully!",
// 	})
// }

