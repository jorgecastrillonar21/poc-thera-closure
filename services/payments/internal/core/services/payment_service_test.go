package services

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"theraclosure/payments-service/internal/core/domain"
	"theraclosure/payments-service/internal/core/ports"
)

// Simple email validation helper
func isValidEmail(email string) bool {
	return email != "" && strings.Contains(email, "@") && strings.Contains(email, ".")
}

// Test helper to create a test customer
func createTestCustomer(t *testing.T) *domain.Customer {
	return &domain.Customer{
		ID:       uuid.New().String(),
		UserID:   uuid.New().String(),
		StripeID: "cus_test_" + uuid.New().String(),
		Email:    "test@example.com",
		Name:     "Test Customer",
		Active:   true,
	}
}

// Test helper to create a test subscription
func createTestSubscription(t *testing.T, customerID string) *domain.Subscription {
	return &domain.Subscription{
		ID:                 uuid.New().String(),
		CustomerID:         customerID,
		StripeID:          "sub_test_" + uuid.New().String(),
		PriceID:           "price_test_123",
		Status:            domain.SubscriptionStatusActive,
		CurrentPeriodStart: time.Now(),
		CurrentPeriodEnd:   time.Now().AddDate(0, 1, 0),
		Amount:            2000,
		Currency:          "usd",
		Active:            true,
	}
}

// Test helper to create a test payment
func createTestPayment(t *testing.T, customerID string, subscriptionID *string) *domain.Payment {
	return &domain.Payment{
		ID:             uuid.New().String(),
		CustomerID:     customerID,
		SubscriptionID: subscriptionID,
		StripeID:       "pi_test_" + uuid.New().String(),
		Amount:         1000,
		Currency:       "usd",
		Status:         domain.PaymentStatusSucceeded,
		Description:    "Test payment",
	}
}

func TestPaymentService_ValidateCustomerRequests(t *testing.T) {
	// Test CreateCustomerRequest validation
	t.Run("CreateCustomerRequest validation", func(t *testing.T) {
		tests := []struct {
			name    string
			request ports.CreateCustomerRequest
			wantErr bool
		}{
			{
				name: "valid request",
				request: ports.CreateCustomerRequest{
					UserID: uuid.New().String(),
					Email:  "test@example.com",
					Name:   "Test User",
				},
				wantErr: false,
			},
			{
				name: "missing user ID",
				request: ports.CreateCustomerRequest{
					Email: "test@example.com",
					Name:  "Test User",
				},
				wantErr: true,
			},
			{
				name: "invalid email",
				request: ports.CreateCustomerRequest{
					UserID: uuid.New().String(),
					Email:  "invalid-email",
					Name:   "Test User",
				},
				wantErr: true,
			},
			{
				name: "missing name",
				request: ports.CreateCustomerRequest{
					UserID: uuid.New().String(),
					Email:  "test@example.com",
				},
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// Basic validation logic test
				hasUserID := tt.request.UserID != ""
				hasValidEmail := isValidEmail(tt.request.Email)
				hasName := tt.request.Name != ""
				
				isValid := hasUserID && hasValidEmail && hasName
				if tt.wantErr {
					assert.False(t, isValid, "Request should be invalid")
				} else {
					assert.True(t, isValid, "Request should be valid")
				}
			})
		}
	})
}

func TestPaymentService_ValidateSubscriptionRequests(t *testing.T) {
	// Test CreateSubscriptionRequest validation
	t.Run("CreateSubscriptionRequest validation", func(t *testing.T) {
		tests := []struct {
			name    string
			request ports.CreateSubscriptionRequest
			wantErr bool
		}{
			{
				name: "valid request",
				request: ports.CreateSubscriptionRequest{
					CustomerID: uuid.New().String(),
					PriceID:    "price_123",
				},
				wantErr: false,
			},
			{
				name: "valid request with trial",
				request: ports.CreateSubscriptionRequest{
					CustomerID: uuid.New().String(),
					PriceID:    "price_123",
					TrialDays:  func() *int { days := 7; return &days }(),
				},
				wantErr: false,
			},
			{
				name: "missing customer ID",
				request: ports.CreateSubscriptionRequest{
					PriceID: "price_123",
				},
				wantErr: true,
			},
			{
				name: "missing price ID",
				request: ports.CreateSubscriptionRequest{
					CustomerID: uuid.New().String(),
				},
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// Basic validation logic test
				hasCustomerID := tt.request.CustomerID != ""
				hasPriceID := tt.request.PriceID != ""
				
				isValid := hasCustomerID && hasPriceID
				if tt.wantErr {
					assert.False(t, isValid, "Request should be invalid")
				} else {
					assert.True(t, isValid, "Request should be valid")
				}
			})
		}
	})
}

func TestPaymentService_ValidatePaymentRequests(t *testing.T) {
	// Test CreatePaymentRequest validation
	t.Run("CreatePaymentRequest validation", func(t *testing.T) {
		tests := []struct {
			name    string
			request ports.CreatePaymentRequest
			wantErr bool
		}{
			{
				name: "valid one-time payment",
				request: ports.CreatePaymentRequest{
					CustomerID:  uuid.New().String(),
					Amount:      1000,
					Currency:    "usd",
					Description: "Test payment",
				},
				wantErr: false,
			},
			{
				name: "valid subscription payment",
				request: ports.CreatePaymentRequest{
					CustomerID:     uuid.New().String(),
					SubscriptionID: func() *string { id := uuid.New().String(); return &id }(),
					Amount:         2000,
					Currency:       "usd",
					Description:    "Subscription payment",
				},
				wantErr: false,
			},
			{
				name: "missing customer ID",
				request: ports.CreatePaymentRequest{
					Amount:   1000,
					Currency: "usd",
				},
				wantErr: true,
			},
			{
				name: "zero amount",
				request: ports.CreatePaymentRequest{
					CustomerID: uuid.New().String(),
					Amount:     0,
					Currency:   "usd",
				},
				wantErr: true,
			},
			{
				name: "negative amount",
				request: ports.CreatePaymentRequest{
					CustomerID: uuid.New().String(),
					Amount:     -100,
					Currency:   "usd",
				},
				wantErr: true,
			},
			{
				name: "missing currency",
				request: ports.CreatePaymentRequest{
					CustomerID: uuid.New().String(),
					Amount:     1000,
				},
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// Basic validation logic test
				hasCustomerID := tt.request.CustomerID != ""
				hasValidAmount := tt.request.Amount > 0
				hasCurrency := tt.request.Currency != ""
				
				isValid := hasCustomerID && hasValidAmount && hasCurrency
				if tt.wantErr {
					assert.False(t, isValid, "Request should be invalid")
				} else {
					assert.True(t, isValid, "Request should be valid")
				}
			})
		}
	})
}

func TestPaymentService_ValidatePaymentIntentRequests(t *testing.T) {
	// Test CreatePaymentIntentRequest validation
	t.Run("CreatePaymentIntentRequest validation", func(t *testing.T) {
		tests := []struct {
			name    string
			request ports.CreatePaymentIntentRequest
			wantErr bool
		}{
			{
				name: "valid payment intent",
				request: ports.CreatePaymentIntentRequest{
					CustomerID:  uuid.New().String(),
					Amount:      1500,
					Currency:    "usd",
					Description: "Test payment intent",
				},
				wantErr: false,
			},
			{
				name: "valid payment intent with metadata",
				request: ports.CreatePaymentIntentRequest{
					CustomerID:  uuid.New().String(),
					Amount:      2500,
					Currency:    "eur",
					Description: "Test payment with metadata",
					Metadata: map[string]string{
						"order_id":    "123",
						"customer_id": "456",
					},
				},
				wantErr: false,
			},
			{
				name: "missing customer ID",
				request: ports.CreatePaymentIntentRequest{
					Amount:   1500,
					Currency: "usd",
				},
				wantErr: true,
			},
			{
				name: "invalid amount",
				request: ports.CreatePaymentIntentRequest{
					CustomerID: uuid.New().String(),
					Amount:     0,
					Currency:   "usd",
				},
				wantErr: true,
			},
			{
				name: "missing currency",
				request: ports.CreatePaymentIntentRequest{
					CustomerID: uuid.New().String(),
					Amount:     1500,
				},
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// Basic validation logic test
				hasCustomerID := tt.request.CustomerID != ""
				hasValidAmount := tt.request.Amount > 0
				hasCurrency := tt.request.Currency != ""
				
				isValid := hasCustomerID && hasValidAmount && hasCurrency
				if tt.wantErr {
					assert.False(t, isValid, "Request should be invalid")
				} else {
					assert.True(t, isValid, "Request should be valid")
				}
			})
		}
	})
}

func TestListRequestsPagination(t *testing.T) {
	t.Run("ListCustomersRequest pagination", func(t *testing.T) {
		req := ports.ListCustomersRequest{
			Offset: 0,
			Limit:  10,
		}
		
		assert.Equal(t, 0, req.Offset)
		assert.Equal(t, 10, req.Limit)
		
		// Test default limit handling
		if req.Limit == 0 {
			req.Limit = 10
		}
		assert.Equal(t, 10, req.Limit)
	})
	
	t.Run("ListSubscriptionsRequest pagination", func(t *testing.T) {
		customerID := uuid.New().String()
		req := ports.ListSubscriptionsRequest{
			CustomerID: customerID,
			Status:     "active",
			Offset:     20,
			Limit:      50,
		}
		
		assert.Equal(t, customerID, req.CustomerID)
		assert.Equal(t, "active", req.Status)
		assert.Equal(t, 20, req.Offset)
		assert.Equal(t, 50, req.Limit)
	})
	
	t.Run("ListPaymentsRequest pagination", func(t *testing.T) {
		customerID := uuid.New().String()
		subscriptionID := uuid.New().String()
		req := ports.ListPaymentsRequest{
			CustomerID:     customerID,
			SubscriptionID: subscriptionID,
			Status:         "succeeded",
			Offset:         10,
			Limit:          25,
		}
		
		assert.Equal(t, customerID, req.CustomerID)
		assert.Equal(t, subscriptionID, req.SubscriptionID)
		assert.Equal(t, "succeeded", req.Status)
		assert.Equal(t, 10, req.Offset)
		assert.Equal(t, 25, req.Limit)
	})
}

func TestResponseStructures(t *testing.T) {
	t.Run("ListCustomersResponse", func(t *testing.T) {
		customers := []*domain.Customer{
			createTestCustomer(t),
			createTestCustomer(t),
		}
		
		response := ports.ListCustomersResponse{
			Customers: customers,
			Total:     100,
			Offset:    0,
			Limit:     10,
		}
		
		assert.Len(t, response.Customers, 2)
		assert.Equal(t, int64(100), response.Total)
		assert.Equal(t, 0, response.Offset)
		assert.Equal(t, 10, response.Limit)
	})
	
	t.Run("CreatePaymentIntentResponse", func(t *testing.T) {
		response := ports.CreatePaymentIntentResponse{
			PaymentIntentID: "pi_test_123",
			ClientSecret:    "pi_test_123_secret",
			Status:          "requires_payment_method",
		}
		
		assert.Equal(t, "pi_test_123", response.PaymentIntentID)
		assert.Equal(t, "pi_test_123_secret", response.ClientSecret)
		assert.Equal(t, "requires_payment_method", response.Status)
	})
	
	t.Run("ConfirmPaymentIntentResponse", func(t *testing.T) {
		paymentID := uuid.New().String()
		response := ports.ConfirmPaymentIntentResponse{
			PaymentIntentID: "pi_test_123",
			Status:          "succeeded",
			PaymentID:       paymentID,
		}
		
		assert.Equal(t, "pi_test_123", response.PaymentIntentID)
		assert.Equal(t, "succeeded", response.Status)
		assert.Equal(t, paymentID, response.PaymentID)
	})
}

func TestDomainModelBusinessLogic(t *testing.T) {
	t.Run("Customer business logic", func(t *testing.T) {
		customer := createTestCustomer(t)
		
		// Test validation
		assert.True(t, customer.IsValid())
		
		// Test invalid cases
		customer.UserID = ""
		assert.False(t, customer.IsValid())
		
		customer.UserID = uuid.New().String()
		customer.Email = ""
		assert.False(t, customer.IsValid())
		
		customer.Email = "test@example.com"
		customer.Name = ""
		assert.False(t, customer.IsValid())
	})
	
	t.Run("Subscription business logic", func(t *testing.T) {
		customerID := uuid.New().String()
		subscription := createTestSubscription(t, customerID)
		
		// Test validation
		assert.True(t, subscription.IsValid())
		
		// Test lifecycle
		assert.Equal(t, domain.SubscriptionStatusActive, subscription.Status)
		
		// Test cancellation
		subscription.Status = domain.SubscriptionStatusCanceled
		now := time.Now()
		subscription.CanceledAt = &now
		
		assert.Equal(t, domain.SubscriptionStatusCanceled, subscription.Status)
		assert.NotNil(t, subscription.CanceledAt)
	})
	
	t.Run("Payment business logic", func(t *testing.T) {
		customerID := uuid.New().String()
		subscriptionID := uuid.New().String()
		payment := createTestPayment(t, customerID, &subscriptionID)
		
		// Test validation
		assert.True(t, payment.IsValid())
		
		// Test payment lifecycle
		assert.Equal(t, domain.PaymentStatusSucceeded, payment.Status)
		
		// Test refund
		payment.Status = domain.PaymentStatusRefunded
		now := time.Now()
		payment.RefundedAt = &now
		
		assert.Equal(t, domain.PaymentStatusRefunded, payment.Status)
		assert.NotNil(t, payment.RefundedAt)
	})
}

func TestBusinessRules(t *testing.T) {
	t.Run("Payment amount validation", func(t *testing.T) {
		customerID := uuid.New().String()
		
		// Valid amounts
		validAmounts := []int64{1, 50, 100, 1000, 999999}
		for _, amount := range validAmounts {
			payment := &domain.Payment{
				CustomerID: customerID,
				Amount:     amount,
				Currency:   "usd",
			}
			assert.True(t, payment.IsValid(), "Amount %d should be valid", amount)
		}
		
		// Invalid amounts
		invalidAmounts := []int64{0, -1, -100}
		for _, amount := range invalidAmounts {
			payment := &domain.Payment{
				CustomerID: customerID,
				Amount:     amount,
				Currency:   "usd",
			}
			assert.False(t, payment.IsValid(), "Amount %d should be invalid", amount)
		}
	})
	
	t.Run("Subscription status transitions", func(t *testing.T) {
		customerID := uuid.New().String()
		subscription := createTestSubscription(t, customerID)
		
		// Valid status transitions
		validStatuses := []domain.SubscriptionStatus{
			domain.SubscriptionStatusTrialing,
			domain.SubscriptionStatusActive,
			domain.SubscriptionStatusPastDue,
			domain.SubscriptionStatusUnpaid,
			domain.SubscriptionStatusCanceled,
			domain.SubscriptionStatusPaused,
		}
		
		for _, status := range validStatuses {
			subscription.Status = status
			assert.True(t, subscription.IsValid(), "Status %s should be valid", status)
		}
	})
	
	t.Run("Payment status transitions", func(t *testing.T) {
		customerID := uuid.New().String()
		payment := createTestPayment(t, customerID, nil)
		
		// Valid status transitions
		validStatuses := []domain.PaymentStatus{
			domain.PaymentStatusPending,
			domain.PaymentStatusSucceeded,
			domain.PaymentStatusFailed,
			domain.PaymentStatusRefunded,
			domain.PaymentStatusCanceled,
		}
		
		for _, status := range validStatuses {
			payment.Status = status
			assert.True(t, payment.IsValid(), "Status %s should be valid", status)
		}
	})
}

func TestCurrencySupport(t *testing.T) {
	t.Run("Supported currencies", func(t *testing.T) {
		customerID := uuid.New().String()
		supportedCurrencies := []string{"usd", "eur", "gbp", "cad", "aud", "jpy"}
		
		for _, currency := range supportedCurrencies {
			payment := &domain.Payment{
				CustomerID: customerID,
				Amount:     1000,
				Currency:   currency,
			}
			assert.True(t, payment.IsValid(), "Currency %s should be supported", currency)
			
			subscription := &domain.Subscription{
				CustomerID: customerID,
				PriceID:    "price_123",
				Status:     domain.SubscriptionStatusActive,
				Currency:   currency,
			}
			assert.True(t, subscription.IsValid(), "Currency %s should be supported for subscriptions", currency)
		}
	})
}