package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"theraclosure/payments-service/internal/adapters/logging"
)

func TestRateLimiter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create test logger
	logger := logging.NewLogger("debug")

	// Create rate limiter with very low limits for testing
	config := &RateLimiterConfig{
		RequestsPerWindow: 2,
		WindowDuration:    time.Second,
		SkipPaths:         []string{"/health"},
	}

	rateLimiter := NewRateLimiter(config, nil, logger) // No Redis for test

	// Create test router
	router := gin.New()
	router.Use(rateLimiter.Middleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Test allowed requests
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	// Test rate limited request
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusTooManyRequests, w.Code)

	// Test skipped path
	req = httptest.NewRequest("GET", "/health", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRequestValidator(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create test logger
	logger := logging.NewLogger("debug")

	// Create validator with test config
	config := &ValidationConfig{
		MaxBodySize:         100, // Very small for testing
		AllowedContentTypes: []string{"application/json"},
		AllowedMethods:      []string{"GET", "POST"},
		BlockedUserAgents:   []string{"(?i)bot"},
	}

	validator := NewRequestValidator(config, logger)

	// Create test router
	router := gin.New()
	router.Use(validator.Middleware())
	router.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Test allowed request
	req := httptest.NewRequest("POST", "/test", bytes.NewReader([]byte(`{"test": true}`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "test-client")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Test blocked user agent
	req = httptest.NewRequest("POST", "/test", bytes.NewReader([]byte(`{"test": true}`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "test-bot")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	// Test invalid content type
	req = httptest.NewRequest("POST", "/test", bytes.NewReader([]byte(`{"test": true}`)))
	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("User-Agent", "test-client")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnsupportedMediaType, w.Code)

	// Test disallowed method
	req = httptest.NewRequest("DELETE", "/test", nil)
	req.Header.Set("User-Agent", "test-client")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestWebhookVerifier(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create test logger
	logger := logging.NewLogger("debug")

	// Test webhook secret
	secret := "test_webhook_secret"
	verifier := NewWebhookVerifier(secret, logger)

	// Create test router
	router := gin.New()
	router.Use(verifier.VerifyStripeSignature())
	router.POST("/api/v1/webhooks/stripe", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "webhook processed"})
	})
	router.POST("/api/v1/other", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "other endpoint"})
	})

	// Test non-webhook endpoint (should skip verification)
	req := httptest.NewRequest("POST", "/api/v1/other", bytes.NewReader([]byte(`{"test": true}`)))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Test webhook endpoint without signature (should fail)
	req = httptest.NewRequest("POST", "/api/v1/webhooks/stripe", bytes.NewReader([]byte(`{"test": true}`)))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthenticator(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create test logger
	logger := logging.NewLogger("debug")

	// Create auth config
	config := &AuthConfig{
		JWTSecret: "test-secret",
		APIKeys: map[string]string{
			"test-service": "test-api-key",
		},
		BasicAuthUsers: map[string]string{
			"testuser": "testpass",
		},
		SkipPaths: []string{"/health"},
	}

	authenticator := NewAuthenticator(config, logger)

	// Create test router with optional auth
	router := gin.New()
	router.Use(authenticator.OptionalAuthMiddleware())
	router.GET("/test", func(c *gin.Context) {
		authenticated, _ := c.Get("authenticated")
		c.JSON(http.StatusOK, gin.H{"authenticated": authenticated})
	})
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Test unauthenticated request (should pass with optional auth)
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Test API key authentication
	req = httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-API-Key", "test-api-key")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Test skipped path
	req = httptest.NewRequest("GET", "/health", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
