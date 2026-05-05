package core_http_middleware

import "net/http"

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
