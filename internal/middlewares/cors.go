package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func CORSMiddleware(allowedOrigins []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			origin := c.Request().Header.Get("Origin")
			if origin != "" {
				allowed := false
				for _, o := range allowedOrigins {
					if o == "*" || o == origin {
						allowed = true
						break
					}
				}
				if allowed {
					resp := c.Response().(*echo.Response)
					if allowedOrigins[0] == "*" {
						resp.Header().Set("Access-Control-Allow-Origin", "*")
					} else {
						resp.Header().Set("Access-Control-Allow-Origin", origin)
						resp.Header().Set("Vary", "Origin")
					}
					resp.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
					resp.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-Request-ID")
					resp.Header().Set("Access-Control-Max-Age", "86400")
				}
			}

			if c.Request().Method == http.MethodOptions {
				return c.NoContent(http.StatusNoContent)
			}

			return next(c)
		}
	}
}
