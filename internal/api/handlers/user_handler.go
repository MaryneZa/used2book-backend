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
	"strings"
)

type UserHandler struct {
	UserService  *services.UserService
	UploadService *services.UploadService
}

type CreatePaymentRequest struct {
    ListingID int `json:"listingId"`
    // Possibly other fields, e.g. quantity, shipping, etc.
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

func (uh *UserHandler) SetUserPreferredGenresHandler(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context
	userID, ok := r.Context().Value("user_id").(int)
	if !ok || userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse request body
	var request struct {
		GenreIDs []int `json:"genre_ids"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil || len(request.GenreIDs) == 0 {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Call service to update preferences
	err = uh.UserService.SetUserPreferredGenres(r.Context(), userID, request.GenreIDs)
	if err != nil {
		http.Error(w, "Failed to update preferences", http.StatusInternalServerError)
		return
	}

	// Send success response
	sendSuccessResponse(w, map[string]interface{}{
		"message": "User preferred genres updated successfully",
	})
}

func (uh *UserHandler) GetUserPreferencesHandler(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by authentication middleware)
	userID, ok := r.Context().Value("user_id").(int)
	if !ok || userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	preferredGenres, err := uh.UserService.GetUserPreferences(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to retrieve user preferences", http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(w, map[string]interface{}{
		"preferred_genres": preferredGenres,
	})
}

func (uh *UserHandler) GetGenderHandler(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by authentication middleware)
	userID, ok := r.Context().Value("user_id").(int)
	if !ok || userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	gender, err := uh.UserService.GetGender(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to retrieve user gender", http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(w, map[string]interface{}{
		"gender": gender,
	})
}

func (uh *UserHandler) UpdateGenderHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)

	if !ok || userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

    var req struct {
        Gender string `json:"gender"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    if req.Gender != "male" && req.Gender != "female" && req.Gender != "other" {
        http.Error(w, "Invalid gender value", http.StatusBadRequest)
        return
    }

    err := uh.UserService.UpdateGender(r.Context(), userID, req.Gender)
    if err != nil {
        http.Error(w, "Failed to update gender", http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Gender updated successfully"})
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

func (uh *UserHandler) GetCartHandler(w http.ResponseWriter, r *http.Request){
	userID := r.Context().Value("user_id").(int)

	// Call the UserService to get the user's wishlist
	carts, err := uh.UserService.GetCart(r.Context(), userID)
	if err != nil {
		// Handle the error, e.g., return a 500 Internal Server Error
		http.Error(w, "Failed to fetch user cart", http.StatusInternalServerError)
		return
	}

	// Send the wishlist as a JSON response
	sendSuccessResponse(w, map[string]interface{}{
		"carts": carts,
	})

}

func (uh *UserHandler) AddToCartHandler(w http.ResponseWriter, r *http.Request){

	var req struct {
		ListingId int `json:"listingId"`
	}
	log.Printf("Decoded request: %+v\n", req)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		sendErrorResponse(w, http.StatusUnauthorized, "User ID missing")
		return
	}

	log.Println("listingID:", req.ListingId)

	// ✅ Call service to toggle wishlist and get status
	_, err := uh.UserService.AddToCart(r.Context(), userID, req.ListingId)
	if err != nil {
		log.Println("❌ Add listing to cart Error:", err)
		sendErrorResponse(w, http.StatusConflict, "cart error: "+err.Error()) 
		return
	}

	// ✅ Return updated wishlist status
	log.Println("✅ Add listing to cart successfully")
	sendSuccessResponse(w, map[string]interface{}{
		"success": true,
		"message": "addded cart successfully!",
	})

}

func (uh *UserHandler) RemoveFromCartHandler(w http.ResponseWriter, r *http.Request){
	var req struct {
		ListingId int `json:"listingId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		sendErrorResponse(w, http.StatusUnauthorized, "User ID missing")
		return
	}

	// ✅ Call service to toggle wishlist and get status
	err := uh.UserService.RemoveFromCart(r.Context(), userID, req.ListingId)
	if err != nil {
		log.Println("❌ Wishlist Error:", err)
		sendErrorResponse(w, http.StatusConflict, "Wishlist error: "+err.Error()) 
		return
	}

	// ✅ Return updated wishlist status
	log.Println("✅ Wishlist updated successfully")
	sendSuccessResponse(w, map[string]interface{}{
		"success": true,
		"message": "delete listing in cart successfully!",
	})

}

func (uh *UserHandler) UploadPostImagesHandler(w http.ResponseWriter, r *http.Request) {
    err := r.ParseMultipartForm(10 << 20) // 10MB max
    if err != nil {
        sendErrorResponse(w, http.StatusBadRequest, "File too large")
        return
    }

    files := r.MultipartForm.File["images"] // Get all files under "images" key
    if len(files) == 0 {
        sendErrorResponse(w, http.StatusBadRequest, "No images provided")
        return
    }

    var uploadURLs []string
    for _, handler := range files {
        file, err := handler.Open()
        if err != nil {
            sendErrorResponse(w, http.StatusInternalServerError, "Error opening file")
            return
        }
        defer file.Close()

        url, err := uh.UploadService.UploadPostImageURL(file, handler.Filename) // Reuse existing service
        if err != nil {
            sendErrorResponse(w, http.StatusInternalServerError, "Image upload failed: "+err.Error())
            return
        }
        uploadURLs = append(uploadURLs, url)
    }

    sendSuccessResponse(w, map[string]interface{}{
        "success":    true,
        "image_urls": uploadURLs,
    })
}




// // CreateExpressAccountHandler creates a Stripe Connect Express account for the user
// func (uh *UserHandler) CreateExpressAccountHandler(w http.ResponseWriter, r *http.Request) {
//     // Typically, you get userID from JWT or context
//     userID, ok := r.Context().Value("user_id").(int)
//     if !ok || userID == 0 {
//         sendErrorResponse(w, http.StatusUnauthorized, "Unauthorized (missing user_id)")
//         return
//     }

//     // 1) Get user from DB (we need to confirm they don't already have an account)
//     user, err := uh.UserService.GetMe(r.Context(), userID)
//     if err != nil {
//         sendErrorResponse(w, http.StatusNotFound, "User not found")
//         return
//     }

//     // If user already has a stripe_account_id, skip creation (or create again if you prefer)
//     if user.StripeAccountID.String != "" {
//         // Possibly just return the existing account link or an error
//         sendErrorResponse(w, http.StatusConflict, "User already has a Stripe account")
//         return
//     }

//     // 2) Create a brand-new Express account
//     accParams := &stripe.AccountParams{
//         Type: stripe.String(string(stripe.AccountTypeExpress)),
//         // If you know the user's country, set it:
//         // Country: stripe.String("TH"), // e.g. for Thailand
//     }
//     acc, err := account.New(accParams)
//     if err != nil {
//         log.Println("Error creating stripe account:", err)
//         sendErrorResponse(w, http.StatusInternalServerError, "Failed to create stripe account")
//         return
//     }

//     // 3) Store the new acc.ID in DB
//     //    e.g. update user with user.StripeAccountID = acc.ID
//     err = uh.UserService.UpdateStripeAccountID(r.Context(), userID, acc.ID)
//     if err != nil {
//         log.Println("Error saving stripe account ID to DB:", err)
//         sendErrorResponse(w, http.StatusInternalServerError, "Failed to store stripe account ID")
//         return
//     }

//     // 4) Generate an Account Link so the user can fill out bank info, etc.
//     linkParams := &stripe.AccountLinkParams{
//         Account:    stripe.String(acc.ID),
//         RefreshURL: stripe.String("http://localhost:3000/seller-onboarding?refresh=1"),
//         ReturnURL:  stripe.String("http://localhost:3000/seller-onboarding-complete"),
//         Type:       stripe.String("account_onboarding"),
//     }
//     link, err := accountlink.New(linkParams)
//     if err != nil {
//         log.Println("Error creating account link:", err)
//         sendErrorResponse(w, http.StatusInternalServerError, "Failed to create account link")
//         return
//     }

//     // 5) Return the onboarding link to the frontend
//     resp := map[string]interface{}{
//         "url":   link.URL,
//         "accID": acc.ID, // optional
//     }
//     sendSuccessResponse(w, resp)
// }

// // CreatePaymentIntentHandler handles buyer purchase requests
// func (uh *UserHandler) CreatePaymentIntentHandler(w http.ResponseWriter, r *http.Request) {
//     // 1) Parse the JSON body to get the listingId
//     body, err := io.ReadAll(r.Body)
//     if err != nil {
//         sendErrorResponse(w, http.StatusBadRequest, "Invalid body")
//         return
//     }
//     var req CreatePaymentRequest
//     if err := json.Unmarshal(body, &req); err != nil {
//         sendErrorResponse(w, http.StatusBadRequest, "JSON decode error")
//         return
//     }

//     // 2) Look up the listing in DB
//     listing, err := uh.UserService.GetListingByID(r.Context(), req.ListingID)
//     if err != nil || listing == nil {
//         sendErrorResponse(w, http.StatusNotFound, "Listing not found")
//         return
//     }

//     // 3) Find the seller's stripe_account_id
//     sellerUser, err := uh.UserService.GetMe(r.Context(), listing.SellerID)
//     if err != nil {
//         sendErrorResponse(w, http.StatusNotFound, "Seller not found")
//         return
//     }
//     if sellerUser.StripeAccountID.String == "" {
//         sendErrorResponse(w, http.StatusBadRequest, "Seller has no stripe account")
//         return
//     }

//     sellerAccID := sellerUser.StripeAccountID.String

//     // Convert listing.Price to integer in cents
//     // e.g. if listing.Price = 12.99 => 1299
//     // If using THB, also multiply by 100 for satang.
//     amount := int64(listing.Price)

//     // 4) Create the PaymentIntent
//     params := &stripe.PaymentIntentParams{
//         Amount:   stripe.Int64(amount),
//         Currency: stripe.String("thb"), // or "thb" if in Thailand
//         OnBehalfOf: stripe.String(sellerAccID),
//         TransferData: &stripe.PaymentIntentTransferDataParams{
//             Destination: stripe.String(sellerAccID),
//         },
//         // If you want to take a 10% platform fee:
//         ApplicationFeeAmount: stripe.Int64(amount / 10),
//         PaymentMethodTypes:    stripe.StringSlice([]string{"card"}), 
//     }

//     pi, err := paymentintent.New(params)
//     if err != nil {
//         log.Println("Error creating payment intent:", err)
//         sendErrorResponse(w, http.StatusInternalServerError, "Failed to create payment intent")
//         return
//     }

//     // 5) Return the client secret to the front end
//     resp := map[string]interface{}{
//         "clientSecret": pi.ClientSecret,
//     }
//     sendSuccessResponse(w, resp)
// }


func (uh *UserHandler) MarkListingAsSoldHandler(w http.ResponseWriter, r *http.Request) {
    // Get buyer ID from JWT context
    buyerID, ok := r.Context().Value("user_id").(int)
    if !ok || buyerID == 0 {
        sendErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
        return
    }

    // Parse request body
    var request struct {
        ListingID int     `json:"listing_id"`
        Amount    float64 `json:"amount"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        sendErrorResponse(w, http.StatusBadRequest, "Invalid request format")
        return
    }

    // Call service to mark listing as sold
    err := uh.UserService.MarkListingAsSold(r.Context(), request.ListingID, buyerID, request.Amount)
    if err != nil {
        sendErrorResponse(w, http.StatusInternalServerError, "Failed to mark listing as sold: "+err.Error())
        return
    }

    sendSuccessResponse(w, map[string]interface{}{
        "success": true,
        "message": "Listing marked as sold successfully!",
    })
}


func (uh *UserHandler) GetAllUserReview(w http.ResponseWriter, r *http.Request) {
	// Call the UserService to get the user's wishlist
	reviews, err := uh.UserService.GetAllUserReview(r.Context())
	if err != nil {
		log.Printf("Failed to fetch user reviews: %v", err) // Log the error
		http.Error(w, "Failed to fetch user review", http.StatusInternalServerError)
		return
	}

	// Send the wishlist as a JSON response
	sendSuccessResponse(w, map[string]interface{}{
		"reviews": reviews,
	})
}


func (uh *UserHandler) GetAllUserPreferred(w http.ResponseWriter, r *http.Request) {
	// Call the UserService to get the user's wishlist
	user_preferred_genres, err := uh.UserService.GetAllUserPreferred(r.Context())
	if err != nil {
		log.Printf("Failed to fetch user user_preferred_genres: %v", err) // Log the error
		http.Error(w, "Failed to fetch user user_preferred_genres", http.StatusInternalServerError)
		return
	}

	// Send the wishlist as a JSON response
	sendSuccessResponse(w, map[string]interface{}{
		"user_preferred_genres": user_preferred_genres,
	})
}

func (uh *UserHandler) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
    // Parse multipart form (10MB max)
    err := r.ParseMultipartForm(10 << 20)
    if err != nil {
        log.Println("Parse Form Error:", err)
        sendErrorResponse(w, http.StatusBadRequest, "Invalid form data")
        return
    }

    // Log the form data for debugging
    log.Println("Form Data:", r.Form)

    // Get content
    content := r.FormValue("content")
    if content == "" {
        sendErrorResponse(w, http.StatusBadRequest, "Content is required")
        return
    }

    // Get image URLs (optional, from form field)
    imageURLs := r.Form["image_urls"] // Multiple values if sent as array

    // Get userID from context
    userID, ok := r.Context().Value("user_id").(int)
    if !ok {
        sendErrorResponse(w, http.StatusUnauthorized, "User ID missing")
        return
    }

    // Create post
    post, err := uh.UserService.CreatePost(r.Context(), userID, content, imageURLs)
    if err != nil {
        sendErrorResponse(w, http.StatusInternalServerError, "Failed to create post: "+err.Error())
        return
    }

    // Send success response with the created post
    sendSuccessResponse(w, map[string]interface{}{
        "success": true,
        "message": "Post created successfully!",
        "post":    post,
    })
}

// GetAllPostsHandler returns all posts
func (uh *UserHandler) GetAllPostsHandler(w http.ResponseWriter, r *http.Request) {
    posts, err := uh.UserService.GetAllPosts(r.Context())
    if err != nil {
        sendErrorResponse(w, http.StatusInternalServerError, "Failed to fetch posts: "+err.Error())
        return
    }

    sendSuccessResponse(w, map[string]interface{}{
        "success": true,
        "posts":   posts,
    })
}

// GetPostsByUserIDHandler returns posts for a specific user
func (uh *UserHandler) GetPostsByUserIDHandler(w http.ResponseWriter, r *http.Request) {

	userIDStr := chi.URLParam(r, "userID")

    userID, err := strconv.Atoi(userIDStr) // Convert to int
	log.Println("userID: ", userID)
    if err != nil {
        sendErrorResponse(w, http.StatusBadRequest, "Invalid user ID")
        return
    }

    posts, err := uh.UserService.GetPostsByUserID(r.Context(), userID)
    if err != nil {
        sendErrorResponse(w, http.StatusInternalServerError, "Failed to fetch posts: "+err.Error())
        return
    }

    sendSuccessResponse(w, map[string]interface{}{
        "success": true,
        "posts":   posts,
    })
}

// GetPostByPostIDHandler returns a single post by its ID
func (uh *UserHandler) GetPostByPostIDHandler(w http.ResponseWriter, r *http.Request) {
    postIDStr := r.URL.Query().Get("post_id") // Assume passed as query param
    postID, err := strconv.Atoi(postIDStr)
    if err != nil || postID <= 0 {
        sendErrorResponse(w, http.StatusBadRequest, "Invalid post ID")
        return
    }

    post, err := uh.UserService.GetPostByPostID(r.Context(), postID)
    if err != nil {
        if strings.Contains(err.Error(), "not found") {
            sendErrorResponse(w, http.StatusNotFound, err.Error())
        } else {
            sendErrorResponse(w, http.StatusInternalServerError, "Failed to fetch post: "+err.Error())
        }
        return
    }

    sendSuccessResponse(w, map[string]interface{}{
        "success": true,
        "post":    post,
    })
}

// CreateCommentHandler handles comment creation
func (uh *UserHandler) CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
    err := r.ParseMultipartForm(10 << 20)
    if err != nil {
        sendErrorResponse(w, http.StatusBadRequest, "Invalid form data")
        return
    }

    postIDStr := r.FormValue("post_id")
    content := r.FormValue("content")
    if postIDStr == "" || content == "" {
        sendErrorResponse(w, http.StatusBadRequest, "Post ID and content are required")
        return
    }

    postID, err := strconv.Atoi(postIDStr)
    if err != nil {
        sendErrorResponse(w, http.StatusBadRequest, "Invalid post ID")
        return
    }

    userID, ok := r.Context().Value("user_id").(int)
    if !ok {
        sendErrorResponse(w, http.StatusUnauthorized, "User ID missing")
        return
    }

    comment, err := uh.UserService.CreateComment(r.Context(), postID, userID, content)
    if err != nil {
        sendErrorResponse(w, http.StatusInternalServerError, "Failed to create comment: "+err.Error())
        return
    }

    sendSuccessResponse(w, map[string]interface{}{
        "success": true,
        "comment": comment,
    })
}

// GetCommentsByPostIDHandler fetches comments for a post
func (uh *UserHandler) GetCommentsByPostIDHandler(w http.ResponseWriter, r *http.Request) {

    postIDStr := chi.URLParam(r, "postID")
    postID, err := strconv.Atoi(postIDStr)
    if err != nil || postID <= 0 {
        sendErrorResponse(w, http.StatusBadRequest, "Invalid post ID")
        return
    }

    comments, err := uh.UserService.GetCommentsByPostID(r.Context(), postID)
    if err != nil {
        sendErrorResponse(w, http.StatusInternalServerError, "Failed to fetch comments: "+err.Error())
        return
    }

    sendSuccessResponse(w, map[string]interface{}{
        "success":  true,
        "comments": comments,
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

