package middleware

import (
    "context"
    "net/http"
    "used2book-backend/internal/utils"
	"database/sql"
    "strings"
)


func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // ✅ Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing access token", http.StatusUnauthorized)
			return
		}

		// ✅ Extract token from "Bearer TOKEN"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}
		accessToken := tokenParts[1]

		// ✅ Verify token using backend secret
		userID, err := utils.VerifyToken(accessToken, "access")
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// ✅ Attach user ID to request context
		ctx := context.WithValue(r.Context(), "user_id", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AdminMiddleware(db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// ✅ Get user ID from request context (set by AuthMiddleware)
			userID := r.Context().Value("user_id")
			if userID == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// ✅ Get user role from database
			var role string
			err := db.QueryRow("SELECT role FROM users WHERE id = ?", userID).Scan(&role)
			if err != nil {
				http.Error(w, "User not found", http.StatusUnauthorized)
				return
			}

			// ✅ Allow only admins
			if role != "admin" {
				http.Error(w, "Forbidden: Admins only", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}



