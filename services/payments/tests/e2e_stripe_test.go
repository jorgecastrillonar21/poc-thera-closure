package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"theraclosure/payments-service/internal/adapters/config"
	httpAdapter "theraclosure/payments-service/internal/adapters/http"
	"theraclosure/payments-service/internal/adapters/logging"
	"theraclosure/payments-service/internal/adapters/persistence"
	"theraclosure/payments-service/internal/adapters/stripe"
	"theraclosure/payments-service/internal/core/domain"
	"theraclosure/payments-service/internal/core/services"
)

type E2ETestSuite struct {
	server       *httpAdapter.Server
	db           *gorm.DB
	stripeClient *stripe.StripeClient
	testEmail    string
}

func setupE2ETestSuite() (*E2ETestSuite, error) {
	// Use in-memory SQLite for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Run migrations
	err = db.AutoMigrate(
		&domain.Customer{},
		&domain.Subscription{},
		&domain.Payment{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	// Create test configuration
	cfg := &config.Config{}
	cfg.App.Name = "test-payments-service"
	cfg.App.Version = "test"
	cfg.App.LogLevel = "info"
	cfg.Server.Host = "localhost"
	cfg.Server.Port = "0" // Let the system choose a port
	cfg.Stripe.SecretKey = getTestStripeKey()
	cfg.Security.EnableAuthentication = false // Disable for testing
	cfg.Security.MaxRequestSize = 10485760
	cfg.Security.RateLimitRPS = 1000
	cfg.Security.RateLimitWindow = "1m"

	// Initialize logger
	logger := logging.NewLogger("info")

	// Initialize repositories
	dbWrapper := &persistence.Database{}
	// We'll need to use a different approach since Database fields are private
	// For E2E tests, we can use the NewDatabase constructor with a custom config
	customerRepo := persistence.NewCustomerRepository(dbWrapper)
	subscriptionRepo := persistence.NewSubscriptionRepository(dbWrapper)
	paymentRepo := persistence.NewPaymentRepository(dbWrapper)

	// Initialize Stripe client
	stripeClient := stripe.NewStripeClient(cfg.Stripe.SecretKey)

	// Initialize service
	paymentService := services.NewPaymentService(
		customerRepo,
		subscriptionRepo,
		paymentRepo,
		stripeClient,
	)

	// Initialize HTTP server
	server := httpAdapter.NewServer(paymentService, cfg, db, nil, logger)

	return &E2ETestSuite{
		server:       server,
		db:           db,
		stripeClient: stripeClient,
		testEmail:    fmt.Sprintf("e2e-test-%d@theraclosure.com", time.Now().Unix()),
	}, nil
}

func getTestStripeKey() string {
	if key := os.Getenv("STRIPE_TEST_SECRET_KEY"); key != "" {
		return key
	}
	return "sk_test_fake_key_for_unit_tests" // Placeholder for tests
}

func TestE2E_CompletePaymentWorkflow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	suite, err := setupE2ETestSuite()
	require.NoError(t, err)

	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Test the complete workflow: Customer -> Payment Intent -> Payment
	t.Run("complete_payment_workflow", func(t *testing.T) {
		// Step 1: Create a customer
		customerReq := map[string]interface{}{
			"email":      suite.testEmail,
			"name":       "E2E Test Customer",
			"phone":      "+1234567890",
			"user_id":    fmt.Sprintf("user_%d", time.Now().Unix()),
			"is_active":  true,
		}

		customerJSON, err := json.Marshal(customerReq)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/api/v1/customers", bytes.NewBuffer(customerJSON))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		suite.server.Router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)

		var customerResponse map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &customerResponse)
		require.NoError(t, err)

		customerData := customerResponse["data"].(map[string]interface{})
		customerID := customerData["id"].(string)
		stripeCustomerID := customerData["stripe_customer_id"].(string)

		assert.NotEmpty(t, customerID)
		assert.NotEmpty(t, stripeCustomerID)
		assert.Contains(t, stripeCustomerID, "cus_")

		// Step 2: Create a payment intent
		paymentIntentReq := map[string]interface{}{
			"amount":      2999, // $29.99
			"currency":    "usd",
			"customer_id": customerID,
			"metadata": map[string]string{
				"test_key": "test_value",
			},
		}

		paymentIntentJSON, err := json.Marshal(paymentIntentReq)
		require.NoError(t, err)

		req, err = http.NewRequest("POST", "/api/v1/payment-intents", bytes.NewBuffer(paymentIntentJSON))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		rr = httptest.NewRecorder()
		suite.server.Router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)

		var paymentIntentResponse map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &paymentIntentResponse)
		require.NoError(t, err)

		paymentIntentData := paymentIntentResponse["data"].(map[string]interface{})
		paymentIntentID := paymentIntentData["id"].(string)

		assert.NotEmpty(t, paymentIntentID)
		assert.Contains(t, paymentIntentID, "pi_")
		assert.Equal(t, float64(2999), paymentIntentData["amount"].(float64))
		assert.Equal(t, "usd", paymentIntentData["currency"])

		// Step 3: Create a payment record
		paymentReq := map[string]interface{}{
			"customer_id":        customerID,
			"amount":            29.99,
			"currency":          "USD",
			"payment_method":    "card",
			"status":            "pending",
			"stripe_payment_id": paymentIntentID,
			"metadata": map[string]string{
				"source": "e2e_test",
			},
		}

		paymentJSON, err := json.Marshal(paymentReq)
		require.NoError(t, err)

		req, err = http.NewRequest("POST", "/api/v1/payments", bytes.NewBuffer(paymentJSON))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		rr = httptest.NewRecorder()
		suite.server.Router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)

		var paymentResponse map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &paymentResponse)
		require.NoError(t, err)

		paymentData := paymentResponse["data"].(map[string]interface{})
		paymentID := paymentData["id"].(string)

		assert.NotEmpty(t, paymentID)
		assert.Equal(t, 29.99, paymentData["amount"])
		assert.Equal(t, "USD", paymentData["currency"])
		assert.Equal(t, "pending", paymentData["status"])

		// Step 4: Verify customer has the payment
		req, err = http.NewRequest("GET", fmt.Sprintf("/api/v1/customers/%s", customerID), nil)
		require.NoError(t, err)

		rr = httptest.NewRecorder()
		suite.server.Router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		// Step 5: List payments for the customer
		req, err = http.NewRequest("GET", fmt.Sprintf("/api/v1/payments?customer_id=%s", customerID), nil)
		require.NoError(t, err)

		rr = httptest.NewRecorder()
		suite.server.Router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var paymentsResponse map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &paymentsResponse)
		require.NoError(t, err)

		paymentsData := paymentsResponse["data"].(map[string]interface{})
		payments := paymentsData["payments"].([]interface{})

		assert.Len(t, payments, 1)
		payment := payments[0].(map[string]interface{})
		assert.Equal(t, paymentID, payment["id"])
	})
}

func TestE2E_SubscriptionWorkflow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	suite, err := setupE2ETestSuite()
	require.NoError(t, err)

	gin.SetMode(gin.TestMode)

	t.Run("subscription_workflow", func(t *testing.T) {
		// Step 1: Create a customer
		customerReq := map[string]interface{}{
			"email":     fmt.Sprintf("sub-test-%d@theraclosure.com", time.Now().Unix()),
			"name":      "Subscription Test Customer",
			"user_id":   fmt.Sprintf("user_sub_%d", time.Now().Unix()),
			"is_active": true,
		}

		customerJSON, err := json.Marshal(customerReq)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/api/v1/customers", bytes.NewBuffer(customerJSON))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		suite.server.Router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)

		var customerResponse map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &customerResponse)
		require.NoError(t, err)

		customerData := customerResponse["data"].(map[string]interface{})
		customerID := customerData["id"].(string)

		// Step 2: Attempt to create a subscription (will fail with invalid price, but tests the endpoint)
		subscriptionReq := map[string]interface{}{
			"customer_id": customerID,
			"price_id":    "price_invalid", // Using invalid price to test error handling
			"status":      "active",
		}

		subscriptionJSON, err := json.Marshal(subscriptionReq)
		require.NoError(t, err)

		req, err = http.NewRequest("POST", "/api/v1/subscriptions", bytes.NewBuffer(subscriptionJSON))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		rr = httptest.NewRecorder()
		suite.server.Router.ServeHTTP(rr, req)

		// Expect error due to invalid price ID
		assert.Equal(t, http.StatusInternalServerError, rr.Code)

		var errorResponse map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse["error"], "Failed to create subscription")
	})
}

func TestE2E_ErrorHandling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	suite, err := setupE2ETestSuite()
	require.NoError(t, err)

	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		method         string
		path           string
		body           interface{}
		expectedStatus int
	}{
		{
			name:           "invalid_customer_creation",
			method:         "POST",
			path:           "/api/v1/customers",
			body:           map[string]interface{}{"invalid": "data"},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "get_nonexistent_customer",
			method:         "GET",
			path:           "/api/v1/customers/nonexistent",
			body:           nil,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "invalid_payment_intent",
			method:         "POST",
			path:           "/api/v1/payment-intents",
			body:           map[string]interface{}{"amount": -100},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "get_nonexistent_payment",
			method:         "GET",
			path:           "/api/v1/payments/nonexistent",
			body:           nil,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			var err error

			if tt.body != nil {
				bodyJSON, _ := json.Marshal(tt.body)
				req, err = http.NewRequest(tt.method, tt.path, bytes.NewBuffer(bodyJSON))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req, err = http.NewRequest(tt.method, tt.path, nil)
			}
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			suite.server.Router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			// Verify error response has proper structure
			var errorResponse map[string]interface{}
			err = json.Unmarshal(rr.Body.Bytes(), &errorResponse)
			require.NoError(t, err)
			assert.Contains(t, errorResponse, "error")
		})
	}
}

func TestE2E_HealthChecks(t *testing.T) {
	suite, err := setupE2ETestSuite()
	require.NoError(t, err)

	gin.SetMode(gin.TestMode)

	healthEndpoints := []string{
		"/health",
		"/health/detailed",
		"/health/ready",
		"/health/live",
	}

	for _, endpoint := range healthEndpoints {
		t.Run(fmt.Sprintf("health_check_%s", endpoint), func(t *testing.T) {
			req, err := http.NewRequest("GET", endpoint, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			suite.server.Router.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusOK, rr.Code)

			var response map[string]interface{}
			err = json.Unmarshal(rr.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.Equal(t, "ok", response["status"])
		})
	}
}

func TestE2E_MetricsEndpoint(t *testing.T) {
	suite, err := setupE2ETestSuite()
	require.NoError(t, err)

	gin.SetMode(gin.TestMode)

	req, err := http.NewRequest("GET", "/metrics", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	suite.server.Router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Header().Get("Content-Type"), "text/plain")

	// Check for Prometheus metrics format
	body := rr.Body.String()
	assert.Contains(t, body, "# HELP")
	assert.Contains(t, body, "# TYPE")
}