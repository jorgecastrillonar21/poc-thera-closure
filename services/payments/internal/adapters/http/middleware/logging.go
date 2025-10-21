package middleware

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"theraclosure/payments-service/internal/adapters/errors"
	"theraclosure/payments-service/internal/adapters/logging"
)

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if request ID is already present (from headers)
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Add to context
		ctx := context.WithValue(c.Request.Context(), "request_id", requestID)
		c.Request = c.Request.WithContext(ctx)

		// Add to response headers
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}

// LoggingMiddleware provides structured request/response logging
func LoggingMiddleware(logger *logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		method := c.Request.Method

		// Capture request body for logging (if needed)
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// Create response writer wrapper to capture response
		writer := &responseWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = writer

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Get request ID from context
		requestID := ""
		if id := c.Request.Context().Value("request_id"); id != nil {
			requestID = id.(string)
		}

		// Build log fields
		fields := logging.Fields{
			{Key: "method", Value: method},
			{Key: "path", Value: path},
			{Key: "query", Value: raw},
			{Key: "status_code", Value: c.Writer.Status()},
			{Key: "duration_ms", Value: duration.Milliseconds()},
			{Key: "request_size", Value: len(bodyBytes)},
			{Key: "response_size", Value: c.Writer.Size()},
			{Key: "user_agent", Value: c.GetHeader("User-Agent")},
			{Key: "remote_addr", Value: c.ClientIP()},
			{Key: "request_id", Value: requestID},
		}

		// Add user context if available
		if userID := c.GetString("user_id"); userID != "" {
			fields = append(fields, logging.Field{Key: "user_id", Value: userID})
		}

		// Log based on status code
		logEntry := logger.WithFields(fields)

		if c.Writer.Status() >= 500 {
			logEntry.Error("Request processed with server error")
		} else if c.Writer.Status() >= 400 {
			logEntry.Warn("Request processed with client error")
		} else {
			logEntry.Info("Request processed successfully")
		}

		// Log slow requests separately
		if duration > 1*time.Second {
			logger.LogPerformanceMetric(
				c.Request.Context(),
				"slow_request",
				duration.Seconds(),
				"seconds",
				logging.Fields{
					{Key: "endpoint", Value: path},
					{Key: "method", Value: method},
					{Key: "status_code", Value: c.Writer.Status()},
				},
			)
		}
	}
}

// ErrorHandlingMiddleware provides centralized error handling
func ErrorHandlingMiddleware(logger *logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Handle any errors that occurred during request processing
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// Get request ID from context
			requestID := ""
			if id := c.Request.Context().Value("request_id"); id != nil {
				requestID = id.(string)
			}

			var serviceErr *errors.ServiceError

			// Check if it's already a ServiceError
			if se, ok := err.(*errors.ServiceError); ok {
				serviceErr = se.WithRequestID(requestID)
			} else {
				// Wrap unknown errors as internal errors
				serviceErr = errors.NewServiceError(
					errors.ErrCodeInternal,
					"Internal server error",
				).WithCause(err).WithRequestID(requestID)
			}

			// Log the error using the logging package method
			errorFields := logging.Fields{
				{Key: "error_code", Value: string(serviceErr.Code)},
				{Key: "http_status", Value: serviceErr.HTTPStatus},
				{Key: "endpoint", Value: c.Request.URL.Path},
				{Key: "method", Value: c.Request.Method},
			}
			logger.WithFields(errorFields).WithError(err).Error("Request failed with error")

			// Return structured error response
			c.JSON(serviceErr.HTTPStatus, serviceErr.ToErrorResponse())
			return
		}
	}
}

// RecoveryMiddleware provides panic recovery with structured logging
func RecoveryMiddleware(logger *logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recovery := recover(); recovery != nil {
				// Get request ID from context
				requestID := ""
				if id := c.Request.Context().Value("request_id"); id != nil {
					requestID = id.(string)
				}

				// Log the panic using the logging package method
				panicFields := logging.Fields{
					{Key: "panic", Value: recovery},
					{Key: "endpoint", Value: c.Request.URL.Path},
					{Key: "method", Value: c.Request.Method},
					{Key: "request_id", Value: requestID},
				}
				logger.WithFields(panicFields).Error("Panic recovered")

				// Create service error for panic
				serviceErr := errors.NewServiceError(
					errors.ErrCodeInternal,
					"Internal server error",
				).WithRequestID(requestID).
					WithDetails("An unexpected error occurred")

				c.JSON(http.StatusInternalServerError, serviceErr.ToErrorResponse())
				c.Abort()
			}
		}()

		c.Next()
	}
}

// SecurityLoggingMiddleware logs security-related events
func SecurityLoggingMiddleware(logger *logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for suspicious patterns
		userAgent := c.GetHeader("User-Agent")
		if isSuspiciousUserAgent(userAgent) {
			logger.LogSecurityEvent(
				c.Request.Context(),
				"suspicious_user_agent",
				"medium",
				logging.Fields{
					{Key: "user_agent", Value: userAgent},
					{Key: "remote_addr", Value: c.ClientIP()},
					{Key: "endpoint", Value: c.Request.URL.Path},
				},
			)
		}

		// Log authentication attempts
		if authHeader := c.GetHeader("Authorization"); authHeader != "" {
			logger.LogSecurityEvent(
				c.Request.Context(),
				"authentication_attempt",
				"low",
				logging.Fields{
					{Key: "has_auth_header", Value: true},
					{Key: "remote_addr", Value: c.ClientIP()},
				},
			)
		}

		c.Next()
	}
}

// responseWriter wraps gin.ResponseWriter to capture response body
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Helper functions

func isSuspiciousUserAgent(userAgent string) bool {
	suspiciousPatterns := []string{
		"sqlmap",
		"nikto",
		"nmap",
		"masscan",
		"<script>",
		"curl", // Might be too broad for production
	}

	for _, pattern := range suspiciousPatterns {
		if contains(userAgent, pattern) {
			return true
		}
	}

	return false
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			(len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr ||
					containsMiddle(s, substr))))
}

func containsMiddle(s, substr string) bool {
	for i := 1; i < len(s)-len(substr)+1; i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// RequestSizeMiddleware limits request body size
func RequestSizeMiddleware(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > maxSize {
			serviceErr := errors.NewServiceError(
				errors.ErrCodeValidation,
				"Request body too large",
			).WithDetails("Request body exceeds maximum allowed size of " + strconv.FormatInt(maxSize, 10) + " bytes")

			c.JSON(serviceErr.HTTPStatus, serviceErr.ToErrorResponse())
			c.Abort()
			return
		}

		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)
		c.Next()
	}
}
