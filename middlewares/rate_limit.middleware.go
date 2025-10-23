package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type RateLimiter struct {
	visitors map[string]*Visitor
	mu       sync.RWMutex
	rate     int           // requests per window
	window   time.Duration // time window
}

type Visitor struct {
	lastSeen  time.Time
	count     int
	windowEnd time.Time
}

// NewRateLimiter creates a new rate limiter
// rate: number of requests allowed per window
// window: time window duration (e.g., 1 minute)
func NewRateLimiter(rate int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*Visitor),
		rate:     rate,
		window:   window,
	}

	// Clean up old visitors every 5 minutes
	go rl.cleanupVisitors()

	return rl
}

// RateLimitMiddleware creates a rate limiting middleware
func (rl *RateLimiter) RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		rl.mu.Lock()
		visitor, exists := rl.visitors[ip]
		now := time.Now()

		if !exists {
			// New visitor
			rl.visitors[ip] = &Visitor{
				lastSeen:  now,
				count:     1,
				windowEnd: now.Add(rl.window),
			}
			rl.mu.Unlock()
			c.Next()
			return
		}

		// Check if window has expired
		if now.After(visitor.windowEnd) {
			// Reset counter for new window
			visitor.count = 1
			visitor.windowEnd = now.Add(rl.window)
			visitor.lastSeen = now
			rl.mu.Unlock()
			c.Next()
			return
		}

		// Within window - check rate limit
		if visitor.count >= rl.rate {
			rl.mu.Unlock()
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Rate limit exceeded",
				"message": "Too many requests. Please try again later.",
			})
			c.Abort()
			return
		}

		// Increment counter
		visitor.count++
		visitor.lastSeen = now
		rl.mu.Unlock()

		c.Next()
	}
}

// cleanupVisitors removes old visitors to prevent memory leaks
func (rl *RateLimiter) cleanupVisitors() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, visitor := range rl.visitors {
			// Remove visitors not seen in the last 10 minutes
			if now.Sub(visitor.lastSeen) > 10*time.Minute {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// Global rate limiters with different configurations
var (
	// Strict rate limiter: 10 requests per minute (for sensitive endpoints like login)
	StrictRateLimiter = NewRateLimiter(10, 1*time.Minute)

	// Normal rate limiter: 100 requests per minute (for regular API endpoints)
	NormalRateLimiter = NewRateLimiter(100, 1*time.Minute)

	// Generous rate limiter: 300 requests per minute (for read-heavy endpoints)
	GenerousRateLimiter = NewRateLimiter(300, 1*time.Minute)
)

// Convenience functions for middleware
func StrictRateLimit() gin.HandlerFunc {
	return StrictRateLimiter.RateLimitMiddleware()
}

func NormalRateLimit() gin.HandlerFunc {
	return NormalRateLimiter.RateLimitMiddleware()
}

func GenerousRateLimit() gin.HandlerFunc {
	return GenerousRateLimiter.RateLimitMiddleware()
}
