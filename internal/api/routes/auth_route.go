package routes

import (
	"database/sql"
	"net/http"
	"used2book-backend/internal/api/handlers"
	"used2book-backend/internal/middleware"
	"used2book-backend/internal/repository/mysql"
	"used2book-backend/internal/services"

	"github.com/go-chi/chi/v5"
	_ "github.com/go-sql-driver/mysql" // Import the MySQL driver
)

func AuthRoutes(db *sql.DB) http.Handler {

	// Initialize repository, service, and handler for user-related actions
	userRepo := mysql.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo)
	tokenRepo := mysql.NewTokenRepository(db)
	tokenService := services.NewTokenService(tokenRepo, userRepo)
	authHandler := &handlers.AuthHandler{
		AuthService:  authService,
		TokenService: tokenService,
	}
	twiliootpHandler := &handlers.TwilioOTPHandler{
		UserService : userService,
	}

	
	r := chi.NewRouter()
	
	r.Get("/google-provider", authHandler.GoogleHandler)
	r.Get("/callback", authHandler.GoogleCallbackHandler)
	
	r.Post("/login/email", authHandler.LoginEmailHandler)
	r.Post("/signup/email", authHandler.SignupEmailHandler)
	
	r.With(middleware.AuthMiddleware).Post("/send-otp", twiliootpHandler.SendOTPHandler)
	r.With(middleware.AuthMiddleware).Post("/verify-otp", twiliootpHandler.VerifyOTPHandler)
	r.With(middleware.AuthMiddleware).Post("/resend-otp", twiliootpHandler.ResendOTPHandler)




	return r

}
