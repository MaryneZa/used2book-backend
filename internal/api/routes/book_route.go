package routes

import (
	"database/sql"
	"net/http"
	"used2book-backend/internal/api/handlers"
	"used2book-backend/internal/repository/mysql"
	"used2book-backend/internal/services"
	"used2book-backend/internal/middleware"

	"github.com/go-chi/chi/v5"
	_ "github.com/go-sql-driver/mysql" // Import the MySQL driver
)

// BookRoutes sets up routes for book-related operations
func BookRoutes(db *sql.DB) http.Handler {
	// Initialize Repositories
	bookRepo := mysql.NewBookRepository(db)

	// Initialize Services
	bookService := services.NewBookService(bookRepo)

	userRepo := mysql.NewUserRepository(db)
	userService := services.NewUserService(userRepo)

	// Initialize Handlers
	bookHandler := &handlers.BookHandler{
		BookService: bookService,
		UserService: userService,
	}

	r := chi.NewRouter()

	// // Public Routes
	// r.Get("/{id:[0-9]+}", bookHandler.GetBookWithRatings) // Get book details with ratings
	// // Protected Routes (Require Authentication)
	// r.With(middleware.AuthMiddleware).Post("/{id:[0-9]+}/rate", bookHandler.AddOrUpdateUserRating)
	// r.With(middleware.AuthMiddleware).Delete("/{id:[0-9]+}/rate", bookHandler.DeleteUserRating)


	r.Get("/all-books", bookHandler.GetAllBooks) // Get book details with ratings
	r.Get("/get-book/{id:[0-9]+}", bookHandler.GetBookByID)
	r.Get("/get-book-genres/{id:[0-9]+}", bookHandler.GetGenresByBookID)


	r.Post("/sync", bookHandler.SyncBooksFromGoogleSheets) // Sync books from Google Sheets
	r.With(middleware.AuthMiddleware).With(middleware.AdminMiddleware(db)).Get("/book-count", bookHandler.GetBookCount) // Sync books from Google Sheets



	r.With(middleware.AuthMiddleware).Get("/{id:[0-9]+}/listings", bookHandler.GetAllListingsByBookID)
	r.With(middleware.AuthMiddleware).Get("/{id:[0-9]+}/get-reviews", bookHandler.GetReviewsByBookIDHandler)


	r.With(middleware.AuthMiddleware).Post("/add-review", bookHandler.AddBookReviewHandler)
	
	r.With(middleware.AuthMiddleware).Get("/all-genres", bookHandler.GetAllGenres)

	r.Get("/all-book-genres", bookHandler.GetAllBookGenres)

	r.With(middleware.AuthMiddleware).Get("/recommended-books", bookHandler.GetRecommendedBooks)



	return r
}
