package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"used2book-backend/internal/services"
	"log"
	"github.com/gorilla/mux"
	"github.com/go-chi/chi/v5"

)

// BookHandler handles book-related HTTP requests
type BookHandler struct {
	BookService *services.BookService
}

// NewBookHandler initializes BookHandler
func NewBookHandler(BookService *services.BookService) *BookHandler {
	return &BookHandler{BookService: BookService}
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

// GetBookWithRatings retrieves a book and its ratings
func (bh *BookHandler) GetBookWithRatings(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	book, err := bh.BookService.GetBookWithRatings(context.Background(), bookID)
	if err != nil {
		http.Error(w, "Error fetching book", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

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



// AddOrUpdateUserRating allows a user to add or update their rating
func (bh *BookHandler) AddOrUpdateUserRating(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("user_id").(int) // Assume User-ID is passed in header
	var payload struct {
		Rating float64 `json:"rating"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = bh.BookService.AddOrUpdateUserRating(context.Background(), userID, bookID, payload.Rating)
	if err != nil {
		http.Error(w, "Error updating rating", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Rating updated successfully"}`))
}

// DeleteUserRating allows a user to delete their rating
func (bh *BookHandler) DeleteUserRating(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("user_id").(int) // Assume User-ID is passed in header

	err = bh.BookService.DeleteUserRating(context.Background(), userID, bookID)
	if err != nil {
		http.Error(w, "Error deleting rating", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Rating deleted successfully"}`))
}

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

