package middleware

import (
    "context"
    "net/http"
    "used2book-backend/internal/utils"
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



// func AdminMiddleware(next http.Handler) http.Handler {
//     return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//         userID := r.Context().Value("user_id").(int)
        
//         db := utils.GetDB()
//         var role string
//         err := db.Get(&role, `SELECT role FROM users WHERE id = ?`, userID)
//         if err != nil || role != "admin" {
//             http.Error(w, "Unauthorized", http.StatusForbidden)
//             return
//         }

//         next.ServeHTTP(w, r)
//     })
// }

// func AuthMiddleware(next http.Handler) http.Handler {
//     return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//         cookie, err := r.Cookie("access_token")
//         if err != nil {
//             http.Error(w, "Access token not provided", http.StatusUnauthorized)
//             return
//         }

//         accessToken := cookie.Value
//         userID, err := utils.VerifyToken(accessToken, "access")
//         if err != nil {
//             http.Error(w, "Invalid or expired access token", http.StatusUnauthorized)
//             return
//         }

//         // Attach user ID to request context
//         ctx := context.WithValue(r.Context(), "user_id", userID)
//         next.ServeHTTP(w, r.WithContext(ctx))
//     })
// }
