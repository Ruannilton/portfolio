package jwt

import (
	"context"
	"net/http"
	"strings"
)

const UserIDKey string = "userID"
const UserEmailKey string = "userEmail"

func AuthMiddleware(tokenService *TokenService, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extrai o token do header Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Espera formato "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// Valida o token
		userID, err := tokenService.GetUserIDFromToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		userEmail, err := tokenService.GetUserEmailFromToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Adiciona o userID e userEmail ao contexto
		ctx := context.WithValue(r.Context(), UserIDKey, *userID)
		ctx = context.WithValue(ctx, UserEmailKey, *userEmail)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
