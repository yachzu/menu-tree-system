package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/xos/Projects/stk-project/backend/internal/dto"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "internal_error",
				Message: err.Error(),
			})
		}
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		c.Next()
	}
}

type visitor struct {
	tokens int
	last   time.Time
}

type rateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	rate     int
	burst    int
}

func newRateLimiter(rate, burst int) *rateLimiter {
	rl := &rateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		burst:    burst,
	}

	go func() {
		for {
			time.Sleep(5 * time.Minute)
			rl.mu.Lock()
			for ip, v := range rl.visitors {
				if time.Since(v.last) > 10*time.Minute {
					delete(rl.visitors, ip)
				}
			}
			rl.mu.Unlock()
		}
	}()

	return rl
}

func (rl *rateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	v, exists := rl.visitors[ip]
	if !exists {
		rl.visitors[ip] = &visitor{tokens: rl.burst - 1, last: now}
		return true
	}

	elapsed := now.Sub(v.last)
	v.tokens += int(elapsed.Seconds()) * rl.rate
	if v.tokens > rl.burst {
		v.tokens = rl.burst
	}
	v.last = now

	if v.tokens <= 0 {
		return false
	}
	v.tokens--
	return true
}

func RateLimiter() gin.HandlerFunc {
	rl := newRateLimiter(10, 20)
	return func(c *gin.Context) {
		if !rl.allow(c.ClientIP()) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, dto.ErrorResponse{
				Error:   "rate_limit_exceeded",
				Message: "too many requests, please try again later",
			})
			return
		}
		c.Next()
	}
}
