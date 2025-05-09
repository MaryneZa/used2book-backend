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
	"math/rand"
	"github.com/go-chi/chi/v5"

	// "github.com/gorilla/mux"
)

// BookHandler handles book-related HTTP requests
type BookHandler struct {
	BookService *services.BookService
	UserService *services.UserService
	UploadService *services.UploadService

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

func (bh *BookHandler) GetAllGenres(w http.ResponseWriter, r *http.Request) {

	// Call the BookService method to get the total book count
	genres, err := bh.BookService.GetAllGenres(r.Context())
	log.Println(genres)
	if err != nil {
		// Handle the error, e.g., return a 500 Internal Server Error
		http.Error(w, "Failed to get all genres", http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(w, map[string]interface{}{
		"genres": genres,
	})
	
}

func (bh *BookHandler) GetAllAuthors(w http.ResponseWriter, r *http.Request) {

	// Call the BookService method to get the total book count
	authors, err := bh.BookService.GetAllAuthors(r.Context())
	log.Println(authors)
	if err != nil {
		// Handle the error, e.g., return a 500 Internal Server Error
		http.Error(w, "Failed to get all authors", http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(w, map[string]interface{}{
		"authors": authors,
	})
	
}

func (bh *BookHandler) GetAllBookAuthors(w http.ResponseWriter, r *http.Request) {

	// Call the BookService method to get the total book count
	book_authors, err := bh.BookService.GetAllBookAuthors(r.Context())
	log.Println(book_authors)
	if err != nil {
		// Handle the error, e.g., return a 500 Internal Server Error
		http.Error(w, "Failed to get all book_authors", http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(w, map[string]interface{}{
		"book_authors": book_authors,
	})
	
}




func (bh *BookHandler) GetRecommendedBooks(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		sendErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Fetch recommendations from the Flask API
	resp, err := http.Get("http://localhost:5005/recommendations?user_id=" + strconv.Itoa(userID))
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to fetch recommendations: "+err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		sendErrorResponse(w, http.StatusInternalServerError, "Recommendation API returned: "+resp.Status)
		return
	}

	// Decode the recommendation response
	var recs struct {
		Recommendations []models.Recommendation `json:"recommendations"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&recs); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to parse recommendations: "+err.Error())
		return
	}

	// Fetch book details for each recommended ID
	var books []models.Book
	for _, rec := range recs.Recommendations {
		book, err := bh.BookService.GetBookByID(context.Background(), rec.ID)
		if err != nil {
			log.Printf("Failed to fetch book ID %d: %v", rec.ID, err)
			continue // Skip failed books, or handle differently
		}
		books = append(books, *book)
	}

	perm := rand.Perm(len(books))
	shuffledBooks := make([]models.Book, len(books))
	for i, p := range perm {
		shuffledBooks[i] = books[p]
	}
	books = shuffledBooks

	// Send response
	sendSuccessResponse(w, map[string]interface{}{
		"books": books,
	})
}

func (bh *BookHandler) GetAllBookGenres(w http.ResponseWriter, r *http.Request) {

	// Call the BookService method to get the total book count
	book_genres, err := bh.BookService.GetAllBookGenres(r.Context())
	// log.Println(book_genres)
	if err != nil {
		// Handle the error, e.g., return a 500 Internal Server Error
		http.Error(w, "Failed to get all genres", http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(w, map[string]interface{}{
		"book_genres": book_genres,
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

func (bh *BookHandler) GetReviewsByUserIDHandler(w http.ResponseWriter, r *http.Request) {
	
	body, _ := io.ReadAll(r.Body)
	log.Println("Request Body:", string(body)) 

	userIDStr := chi.URLParam(r, "userID")
	log.Println("userIDStr:", userIDStr)

	if userIDStr == "" {
		sendErrorResponse(w, http.StatusBadRequest, "Book ID is required")
		return
	}

	userID, err := strconv.Atoi(userIDStr)

	// Call the BookService method to get the total book count
	reviews, err := bh.BookService.GetReviewsByUserID(r.Context(), userID)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to get reviews"+err.Error())
		return
	}

	sendSuccessResponse(w, map[string]interface{}{
		"reviews": reviews,
	})
	
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



// InsertBookHandler handles book insertion with optional cover image upload
func (bh *BookHandler) InsertBookHandler(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form (10MB max)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Println("Parse Form Error:", err)
		sendErrorResponse(w, http.StatusBadRequest, "Invalid form data")
		return
	}

	// Get JSON data from form field 'data'
	jsonData := r.FormValue("data")
	var bookForm models.BookForm
	if err := json.Unmarshal([]byte(jsonData), &bookForm); err != nil {
		log.Println("JSON Decode Error:", err)
		sendErrorResponse(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// Validate required fields
	if bookForm.Title == "" || len(bookForm.Author) == 0 {
		sendErrorResponse(w, http.StatusBadRequest, "Title and Author are required")
		return
	}
	if len(bookForm.Genres) == 0 {
		sendErrorResponse(w, http.StatusBadRequest, "At least one genre is required")
		return
	}

	// Handle cover image upload (optional)
	var coverImageURL string
	files := r.MultipartForm.File["cover_image"]
	if len(files) > 0 {
		if len(files) > 1 {
			sendErrorResponse(w, http.StatusBadRequest, "Only one cover image is allowed")
			return
		}
		fileHandler := files[0]
		file, err := fileHandler.Open()
		if err != nil {
			sendErrorResponse(w, http.StatusInternalServerError, "Error opening file")
			return
		}
		defer file.Close()

		coverImageURL, err = bh.UploadService.UploadImageURL(file, fileHandler.Filename)
		if err != nil {
			sendErrorResponse(w, http.StatusInternalServerError, "Image upload failed: "+err.Error())
			return
		}
	}

	// Create book model directly from bookForm (PublishDate is already time.Time)
	book := models.Book{
		Title:         bookForm.Title,
		Author:        bookForm.Author,
		Description:   bookForm.Description,
		Language:      bookForm.Language,
		ISBN:          bookForm.ISBN,
		Publisher:     bookForm.Publisher,
		PublishDate:   bookForm.PublishDate, // Use directly, no parsing needed
		CoverImageURL: coverImageURL,
	}

	// Insert book into database
	bookID, err := bh.BookService.InsertBook(r.Context(), book)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to insert book: "+err.Error())
		return
	}

	// Insert or get genres and associate with book
	for _, genreName := range bookForm.Genres {
		genreID, err := bh.BookService.GetOrInsertGenre(r.Context(), genreName)
		if err != nil {
			sendErrorResponse(w, http.StatusInternalServerError, "Failed to process genre "+genreName+": "+err.Error())
			return
		}
		err = bh.BookService.AssociateBookWithGenre(r.Context(), bookID, genreID)
		if err != nil {
			sendErrorResponse(w, http.StatusInternalServerError, "Failed to associate genre "+genreName+": "+err.Error())
			return
		}
	}

	// Success response
	sendSuccessResponse(w, map[string]interface{}{
		"success": true,
		"message": "Book inserted successfully",
		"book_id": bookID,
	})
}

func (bh *BookHandler) UpdateBookHandler(w http.ResponseWriter, r *http.Request) {
	bookIDStr := chi.URLParam(r, "bookID")
	bookID, err := strconv.Atoi(bookIDStr)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid book ID")
		return
	}

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid form data")
		return
	}

	jsonData := r.FormValue("data")
	var bookForm models.BookForm
	if err := json.Unmarshal([]byte(jsonData), &bookForm); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	if bookForm.Title == "" || len(bookForm.Author) == 0 {
		sendErrorResponse(w, http.StatusBadRequest, "Title and Author are required")
		return
	}

	if len(bookForm.Genres) == 0 {
		sendErrorResponse(w, http.StatusBadRequest, "At least one genre is required")
		return
	}

	// Optional cover image upload
	var coverImageURL string
	files := r.MultipartForm.File["cover_image"]
	if len(files) > 0 {
		fileHandler := files[0]
		file, err := fileHandler.Open()
		if err != nil {
			sendErrorResponse(w, http.StatusInternalServerError, "Error opening file")
			return
		}
		defer file.Close()

		coverImageURL, err = bh.UploadService.UploadImageURL(file, fileHandler.Filename)
		if err != nil {
			sendErrorResponse(w, http.StatusInternalServerError, "Image upload failed: "+err.Error())
			return
		}
	}

	// If no new image uploaded, use existing (client should send it back)
	if coverImageURL == "" {
		coverImageURL = bookForm.CoverImageURL
	}

	book := models.Book{
		Title:         bookForm.Title,
		Author:        bookForm.Author,
		Description:   bookForm.Description,
		Language:      bookForm.Language,
		ISBN:          bookForm.ISBN,
		Publisher:     bookForm.Publisher,
		PublishDate:   bookForm.PublishDate,
		CoverImageURL: coverImageURL,
	}

	err = bh.BookService.UpdateBook(r.Context(), bookID, book)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to update book: "+err.Error())
		return
	}

	// Handle genres
	err = bh.BookService.UpdateBookGenres(r.Context(), bookID, bookForm.Genres)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to update genres: "+err.Error())
		return
	}

	sendSuccessResponse(w, map[string]interface{}{
		"success": true,
		"message": "Book updated successfully",
	})
}

