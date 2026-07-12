package middlewares

import (
	"net/http"
	"sync"

	"github.com/labstack/echo/v5"
	"golang.org/x/time/rate"
)

type RateLimiter struct {
	mu      sync.Mutex
	clients map[string]*rate.Limiter
	rate    rate.Limit
	burst   int
}

func NewRateLimiter(r rate.Limit, burst int) *RateLimiter {
	return &RateLimiter{
		clients: make(map[string]*rate.Limiter),
		rate:    r,
		burst:   burst,
	}
}

func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, ok := rl.clients[ip]
	if !ok {
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.clients[ip] = limiter
	}
	return limiter
}

func RateLimitMiddleware(rl *RateLimiter) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			ip := c.RealIP()
			limiter := rl.getLimiter(ip)
			if !limiter.Allow() {
				return echo.NewHTTPError(http.StatusTooManyRequests, "rate limit exceeded")
			}
			return next(c)
		}
	}
}
