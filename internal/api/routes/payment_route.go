package routes

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"net/http"
	"used2book-backend/internal/api/handlers"
	"used2book-backend/internal/middleware"
	"used2book-backend/internal/repository/mysql"
	"used2book-backend/internal/services"
	"github.com/streadway/amqp"
	
)

// PaymentRoutes initializes all payment-related routes
func PaymentRoutes(db *sql.DB, rabbitConn *amqp.Connection) http.Handler {
	// Initialize required services and repositories
	omiseService := services.NewOmiseService()
	userRepo := mysql.NewUserRepository(db)
	userService := services.NewUserService(userRepo)


	// Initialize payment handler
	paymentHandler := &handlers.PaymentHandler{
		OmiseService: omiseService,
		UserService:  userService,
		RabbitMQConn: rabbitConn,
	}

	// Create a new router
	r := chi.NewRouter()

	// Define payment-related API endpoints
	r.With(middleware.AuthMiddleware).Post("/api/omise/create-account", paymentHandler.CreateOrUpdateOmiseAccountHandler) // Create Omise account for sellers
	r.With(middleware.AuthMiddleware).Post("/charge", paymentHandler.ChargeHandler)
	r.With(middleware.AuthMiddleware).Post("/omise/create-account", paymentHandler.CreateOrUpdateOmiseAccountHandler)
	r.Post("/webhook", paymentHandler.WebhookHandler)
	r.With(middleware.AuthMiddleware).Get("/omise/bank-account", paymentHandler.GetBankAccountInfoHandler)

	// Future payment-related routes can be added here

	return r
}
