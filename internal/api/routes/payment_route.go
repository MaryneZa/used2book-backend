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
	userRepo := mysql.NewUserRepository(db)
	userService := services.NewUserService(userRepo)


	// Initialize payment handler
	paymentHandler := &handlers.PaymentHandler{
		UserService:  userService,
		RabbitMQConn: rabbitConn,
	}

	// Create a new router
	r := chi.NewRouter()

	r.With(middleware.AuthMiddleware).Post("/check-out", paymentHandler.CheckOutHandler)
	r.Post("/webhook", paymentHandler.WebhookHandler)


	return r
}
