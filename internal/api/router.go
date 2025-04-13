package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"used2book-backend/internal/api/routes"
	"used2book-backend/internal/config"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/streadway/amqp"
)

func SetupRouter(db *sql.DB, rabbitConn *amqp.Connection) http.Handler {
	config.InitOAuth()

	r := chi.NewRouter()

	// Basic CORS
  	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
  	r.Use(cors.Handler(cors.Options{
    // AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
    AllowedOrigins:   []string{"https://*", "http://*"},
    // AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
    AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
    ExposedHeaders:   []string{"Link"},
    AllowCredentials: true,
    MaxAge:           300, // Maximum value not ignored by any of major browsers
  	}))

	// ✅ Add logging middleware for debugging
	r.Use(middleware.Logger)

	// ✅ Register API routes correctly
	r.Mount("/auth", routes.AuthRoutes(db))
	r.Mount("/user", routes.UserRoutes(db, rabbitConn))
	r.Mount("/book", routes.BookRoutes(db))
	r.Mount("/auth-token", routes.TokenRoutes(db))
	r.Mount("/payment", routes.PaymentRoutes(db, rabbitConn))

	// ✅ Debugging: Print all registered routes
	fmt.Println("🔍 Registered Routes:")
	_ = chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("%s %s\n", method, route)
		return nil
	})

	return r
}
