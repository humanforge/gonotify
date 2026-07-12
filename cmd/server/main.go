package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"notification-service/internal/notification"
	"notification-service/internal/platform/config"
	"notification-service/internal/platform/httpserver"
	"notification-service/internal/platform/logging"
	"notification-service/internal/platform/postgres"

	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}

	log, err := logging.New(cfg)
	if err != nil {
		os.Stderr.WriteString("failed to initialize logger: " + err.Error() + "\n")
		os.Exit(1)
	}

	db, err := postgres.NewDBConnection(cfg)
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}

	notifStore := notification.NewStore(db.DB)
	notifSvc := notification.NewService(notifStore)
	notifHandler := notification.NewHandler(notifSvc, log)

	srv := httpserver.New(cfg, log, db)
	notifHandler.RegisterRoutes(srv.Echo.Group("/api/v1"))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := srv.Start(ctx); err != nil {
		log.Error("server error", zap.Error(err))
	}

	if err := db.Close(); err != nil {
		log.Error("db close error", zap.Error(err))
	}

	_ = log.Sync()
}
