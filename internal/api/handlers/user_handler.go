package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
	"strings"
	"used2book-backend/internal/models"
	"used2book-backend/internal/services"
	"github.com/streadway/amqp"
	"time"
)

type UserHandler struct {
	UserService   *services.UserService
	UploadService *services.UploadService
	RabbitMQConn *amqp.Connection
}

type CreatePaymentRequest struct {
	ListingID int `json:"listingId"`
	// Possibly other fields, e.g. quantity, shipping, etc.
}

func (uh *UserHandler) TestPort(w http.ResponseWriter, r *http.Request) {
	sendSuccessResponse(w, map[string]interface{}{
		"success": true,
	})
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

func (uh * UserHandler) CreateBankAccountHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok || userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse request body
	var request *models.BankAccount

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	request.UserID = userID

	// Call service to update preferences
	_, err = uh.UserService.CreateBankAccount(r.Context(), request)
	if err != nil {
		http.Error(w, "Failed to update preferences", http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(w, map[string]interface{}{
		"success": true,
	})

}

func (uh * UserHandler) CreateBookRequestHandle(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok || userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse request body
	var req *models.BookRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	req.UserID = userID

	// Call service to update preferences
	status, err := uh.UserService.CreateBookRequest(r.Context(), req)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to update preferences",)
		return
	}

	sendSuccessResponse(w, map[string]interface{}{
		"success": status,
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
		"success":           true,
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
		"success":              true,
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

func (uh *UserHandler) EditProfileHandler(w http.ResponseWriter, r *http.Request) {

	// 1. Parse JSON body for email, password, name, etc.
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	userID := r.Context().Value("user_id").(int)

	log.Println("first name:", user.FirstName)
	log.Println("last name:", user.LastName)
	log.Println("address:", user.Address)
	log.Println("quote:", user.Quote)
	log.Println("bio:", user.Bio)


	// 2. Check if user with same email already exists
	err := uh.UserService.EditProfile(r.Context(), userID, user.FirstName, user.LastName, user.Address, user.Quote, user.Bio)
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
	var req struct {
		BookID int `json:"book_id"`
		ReadingStatus int `json:"reading_status"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	log.Printf("Decoded request: %+v\n", req)

	userID, ok := r.Context().Value("user_id").(int)
    if !ok {
        sendErrorResponse(w, http.StatusUnauthorized, "User ID missing")
        return
    }

    _ , err := uh.UserService.AddBookToLibrary(r.Context(), userID, req.BookID, req.ReadingStatus)
    if err != nil {
        sendErrorResponse(w, http.StatusConflict, "Failed to process book: "+err.Error())
        return
    }


    sendSuccessResponse(w, map[string]interface{}{
        "success": true,
        "message": "add to library successfully",
    })
}

func (uh *UserHandler) AddBookToListingHandler(w http.ResponseWriter, r *http.Request) {
    // Parse multipart form (for JSON and files)
    err := r.ParseMultipartForm(10 << 20) // 10MB max
    if err != nil {
        log.Println("Parse Form Error:", err)
        sendErrorResponse(w, http.StatusBadRequest, "Invalid form data")
        return
    }

    // Get JSON data from form field 'data'
    jsonData := r.FormValue("data")
    var user models.UserAddListingForm
    if err := json.Unmarshal([]byte(jsonData), &user); err != nil {
        log.Println("JSON Decode Error:", err)
        sendErrorResponse(w, http.StatusBadRequest, "Invalid JSON format")
        return
    }

    userID, ok := r.Context().Value("user_id").(int)
    if !ok {
        sendErrorResponse(w, http.StatusUnauthorized, "User ID missing")
        return
    }

    files := r.MultipartForm.File["images"]

    var uploadURLs []string
    if len(files) > 0 { // Only process images if provided
        for _, handler := range files {
            file, err := handler.Open()
            if err != nil {
                sendErrorResponse(w, http.StatusInternalServerError, "Error opening file")
                return
            }
            defer file.Close()

            url, err := uh.UploadService.UploadImageURL(file, handler.Filename)
            if err != nil {
                sendErrorResponse(w, http.StatusInternalServerError, "Image upload failed: "+err.Error())
                return
            }
            uploadURLs = append(uploadURLs, url)
        }
    }

    _, err = uh.UserService.AddBookToListing(r.Context(), userID, user.BookID, user.Price, user.AllowOffer, uploadURLs, user.SellerNote, user.PhoneNumber)
    if err != nil {
        sendErrorResponse(w, http.StatusConflict, "Failed to process book: "+err.Error())
        return
    }

    sendSuccessResponse(w, map[string]interface{}{
        "success": true,
        "message": "Book listing added successfully!",
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

	in_wishlist, err := uh.UserService.AddBookToWishlist(r.Context(), userID, bookID)
	if err != nil {
		log.Println("❌ Wishlist Error:", err)
		sendErrorResponse(w, http.StatusConflict, "Wishlist error: "+err.Error())
		return
	}

	log.Println("✅ Wishlist updated successfully")
	sendSuccessResponse(w, map[string]interface{}{
		"success":     true,
		"message":     "Wishlist updated successfully!",
		"in_wishlist": in_wishlist,
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

func (uh *UserHandler) GetMyPurchasedListingsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	listing, err := uh.UserService.GetPurchasedListingsByUserID(r.Context(), userID)
	if err != nil {
		// Handle the error, e.g., return a 500 Internal Server Error
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to get my purchase listing: "+err.Error()) // 409 Conflict if user exists
		return
	}

	sendSuccessResponse(w, map[string]interface{}{
		"listing": listing,
	})

}

func (uh *UserHandler) GetMyOrdersHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	orders, err := uh.UserService.GetMyOrders(r.Context(), userID)
	if err != nil {
		// Handle the error, e.g., return a 500 Internal Server Error
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to get my purchase listing: "+err.Error()) // 409 Conflict if user exists
		return
	}

	sendSuccessResponse(w, map[string]interface{}{
		"orders": orders,
	})

}
func (uh *UserHandler) GetUsersByWishlistBookIDHandler(w http.ResponseWriter, r *http.Request) {
	bookIDStr := chi.URLParam(r, "bookID")
	bookID, err := strconv.Atoi(bookIDStr)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid book ID")
		return
	}

	users, err := uh.UserService.GetUsersByBookInWishlist(r.Context(), bookID)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to fetch wishlist users")
		return
	}

	sendSuccessResponse(w, map[string]interface{}{
		"users": users,
	})
}

func (uh *UserHandler) DeleteUserLibraryByIDHandler(w http.ResponseWriter, r *http.Request) {
	IDStr := chi.URLParam(r, "id")
	ID, err := strconv.Atoi(IDStr)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid book ID")
		return
	}

	delete_status, err := uh.UserService.DeleteUserLibraryByID(r.Context(), ID)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to delete book from library")
		return
	}

	sendSuccessResponse(w, map[string]interface{}{
		"success": delete_status,
	})
}



func (uh *UserHandler) GetAllListingsHandler(w http.ResponseWriter, r *http.Request) {

	listing, err := uh.UserService.GetAllListings(r.Context())
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

	wishlist, err := uh.UserService.GetWishlistByUserID(r.Context(), userID)
	if err != nil {
		// Handle the error, e.g., return a 500 Internal Server Error
		http.Error(w, "Failed to fetch user wishlist", http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(w, map[string]interface{}{
		"wishlist": wishlist,
	})
}

func (uh *UserHandler) GetListingByIDHandler(w http.ResponseWriter, r *http.Request) {

	listingIDStr := chi.URLParam(r, "listingID")

	listingID, err := strconv.Atoi(listingIDStr) // Convert to int
	log.Println("listingID: ", listingID)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid listing ID")
		return
	}

	listing, err := uh.UserService.GetListingByID(r.Context(), listingID)
	if err != nil {
		// Handle the error, e.g., return a 500 Internal Server Error
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to fetch listing")
		return
	}

	sendSuccessResponse(w, map[string]interface{}{
		"listing": listing,
	})
}

func (uh *UserHandler) GetCartHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	carts, err := uh.UserService.GetCart(r.Context(), userID)
	if err != nil {
		// Handle the error, e.g., return a 500 Internal Server Error
		http.Error(w, "Failed to fetch user cart", http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(w, map[string]interface{}{
		"carts": carts,
	})

}

func (uh *UserHandler) RemoveListingHandler(w http.ResponseWriter, r *http.Request) {
    userID, ok := r.Context().Value("user_id").(int)
    if !ok {
        sendErrorResponse(w, http.StatusUnauthorized, "User ID missing")
        return
    }

	listingIDStr := chi.URLParam(r, "listingID")

	listingID, err := strconv.Atoi(listingIDStr) // Convert to int
	log.Println("listingID: ", listingID)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid listing ID")
		return
	}

    err = uh.UserService.RemoveListing(r.Context(), userID, listingID)
    if err != nil {
        sendErrorResponse(w, http.StatusBadRequest, err.Error())
        return
    }

    sendSuccessResponse(w, map[string]interface{}{
        "success": true,
        "message": "Listing removed successfully",
    })
}

func (uh *UserHandler) AddToCartHandler(w http.ResponseWriter, r *http.Request) {

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

	_, err := uh.UserService.AddToCart(r.Context(), userID, req.ListingId)
	if err != nil {
		log.Println("❌ Add listing to cart Error:", err)
		sendErrorResponse(w, http.StatusConflict, "cart error: "+err.Error())
		return
	}

	log.Println("✅ Add listing to cart successfully")
	sendSuccessResponse(w, map[string]interface{}{
		"success": true,
		"message": "addded cart successfully!",
	})

}

func (uh *UserHandler) RemoveFromCartHandler(w http.ResponseWriter, r *http.Request) {
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

	err := uh.UserService.RemoveFromCart(r.Context(), userID, req.ListingId)
	if err != nil {
		log.Println("❌ Wishlist Error:", err)
		sendErrorResponse(w, http.StatusConflict, "Wishlist error: "+err.Error())
		return
	}

	log.Println("✅ Wishlist updated successfully")
	sendSuccessResponse(w, map[string]interface{}{
		"success": true,
		"message": "delete listing in cart successfully!",
	})

}

func (uh *UserHandler) AddToOffersHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        ListingID    int     `json:"listingId"`
        OfferedPrice float64 `json:"offeredPrice"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    buyerID, ok := r.Context().Value("user_id").(int)
    if !ok {
        sendErrorResponse(w, http.StatusUnauthorized, "User ID missing")
        return
    }

    log.Println("listingID:", req.ListingID, "offeredPrice:", req.OfferedPrice)

    id, err := uh.UserService.AddToOffers(r.Context(), buyerID, req.ListingID, req.OfferedPrice)
    if err != nil {
        log.Println("❌ Add offer error:", err)
        sendErrorResponse(w, http.StatusConflict, "Offer error: "+err.Error())
        return
    }

	listing, err := uh.UserService.GetListingByID(r.Context(), req.ListingID)
	if err != nil {
		// Handle the error, e.g., return a 500 Internal Server Error
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to fetch listing")
		return
	}

	// Publish to RabbitMQ
	ch, err := uh.RabbitMQConn.Channel()
	if err != nil {
		log.Println("❌ RabbitMQ Channel Error:", err)
		return // Don’t fail webhook response
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"offer_queue", // New queue for offer_queue
		true, false, false, false, nil,
	)
	if err != nil {
		log.Println("❌ Queue Declare Error:", err)
		return
	}

	noti := map[string]interface{}{
		"user_id":   int(listing.SellerID),
		"type":       "offer",
		"created_at": time.Now(),
	}

	body, _ := json.Marshal(noti)

	err = ch.Publish(
		"", q.Name, false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Println("❌ Publish Error:", err)
	}


    log.Println("✅ Added offer successfully with ID:", id)
    sendSuccessResponse(w, map[string]interface{}{
        "success": true,
        "message": "Offer added successfully!",
        "offerId": id,
    })
}

func (uh *UserHandler) GetBuyerOffersHandler(w http.ResponseWriter, r *http.Request) {
    buyerID := r.Context().Value("user_id").(int)

    offers, err := uh.UserService.GetBuyerOffers(r.Context(), buyerID)
    if err != nil {
        http.Error(w, "Failed to fetch user offers", http.StatusInternalServerError)
        return
    }

    sendSuccessResponse(w, map[string]interface{}{
        "offers": offers,
    })
}

func (uh *UserHandler) GetSellerOffersHandler(w http.ResponseWriter, r *http.Request) {
    sellerID, ok := r.Context().Value("user_id").(int)
    if !ok {
        sendErrorResponse(w, http.StatusUnauthorized, "User ID missing")
        return
    }

    offers, err := uh.UserService.GetSellerOffers(r.Context(), sellerID)
    if err != nil {
        http.Error(w, "Failed to fetch seller offers", http.StatusInternalServerError)
        return
    }

    sendSuccessResponse(w, map[string]interface{}{
        "offers": offers,
    })
}

func (uh *UserHandler) RemoveFromOffersHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        ListingID int `json:"listingId"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    buyerID, ok := r.Context().Value("user_id").(int)
    if !ok {
        sendErrorResponse(w, http.StatusUnauthorized, "User ID missing")
        return
    }

    err := uh.UserService.RemoveFromOffers(r.Context(), buyerID, req.ListingID)
    if err != nil {
        log.Println("❌ Remove offer error:", err)
        sendErrorResponse(w, http.StatusConflict, "Offer error: "+err.Error())
        return
    }

    log.Println("✅ Offer removed successfully")
    sendSuccessResponse(w, map[string]interface{}{
        "success": true,
        "message": "Offer removed successfully!",
    })
}

func (uh *UserHandler) AcceptOfferHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        OfferID int `json:"offerId"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    sellerID, ok := r.Context().Value("user_id").(int)
    if !ok {
        sendErrorResponse(w, http.StatusUnauthorized, "User ID missing")
        return
    }

    buyer_id, err := uh.UserService.AcceptOffer(r.Context(), sellerID, req.OfferID)
    if err != nil {
        log.Println("❌ Accept offer error:", err)
        sendErrorResponse(w, http.StatusConflict, "Offer error: "+err.Error())
        return
    }

	ch, err := uh.RabbitMQConn.Channel()
	if err != nil {
		log.Println("❌ RabbitMQ Channel Error:", err)
		return // Don’t fail webhook response
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"offer_queue", // New queue for offer_queue
		true, false, false, false, nil,
	)
	if err != nil {
		log.Println("❌ Queue Declare Error:", err)
		return
	}

	noti := map[string]interface{}{
		"user_id":  buyer_id,
		"type":      "offer",
		"created_at": time.Now(),
	}

	body, _ := json.Marshal(noti)

	err = ch.Publish(
		"", q.Name, false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Println("❌ Publish Error:", err)
	}



    sendSuccessResponse(w, map[string]interface{}{
        "success": true,
        "message": "Offer accepted successfully!",
    })
}

func (uh *UserHandler) RejectOfferHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        OfferID int `json:"offerId"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    sellerID, ok := r.Context().Value("user_id").(int)
    if !ok {
        sendErrorResponse(w, http.StatusUnauthorized, "User ID missing")
        return
    }

    buyer_id, err := uh.UserService.RejectOffer(r.Context(), sellerID, req.OfferID)
    if err != nil {
        log.Println("❌ Reject offer error:", err)
        sendErrorResponse(w, http.StatusConflict, "Offer error: "+err.Error())
        return
    }

	ch, err := uh.RabbitMQConn.Channel()
	if err != nil {
		log.Println("❌ RabbitMQ Channel Error:", err)
		return // Don’t fail webhook response
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"offer_queue", // New queue for offer_queue
		true, false, false, false, nil,
	)
	if err != nil {
		log.Println("❌ Queue Declare Error:", err)
		return
	}

	noti := map[string]interface{}{
		"user_id":  buyer_id,
		"type":      "offer",
		"created_at": time.Now(),
	}

	body, _ := json.Marshal(noti)

	err = ch.Publish(
		"", q.Name, false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Println("❌ Publish Error:", err)
	}

    sendSuccessResponse(w, map[string]interface{}{
        "success": true,
        "message": "Offer rejected successfully!",
    })
}

// handler/user_handler.go
func (uh *UserHandler) GetAcceptedOfferHandler(w http.ResponseWriter, r *http.Request) {
    // offerIDStr := r.URL.Path[len("/offers/") : strings.Index(r.URL.Path, "/payment")]
    // offerID, err := strconv.Atoi(offerIDStr)
    // if err != nil {
    //     sendErrorResponse(w, http.StatusBadRequest, "Invalid offer ID")
    //     return
    // }
	offerIDStr := chi.URLParam(r, "offerID")

	offerID, err := strconv.Atoi(offerIDStr) // Convert to int
	log.Println("offerID: ", offerID)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid listing ID")
		return
	}

    buyerID, ok := r.Context().Value("user_id").(int)
    if !ok {
        sendErrorResponse(w, http.StatusUnauthorized, "User ID missing")
        return
    }

    offer, err := uh.UserService.GetAcceptedOffer(r.Context(), offerID)
    if err != nil {
        sendErrorResponse(w, http.StatusNotFound, "Offer not found or not accepted")
        return
    }

    if offer.BuyerID != buyerID {
        sendErrorResponse(w, http.StatusForbidden, "You are not the buyer of this offer")
        return
    }

    sendSuccessResponse(w, map[string]interface{}{
        "offer": offer,
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

		url, err := uh.UploadService.UploadImageURL(file, handler.Filename) // Reuse existing service
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
	err := uh.UserService.MarkListingAsSold(r.Context(), request.ListingID, buyerID)
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
	reviews, err := uh.UserService.GetAllUserReview(r.Context())
	if err != nil {
		log.Printf("Failed to fetch user reviews: %v", err) // Log the error
		http.Error(w, "Failed to fetch user review", http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(w, map[string]interface{}{
		"reviews": reviews,
	})
}

func (uh *UserHandler) GetAllUserPreferred(w http.ResponseWriter, r *http.Request) {
	user_preferred_genres, err := uh.UserService.GetAllUserPreferred(r.Context())
	if err != nil {
		log.Printf("Failed to fetch user user_preferred_genres: %v", err) // Log the error
		http.Error(w, "Failed to fetch user user_preferred_genres", http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(w, map[string]interface{}{
		"user_preferred_genres": user_preferred_genres,
	})
}

// handler/user_handler.go
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

    // Get genre_id and book_id (optional)
    var genreID *int
    var bookID *int
    if genreStr := r.FormValue("genre_id"); genreStr != "" {
        id, err := strconv.Atoi(genreStr)
        if err != nil {
            sendErrorResponse(w, http.StatusBadRequest, "Invalid genre_id")
            return
        }
        genreID = &id
    }
    if bookStr := r.FormValue("book_id"); bookStr != "" {
        id, err := strconv.Atoi(bookStr)
        if err != nil {
            sendErrorResponse(w, http.StatusBadRequest, "Invalid book_id")
            return
        }
        bookID = &id
    }

    // Enforce 1-to-1 relationship: only one of genre_id or book_id can be non-null
    if genreID != nil && bookID != nil {
        sendErrorResponse(w, http.StatusBadRequest, "Post can only reference either a genre OR a book, not both")
        return
    }

    // Get userID from context
    userID, ok := r.Context().Value("user_id").(int)
    if !ok {
        sendErrorResponse(w, http.StatusUnauthorized, "User ID missing")
        return
    }

    // Create post
    post, err := uh.UserService.CreatePost(r.Context(), userID, content, imageURLs, genreID, bookID)
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
// handler/user_handler.go
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

// ToggleLikeHandler toggles a like (like/unlike) on a post
func (uh *UserHandler) ToggleLikeHandler(w http.ResponseWriter, r *http.Request) {
	postIDStr := chi.URLParam(r, "postID")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil || postID <= 0 {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid post ID")
		return
	}
	if postID <= 0 {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		sendErrorResponse(w, http.StatusUnauthorized, "User ID missing")
		return
	}

	isLiked, err := uh.UserService.IsPostLikedByUser(r.Context(), postID, userID)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to check like status: "+err.Error())
		return
	}

	if isLiked {
		err = uh.UserService.RemoveLike(r.Context(), postID, userID)
		if err != nil {
			sendErrorResponse(w, http.StatusInternalServerError, "Failed to unlike: "+err.Error())
			return
		}
	} else {
		_, err = uh.UserService.CreateLike(r.Context(), postID, userID)
		if err != nil {
			sendErrorResponse(w, http.StatusInternalServerError, "Failed to like: "+err.Error())
			return
		}
	}

	sendSuccessResponse(w, map[string]interface{}{
		"success":  true,
		"is_liked": !isLiked,
	})
}

func (uh *UserHandler) IsPostLikedHandler(w http.ResponseWriter, r *http.Request) {
	postIDStr := chi.URLParam(r, "postID")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil || postID <= 0 {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		sendErrorResponse(w, http.StatusUnauthorized, "User ID missing")
		return
	}

	isLiked, err := uh.UserService.IsPostLikedByUser(r.Context(), postID, userID)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to check like status: "+err.Error())
		return
	}

	sendSuccessResponse(w, map[string]interface{}{
		"success":  true,
		"is_liked": isLiked,
	})
}

// GetLikeCountHandler returns the number of likes for a post
func (uh *UserHandler) GetLikeCountHandler(w http.ResponseWriter, r *http.Request) {
	// Get post_id from query parameter
	postIDStr := chi.URLParam(r, "postID")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil || postID <= 0 {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	// Fetch like count from service
	likeCount, err := uh.UserService.GetLikeCountByPostID(r.Context(), postID)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to get like count: "+err.Error())
		return
	}

	// Send success response
	sendSuccessResponse(w, map[string]interface{}{
		"success":    true,
		"like_count": likeCount,
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
