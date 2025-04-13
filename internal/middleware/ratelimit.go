package middleware

import (
	"sync"
	"time"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var (
	limiters sync.Map
	rate_limit = 100 // requests per minute
	burst = 5 // maximum burst size
	windowSize = time.Minute
)

// getRateLimiter returns a rate limiter for the given IP address
func getRateLimiter(ip string) *rate.Limiter {
	limiter, exists := limiters.Load(ip)
	if !exists {
		limiter = rate.NewLimiter(rate.Every(windowSize/time.Duration(rate_limit)), burst)
		limiters.Store(ip, limiter)
	}
	return limiter.(*rate.Limiter)
}

// RateLimitMiddleware limits the number of requests per client
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get client IP
		ip := c.ClientIP()
		
		// Get rate limiter for this IP
		limiter := getRateLimiter(ip)

		// Try to allow request
		if !limiter.Allow() {
			c.JSON(429, gin.H{
				"status": "error",
				"message": "Too many requests",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}