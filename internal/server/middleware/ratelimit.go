package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v5"
	"golang.org/x/time/rate"
)

// ipLimiter associates a token-bucket rate limiter with a timestamp for cleanup.
type ipLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// NewRateLimit creates a per-IP rate limiter with the given rate and window.
// maxRequests is the maximum number of requests allowed per window.
// Each middleware instance keeps its own limiter map and cleanup goroutine.
func NewRateLimit(maxRequests int, window time.Duration) echo.MiddlewareFunc {
	interval := window / time.Duration(maxRequests)

	var (
		mu       sync.Mutex
		limiters = make(map[string]*ipLimiter)
	)

	go func() {
		for {
			time.Sleep(5 * time.Minute)
			mu.Lock()
			for ip, il := range limiters {
				if time.Since(il.lastSeen) > 5*time.Minute {
					delete(limiters, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			ip := c.RealIP()

			mu.Lock()
			il, ok := limiters[ip]
			if !ok {
				il = &ipLimiter{limiter: rate.NewLimiter(rate.Every(interval), maxRequests)}
				limiters[ip] = il
			}
			il.lastSeen = time.Now()
			allowed := il.limiter.Allow()
			mu.Unlock()

			if !allowed {
				return c.JSON(http.StatusTooManyRequests, map[string]string{
					"error": "Too many requests. Please try again later.",
				})
			}
			return next(c)
		}
	}
}
