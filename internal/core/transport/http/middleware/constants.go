package core_http_middleware

type contextKey string

const (
	UserIDKey contextKey = "user_id"
	RoleIDKey contextKey = "role_id"
	LoggerKey contextKey = "logger"
)
