package core_http_middleware

import (
	"context"
	"net/http"
	"strings"

	auth_service "github.com/emount4/concert_reviews/internal/features/auth/service"
)

func Auth(jwtManager auth_service.JWTManager) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				http.Error(w, "missing auth header", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")

			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid auth header", http.StatusUnauthorized)
				return
			}

			claims, err := jwtManager.Parse(parts[1])

			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
			ctx = context.WithValue(ctx, "role_id", claims.RoleID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
