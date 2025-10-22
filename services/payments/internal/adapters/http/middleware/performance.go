package middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	"theraclosure/payments-service/internal/adapters/logging"
	"theraclosure/payments-service/internal/adapters/monitoring"
)

// PerformanceMonitor provides performance monitoring middleware
type PerformanceMonitor struct {
	metrics *monitoring.Metrics
	logger  *logging.Logger
}

// NewPerformanceMonitor creates a new performance monitoring middleware
func NewPerformanceMonitor(metrics *monitoring.Metrics, logger *logging.Logger) *PerformanceMonitor {
	return &PerformanceMonitor{
		metrics: metrics,
		logger:  logger,
	}
}

// Middleware returns a Gin middleware function for performance monitoring
func (pm *PerformanceMonitor) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Increment active requests
		pm.metrics.IncActiveRequests()

		// Process request
		c.Next()

		// Decrement active requests
		pm.metrics.DecActiveRequests()

		// Calculate metrics
		duration := time.Since(start)
		endpoint := c.FullPath()
		if endpoint == "" {
			endpoint = c.Request.URL.Path
		}

		method := c.Request.Method
		statusCode := c.Writer.Status()
		responseSize := int64(c.Writer.Size())

		// Record HTTP metrics
		pm.metrics.RecordHTTPRequest(method, endpoint, statusCode, duration, responseSize)

		// Log performance data for slow requests
		if duration > 1*time.Second {
			pm.logger.WithFields(logging.Fields{
				{Key: "method", Value: method},
				{Key: "endpoint", Value: endpoint},
				{Key: "status_code", Value: statusCode},
				{Key: "duration_ms", Value: duration.Milliseconds()},
				{Key: "response_size", Value: responseSize},
				{Key: "ip", Value: c.ClientIP()},
			}).Warn("Slow request detected")
		}

		// Add performance headers to response
		c.Header("X-Response-Time", duration.String())
		c.Header("X-Request-ID", c.GetHeader("X-Request-ID"))
	}
}

// BusinessEventMiddleware creates middleware for tracking business events
func (pm *PerformanceMonitor) BusinessEventMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Track business events based on endpoint and status
		endpoint := c.FullPath()
		statusCode := c.Writer.Status()

		if statusCode >= 200 && statusCode < 300 {
			switch endpoint {
			case "/api/v1/customers":
				if c.Request.Method == "POST" {
					pm.metrics.RecordBusinessEvent("customer_created", map[string]string{
						"status": "success",
					})
				}
			case "/api/v1/subscriptions":
				if c.Request.Method == "POST" {
					pm.metrics.RecordBusinessEvent("subscription_created", map[string]string{
						"status": "success",
					})
				}
			case "/api/v1/payments":
				if c.Request.Method == "POST" {
					pm.metrics.RecordBusinessEvent("payment_created", map[string]string{
						"status":   "success",
						"currency": "usd",  // Could be extracted from request
						"amount":   "1000", // Could be extracted from request
					})
				}
			}
		} else if statusCode >= 400 {
			switch endpoint {
			case "/api/v1/customers":
				if c.Request.Method == "POST" {
					pm.metrics.RecordBusinessEvent("customer_created", map[string]string{
						"status": "failed",
					})
				}
			case "/api/v1/subscriptions":
				if c.Request.Method == "POST" {
					pm.metrics.RecordBusinessEvent("subscription_created", map[string]string{
						"status": "failed",
					})
				}
			case "/api/v1/payments":
				if c.Request.Method == "POST" {
					pm.metrics.RecordBusinessEvent("payment_created", map[string]string{
						"status":   "failed",
						"currency": "usd",
					})
				}
			}
		}
	}
}

// SecurityMetricsMiddleware tracks security-related metrics
func (pm *PerformanceMonitor) SecurityMetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for authentication
		if authHeader := c.GetHeader("Authorization"); authHeader != "" {
			var method string
			if c.GetHeader("X-API-Key") != "" {
				method = "api_key"
			} else if authHeader[:6] == "Bearer" {
				method = "jwt"
			} else if authHeader[:5] == "Basic" {
				method = "basic"
			}

			// Record authentication attempt before processing
			pm.metrics.RecordSecurityEvent("authentication_attempt", map[string]string{
				"method": method,
				"status": "attempted",
			})
		}

		c.Next()

		// Record final authentication status
		if authHeader := c.GetHeader("Authorization"); authHeader != "" {
			var method string
			if c.GetHeader("X-API-Key") != "" {
				method = "api_key"
			} else if authHeader[:6] == "Bearer" {
				method = "jwt"
			} else if authHeader[:5] == "Basic" {
				method = "basic"
			}

			status := "success"
			if c.Writer.Status() == 401 || c.Writer.Status() == 403 {
				status = "failed"
			}

			pm.metrics.RecordSecurityEvent("authentication_attempt", map[string]string{
				"method": method,
				"status": status,
			})
		}

		// Record rate limit hits
		if c.Writer.Status() == 429 {
			endpoint := c.FullPath()
			clientType := "anonymous"
			if userID, exists := c.Get("user_id"); exists && userID != nil {
				clientType = "authenticated"
			}

			pm.metrics.RecordSecurityEvent("rate_limit_hit", map[string]string{
				"endpoint":    endpoint,
				"client_type": clientType,
			})
		}

		// Record security violations
		if c.Writer.Status() == 400 || c.Writer.Status() == 403 {
			endpoint := c.FullPath()
			violationType := "validation_error"
			if c.Writer.Status() == 403 {
				violationType = "access_denied"
			}

			pm.metrics.RecordSecurityEvent("security_violation", map[string]string{
				"violation_type": violationType,
				"endpoint":       endpoint,
			})
		}
	}
}
