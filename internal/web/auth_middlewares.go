package web

import (
	"context"
	"net/http"
	"portfolio/internal/jwt"
)

func (m *WebModule) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		autenticatedUser := jwt.GetAutenticatedUserFromRequest(r, m.jwtService)

		if autenticatedUser == nil {
			http.Redirect(w, r, "/app/login", http.StatusFound)
			return
		}

		ctx := context.WithValue(r.Context(), jwt.AutenticatedUserKey, *autenticatedUser)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (m *WebModule) optionalAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		autenticatedUser := jwt.GetAutenticatedUserFromRequest(r, m.jwtService)

		if autenticatedUser == nil {
			next.ServeHTTP(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), jwt.AutenticatedUserKey, *autenticatedUser)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
