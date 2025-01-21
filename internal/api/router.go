package api

import (
    "net/http"

    "github.com/go-chi/chi"
	"used2book-backend/internal/api/handlers"
    "used2book-backend/internal/config"
    "used2book-backend/internal/repository/mysql"
    "used2book-backend/internal/services"
	"database/sql"
    _ "github.com/go-sql-driver/mysql" // Import the MySQL driver
)

func SetupRouter(db *sql.DB) http.Handler {
	
	// Initialize OAuth configuration
    config.InitOAuth()
	
    // Initialize repository, service, and handler for user-related actions
    userRepo := mysql.NewUserRepository(db)
    userService := services.NewUserService(userRepo)
    authHandler := &handlers.AuthHandler{
		UserService: userService,
    }
	
	r := chi.NewRouter()
    r.Get("/", authHandler.IndexHandler)
    r.Get("/login", authHandler.LoginHandler)
    r.Get("/signup", authHandler.SignupHandler)

    r.Get("/login/callback", authHandler.LoginCallbackHandler)
    r.Get("/signup/callback", authHandler.SignupCallbackHandler)

    return r
}
