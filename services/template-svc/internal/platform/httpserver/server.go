package httpserver

import (
	"context"

	"notification-service/services/template-svc/internal/platform/config"
	"notification-service/services/template-svc/internal/platform/postgres"

	"github.com/labstack/echo/v5"
	"go.uber.org/zap"
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
		db:     db,
	}
	s.registerMiddlewares()
	s.registerErrorHandler()
	return s
}

func (s *Server) registerMiddlewares() {
	s.Echo.Use(RequestIDMiddleware())
	s.Echo.Use(RequestLoggerMiddleware(s.logger))
}

func (s *Server) registerErrorHandler() {
	s.Echo.HTTPErrorHandler = NewErrorHandler(s.logger, s.cfg.Env == "production")
}

func (s *Server) Start(ctx context.Context) error {
	s.logger.Info("Starting server", zap.String("port", s.cfg.Port))
	sc := echo.StartConfig{
		Address:    s.cfg.Port,
		HideBanner: true,
		HidePort:   true,
	}
	return sc.Start(ctx, s.Echo)
}
