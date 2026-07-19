// Package template
package template

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
	"go.uber.org/zap"
)

type Handler struct {
	svc *Service
	log *zap.Logger
}

func NewHandler(svc *Service, log *zap.Logger) *Handler {
	return &Handler{svc: svc, log: log}
}

func (h *Handler) RegisterRoutes(g *echo.Group) {
	g.POST("/templates", h.Create)
	g.PUT("/templates/:id", h.Update)
	g.GET("/templates/:id", h.GetActive)
	g.GET("/templates/:id/versions/:version", h.GetVersion)
	g.GET("/templates", h.List)
	g.DELETE("/templates/:id", h.Deactivate)
}

func (h *Handler) Create(c *echo.Context) error {
	var req CreateTemplateRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	resp, err := h.svc.Create(c.Request().Context(), req)
	if err != nil {
		h.log.Error("create template failed", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "create template failed")
	}

	return c.JSON(http.StatusCreated, resp)
}

func (h *Handler) Update(c *echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "template id is required")
	}

	var req UpdateTemplateRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	resp, err := h.svc.Update(c.Request().Context(), id, req)
	if err != nil {
		if errors.Is(err, ErrTemplateNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "template not found")
		}
		if errors.Is(err, ErrConcurrentUpdate) {
			return echo.NewHTTPError(http.StatusConflict, "concurrent update detected, retry")
		}
		h.log.Error("update template failed", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "update template failed")
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetActive(c *echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "template id is required")
	}

	resp, err := h.svc.GetActive(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, ErrTemplateNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "template not found")
		}
		h.log.Error("get template failed", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "get template failed")
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetVersion(c *echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "template id is required")
	}

	version, err := strconv.Atoi(c.Param("version"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid version")
	}

	resp, err := h.svc.GetVersion(c.Request().Context(), id, version)
	if err != nil {
		if errors.Is(err, ErrTemplateNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "template not found")
		}
		h.log.Error("get template version failed", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "get template version failed")
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) List(c *echo.Context) error {
	limit := parseInt(c.QueryParam("limit"), 20)
	cursor := c.QueryParam("cursor")

	resp, err := h.svc.List(c.Request().Context(), limit, cursor)
	if err != nil {
		h.log.Error("list templates failed", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "list templates failed")
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) Deactivate(c *echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "template id is required")
	}

	if err := h.svc.Deactivate(c.Request().Context(), id); err != nil {
		if errors.Is(err, ErrTemplateNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "template not found")
		}
		h.log.Error("deactivate template failed", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "deactivate template failed")
	}

	return c.NoContent(http.StatusNoContent)
}

func parseInt(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return n
}
