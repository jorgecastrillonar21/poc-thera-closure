package middleware

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"

	"theraclosure/payments-service/internal/adapters/logging"
)

// ValidationConfig holds configuration for request validation
type ValidationConfig struct {
	// Maximum request body size in bytes
	MaxBodySize int64
	// Allowed content types
	AllowedContentTypes []string
	// Required headers for specific endpoints
	RequiredHeaders map[string][]string
	// Blocked user agents (regex patterns)
	BlockedUserAgents []string
	// Allowed HTTP methods
	AllowedMethods []string
}

// RequestValidator provides request validation functionality
type RequestValidator struct {
	config           *ValidationConfig
	logger           *logging.Logger
	blockedUARegexes []*regexp.Regexp
}

// NewRequestValidator creates a new request validator
func NewRequestValidator(config *ValidationConfig, logger *logging.Logger) *RequestValidator {
	validator := &RequestValidator{
		config: config,
		logger: logger,
	}

	// Compile blocked user agent regexes
	for _, pattern := range config.BlockedUserAgents {
		if regex, err := regexp.Compile(pattern); err == nil {
			validator.blockedUARegexes = append(validator.blockedUARegexes, regex)
		} else {
			logger.WithFields(logging.Fields{
				{Key: "pattern", Value: pattern},
				{Key: "error", Value: err.Error()},
			}).Error("Failed to compile user agent regex")
		}
	}

	return validator
}

// Middleware returns a Gin middleware function for request validation
func (rv *RequestValidator) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check HTTP method
		if !rv.isMethodAllowed(c.Request.Method) {
			rv.logger.WithFields(logging.Fields{
				{Key: "method", Value: c.Request.Method},
				{Key: "path", Value: c.Request.URL.Path},
				{Key: "ip", Value: c.ClientIP()},
			}).Warn("Blocked request with disallowed HTTP method")

			c.JSON(http.StatusMethodNotAllowed, gin.H{
				"error":   "METHOD_NOT_ALLOWED",
				"message": "HTTP method not allowed for this endpoint",
			})
			c.Abort()
			return
		}

		// Check user agent
		if rv.isUserAgentBlocked(c.GetHeader("User-Agent")) {
			rv.logger.WithFields(logging.Fields{
				{Key: "user_agent", Value: c.GetHeader("User-Agent")},
				{Key: "ip", Value: c.ClientIP()},
			}).Warn("Blocked request with suspicious user agent")

			c.JSON(http.StatusForbidden, gin.H{
				"error":   "FORBIDDEN",
				"message": "Request blocked",
			})
			c.Abort()
			return
		}

		// Check content type for requests with body
		if rv.hasBody(c.Request.Method) {
			if !rv.isContentTypeAllowed(c.GetHeader("Content-Type")) {
				rv.logger.WithFields(logging.Fields{
					{Key: "content_type", Value: c.GetHeader("Content-Type")},
					{Key: "path", Value: c.Request.URL.Path},
					{Key: "ip", Value: c.ClientIP()},
				}).Warn("Blocked request with invalid content type")

				c.JSON(http.StatusUnsupportedMediaType, gin.H{
					"error":   "UNSUPPORTED_MEDIA_TYPE",
					"message": "Content-Type not supported",
				})
				c.Abort()
				return
			}
		}

		// Check required headers for specific endpoints
		if !rv.hasRequiredHeaders(c) {
			rv.logger.WithFields(logging.Fields{
				{Key: "path", Value: c.Request.URL.Path},
				{Key: "ip", Value: c.ClientIP()},
			}).Warn("Request missing required headers")

			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "MISSING_REQUIRED_HEADERS",
				"message": "Required headers are missing",
			})
			c.Abort()
			return
		}

		// Validate request size (this is handled by RequestSizeMiddleware but we double-check)
		if c.Request.ContentLength > rv.config.MaxBodySize {
			rv.logger.WithFields(logging.Fields{
				{Key: "content_length", Value: c.Request.ContentLength},
				{Key: "max_size", Value: rv.config.MaxBodySize},
				{Key: "ip", Value: c.ClientIP()},
			}).Warn("Request body too large")

			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"error":   "REQUEST_TOO_LARGE",
				"message": "Request body exceeds maximum allowed size",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// isMethodAllowed checks if the HTTP method is allowed
func (rv *RequestValidator) isMethodAllowed(method string) bool {
	if len(rv.config.AllowedMethods) == 0 {
		return true // Allow all if not configured
	}

	for _, allowed := range rv.config.AllowedMethods {
		if strings.EqualFold(method, allowed) {
			return true
		}
	}
	return false
}

// isUserAgentBlocked checks if the user agent is blocked
func (rv *RequestValidator) isUserAgentBlocked(userAgent string) bool {
	if userAgent == "" {
		return false // Allow empty user agents
	}

	for _, regex := range rv.blockedUARegexes {
		if regex.MatchString(userAgent) {
			return true
		}
	}
	return false
}

// isContentTypeAllowed checks if the content type is allowed
func (rv *RequestValidator) isContentTypeAllowed(contentType string) bool {
	if len(rv.config.AllowedContentTypes) == 0 {
		return true // Allow all if not configured
	}

	// Extract the main content type (ignore charset, boundary, etc.)
	mainType := strings.Split(contentType, ";")[0]
	mainType = strings.TrimSpace(mainType)

	for _, allowed := range rv.config.AllowedContentTypes {
		if strings.EqualFold(mainType, allowed) {
			return true
		}
	}
	return false
}

// hasRequiredHeaders checks if all required headers are present
func (rv *RequestValidator) hasRequiredHeaders(c *gin.Context) bool {
	path := c.Request.URL.Path

	// Check for exact path match first
	if headers, exists := rv.config.RequiredHeaders[path]; exists {
		return rv.checkHeaders(c, headers)
	}

	// Check for pattern matches
	for pattern, headers := range rv.config.RequiredHeaders {
		if matched, _ := regexp.MatchString(pattern, path); matched {
			return rv.checkHeaders(c, headers)
		}
	}

	return true // No required headers for this path
}

// checkHeaders verifies that all required headers are present
func (rv *RequestValidator) checkHeaders(c *gin.Context, requiredHeaders []string) bool {
	for _, header := range requiredHeaders {
		if c.GetHeader(header) == "" {
			return false
		}
	}
	return true
}

// hasBody checks if the HTTP method typically has a request body
func (rv *RequestValidator) hasBody(method string) bool {
	bodyMethods := []string{"POST", "PUT", "PATCH"}
	for _, m := range bodyMethods {
		if strings.EqualFold(method, m) {
			return true
		}
	}
	return false
}

// DefaultValidationConfig returns a sensible default configuration
func DefaultValidationConfig() *ValidationConfig {
	return &ValidationConfig{
		MaxBodySize: 10 * 1024 * 1024, // 10MB
		AllowedContentTypes: []string{
			"application/json",
			"application/x-www-form-urlencoded",
			"multipart/form-data",
		},
		AllowedMethods: []string{
			"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD",
		},
		BlockedUserAgents: []string{
			"(?i)bot",
			"(?i)crawler",
			"(?i)spider",
			"(?i)scraper",
			"curl/7\\.[0-4]\\.", // Block old curl versions
		},
		RequiredHeaders: map[string][]string{
			"/api/v1/webhooks/stripe": {"Stripe-Signature"},
		},
	}
}
