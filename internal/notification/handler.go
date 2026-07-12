package notification

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type notificationService interface {
	Send(ctx context.Context, userID, title, body string) error
	GetByID(ctx context.Context, id string) (*Notification, error)
	ListByUser(ctx context.Context, userID string) ([]Notification, error)
	MarkRead(ctx context.Context, id string) error
}

type Handler struct {
	svc notificationService
	log *zap.Logger
}

func NewHandler(svc notificationService, log *zap.Logger) *Handler {
	return &Handler{svc: svc, log: log}
}

func (h *Handler) Create(c *echo.Context) error {
	var req CreateNotificationRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	if req.UserID == "" || req.Title == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "user_id and title are required")
	}

	if err := h.svc.Send(c.Request().Context(), req.UserID, req.Title, req.Body); err != nil {
		h.log.Error("failed to send notification", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to send notification")
	}

	return c.JSON(http.StatusCreated, map[string]string{"status": "sent"})
}

func (h *Handler) GetByID(c *echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "id is required")
	}

	n, err := h.svc.GetByID(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "notification not found")
		}
		h.log.Error("failed to get notification", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get notification")
	}

	return c.JSON(http.StatusOK, toNotificationResponse(n))
}

func (h *Handler) ListByUser(c *echo.Context) error {
	userID := c.QueryParam("user_id")
	if userID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "user_id query parameter is required")
	}

	notifications, err := h.svc.ListByUser(c.Request().Context(), userID)
	if err != nil {
		h.log.Error("failed to list notifications", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to list notifications")
	}

	if notifications == nil {
		notifications = []Notification{}
	}

	items := make([]NotificationResponse, len(notifications))
	for i, n := range notifications {
		items[i] = toNotificationResponse(&n)
	}

	return c.JSON(http.StatusOK, ListNotificationsResponse{Notifications: items})
}

func (h *Handler) MarkRead(c *echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "id is required")
	}

	if err := h.svc.MarkRead(c.Request().Context(), id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "notification not found")
		}
		h.log.Error("failed to mark notification as read", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to mark notification as read")
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "updated"})
}

func (h *Handler) RegisterRoutes(g *echo.Group) {
	g.POST("/notifications", h.Create)
	g.GET("/notifications/:id", h.GetByID)
	g.GET("/notifications", h.ListByUser)
	g.PATCH("/notifications/:id/read", h.MarkRead)
}

func toNotificationResponse(n *Notification) NotificationResponse {
	return NotificationResponse{
		ID:        n.ID,
		UserID:    n.UserID,
		Title:     n.Title,
		Body:      n.Body,
		Read:      n.Read,
		CreatedAt: n.CreatedAt,
	}
}
