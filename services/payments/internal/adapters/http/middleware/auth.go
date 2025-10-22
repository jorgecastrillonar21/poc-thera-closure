package middleware

import (
	"crypto/subtle"
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"theraclosure/payments-service/internal/adapters/logging"
)

// AuthConfig holds authentication configuration
type AuthConfig struct {
	// JWT secret for token verification
	JWTSecret string
	// API keys for service-to-service authentication
	APIKeys map[string]string
	// Basic auth credentials for simple authentication
	BasicAuthUsers map[string]string
	// Skip authentication for these paths
	SkipPaths []string
	// Token expiration time
	TokenExpiration time.Duration
}

// Authenticator provides authentication functionality
type Authenticator struct {
	config *AuthConfig
	logger *logging.Logger
}

// NewAuthenticator creates a new authenticator
func NewAuthenticator(config *AuthConfig, logger *logging.Logger) *Authenticator {
	return &Authenticator{
		config: config,
		logger: logger,
	}
}

// JWTMiddleware provides JWT token authentication
func (a *Authenticator) JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip authentication for specified paths
		if a.shouldSkipAuth(c.Request.URL.Path) {
			c.Next()
			return
		}

		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			a.logUnauthorized(c, "Missing Authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "MISSING_AUTHORIZATION",
				"message": "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Extract bearer token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			a.logUnauthorized(c, "Invalid Authorization format")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "INVALID_AUTHORIZATION_FORMAT",
				"message": "Authorization must be Bearer token",
			})
			c.Abort()
			return
		}

		// Verify JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(a.config.JWTSecret), nil
		})

		if err != nil {
			a.logUnauthorized(c, "Invalid JWT token: "+err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "INVALID_TOKEN",
				"message": "JWT token verification failed",
			})
			c.Abort()
			return
		}

		// Extract claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Set user information in context
			if userID, exists := claims["user_id"]; exists {
				c.Set("user_id", userID)
			}
			if role, exists := claims["role"]; exists {
				c.Set("user_role", role)
			}
			if customerID, exists := claims["customer_id"]; exists {
				c.Set("customer_id", customerID)
			}

			a.logger.WithFields(logging.Fields{
				{Key: "user_id", Value: claims["user_id"]},
				{Key: "path", Value: c.Request.URL.Path},
			}).Debug("JWT authentication successful")

			c.Next()
		} else {
			a.logUnauthorized(c, "Invalid JWT claims")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "INVALID_CLAIMS",
				"message": "JWT token claims are invalid",
			})
			c.Abort()
		}
	}
}

// APIKeyMiddleware provides API key authentication
func (a *Authenticator) APIKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip authentication for specified paths
		if a.shouldSkipAuth(c.Request.URL.Path) {
			c.Next()
			return
		}

		// Check for API key in header
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			// Also check Authorization header for API key
			authHeader := c.GetHeader("Authorization")
			if strings.HasPrefix(authHeader, "ApiKey ") {
				apiKey = strings.TrimPrefix(authHeader, "ApiKey ")
			}
		}

		if apiKey == "" {
			a.logUnauthorized(c, "Missing API key")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "MISSING_API_KEY",
				"message": "API key is required",
			})
			c.Abort()
			return
		}

		// Verify API key
		serviceName := a.verifyAPIKey(apiKey)
		if serviceName == "" {
			a.logUnauthorized(c, "Invalid API key")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "INVALID_API_KEY",
				"message": "API key verification failed",
			})
			c.Abort()
			return
		}

		// Set service information in context
		c.Set("service_name", serviceName)
		c.Set("authenticated_via", "api_key")

		a.logger.WithFields(logging.Fields{
			{Key: "service_name", Value: serviceName},
			{Key: "path", Value: c.Request.URL.Path},
		}).Debug("API key authentication successful")

		c.Next()
	}
}

// BasicAuthMiddleware provides basic authentication
func (a *Authenticator) BasicAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip authentication for specified paths
		if a.shouldSkipAuth(c.Request.URL.Path) {
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Basic ") {
			a.requestBasicAuth(c)
			return
		}

		// Decode base64 credentials
		encoded := strings.TrimPrefix(authHeader, "Basic ")
		decoded, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			a.logUnauthorized(c, "Invalid base64 encoding")
			a.requestBasicAuth(c)
			return
		}

		// Extract username and password
		credentials := strings.SplitN(string(decoded), ":", 2)
		if len(credentials) != 2 {
			a.logUnauthorized(c, "Invalid credentials format")
			a.requestBasicAuth(c)
			return
		}

		username, password := credentials[0], credentials[1]

		// Verify credentials
		if !a.verifyBasicAuth(username, password) {
			a.logUnauthorized(c, "Invalid basic auth credentials")
			a.requestBasicAuth(c)
			return
		}

		// Set user information in context
		c.Set("username", username)
		c.Set("authenticated_via", "basic_auth")

		a.logger.WithFields(logging.Fields{
			{Key: "username", Value: username},
			{Key: "path", Value: c.Request.URL.Path},
		}).Debug("Basic authentication successful")

		c.Next()
	}
}

// OptionalAuthMiddleware provides optional authentication (doesn't block unauthenticated requests)
func (a *Authenticator) OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		// Try JWT first
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(a.config.JWTSecret), nil
			})

			if err == nil {
				if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
					c.Set("user_id", claims["user_id"])
					c.Set("authenticated", true)
				}
			}
		}

		// Try API key
		apiKey := c.GetHeader("X-API-Key")
		if apiKey != "" {
			serviceName := a.verifyAPIKey(apiKey)
			if serviceName != "" {
				c.Set("service_name", serviceName)
				c.Set("authenticated", true)
			}
		}

		c.Next()
	}
}

// shouldSkipAuth checks if authentication should be skipped for the given path
func (a *Authenticator) shouldSkipAuth(path string) bool {
	for _, skipPath := range a.config.SkipPaths {
		if strings.HasPrefix(path, skipPath) {
			return true
		}
	}
	return false
}

// verifyAPIKey verifies an API key and returns the associated service name
func (a *Authenticator) verifyAPIKey(apiKey string) string {
	for serviceName, validKey := range a.config.APIKeys {
		if subtle.ConstantTimeCompare([]byte(apiKey), []byte(validKey)) == 1 {
			return serviceName
		}
	}
	return ""
}

// verifyBasicAuth verifies basic auth credentials
func (a *Authenticator) verifyBasicAuth(username, password string) bool {
	validPassword, exists := a.config.BasicAuthUsers[username]
	if !exists {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(password), []byte(validPassword)) == 1
}

// requestBasicAuth sends a 401 response with WWW-Authenticate header
func (a *Authenticator) requestBasicAuth(c *gin.Context) {
	c.Header("WWW-Authenticate", "Basic realm=\"Payments Service\"")
	c.JSON(http.StatusUnauthorized, gin.H{
		"error":   "AUTHENTICATION_REQUIRED",
		"message": "Basic authentication is required",
	})
	c.Abort()
}

// logUnauthorized logs unauthorized access attempts
func (a *Authenticator) logUnauthorized(c *gin.Context, reason string) {
	a.logger.WithFields(logging.Fields{
		{Key: "ip", Value: c.ClientIP()},
		{Key: "path", Value: c.Request.URL.Path},
		{Key: "method", Value: c.Request.Method},
		{Key: "user_agent", Value: c.GetHeader("User-Agent")},
		{Key: "reason", Value: reason},
	}).Warn("Unauthorized access attempt")
}

// DefaultAuthConfig returns a default authentication configuration
func DefaultAuthConfig() *AuthConfig {
	return &AuthConfig{
		JWTSecret:       "your-jwt-secret-key", // Should be loaded from environment
		APIKeys:         make(map[string]string),
		BasicAuthUsers:  make(map[string]string),
		TokenExpiration: 24 * time.Hour,
		SkipPaths: []string{
			"/health",
			"/swagger",
			"/api/v1/webhooks", // Webhooks use different auth
		},
	}
}
