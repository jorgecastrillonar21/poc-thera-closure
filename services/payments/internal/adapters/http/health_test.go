package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"theraclosure/payments-service/internal/adapters/config"
	"theraclosure/payments-service/internal/adapters/logging"
)

func TestHealthCheckEndpoints_WithDatabase(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create in-memory SQLite database for testing (without UUID dependencies)
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Create test config
	cfg := &config.Config{}
	cfg.App.Name = "test-payments-service"
	cfg.Stripe.SecretKey = "sk_test_..."
	cfg.CORS.AllowedOrigins = []string{"*"}
	cfg.CORS.AllowedMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	cfg.CORS.AllowedHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}

	// Create test logger
	logger := logging.NewLogger("debug")

	// Create test server with database
	server := NewServer(&MockPaymentService{}, cfg, db, nil, logger)

	t.Run("Basic Health Check", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/health", nil)
		assert.NoError(t, err)
		
		w := httptest.NewRecorder()
		server.Router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "healthy")
	})

	t.Run("Detailed Health Check", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/health/detailed", nil)
		w := httptest.NewRecorder()
		server.Router.ServeHTTP(w, req)

		// Should return Service Unavailable due to invalid Stripe key
		assert.Equal(t, http.StatusServiceUnavailable, w.Code)
		assert.Contains(t, w.Body.String(), "database")
		assert.Contains(t, w.Body.String(), "redis")
		assert.Contains(t, w.Body.String(), "stripe")
		assert.Contains(t, w.Body.String(), "unhealthy")
	})

	t.Run("Readiness Probe", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/health/ready", nil)
		w := httptest.NewRecorder()
		server.Router.ServeHTTP(w, req)

		// Readiness should fail due to invalid Stripe key
		assert.Equal(t, http.StatusServiceUnavailable, w.Code)
		assert.Contains(t, w.Body.String(), "not_ready")
	})

	t.Run("Liveness Probe", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/health/live", nil)
		w := httptest.NewRecorder()
		server.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "alive")
	})
}

func TestHealthChecker_Components(t *testing.T) {
	// Create in-memory SQLite database for testing (without domain migrations)
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	checker := NewHealthChecker(db, nil, "sk_test_invalid")

	t.Run("Database Health Check", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		health := checker.checkDatabase(ctx)
		assert.Equal(t, HealthStatusHealthy, health.Status)
		assert.Contains(t, health.Message, "successful")
		assert.True(t, health.LastChecked.After(time.Now().Add(-time.Minute)))
	})

	t.Run("Redis Health Check - Not Configured", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		health := checker.checkRedis(ctx)
		assert.Equal(t, HealthStatusDegraded, health.Status)
		assert.Contains(t, health.Message, "not configured")
	})

	t.Run("Stripe Health Check - Invalid Key", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		health := checker.checkStripe(ctx)
		assert.Equal(t, HealthStatusUnhealthy, health.Status)
		assert.Contains(t, health.Message, "failed")
	})
}

func TestHealthStatus_Constants(t *testing.T) {
	assert.Equal(t, "healthy", string(HealthStatusHealthy))
	assert.Equal(t, "unhealthy", string(HealthStatusUnhealthy))
	assert.Equal(t, "degraded", string(HealthStatusDegraded))
}
