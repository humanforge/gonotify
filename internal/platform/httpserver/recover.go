package httpserver

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"go.uber.org/zap"
)

func RecoverMiddleware(logger *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					recErr, ok := r.(error)
					if !ok {
						recErr = echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
					}
					_ = recErr
					logger.Error("panic recovered",
						zap.Any("panic", r),
						zap.String("path", c.Path()),
						zap.String("method", c.Request().Method),
						zap.String("request_id", GetRequestID(c)),
					)
					c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
				}
			}()
			return next(c)
		}
	}
}
