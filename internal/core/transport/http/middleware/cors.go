package core_http_middleware

import (
	"net/http"
	"strconv"
	"strings"
)

const (
	headerOrigin                   = "Origin"
	headerVary                     = "Vary"
	headerAccessControlRequestHdrs = "Access-Control-Request-Headers"
	headerAccessControlRequestMth  = "Access-Control-Request-Method"

	headerAllowOrigin      = "Access-Control-Allow-Origin"
	headerAllowMethods     = "Access-Control-Allow-Methods"
	headerAllowHeaders     = "Access-Control-Allow-Headers"
	headerAllowCredentials = "Access-Control-Allow-Credentials"
	headerExposeHeaders    = "Access-Control-Expose-Headers"
	headerMaxAge           = "Access-Control-Max-Age"
)

// CORSFromCSV enables CORS for explicit origins from a comma-separated list.
// Example: "http://127.0.0.1:5174,https://app.example.com"
func CORSFromCSV(originsCSV string, allowCredentials bool, maxAgeSeconds int) Middleware {
	allowedOrigins := parseOrigins(originsCSV)
	allowAnyOrigin := len(allowedOrigins) == 1 && allowedOrigins[0] == "*"

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get(headerOrigin)
			if origin == "" {
				next.ServeHTTP(w, r)
				return
			}

			isOriginAllowed := allowAnyOrigin || containsOrigin(allowedOrigins, origin)
			if !isOriginAllowed {
				if isPreflight(r) {
					http.Error(w, "cors origin is not allowed", http.StatusForbidden)
					return
				}
				next.ServeHTTP(w, r)
				return
			}

			w.Header().Add(headerVary, headerOrigin)
			if allowAnyOrigin && !allowCredentials {
				w.Header().Set(headerAllowOrigin, "*")
			} else {
				w.Header().Set(headerAllowOrigin, origin)
			}

			if allowCredentials {
				w.Header().Set(headerAllowCredentials, "true")
			}
			w.Header().Set(headerExposeHeaders, "X-Request-ID")

			if isPreflight(r) {
				w.Header().Add(headerVary, headerAccessControlRequestMth)
				w.Header().Add(headerVary, headerAccessControlRequestHdrs)

				w.Header().Set(headerAllowMethods, "GET, POST, PUT, PATCH, DELETE, OPTIONS")

				requestedHeaders := r.Header.Get(headerAccessControlRequestHdrs)
				if requestedHeaders == "" {
					w.Header().Set(headerAllowHeaders, "Authorization, Content-Type, X-Request-ID")
				} else {
					w.Header().Set(headerAllowHeaders, requestedHeaders)
				}

				if maxAgeSeconds > 0 {
					w.Header().Set(headerMaxAge, strconv.Itoa(maxAgeSeconds))
				}

				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func parseOrigins(originsCSV string) []string {
	parts := strings.Split(originsCSV, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		origin := strings.TrimSpace(part)
		if origin == "" {
			continue
		}
		result = append(result, origin)
	}
	return result
}

func containsOrigin(allowed []string, origin string) bool {
	for _, item := range allowed {
		if item == origin {
			return true
		}
	}
	return false
}

func isPreflight(r *http.Request) bool {
	return r.Method == http.MethodOptions && r.Header.Get(headerAccessControlRequestMth) != ""
}
