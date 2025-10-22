package monitoring

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMetrics(t *testing.T) {
	metrics := NewMetrics("test_service", "1.0.0")
	
	assert.NotNil(t, metrics)
	assert.NotNil(t, metrics.HTTPRequestsTotal)
	assert.NotNil(t, metrics.HTTPRequestDuration)
	assert.NotNil(t, metrics.CustomersTotal)
	assert.NotNil(t, metrics.DatabaseQueriesTotal)
	assert.NotNil(t, metrics.SystemMemoryUsage)
}

func TestMetricsHTTPHandler(t *testing.T) {
	metrics := NewMetrics("test_service_handler", "1.0.0")
	
	handler := metrics.GetHandler()
	assert.NotNil(t, handler)
}

func TestDatabaseMonitor(t *testing.T) {
	// Skip this test as it requires a real database connection
	t.Skip("Database monitor test requires GORM DB instance")
}

func TestGormMetricsPlugin(t *testing.T) {
	metrics := NewMetrics("test_service_gorm", "1.0.0")
	plugin := &gormMetricsPlugin{metrics: metrics}
	
	// Test plugin name
	name := plugin.Name()
	assert.Equal(t, "metrics", name)
}

func TestRedisMonitor(t *testing.T) {
	metrics := NewMetrics("test_service_redis", "1.0.0")
	redisMonitor := NewRedisMonitor(metrics)
	
	assert.NotNil(t, redisMonitor)
	assert.Equal(t, metrics, redisMonitor.metrics)
}

func TestStripeMonitor(t *testing.T) {
	metrics := NewMetrics("test_service_stripe", "1.0.0")
	stripeMonitor := NewStripeMonitor(metrics)
	
	assert.NotNil(t, stripeMonitor)
	assert.Equal(t, metrics, stripeMonitor.metrics)
}