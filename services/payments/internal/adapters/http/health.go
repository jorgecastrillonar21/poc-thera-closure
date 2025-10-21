package http

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/account"
	"gorm.io/gorm"
)

type HealthChecker struct {
	db        *gorm.DB
	redis     *redis.Client
	stripeKey string
}

type HealthStatus string

const (
	HealthStatusHealthy   HealthStatus = "healthy"
	HealthStatusUnhealthy HealthStatus = "unhealthy"
	HealthStatusDegraded  HealthStatus = "degraded"
)

type ComponentHealth struct {
	Status       HealthStatus `json:"status"`
	Message      string       `json:"message,omitempty"`
	LastChecked  time.Time    `json:"last_checked"`
	ResponseTime string       `json:"response_time,omitempty"`
}

type HealthResponse struct {
	Status     HealthStatus               `json:"status"`
	Version    string                     `json:"version"`
	Timestamp  time.Time                  `json:"timestamp"`
	Components map[string]ComponentHealth `json:"components"`
}

func NewHealthChecker(db *gorm.DB, redis *redis.Client, stripeKey string) *HealthChecker {
	return &HealthChecker{
		db:        db,
		redis:     redis,
		stripeKey: stripeKey,
	}
}

// Basic health check - lightweight endpoint for load balancers
func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
		"service":   "payments-service",
	})
}

// Detailed health check with component status
func (s *Server) detailedHealthCheck(c *gin.Context) {
	checker := NewHealthChecker(s.db, s.redis, s.stripeKey)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	components := make(map[string]ComponentHealth)
	overallStatus := HealthStatusHealthy

	// Check database connectivity
	dbHealth := checker.checkDatabase(ctx)
	components["database"] = dbHealth
	if dbHealth.Status != HealthStatusHealthy {
		overallStatus = HealthStatusDegraded
	}

	// Check Redis connectivity (optional service)
	redisHealth := checker.checkRedis(ctx)
	components["redis"] = redisHealth
	if redisHealth.Status == HealthStatusUnhealthy {
		// Redis is optional, so degraded instead of unhealthy
		if overallStatus == HealthStatusHealthy {
			overallStatus = HealthStatusDegraded
		}
	}

	// Check Stripe API connectivity
	stripeHealth := checker.checkStripe(ctx)
	components["stripe"] = stripeHealth
	if stripeHealth.Status != HealthStatusHealthy {
		overallStatus = HealthStatusUnhealthy
	}

	response := HealthResponse{
		Status:     overallStatus,
		Version:    "1.0.0", // TODO: Get from build info
		Timestamp:  time.Now().UTC(),
		Components: components,
	}

	statusCode := http.StatusOK
	if overallStatus == HealthStatusUnhealthy {
		statusCode = http.StatusServiceUnavailable
	} else if overallStatus == HealthStatusDegraded {
		statusCode = http.StatusPartialContent
	}

	c.JSON(statusCode, response)
}

// Check database connectivity and basic query performance
func (hc *HealthChecker) checkDatabase(ctx context.Context) ComponentHealth {
	start := time.Now()

	// Test database connection with a simple query
	var result int
	err := hc.db.WithContext(ctx).Raw("SELECT 1").Scan(&result).Error

	responseTime := time.Since(start)

	if err != nil {
		return ComponentHealth{
			Status:       HealthStatusUnhealthy,
			Message:      "Database connection failed: " + err.Error(),
			LastChecked:  time.Now().UTC(),
			ResponseTime: responseTime.String(),
		}
	}

	// Check if response time is acceptable (< 100ms for simple query)
	status := HealthStatusHealthy
	message := "Database connection successful"

	if responseTime > 100*time.Millisecond {
		status = HealthStatusDegraded
		message = "Database responding slowly"
	}

	return ComponentHealth{
		Status:       status,
		Message:      message,
		LastChecked:  time.Now().UTC(),
		ResponseTime: responseTime.String(),
	}
}

// Check Redis connectivity (optional service)
func (hc *HealthChecker) checkRedis(ctx context.Context) ComponentHealth {
	if hc.redis == nil {
		return ComponentHealth{
			Status:      HealthStatusDegraded,
			Message:     "Redis not configured",
			LastChecked: time.Now().UTC(),
		}
	}

	start := time.Now()

	// Test Redis connection with ping
	err := hc.redis.Ping(ctx).Err()

	responseTime := time.Since(start)

	if err != nil {
		return ComponentHealth{
			Status:       HealthStatusUnhealthy,
			Message:      "Redis connection failed: " + err.Error(),
			LastChecked:  time.Now().UTC(),
			ResponseTime: responseTime.String(),
		}
	}

	return ComponentHealth{
		Status:       HealthStatusHealthy,
		Message:      "Redis connection successful",
		LastChecked:  time.Now().UTC(),
		ResponseTime: responseTime.String(),
	}
}

// Check Stripe API connectivity
func (hc *HealthChecker) checkStripe(ctx context.Context) ComponentHealth {
	start := time.Now()

	// Test Stripe API with account retrieval
	stripe.Key = hc.stripeKey

	// Use a lightweight API call to check connectivity
	_, err := account.Get()

	responseTime := time.Since(start)

	if err != nil {
		return ComponentHealth{
			Status:       HealthStatusUnhealthy,
			Message:      "Stripe API connection failed: " + err.Error(),
			LastChecked:  time.Now().UTC(),
			ResponseTime: responseTime.String(),
		}
	}

	// Check if response time is acceptable (< 2s for external API)
	status := HealthStatusHealthy
	message := "Stripe API connection successful"

	if responseTime > 2*time.Second {
		status = HealthStatusDegraded
		message = "Stripe API responding slowly"
	}

	return ComponentHealth{
		Status:       status,
		Message:      message,
		LastChecked:  time.Now().UTC(),
		ResponseTime: responseTime.String(),
	}
}

// Readiness probe - checks if service is ready to handle requests
func (s *Server) readinessProbe(c *gin.Context) {
	checker := NewHealthChecker(s.db, s.redis, s.stripeKey)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// For readiness, we only check essential services
	dbHealth := checker.checkDatabase(ctx)
	stripeHealth := checker.checkStripe(ctx)

	if dbHealth.Status == HealthStatusUnhealthy || stripeHealth.Status == HealthStatusUnhealthy {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "not_ready",
			"message": "Essential services are not available",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "ready",
		"message": "Service is ready to handle requests",
	})
}

// Liveness probe - checks if service is alive (minimal check)
func (s *Server) livenessProbe(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "alive",
		"timestamp": time.Now().UTC(),
	})
}
