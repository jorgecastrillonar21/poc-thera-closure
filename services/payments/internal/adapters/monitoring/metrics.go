package monitoring

import (
	"net/http"
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics holds all Prometheus metrics for the payments service
type Metrics struct {
	// HTTP Metrics
	HTTPRequestsTotal   *prometheus.CounterVec
	HTTPRequestDuration *prometheus.HistogramVec
	HTTPResponseSize    *prometheus.HistogramVec
	HTTPActiveRequests  prometheus.Gauge
	HTTPErrorsTotal     *prometheus.CounterVec

	// Business Metrics
	CustomersTotal        prometheus.Counter
	CustomersCreated      *prometheus.CounterVec
	SubscriptionsTotal    prometheus.Counter
	SubscriptionsCreated  *prometheus.CounterVec
	PaymentsTotal         prometheus.Counter
	PaymentsCreated       *prometheus.CounterVec
	PaymentAmount         *prometheus.HistogramVec
	StripeAPICallsTotal   *prometheus.CounterVec
	StripeAPICallDuration *prometheus.HistogramVec

	// Database Metrics
	DatabaseConnectionsActive prometheus.Gauge
	DatabaseConnectionsIdle   prometheus.Gauge
	DatabaseQueriesTotal      *prometheus.CounterVec
	DatabaseQueryDuration     *prometheus.HistogramVec
	DatabaseErrors            *prometheus.CounterVec

	// Redis Metrics
	RedisConnectionsActive prometheus.Gauge
	RedisCommandsTotal     *prometheus.CounterVec
	RedisCommandDuration   *prometheus.HistogramVec
	RedisErrors            *prometheus.CounterVec

	// Security Metrics
	AuthenticationAttempts *prometheus.CounterVec
	RateLimitHits          *prometheus.CounterVec
	SecurityViolations     *prometheus.CounterVec
	WebhookVerifications   *prometheus.CounterVec

	// System Metrics
	SystemMemoryUsage prometheus.Gauge
	SystemCPUUsage    prometheus.Gauge
	SystemGoroutines  prometheus.Gauge
	SystemGCDuration  prometheus.Gauge
	SystemUptime      prometheus.Gauge

	// Application start time
	startTime time.Time
}

// NewMetrics creates and registers all Prometheus metrics
func NewMetrics(serviceName, version string) *Metrics {
	m := &Metrics{
		startTime: time.Now(),
	}

	// HTTP Metrics
	m.HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
			ConstLabels: prometheus.Labels{
				"service": serviceName,
				"version": version,
			},
		},
		[]string{"method", "endpoint", "status_code"},
	)

	m.HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
			ConstLabels: prometheus.Labels{
				"service": serviceName,
				"version": version,
			},
		},
		[]string{"method", "endpoint"},
	)

	m.HTTPResponseSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "Size of HTTP response in bytes",
			Buckets: []float64{100, 1000, 10000, 100000, 1000000},
			ConstLabels: prometheus.Labels{
				"service": serviceName,
				"version": version,
			},
		},
		[]string{"method", "endpoint"},
	)

	m.HTTPActiveRequests = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_active",
			Help: "Number of active HTTP requests",
			ConstLabels: prometheus.Labels{
				"service": serviceName,
				"version": version,
			},
		},
	)

	m.HTTPErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_errors_total",
			Help: "Total number of HTTP errors",
			ConstLabels: prometheus.Labels{
				"service": serviceName,
				"version": version,
			},
		},
		[]string{"method", "endpoint", "error_type"},
	)

	// Business Metrics
	m.CustomersTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "customers_total",
			Help: "Total number of customers",
			ConstLabels: prometheus.Labels{
				"service": serviceName,
			},
		},
	)

	m.CustomersCreated = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "customers_created_total",
			Help: "Total number of customers created",
			ConstLabels: prometheus.Labels{
				"service": serviceName,
			},
		},
		[]string{"status"},
	)

	m.SubscriptionsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "subscriptions_total",
			Help: "Total number of subscriptions",
			ConstLabels: prometheus.Labels{
				"service": serviceName,
			},
		},
	)

	m.SubscriptionsCreated = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "subscriptions_created_total",
			Help: "Total number of subscriptions created",
			ConstLabels: prometheus.Labels{
				"service": serviceName,
			},
		},
		[]string{"status"},
	)

	m.PaymentsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "payments_total",
			Help: "Total number of payments",
			ConstLabels: prometheus.Labels{
				"service": serviceName,
			},
		},
	)

	m.PaymentsCreated = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "payments_created_total",
			Help: "Total number of payments created",
			ConstLabels: prometheus.Labels{
				"service": serviceName,
			},
		},
		[]string{"status", "currency"},
	)

	m.PaymentAmount = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "payment_amount",
			Help:    "Payment amounts in cents",
			Buckets: []float64{100, 500, 1000, 2500, 5000, 10000, 25000, 50000, 100000, 250000},
			ConstLabels: prometheus.Labels{
				"service": serviceName,
			},
		},
		[]string{"currency"},
	)

	m.StripeAPICallsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "stripe_api_calls_total",
			Help: "Total number of Stripe API calls",
			ConstLabels: prometheus.Labels{
				"service": serviceName,
			},
		},
		[]string{"operation", "status"},
	)

	m.StripeAPICallDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "stripe_api_call_duration_seconds",
			Help:    "Duration of Stripe API calls in seconds",
			Buckets: prometheus.DefBuckets,
			ConstLabels: prometheus.Labels{
				"service": serviceName,
			},
		},
		[]string{"operation"},
	)

	// Database Metrics
	m.DatabaseConnectionsActive = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "database_connections_active",
			Help: "Number of active database connections",
			ConstLabels: prometheus.Labels{
				"service": serviceName,
			},
		},
	)

	m.DatabaseConnectionsIdle = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "database_connections_idle",
			Help: "Number of idle database connections",
			ConstLabels: prometheus.Labels{
				"service": serviceName,
			},
		},
	)

	m.DatabaseQueriesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "database_queries_total",
			Help: "Total number of database queries",
			ConstLabels: prometheus.Labels{
				"service": serviceName,
			},
		},
		[]string{"operation", "table", "status"},
	)

	m.DatabaseQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "database_query_duration_seconds",
			Help:    "Duration of database queries in seconds",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5},
			ConstLabels: prometheus.Labels{
				"service": serviceName,
			},
		},
		[]string{"operation", "table"},
	)

	m.DatabaseErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "database_errors_total",
			Help: "Total number of database errors",
			ConstLabels: prometheus.Labels{
				"service": serviceName,
			},
		},
		[]string{"operation", "error_type"},
	)

	// Redis Metrics
	m.RedisConnectionsActive = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "redis_connections_active",
			Help: "Number of active Redis connections",
			ConstLabels: prometheus.Labels{
				"service": serviceName,
			},
		},
	)

	m.RedisCommandsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "redis_commands_total",
			Help: "Total number of Redis commands",
			ConstLabels: prometheus.Labels{
				"service": serviceName,
			},
		},
		[]string{"command", "status"},
	)

	m.RedisCommandDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "redis_command_duration_seconds",
			Help:    "Duration of Redis commands in seconds",
			Buckets: []float64{0.0001, 0.0005, 0.001, 0.005, 0.01, 0.025, 0.05, 0.1},
			ConstLabels: prometheus.Labels{
				"service": serviceName,
			},
		},
		[]string{"command"},
	)

	m.RedisErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "redis_errors_total",
			Help: "Total number of Redis errors",
			ConstLabels: prometheus.Labels{
				"service": serviceName,
			},
		},
		[]string{"command", "error_type"},
	)

	// Security Metrics
	m.AuthenticationAttempts = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "authentication_attempts_total",
			Help: "Total number of authentication attempts",
			ConstLabels: prometheus.Labels{
				"service": serviceName,
			},
		},
		[]string{"method", "status"},
	)

	m.RateLimitHits = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rate_limit_hits_total",
			Help: "Total number of rate limit hits",
			ConstLabels: prometheus.Labels{
				"service": serviceName,
			},
		},
		[]string{"endpoint", "client_type"},
	)

	m.SecurityViolations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "security_violations_total",
			Help: "Total number of security violations",
			ConstLabels: prometheus.Labels{
				"service": serviceName,
			},
		},
		[]string{"violation_type", "endpoint"},
	)

	m.WebhookVerifications = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "webhook_verifications_total",
			Help: "Total number of webhook verifications",
			ConstLabels: prometheus.Labels{
				"service": serviceName,
			},
		},
		[]string{"provider", "status"},
	)

	// System Metrics
	m.SystemMemoryUsage = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "system_memory_usage_bytes",
			Help: "System memory usage in bytes",
			ConstLabels: prometheus.Labels{
				"service": serviceName,
			},
		},
	)

	m.SystemCPUUsage = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "system_cpu_usage_percent",
			Help: "System CPU usage percentage",
			ConstLabels: prometheus.Labels{
				"service": serviceName,
			},
		},
	)

	m.SystemGoroutines = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "system_goroutines",
			Help: "Number of goroutines",
			ConstLabels: prometheus.Labels{
				"service": serviceName,
			},
		},
	)

	m.SystemGCDuration = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "system_gc_duration_seconds",
			Help: "Last garbage collection duration in seconds",
			ConstLabels: prometheus.Labels{
				"service": serviceName,
			},
		},
	)

	m.SystemUptime = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "system_uptime_seconds",
			Help: "System uptime in seconds",
			ConstLabels: prometheus.Labels{
				"service": serviceName,
			},
		},
	)

	// Start system metrics collection
	go m.collectSystemMetrics()

	return m
}

// collectSystemMetrics collects system-level metrics periodically
func (m *Metrics) collectSystemMetrics() {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	var memStats runtime.MemStats

	for range ticker.C {
		// Memory metrics
		runtime.ReadMemStats(&memStats)
		m.SystemMemoryUsage.Set(float64(memStats.Alloc))

		// Goroutine metrics
		m.SystemGoroutines.Set(float64(runtime.NumGoroutine()))

		// GC metrics
		m.SystemGCDuration.Set(float64(memStats.PauseNs[(memStats.NumGC+255)%256]) / 1e9)

		// Uptime
		m.SystemUptime.Set(time.Since(m.startTime).Seconds())
	}
}

// RecordHTTPRequest records metrics for an HTTP request
func (m *Metrics) RecordHTTPRequest(method, endpoint string, statusCode int, duration time.Duration, responseSize int64) {
	m.HTTPRequestsTotal.WithLabelValues(method, endpoint, string(rune(statusCode))).Inc()
	m.HTTPRequestDuration.WithLabelValues(method, endpoint).Observe(duration.Seconds())
	m.HTTPResponseSize.WithLabelValues(method, endpoint).Observe(float64(responseSize))

	if statusCode >= 400 {
		errorType := "client_error"
		if statusCode >= 500 {
			errorType = "server_error"
		}
		m.HTTPErrorsTotal.WithLabelValues(method, endpoint, errorType).Inc()
	}
}

// IncActiveRequests increments the active requests counter
func (m *Metrics) IncActiveRequests() {
	m.HTTPActiveRequests.Inc()
}

// DecActiveRequests decrements the active requests counter
func (m *Metrics) DecActiveRequests() {
	m.HTTPActiveRequests.Dec()
}

// RecordBusinessEvent records business-specific metrics
func (m *Metrics) RecordBusinessEvent(eventType string, labels map[string]string) {
	switch eventType {
	case "customer_created":
		status := labels["status"]
		m.CustomersCreated.WithLabelValues(status).Inc()
		if status == "success" {
			m.CustomersTotal.Inc()
		}
	case "subscription_created":
		status := labels["status"]
		m.SubscriptionsCreated.WithLabelValues(status).Inc()
		if status == "success" {
			m.SubscriptionsTotal.Inc()
		}
	case "payment_created":
		status := labels["status"]
		currency := labels["currency"]
		m.PaymentsCreated.WithLabelValues(status, currency).Inc()
		if status == "success" {
			m.PaymentsTotal.Inc()
			if amount := labels["amount"]; amount != "" {
				// Parse amount and record
				// Implementation would parse the amount string
				m.PaymentAmount.WithLabelValues(currency).Observe(1000) // Placeholder
			}
		}
	}
}

// RecordStripeAPICall records metrics for Stripe API calls
func (m *Metrics) RecordStripeAPICall(operation, status string, duration time.Duration) {
	m.StripeAPICallsTotal.WithLabelValues(operation, status).Inc()
	m.StripeAPICallDuration.WithLabelValues(operation).Observe(duration.Seconds())
}

// RecordDatabaseOperation records metrics for database operations
func (m *Metrics) RecordDatabaseOperation(operation, table, status string, duration time.Duration) {
	m.DatabaseQueriesTotal.WithLabelValues(operation, table, status).Inc()
	m.DatabaseQueryDuration.WithLabelValues(operation, table).Observe(duration.Seconds())
}

// RecordSecurityEvent records security-related metrics
func (m *Metrics) RecordSecurityEvent(eventType string, labels map[string]string) {
	switch eventType {
	case "authentication_attempt":
		method := labels["method"]
		status := labels["status"]
		m.AuthenticationAttempts.WithLabelValues(method, status).Inc()
	case "rate_limit_hit":
		endpoint := labels["endpoint"]
		clientType := labels["client_type"]
		m.RateLimitHits.WithLabelValues(endpoint, clientType).Inc()
	case "security_violation":
		violationType := labels["violation_type"]
		endpoint := labels["endpoint"]
		m.SecurityViolations.WithLabelValues(violationType, endpoint).Inc()
	case "webhook_verification":
		provider := labels["provider"]
		status := labels["status"]
		m.WebhookVerifications.WithLabelValues(provider, status).Inc()
	}
}

// GetHandler returns the Prometheus HTTP handler for metrics exposition
func (m *Metrics) GetHandler() http.Handler {
	return promhttp.Handler()
}
