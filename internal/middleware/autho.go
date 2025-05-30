package middleware

import (
	"avito_tech/internal/auth"
	"context"
	"net/http"
	"strings"
)

func Middleware(jwt *auth.JwtProvider) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/api/auth" {
				next.ServeHTTP(w, r)
				return
			}

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			claims, err := jwt.ParseWithClaims(tokenString)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
			}

			ctx := context.WithValue(r.Context(), "jwt_claims", claims)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
