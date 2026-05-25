package core_http_middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	auth_service "github.com/emount4/concert_reviews/internal/features/auth/service"
	"github.com/google/uuid"
)

type JwtManager interface {
	Parse(accessToken string) (*auth_service.JWTClaims, error)
}

func Auth(jwtManager JwtManager) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			w.Header().Set("Content-Type", "application/json")

			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				writeJSONError(w, http.StatusUnauthorized, "missing auth header")
				return
			}

			parts := strings.Fields(authHeader)

			if len(parts) != 2 || parts[0] != "Bearer" {
				writeJSONError(w, http.StatusUnauthorized, "invalid auth header")
				return
			}

			claims, err := jwtManager.Parse(parts[1])

			if err != nil {
				writeJSONError(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, RoleIDKey, claims.RoleID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func writeJSONError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": message,
		"error":   message,
	})
}

func GetUserID(ctx context.Context) (uuid.UUID, error) {
	id := ctx.Value(UserIDKey)
	if id == nil {
		return uuid.Nil, fmt.Errorf("user id not found in context")
	}
	userID, ok := id.(uuid.UUID)
	if !ok {
		return uuid.Nil, fmt.Errorf("user id in context is not uuid.UUID, got %T", id)
	}
	return userID, nil
}

func GetRoleID(ctx context.Context) (int, error) {
	id := ctx.Value(RoleIDKey)
	if id == nil {
		return 0, fmt.Errorf("role id not found in context")
	}
	roleID, ok := id.(int)
	if !ok {
		return 0, fmt.Errorf("role id in context is not int, got %T", id)
	}
	return roleID, nil
}
