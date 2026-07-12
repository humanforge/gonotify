package httpserver

import (
	"net/http"

	"notification-service/internal/platform/postgres"

	"github.com/labstack/echo/v5"
	"go.uber.org/zap"
)

type HealthHandler struct {
	db  *postgres.DBConnection
	log *zap.Logger
}

func NewHealthHandler(db *postgres.DBConnection, log *zap.Logger) *HealthHandler {
	return &HealthHandler{db: db, log: log}
}

func (h *HealthHandler) Healthz(c *echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func (h *HealthHandler) Readyz(c *echo.Context) error {
	if err := h.db.Ping(); err != nil {
		h.log.Warn("readiness check failed", zap.Error(err))
		return c.JSON(http.StatusServiceUnavailable, map[string]string{"status": "unavailable"})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}
