package core_http_middleware

import "net/http"

func AdminOnly() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			roleID, ok := r.Context().Value(RoleIDKey).(int)

			if !ok || roleID < 2 {
				http.Error(w, "forbidden: admin rights required", http.StatusForbidden)
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
				http.Error(w, "forbidden: super_admin rights required", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
