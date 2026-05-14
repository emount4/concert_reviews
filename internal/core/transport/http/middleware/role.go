package core_http_middleware

import (
	"context"
	"net/http"
	"strings"
)

func AuthOnly() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			roleID, ok := r.Context().Value(RoleIDKey).(int)

			if !ok || roleID < 1 {
				writeJSONError(w, http.StatusForbidden, "forbidden: auth rights required")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func AdminOnly() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			roleID, ok := r.Context().Value(RoleIDKey).(int)

			if !ok || roleID < 2 {
				writeJSONError(w, http.StatusForbidden, "forbidden: admin rights required")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func SuperAdminOnly() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			roleID, ok := r.Context().Value(RoleIDKey).(int)

			if !ok || roleID < 3 {
				writeJSONError(w, http.StatusForbidden, "forbidden: super_admin rights required")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func OptionalAuth(jwtManager JwtManager) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				next.ServeHTTP(w, r)
				return
			}

			parts := strings.Fields(authHeader)
			if len(parts) != 2 || parts[0] != "Bearer" {
				writeJSONError(w, http.StatusUnauthorized, "invalid auth header format")
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
