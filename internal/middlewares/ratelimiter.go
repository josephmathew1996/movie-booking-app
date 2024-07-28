package middlewares

import (
	"movie-booking-app/users-service/internal/ratelimiter"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// RateLimiterMiddleware is an Echo middleware function for rate limiting
func RateLimiterMiddleware(rateLimiter *ratelimiter.FixedWindowClientsRateLimiter, enabled bool, logger *zap.SugaredLogger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if enabled {
				ip := c.RealIP()
				if allowed, retryAfter := rateLimiter.Allow(ip); !allowed {
					logger.Warn("Rate limit exceeded", zap.String("client_ip", c.RealIP()))
					c.Response().Header().Set("Retry-After", retryAfter.String())
					return c.String(http.StatusTooManyRequests, "Rate limit exceeded")
				}
				return next(c)
			}
			// else dont use rate limiter
			return next(c)
		}
	}
}
