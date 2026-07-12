package middlewares

import (
	"github.com/labstack/echo/v5"
)

func SecureHeadersMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			resp := c.Response().(*echo.Response)
			resp.Header().Set("X-Content-Type-Options", "nosniff")
			resp.Header().Set("X-Frame-Options", "DENY")
			resp.Header().Set("X-XSS-Protection", "1; mode=block")
			resp.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			return next(c)
		}
	}
}
