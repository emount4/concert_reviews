package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_logger "github.com/emount4/concert_reviews/internal/core/logger"
	core_postgres_pool "github.com/emount4/concert_reviews/internal/core/repository/postgres/pool"
	core_postgres_tx "github.com/emount4/concert_reviews/internal/core/repository/postgres/tx"
	core_http_middleware "github.com/emount4/concert_reviews/internal/core/transport/http/middleware"
	core_http_server "github.com/emount4/concert_reviews/internal/core/transport/http/server"
	auth_postgres_repository "github.com/emount4/concert_reviews/internal/features/auth/repository/postgres"
	auth_service "github.com/emount4/concert_reviews/internal/features/auth/service"
	auth_transport_http "github.com/emount4/concert_reviews/internal/features/auth/transport/http"
	user_transport_http "github.com/emount4/concert_reviews/internal/features/user/transport/http"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())

	if err != nil {
		fmt.Println("failed to init app logger: %w", err)
		os.Exit(1)
	}

	defer logger.Close()

	logger.Debug("starting application")

	logger.Debug("initializing postgres pool")
	pool, err := core_postgres_pool.NewConnectionPool(
		ctx,
		core_postgres_pool.NewConfigMust(),
	)
	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("initializing features", zap.String("features", "auth"))
	authRepository := auth_postgres_repository.NewAuthRepository(pool)
	txManager := core_postgres_tx.NewManager(pool)

	usersTransportHTTP := user_transport_http.NewUsersHTTPHandler(nil)
	usersRoutes := usersTransportHTTP.Routes()

	authConfig := auth_service.NewConfigMust()
	hasher := auth_service.NewSHA1Hasher(authConfig.PasswordSalt)
	jwtManager := auth_service.NewManager(authConfig.JWTSigningKey)
	authService := auth_service.NewAuthService(authRepository, txManager, authConfig, hasher, jwtManager)
	authTransportHTTP := auth_transport_http.NewAuthHTTPHandler(authService)
	authRoutes := authTransportHTTP.Routes()

	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RigisterRoutes(usersRoutes...)
	apiVersionRouter.RigisterRoutes(authRoutes...)

	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		//подключение основных мв
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Panic(),
		core_http_middleware.Trace(),
	)
	httpServer.RegisterAPIRouters(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error: %w", zap.Error(err))
	}
}
