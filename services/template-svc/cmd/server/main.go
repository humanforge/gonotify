package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"notification-service/services/template-svc/internal/template"
	"notification-service/services/template-svc/internal/platform/config"
	"notification-service/services/template-svc/internal/platform/httpserver"
	"notification-service/services/template-svc/internal/platform/logging"
	"notification-service/services/template-svc/internal/platform/postgres"

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

	store := template.NewPGStore(postgres.NewDatabase(db.DB))
	svc := template.NewService(store)
	h := template.NewHandler(svc, log)

	srv := httpserver.New(cfg, log, db)
	h.RegisterRoutes(srv.Echo.Group("/api/v1"))

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
