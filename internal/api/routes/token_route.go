package routes

import (
	"database/sql"
	"net/http"
	"used2book-backend/internal/api/handlers"
	"used2book-backend/internal/repository/mysql"
	"used2book-backend/internal/services"

	"github.com/go-chi/chi/v5"
	_ "github.com/go-sql-driver/mysql" // Import the MySQL driver
)

func TokenRoutes(db *sql.DB) http.Handler {

	// Initialize repository, service, and handler for user-related actions
	tokenRepo := mysql.NewTokenRepository(db)
	userRepo := mysql.NewUserRepository(db)
	tokenService := services.NewTokenService(tokenRepo, userRepo)
	userService := services.NewUserService(userRepo)
	tokenHandler := &handlers.TokenHandler{
		TokenService: tokenService,
		UserService:  userService,
	}

	r := chi.NewRouter()

	r.Get("/refresh", tokenHandler.RefreshTokenHandler)

	return r

	// r.Post("/sign-up", user.SignupGoogleHandler)
	// r.Post("/login", user.LoginGoogleHandler)
	// r.Post("/get-access", user.RefreshTokenHandler)
	// r.With(middleware.AdminMiddleware).Post("/get-blacklist", user.AdminMiddleware)
	// r.With(middleware.AuthMiddleware).Get("/users", user.GetUsersHandler)
}
