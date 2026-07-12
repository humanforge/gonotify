package httpserver

import (
	"context"

	"notification-service/internal/platform/config"
	"notification-service/internal/platform/postgres"

	"github.com/labstack/echo/v5"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

type Server struct {
	Echo   *echo.Echo
	logger *zap.Logger
	cfg    *config.Config
	db     *postgres.DBConnection
}

func New(cfg *config.Config, logger *zap.Logger, db *postgres.DBConnection) *Server {
	e := echo.New()

	s := &Server{
		Echo:   e,
		logger: logger,
		cfg:    cfg,
		db:     db,
	}

	s.registerMiddlewares()
	s.registerErrorHandler()
	s.registerRoutes()

	return s
}

func (s *Server) registerMiddlewares() {
	s.Echo.Use(RequestIDMiddleware())
	s.Echo.Use(SecureHeadersMiddleware())
	s.Echo.Use(CORSMiddleware(s.cfg.CORSAllowedOrigins))
	s.Echo.Use(RecoverMiddleware(s.logger))
	s.Echo.Use(RequestLoggerMiddleware(s.logger))

	if s.cfg.DBConnMaxLifetime > 0 {
		s.Echo.Use(TimeoutMiddleware(s.cfg.DBConnMaxLifetime))
	}

	rl := NewRateLimiter(rate.Limit(100), 200)
	s.Echo.Use(RateLimitMiddleware(rl))
}

func (s *Server) registerErrorHandler() {
	s.Echo.HTTPErrorHandler = NewErrorHandler(s.logger, s.cfg.Env == "production")
}

func (s *Server) registerRoutes() {
	health := NewHealthHandler(s.db, s.logger)

	api := s.Echo.Group("/api/v1")

	s.Echo.GET("/healthz", health.Healthz)
	s.Echo.GET("/readyz", health.Readyz)

	api.GET("/healthz", health.Healthz)
	api.GET("/readyz", health.Readyz)
}

func (s *Server) Start(ctx context.Context) error {
	s.logger.Info("starting server", zap.String("port", s.cfg.Port))
	sc := echo.StartConfig{
		Address:    s.cfg.Port,
		HideBanner: true,
		HidePort:   true,
	}
	return sc.Start(ctx, s.Echo)
}
