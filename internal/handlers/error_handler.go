package handlers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"
	"go.uber.org/zap"
)

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error APIError `json:"error"`
}

func NewErrorHandler(log *zap.Logger, isProduction bool) echo.HTTPErrorHandler {
	return func(c *echo.Context, err error) {
		he := &echo.HTTPError{}
		if !errors.As(err, &he) {
			he = echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
		}

		code := he.Code
		message := he.Message

		if isProduction && code >= 500 {
			message = "internal server error"
		}

		log.Warn("request error",
			zap.Int("status", code),
			zap.Error(err),
			zap.String("path", c.Path()),
			zap.String("request_id", GetRequestID(c)),
		)

		resp := c.Response().(*echo.Response)
		if !resp.Committed {
			if err := c.JSON(code, ErrorResponse{Error: APIError{Code: code, Message: message}}); err != nil {
				log.Error("failed to write error response", zap.Error(err))
			}
		}
	}
}

func GetRequestID(c *echo.Context) string {
	if id, ok := c.Get("request_id").(string); ok {
		return id
	}
	return ""
}
