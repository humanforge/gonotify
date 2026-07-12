package middlewares

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/labstack/echo/v5"
)

type contextKey string

const RequestIDKey contextKey = "request_id"

func generateRequestID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func RequestIDMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			req := c.Request()
			id := req.Header.Get("X-Request-ID")
			if id == "" {
				id = generateRequestID()
			}
			c.Set(string(RequestIDKey), id)
			resp := c.Response().(*echo.Response)
			resp.Header().Set("X-Request-ID", id)
			return next(c)
		}
	}
}

func GetRequestID(c *echo.Context) string {
	if id, ok := c.Get(string(RequestIDKey)).(string); ok {
		return id
	}
	return ""
}
