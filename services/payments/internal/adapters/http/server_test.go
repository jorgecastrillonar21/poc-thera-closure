package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"theraclosure/payments-service/internal/adapters/config"
	"theraclosure/payments-service/internal/core/domain"
	"theraclosure/payments-service/internal/core/ports"
)

func setupTestServer() (*gin.Engine, *config.Config) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create test config
	cfg := &config.Config{}
	cfg.App.Name = "test-payments-service"
	cfg.App.LogLevel = "debug"
	cfg.CORS.AllowedOrigins = []string{"*"}
	cfg.CORS.AllowedMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	cfg.CORS.AllowedHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}

	// Create a mock payment service for testing
	mockService := &MockPaymentService{}

	// Create server with mock service and nil dependencies for testing
	server := NewServer(mockService, cfg, nil, nil, nil)

	return server.router, cfg
}

// MockPaymentService is a simple mock for testing HTTP endpoints
type MockPaymentService struct{}

func (m *MockPaymentService) CreateCustomer(ctx context.Context, req ports.CreateCustomerRequest) (*domain.Customer, error) {
	// Return a mock customer for testing
	customer := &domain.Customer{
		ID:     uuid.New().String(),
		UserID: req.UserID,
		Email:  req.Email,
		Name:   req.Name,
		Active: true,
	}
	return customer, nil
}

func (m *MockPaymentService) GetCustomer(ctx context.Context, id string) (*domain.Customer, error) {
	return &domain.Customer{
		ID:     id,
		UserID: uuid.New().String(),
		Email:  "test@example.com",
		Name:   "Test Customer",
		Active: true,
	}, nil
}

func (m *MockPaymentService) GetCustomerByUserID(ctx context.Context, userID string) (*domain.Customer, error) {
	return &domain.Customer{
		ID:     uuid.New().String(),
		UserID: userID,
		Email:  "test@example.com",
		Name:   "Test Customer",
		Active: true,
	}, nil
}

func (m *MockPaymentService) UpdateCustomer(ctx context.Context, id string, req ports.UpdateCustomerRequest) (*domain.Customer, error) {
	return &domain.Customer{
		ID:     id,
		UserID: uuid.New().String(),
		Email:  req.Email,
		Name:   req.Name,
		Active: true,
	}, nil
}

func (m *MockPaymentService) DeleteCustomer(ctx context.Context, id string) error {
	return nil
}

func (m *MockPaymentService) ListCustomers(ctx context.Context, req ports.ListCustomersRequest) (*ports.ListCustomersResponse, error) {
	return &ports.ListCustomersResponse{
		Customers: []*domain.Customer{},
		Total:     0,
		Offset:    req.Offset,
		Limit:     req.Limit,
	}, nil
}

func (m *MockPaymentService) HealthCheck(ctx context.Context) error {
	return nil
}

// Add other required methods as no-ops for testing
func (m *MockPaymentService) CreateSubscription(ctx context.Context, req ports.CreateSubscriptionRequest) (*domain.Subscription, error) {
	return nil, nil
}
func (m *MockPaymentService) GetSubscription(ctx context.Context, id string) (*domain.Subscription, error) {
	return nil, nil
}
func (m *MockPaymentService) UpdateSubscription(ctx context.Context, id string, req ports.UpdateSubscriptionRequest) (*domain.Subscription, error) {
	return nil, nil
}
func (m *MockPaymentService) CancelSubscription(ctx context.Context, id string) (*domain.Subscription, error) {
	return nil, nil
}
func (m *MockPaymentService) ListSubscriptions(ctx context.Context, req ports.ListSubscriptionsRequest) (*ports.ListSubscriptionsResponse, error) {
	return nil, nil
}
func (m *MockPaymentService) GetCustomerSubscriptions(ctx context.Context, customerID string) ([]*domain.Subscription, error) {
	return nil, nil
}
func (m *MockPaymentService) CreatePayment(ctx context.Context, req ports.CreatePaymentRequest) (*domain.Payment, error) {
	return nil, nil
}
func (m *MockPaymentService) GetPayment(ctx context.Context, id string) (*domain.Payment, error) {
	return nil, nil
}
func (m *MockPaymentService) ListPayments(ctx context.Context, req ports.ListPaymentsRequest) (*ports.ListPaymentsResponse, error) {
	return nil, nil
}
func (m *MockPaymentService) RefundPayment(ctx context.Context, id string, amount *int64) (*domain.Payment, error) {
	return nil, nil
}
func (m *MockPaymentService) CreatePaymentIntent(ctx context.Context, req ports.CreatePaymentIntentRequest) (*ports.CreatePaymentIntentResponse, error) {
	return nil, nil
}
func (m *MockPaymentService) ConfirmPaymentIntent(ctx context.Context, paymentIntentID string) (*ports.ConfirmPaymentIntentResponse, error) {
	return nil, nil
}
func (m *MockPaymentService) HandleWebhook(ctx context.Context, payload []byte, signature string) error {
	return nil
}

func TestHealthCheckEndpoints(t *testing.T) {
	router, _ := setupTestServer()

	tests := []struct {
		name     string
		endpoint string
		method   string
		expected int
	}{
		{
			name:     "health check root",
			endpoint: "/health",
			method:   "GET",
			expected: http.StatusOK,
		},
		{
			name:     "health check API v1",
			endpoint: "/api/v1/health",
			method:   "GET",
			expected: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.endpoint, nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expected, w.Code)

			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.Equal(t, "healthy", response["status"])
			assert.NotEmpty(t, response["timestamp"])
		})
	}
}

func TestCustomerEndpoints(t *testing.T) {
	router, _ := setupTestServer()

	t.Run("Create Customer", func(t *testing.T) {
		customerReq := ports.CreateCustomerRequest{
			UserID: uuid.New().String(),
			Email:  "test@example.com",
			Name:   "Test Customer",
		}

		jsonData, err := json.Marshal(customerReq)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/api/v1/customers", bytes.NewBuffer(jsonData))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.NotNil(t, response["data"])
	})

	t.Run("Create Customer - Invalid Request", func(t *testing.T) {
		invalidReq := map[string]interface{}{
			"email": "invalid-email",
		}

		jsonData, err := json.Marshal(invalidReq)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/api/v1/customers", bytes.NewBuffer(jsonData))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Get Customer", func(t *testing.T) {
		customerID := uuid.New().String()

		req, err := http.NewRequest("GET", "/api/v1/customers/"+customerID, nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.NotNil(t, response["data"])
	})

	t.Run("Get Customer - Missing ID", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/v1/customers/", nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should return 404 for missing ID path
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Update Customer", func(t *testing.T) {
		customerID := uuid.New().String()
		updateReq := ports.UpdateCustomerRequest{
			Email: "updated@example.com",
			Name:  "Updated Name",
		}

		jsonData, err := json.Marshal(updateReq)
		require.NoError(t, err)

		req, err := http.NewRequest("PUT", "/api/v1/customers/"+customerID, bytes.NewBuffer(jsonData))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.NotNil(t, response["data"])
	})

	t.Run("Delete Customer", func(t *testing.T) {
		customerID := uuid.New().String()

		req, err := http.NewRequest("DELETE", "/api/v1/customers/"+customerID, nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "Customer deleted successfully", response["message"])
	})

	t.Run("List Customers", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/v1/customers?offset=0&limit=10", nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.NotNil(t, response["data"])
	})

	t.Run("Get Customer By User ID", func(t *testing.T) {
		userID := uuid.New().String()

		req, err := http.NewRequest("GET", "/api/v1/customers/user/"+userID, nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.NotNil(t, response["data"])
	})
}

func TestCORSHeaders(t *testing.T) {
	router, _ := setupTestServer()

	t.Run("CORS Preflight", func(t *testing.T) {
		req, err := http.NewRequest("OPTIONS", "/api/v1/customers", nil)
		require.NoError(t, err)
		req.Header.Set("Origin", "http://localhost:3000")
		req.Header.Set("Access-Control-Request-Method", "POST")
		req.Header.Set("Access-Control-Request-Headers", "Content-Type")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Origin"))
		assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Methods"))
		assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Headers"))
	})
}

func TestInvalidJSONHandling(t *testing.T) {
	router, _ := setupTestServer()

	t.Run("Invalid JSON Body", func(t *testing.T) {
		invalidJSON := `{"invalid": json}`

		req, err := http.NewRequest("POST", "/api/v1/customers", bytes.NewBufferString(invalidJSON))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "Invalid request body", response["error"])
	})
}

func TestEndpointNotFound(t *testing.T) {
	router, _ := setupTestServer()

	t.Run("Non-existent endpoint", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/v1/nonexistent", nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestHTTPMethodsNotAllowed(t *testing.T) {
	router, _ := setupTestServer()

	t.Run("Method not allowed", func(t *testing.T) {
		// Try PATCH on an endpoint that doesn't support it
		req, err := http.NewRequest("PATCH", "/api/v1/customers", nil)
		require.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})
}
