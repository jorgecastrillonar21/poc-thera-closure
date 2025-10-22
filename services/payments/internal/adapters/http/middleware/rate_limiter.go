package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"theraclosure/payments-service/internal/adapters/logging"
)

// RateLimiterConfig holds the configuration for rate limiting
type RateLimiterConfig struct {
	// Requests per window
	RequestsPerWindow int
	// Window duration
	WindowDuration time.Duration
	// Skip rate limiting for these paths
	SkipPaths []string
	// Custom key generator function
	KeyGenerator func(*gin.Context) string
}

// RateLimiter provides rate limiting functionality
type RateLimiter struct {
	config      *RateLimiterConfig
	redis       *redis.Client
	localLimits sync.Map // Fallback for when Redis is unavailable
	logger      *logging.Logger
}

// NewRateLimiter creates a new rate limiter middleware
func NewRateLimiter(config *RateLimiterConfig, redisClient *redis.Client, logger *logging.Logger) *RateLimiter {
	if config.KeyGenerator == nil {
		config.KeyGenerator = defaultKeyGenerator
	}

	return &RateLimiter{
		config: config,
		redis:  redisClient,
		logger: logger,
	}
}

// Middleware returns a Gin middleware function for rate limiting
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip rate limiting for specified paths
		for _, path := range rl.config.SkipPaths {
			if strings.HasPrefix(c.Request.URL.Path, path) {
				c.Next()
				return
			}
		}

		key := rl.config.KeyGenerator(c)
		allowed, remaining, resetTime, err := rl.checkRateLimit(c.Request.Context(), key)

		if err != nil {
			rl.logger.WithFields(logging.Fields{
				{Key: "error", Value: err.Error()},
				{Key: "key", Value: key},
				{Key: "ip", Value: c.ClientIP()},
			}).Error("Rate limiter error")

			// Allow request on error (fail open)
			c.Next()
			return
		}

		// Add rate limit headers
		c.Header("X-RateLimit-Limit", strconv.Itoa(rl.config.RequestsPerWindow))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(resetTime.Unix(), 10))

		if !allowed {
			rl.logger.WithFields(logging.Fields{
				{Key: "key", Value: key},
				{Key: "ip", Value: c.ClientIP()},
				{Key: "path", Value: c.Request.URL.Path},
				{Key: "method", Value: c.Request.Method},
			}).Warn("Rate limit exceeded")

			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "RATE_LIMIT_EXCEEDED",
				"message":     "Too many requests. Please try again later.",
				"retry_after": int(rl.config.WindowDuration.Seconds()),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// checkRateLimit checks if the request is within rate limits
func (rl *RateLimiter) checkRateLimit(ctx context.Context, key string) (allowed bool, remaining int, resetTime time.Time, err error) {
	now := time.Now()
	window := now.Truncate(rl.config.WindowDuration)
	resetTime = window.Add(rl.config.WindowDuration)
	redisKey := fmt.Sprintf("rate_limit:%s:%d", key, window.Unix())

	// Try Redis first
	if rl.redis != nil {
		count, err := rl.redis.Incr(ctx, redisKey).Result()
		if err == nil {
			// Set expiration on first increment
			if count == 1 {
				rl.redis.Expire(ctx, redisKey, rl.config.WindowDuration)
			}

			remaining = rl.config.RequestsPerWindow - int(count)
			if remaining < 0 {
				remaining = 0
			}

			allowed = count <= int64(rl.config.RequestsPerWindow)
			return allowed, remaining, resetTime, nil
		}

		// Log Redis error but continue with local fallback
		rl.logger.WithFields(logging.Fields{
			{Key: "error", Value: err.Error()},
			{Key: "key", Value: key},
		}).Warn("Redis rate limiter failed, falling back to local")
	}

	// Local fallback
	return rl.checkRateLimitLocal(key, window, resetTime)
}

// checkRateLimitLocal provides local rate limiting when Redis is unavailable
func (rl *RateLimiter) checkRateLimitLocal(key string, window time.Time, resetTime time.Time) (bool, int, time.Time, error) {
	localKey := fmt.Sprintf("%s:%d", key, window.Unix())

	val, exists := rl.localLimits.Load(localKey)
	if !exists {
		rl.localLimits.Store(localKey, 1)

		// Clean up expired entries
		go rl.cleanupExpiredEntries(window)

		remaining := rl.config.RequestsPerWindow - 1
		return true, remaining, resetTime, nil
	}

	count, ok := val.(int)
	if !ok {
		count = 0
	}

	count++
	rl.localLimits.Store(localKey, count)

	remaining := rl.config.RequestsPerWindow - count
	if remaining < 0 {
		remaining = 0
	}

	allowed := count <= rl.config.RequestsPerWindow
	return allowed, remaining, resetTime, nil
}

// cleanupExpiredEntries removes expired entries from local storage
func (rl *RateLimiter) cleanupExpiredEntries(currentWindow time.Time) {
	cutoff := currentWindow.Add(-rl.config.WindowDuration).Unix()

	rl.localLimits.Range(func(key, value interface{}) bool {
		keyStr, ok := key.(string)
		if !ok {
			return true
		}

		parts := strings.Split(keyStr, ":")
		if len(parts) < 2 {
			return true
		}

		timestamp, err := strconv.ParseInt(parts[len(parts)-1], 10, 64)
		if err != nil {
			return true
		}

		if timestamp < cutoff {
			rl.localLimits.Delete(key)
		}

		return true
	})
}

// defaultKeyGenerator generates a rate limiting key based on client IP
func defaultKeyGenerator(c *gin.Context) string {
	return c.ClientIP()
}

// UserBasedKeyGenerator generates a rate limiting key based on user ID (if available)
func UserBasedKeyGenerator(c *gin.Context) string {
	// Try to get user ID from context (set by auth middleware)
	if userID, exists := c.Get("user_id"); exists {
		if uid, ok := userID.(string); ok {
			return fmt.Sprintf("user:%s", uid)
		}
	}

	// Fallback to IP-based limiting
	return fmt.Sprintf("ip:%s", c.ClientIP())
}

// EndpointBasedKeyGenerator generates different limits per endpoint
func EndpointBasedKeyGenerator(c *gin.Context) string {
	clientKey := c.ClientIP()
	if userID, exists := c.Get("user_id"); exists {
		if uid, ok := userID.(string); ok {
			clientKey = fmt.Sprintf("user:%s", uid)
		}
	}

	return fmt.Sprintf("%s:%s:%s", clientKey, c.Request.Method, c.FullPath())
}
