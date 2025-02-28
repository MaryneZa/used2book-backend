package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"used2book-backend/internal/models"
	"used2book-backend/internal/services"

	"io"
	"log"
	"strconv"

	"github.com/go-chi/chi/v5"
	// "github.com/gorilla/mux"
)

// BookHandler handles book-related HTTP requests
type BookHandler struct {
	BookService *services.BookService
	UserService *services.UserService
}

func (bh *BookHandler) GetAllBooks(w http.ResponseWriter, r *http.Request) {

	// Call the BookService method to get the total book count
	books, err := bh.BookService.GetAllBooks(r.Context())
	if err != nil {
		// Handle the error, e.g., return a 500 Internal Server Error
		http.Error(w, "Failed to get all books", http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(w, map[string]interface{}{
		"books": books,
	})
	
}

func (bh *BookHandler) GetReviewsByBookIDHandler(w http.ResponseWriter, r *http.Request) {
	
	body, _ := io.ReadAll(r.Body)
	log.Println("Request Body:", string(body)) 

	bookIDStr := chi.URLParam(r, "id")
	log.Println("bookIDStr:", bookIDStr)

	if bookIDStr == "" {
		sendErrorResponse(w, http.StatusBadRequest, "Book ID is required")
		return
	}

	bookID, err := strconv.Atoi(bookIDStr)

	// Call the BookService method to get the total book count
	reviews, err := bh.BookService.GetReviewsByBookID(r.Context(), bookID)
	if err != nil {
		// Handle the error, e.g., return a 500 Internal Server Error
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to get reviews"+err.Error())
		return
	}

	sendSuccessResponse(w, map[string]interface{}{
		"reviews": reviews,
	})
	
}

// // GetBookWithRatings retrieves a book and its ratings
// func (bh *BookHandler) GetBookWithRatings(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	bookID, err := strconv.Atoi(vars["id"])
// 	if err != nil {
// 		http.Error(w, "Invalid book ID", http.StatusBadRequest)
// 		return
// 	}

// 	book, err := bh.BookService.GetBookWithRatings(context.Background(), bookID)
// 	if err != nil {
// 		http.Error(w, "Error fetching book", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(book)
// }

func (bh *BookHandler) GetBookByID(w http.ResponseWriter, r *http.Request) {
	// Use chi's URLParam to get the 'id' parameter
	bookIDStr := chi.URLParam(r, "id")
	log.Println("bookIDStr:", bookIDStr)

	if bookIDStr == "" {
		sendErrorResponse(w, http.StatusBadRequest, "Book ID is required")
		return
	}

	bookID, err := strconv.Atoi(bookIDStr)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid book ID: "+err.Error())
		return
	}

	book, err := bh.BookService.GetBookByID(context.Background(), bookID)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Error fetching book: "+err.Error())
		return
	}

	sendSuccessResponse(w, map[string]interface{}{
		"book": book,
	})
}

func (bh *BookHandler) AddBookReviewHandler(w http.ResponseWriter, r *http.Request) {
	// ✅ Read request body only once
	body, err := io.ReadAll(r.Body)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Failed to read request body")
		return
	}

	// ✅ Log request body for debugging
	log.Println("Request Body:", string(body))

	// ✅ Parse JSON body
	var review models.AddBookReview
	if err := json.Unmarshal(body, &review); err != nil {
		log.Println("Error decoding JSON:", err)
		sendErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// ✅ Get user ID from context
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		sendErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// ✅ Call the service layer to save the review
	err = bh.BookService.AddBookReview(context.Background(), userID, review.BookID, review.Rating, review.Comment)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Error saving review: "+err.Error())
		return
	}

	// ✅ Send success response
	sendSuccessResponse(w, map[string]interface{}{
		"success": true,
	})
}


func (bh *BookHandler) GetGenresByBookID(w http.ResponseWriter, r *http.Request) {
	// Use chi's URLParam to get the 'id' parameter
	bookIDStr := chi.URLParam(r, "id")
	log.Println("bookIDStr:", bookIDStr)

	if bookIDStr == "" {
		sendErrorResponse(w, http.StatusBadRequest, "Book ID is required")
		return
	}

	bookID, err := strconv.Atoi(bookIDStr)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid book ID genres: "+err.Error())
		return
	}

	book_genres, err := bh.BookService.GetGenresByBookID(context.Background(), bookID)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Error fetching book genres: "+err.Error())
		return
	}

	sendSuccessResponse(w, map[string]interface{}{
		"genres": book_genres,
	})
}



// // AddOrUpdateUserRating allows a user to add or update their rating
// func (bh *BookHandler) AddOrUpdateUserRating(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	bookID, err := strconv.Atoi(vars["id"])
// 	if err != nil {
// 		http.Error(w, "Invalid book ID", http.StatusBadRequest)
// 		return
// 	}

// 	userID := r.Context().Value("user_id").(int) // Assume User-ID is passed in header
// 	var payload struct {
// 		Rating float64 `json:"rating"`
// 	}

// 	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	err = bh.BookService.AddOrUpdateUserRating(context.Background(), userID, bookID, payload.Rating)
// 	if err != nil {
// 		http.Error(w, "Error updating rating", http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte(`{"message": "Rating updated successfully"}`))
// }

// // DeleteUserRating allows a user to delete their rating
// func (bh *BookHandler) DeleteUserRating(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	bookID, err := strconv.Atoi(vars["id"])
// 	if err != nil {
// 		http.Error(w, "Invalid book ID", http.StatusBadRequest)
// 		return
// 	}

// 	userID := r.Context().Value("user_id").(int) // Assume User-ID is passed in header

// 	err = bh.BookService.DeleteUserRating(context.Background(), userID, bookID)
// 	if err != nil {
// 		http.Error(w, "Error deleting rating", http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte(`{"message": "Rating deleted successfully"}`))
// }

// SyncBooksFromGoogleSheets triggers syncing books from Google Sheets
func (bh *BookHandler) SyncBooksFromGoogleSheets(w http.ResponseWriter, r *http.Request) {
	sheetID := r.URL.Query().Get("sheet_id")
	apiKey := r.URL.Query().Get("api_key")

	if sheetID == "" || apiKey == "" {
		http.Error(w, "Missing sheet_id or api_key", http.StatusBadRequest)
		return
	}

	err := bh.BookService.SyncBooksFromGoogleSheets(sheetID, apiKey)
	if err != nil {
		http.Error(w, "Error syncing books", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Books synced successfully"}`))
}

func (bh *BookHandler) GetBookCount(w http.ResponseWriter, r *http.Request) {

	// Call the BookService method to get the total book count
	count, err := bh.BookService.CountBooks()
	if err != nil {
		// Handle the error, e.g., return a 500 Internal Server Error
		http.Error(w, "Failed to count books", http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(w, map[string]interface{}{
		"count": count,
	})
	
}


func (bh *BookHandler) GetAllListingsByBookID(w http.ResponseWriter, r *http.Request) {
    // Extract user ID from context

	body, _ := io.ReadAll(r.Body)
	log.Println("Request Body:", string(body)) 

    userID, ok := r.Context().Value("user_id").(int)
    if !ok {
        sendErrorResponse(w, http.StatusUnauthorized, "User ID missing")
        return
    }

	

    // Extract bookID from URL parameters
    bookIDStr := chi.URLParam(r, "id")
	log.Println("bookIDStr listing:", bookIDStr)

    bookID, err := strconv.Atoi(bookIDStr) // Convert to int
    if err != nil {
        sendErrorResponse(w, http.StatusBadRequest, "Invalid book ID")
        return
    }

    // Call service to get listings for the book
    listing, err := bh.UserService.GetAllListingsByBookID(r.Context(), userID, bookID)
    if err != nil {
        sendErrorResponse(w, http.StatusInternalServerError, "Error fetching listings")
        return
    }

    // Return success response
    sendSuccessResponse(w, map[string]interface{}{
        "success":  true,
        "listings": listing,
    })
}

