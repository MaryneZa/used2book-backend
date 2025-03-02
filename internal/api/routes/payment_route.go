package routes

import (
	"database/sql"
	"net/http"
	"used2book-backend/internal/api/handlers"
	"used2book-backend/internal/repository/mysql"
	"used2book-backend/internal/services"
	"used2book-backend/internal/middleware"
	"github.com/go-chi/chi/v5"
)

// PaymentRoutes initializes all payment-related routes
func PaymentRoutes(db *sql.DB) http.Handler {
	// Initialize required services and repositories
	omiseService := services.NewOmiseService()
	userRepo := mysql.NewUserRepository(db)
	userService := services.NewUserService(userRepo)

	// Initialize payment handler
	paymentHandler := &handlers.PaymentHandler{
		OmiseService: omiseService,
		UserService:  userService,
	}

	// Create a new router
	r := chi.NewRouter()

	// Define payment-related API endpoints
	r.With(middleware.AuthMiddleware).Post("/api/omise/create-account", paymentHandler.CreateOmiseAccountHandler) // Create Omise account for sellers

	// Future payment-related routes can be added here

	return r
}
