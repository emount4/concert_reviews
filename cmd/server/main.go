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
	core_redis "github.com/emount4/concert_reviews/internal/core/repository/redis"

	core_s3 "github.com/emount4/concert_reviews/internal/core/repository/s3"
	core_http_middleware "github.com/emount4/concert_reviews/internal/core/transport/http/middleware"
	core_http_server "github.com/emount4/concert_reviews/internal/core/transport/http/server"
	artist_postgres_repository "github.com/emount4/concert_reviews/internal/features/artist/repository/postgres"
	artist_service "github.com/emount4/concert_reviews/internal/features/artist/service"
	artist_transport_http "github.com/emount4/concert_reviews/internal/features/artist/transport/http"
	auth_postgres_repository "github.com/emount4/concert_reviews/internal/features/auth/repository/postgres"
	auth_service "github.com/emount4/concert_reviews/internal/features/auth/service"
	auth_transport_http "github.com/emount4/concert_reviews/internal/features/auth/transport/http"
	city_postgres_repository "github.com/emount4/concert_reviews/internal/features/city/repository/postgres"
	city_service "github.com/emount4/concert_reviews/internal/features/city/service"
	city_transport_http "github.com/emount4/concert_reviews/internal/features/city/transport/http"
	concert_postgres_repository "github.com/emount4/concert_reviews/internal/features/concert/repository/postgres"
	concert_service "github.com/emount4/concert_reviews/internal/features/concert/service"
	concert_transport_http "github.com/emount4/concert_reviews/internal/features/concert/transport/http"
	favorites_postgres_repository "github.com/emount4/concert_reviews/internal/features/favorites/repository/postgres"
	favorites_service "github.com/emount4/concert_reviews/internal/features/favorites/service"
	favorites_transport_http "github.com/emount4/concert_reviews/internal/features/favorites/transport/http"
	review_postgres_repository "github.com/emount4/concert_reviews/internal/features/reviews/repository/postgres"
	review_redis_repository "github.com/emount4/concert_reviews/internal/features/reviews/repository/redis"
	review_service "github.com/emount4/concert_reviews/internal/features/reviews/service"
	review_transport_http "github.com/emount4/concert_reviews/internal/features/reviews/transport/http"
	stats_postgres_repository "github.com/emount4/concert_reviews/internal/features/stats/repository/postgres"
	stats_redis_repository "github.com/emount4/concert_reviews/internal/features/stats/repository/redis"
	stats_service "github.com/emount4/concert_reviews/internal/features/stats/service"
	stats_transport_http "github.com/emount4/concert_reviews/internal/features/stats/transport/http"
	user_postgres_repository "github.com/emount4/concert_reviews/internal/features/user/repository/postgres"
	user_service "github.com/emount4/concert_reviews/internal/features/user/service"
	user_transport_http "github.com/emount4/concert_reviews/internal/features/user/transport/http"
	venue_postgres_repository "github.com/emount4/concert_reviews/internal/features/venues/repository/postgres"
	venue_service "github.com/emount4/concert_reviews/internal/features/venues/service"
	venue_transport_http "github.com/emount4/concert_reviews/internal/features/venues/transport/http"

	media_service "github.com/emount4/concert_reviews/internal/features/media/service"
	media_transport_http "github.com/emount4/concert_reviews/internal/features/media/transport/http"
	moderation_postgres_repository "github.com/emount4/concert_reviews/internal/features/moderation/repository/postgres"
	moderation_service "github.com/emount4/concert_reviews/internal/features/moderation/service"
	moderation_transport_http "github.com/emount4/concert_reviews/internal/features/moderation/transport/http"
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

	txManager := core_postgres_tx.NewManager(pool)

	logger.Debug("initializing s3")
	s3Config := core_s3.NewConfigMust()
	s3Storage, err := core_s3.NewS3Storage(logger, s3Config)
	if err != nil {
		logger.Fatal("failed to init s3 storage", zap.Error(err))
	}

	logger.Debug("initializing redis")
	redisConfig := core_redis.NewConfigMust()
	redis, err := core_redis.NewClient(ctx, redisConfig)
	if err != nil {
		logger.Fatal("failed to init redis storage", zap.Error(err))
	}

	logger.Debug("initializing features", zap.String("features", "auth"))
	authRepository := auth_postgres_repository.NewAuthRepository(pool)
	authConfig := auth_service.NewConfigMust()
	hasher := auth_service.NewSHA1Hasher(authConfig.PasswordSalt)
	jwtManager := auth_service.NewManager(authConfig.JWTSigningKey)
	authService := auth_service.NewAuthService(authRepository, txManager, authConfig, hasher, jwtManager)
	authTransportHTTP := auth_transport_http.NewAuthHTTPHandler(authService)
	authRoutes := authTransportHTTP.Routes()

	logger.Debug("initializing features", zap.String("features", "stats"))
	statsRepository := stats_postgres_repository.NewStatsRepository(pool, txManager)
	statsRedisRepository := stats_redis_repository.NewStatsRedisRepository(redis)
	statsService := stats_service.NewStatsService(statsRepository, statsRedisRepository)
	statsTransportHTTP := stats_transport_http.NewStatsHTTPHandler(statsService)
	statsRoutes := statsTransportHTTP.Routes()

	logger.Debug("initializing features", zap.String("features", "artist"))
	artistRepository := artist_postgres_repository.NewArtistRepository(pool)
	artistService := artist_service.NewArtistService(artistRepository, s3Storage)
	artistsTransportHTTP := artist_transport_http.NewArtistHTTPHandler(artistService)
	artistRoutes := artistsTransportHTTP.Routes()
	applyRouteAccessPolicy(artistRoutes, jwtManager)

	logger.Debug("initializing features", zap.String("features", "city"))
	cityRepository := city_postgres_repository.NewCityRepository(pool)
	cityService := city_service.NewCityService(cityRepository)
	cityTransportHTTP := city_transport_http.NewCityHTTPHandler(cityService)
	cityRoutes := cityTransportHTTP.Routes()
	applyRouteAccessPolicy(cityRoutes, jwtManager)

	logger.Debug("initializing features", zap.String("features", "venue"))
	venueRepository := venue_postgres_repository.NewVenueRepository(pool)
	venueService := venue_service.NewVenueService(venueRepository, s3Storage)
	venueTransportHTTP := venue_transport_http.NewVenueHTTPHandler(venueService)
	venueRoutes := venueTransportHTTP.Routes()
	applyRouteAccessPolicy(venueRoutes, jwtManager)

	logger.Debug("initializing features", zap.String("features", "concerts"))
	concertRepository := concert_postgres_repository.NewConcertRepository(pool, txManager)
	concertService := concert_service.NewConcertService(concertRepository, s3Storage)
	concertTransportHTTP := concert_transport_http.NewConcertHTTPHandler(concertService)
	concertRoutes := concertTransportHTTP.Routes()
	applyRouteAccessPolicy(concertRoutes, jwtManager)

	logger.Debug("initializing features", zap.String("features", "reviews"))
	reviewRepository := review_postgres_repository.NewReviewRepository(pool, txManager)
	reviewRedis := review_redis_repository.NewReviewRedisRepository(redis)
	reviewService := review_service.NewReviewService(reviewRepository, s3Storage, reviewRedis)
	reviewHTTPHandler := review_transport_http.NewReviewHTTPHandler(reviewService)
	reviewRoutes := reviewHTTPHandler.Routes()
	applyRouteAccessPolicy(reviewRoutes, jwtManager)

	logger.Debug("initializing features", zap.String("features", "users"))
	userRepository := user_postgres_repository.NewReviewRepository(pool, txManager)
	userService := user_service.NewUserService(userRepository, reviewRepository, s3Storage)
	userHTTPHandler := user_transport_http.NewUsersHTTPHandler(userService)
	userRoutes := userHTTPHandler.Routes()
	applyRouteAccessPolicy(userRoutes, jwtManager)

	logger.Debug("initializing features", zap.String("features", "favorites"))
	favoritesRepository := favorites_postgres_repository.NewFavoritesRepository(pool, txManager)
	favoritesService := favorites_service.NewFavoritesService(favoritesRepository)
	favoritesTransportHTTP := favorites_transport_http.NewFavoritesHTTPHandler(favoritesService)
	favoritesRoutes := favoritesTransportHTTP.Routes()
	applyRouteAccessPolicy(favoritesRoutes, jwtManager)

	allowedExt := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
		".gif":  true,
		".mp4":  true,
	}
	logger.Debug("initializing features", zap.String("features", "media"))
	mediaService := media_service.NewMediaService(
		s3Storage,
		allowedExt,
		s3Config.UploadMinMB*1024*1024,
		s3Config.UploadMaxMB*1024*1024,
	)
	mediaTransportHTTP := media_transport_http.NewMediaHTTPHandler(mediaService)
	mediaRoutes := mediaTransportHTTP.Routes()
	applyRouteAccessPolicy(mediaRoutes, jwtManager)

	logger.Debug("initializing features", zap.String("features", "moderation"))
	moderationRepository := moderation_postgres_repository.NewModerationRepository(pool)
	moderationService := moderation_service.NewModerationService(moderationRepository)
	moderationTransportHTTP := moderation_transport_http.NewModerationHTTPHandler(moderationService)
	moderationRoutes := moderationTransportHTTP.Routes()
	applyRouteAccessPolicy(moderationRoutes, jwtManager)

	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)

	apiVersionRouter.RigisterRoutes(authRoutes...)
	apiVersionRouter.RigisterRoutes(statsRoutes...)
	apiVersionRouter.RigisterRoutes(cityRoutes...)
	apiVersionRouter.RigisterRoutes(mediaRoutes...)
	apiVersionRouter.RigisterRoutes(artistRoutes...)
	apiVersionRouter.RigisterRoutes(venueRoutes...)
	apiVersionRouter.RigisterRoutes(concertRoutes...)
	apiVersionRouter.RigisterRoutes(reviewRoutes...)
	apiVersionRouter.RigisterRoutes(userRoutes...)
	apiVersionRouter.RigisterRoutes(favoritesRoutes...)
	apiVersionRouter.RigisterRoutes(moderationRoutes...)

	httpConfig := core_http_server.NewConfigMust()

	httpServer := core_http_server.NewHTTPServer(
		httpConfig,
		logger,
		//подключение основных мв
		core_http_middleware.CORSFromCSV(
			httpConfig.CORSAllowedOrigins,
			httpConfig.CORSAllowCredentials,
			httpConfig.CORSMaxAgeSeconds,
		),
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)
	httpServer.RegisterAPIRouters(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error: %w", zap.Error(err))
	}
}

func applyRouteAccessPolicy(routes []core_http_server.Route, jwtManager auth_service.JWTManager) {
	for i := range routes {
		switch routes[i].Access {
		case core_http_server.AccessAdminOnly:
			routes[i].Middleware = append(
				routes[i].Middleware,
				core_http_middleware.Auth(jwtManager),
				core_http_middleware.AdminOnly(),
			)
		case core_http_server.AccessSuperAdminOnly:
			routes[i].Middleware = append(
				routes[i].Middleware,
				core_http_middleware.Auth(jwtManager),
				core_http_middleware.SuperAdminOnly(),
			)
		case core_http_server.AccessAuthOnly:
			routes[i].Middleware = append(
				routes[i].Middleware,
				core_http_middleware.Auth(jwtManager),
			)
		default:
			routes[i].Middleware = append(
				routes[i].Middleware,
				core_http_middleware.OptionalAuth(jwtManager),
			)
		}
	}
}
