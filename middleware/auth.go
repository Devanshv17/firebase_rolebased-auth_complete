package middleware

import (
	"backend/utils"
	"context"
	"net/http"
	"strings"
)

// AuthMiddleware is the middleware to protect routes
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		// Token should be in the format "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			http.Error(w, "Invalid authorization token format", http.StatusUnauthorized)
			return
		}

		// Extract the token
		tokenString := tokenParts[1]

		// Verify the token using the utils.VerifyToken function
		claims, err := utils.VerifyToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Store the UID in context for use in the handler
		ctx := context.WithValue(r.Context(), "uid", claims.UID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
