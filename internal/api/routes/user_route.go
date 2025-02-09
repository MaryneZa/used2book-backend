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

func UserRoutes(db *sql.DB) http.Handler {

	// Initialize repository, service, and handler for user-related actions
	userRepo := mysql.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	tokenRepo := mysql.NewTokenRepository(db)
	tokenService := services.NewTokenService(tokenRepo, userRepo)
	authHandler := &handlers.AuthHandler{
		UserService:  userService,
		TokenService: tokenService,
	}

	r := chi.NewRouter()

	r.Get("/", authHandler.IndexHandler)
	r.Get("/google-provider", authHandler.GoogleHandler)
	r.Get("/callback", authHandler.GoogleCallbackHandler)
	
	r.Post("/login/email", authHandler.LoginEmailHandler)
	r.Post("/signup/email", authHandler.SignupEmailHandler)
	
	// r.Get("/login/callback", authHandler.LoginGoogleCallbackHandler)
	// r.Get("/signup/callback", authHandler.SignupGoogleCallbackHandler)
	// r.Get("/login-google", authHandler.LoginGoogleHandler)
	// r.Get("/signup-google", authHandler.SignupGoogleHandler)

	// ✅ Add protected route for `getMe`
	r.With(middleware.AuthMiddleware).Get("/me", authHandler.GetMeHandler)
	return r

	// r.Post("/sign-up", user.SignupGoogleHandler)
	// r.Post("/login", user.LoginGoogleHandler)
	// r.Post("/get-access", user.RefreshTokenHandler)
	// r.With(middleware.AdminMiddleware).Post("/get-blacklist", user.AdminMiddleware)
	// r.With(middleware.AuthMiddleware).Get("/users", user.GetUsersHandler)
}
