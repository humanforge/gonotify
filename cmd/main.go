package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"notification-service/internal/config"
	"notification-service/internal/database"
	"notification-service/internal/logger"
	"notification-service/internal/server"

	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}

	log, err := logger.New(cfg)
	if err != nil {
		os.Stderr.WriteString("failed to initialize logger: " + err.Error() + "\n")
		os.Exit(1)
	}

	db, err := database.NewDBConnection(cfg)
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}

	srv := server.New(cfg, log, db)

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
