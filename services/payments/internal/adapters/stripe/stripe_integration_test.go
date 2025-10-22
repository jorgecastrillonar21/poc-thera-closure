package stripe

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stripe/stripe-go/v72"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test configuration
const (
	testEmail = "test@theraclosure.com"
	testName  = "Test Customer"
)

// getTestSecretKey returns the test secret key from environment or a placeholder
func getTestSecretKey() string {
	key := os.Getenv("STRIPE_TEST_SECRET_KEY")
	if key == "" {
		return "sk_test_fake_key_for_unit_tests"
	}
	return key
}

func setupTestClient() *StripeClient {
	return NewStripeClient(getTestSecretKey())
}

func TestNewStripeClient(t *testing.T) {
	testKey := getTestSecretKey()
	client := NewStripeClient(testKey)
	
	assert.NotNil(t, client)
	assert.Equal(t, testKey, client.secretKey)
	assert.Equal(t, testKey, stripe.Key)
}

func TestStripeClient_CreateCustomer(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	client := setupTestClient()
	
	tests := []struct {
		name        string
		email       string
		customerName string
		expectError bool
	}{
		{
			name:         "valid_customer",
			email:        fmt.Sprintf("test-%d@theraclosure.com", time.Now().Unix()),
			customerName: "Test Customer",
			expectError:  false,
		},
		{
			name:         "empty_email",
			email:        "",
			customerName: "Test Customer",
			expectError:  false, // Stripe allows empty email
		},
		{
			name:         "empty_name",
			email:        fmt.Sprintf("test-noname-%d@theraclosure.com", time.Now().Unix()),
			customerName: "",
			expectError:  false, // Stripe allows empty name
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			customerID, err := client.CreateCustomer(tt.email, tt.customerName)
			
			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, customerID)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, customerID)
				assert.Contains(t, customerID, "cus_")
			}
		})
	}
}

func TestStripeClient_GetCustomer(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	client := setupTestClient()
	
	// Create a customer first
	email := fmt.Sprintf("get-test-%d@theraclosure.com", time.Now().Unix())
	customerID, err := client.CreateCustomer(email, "Get Test Customer")
	require.NoError(t, err)
	require.NotEmpty(t, customerID)
	
	tests := []struct {
		name        string
		customerID  string
		expectError bool
	}{
		{
			name:        "valid_customer",
			customerID:  customerID,
			expectError: false,
		},
		{
			name:        "invalid_customer",
			customerID:  "cus_invalid",
			expectError: true,
		},
		{
			name:        "empty_customer_id",
			customerID:  "",
			expectError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			customer, err := client.GetCustomer(tt.customerID)
			
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, customer)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, customer)
				assert.Equal(t, customerID, customer["id"])
				assert.Equal(t, email, customer["email"])
			}
		})
	}
}

func TestStripeClient_UpdateCustomer(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	client := setupTestClient()
	
	// Create a customer first
	email := fmt.Sprintf("update-test-%d@theraclosure.com", time.Now().Unix())
	customerID, err := client.CreateCustomer(email, "Update Test Customer")
	require.NoError(t, err)
	require.NotEmpty(t, customerID)
	
	tests := []struct {
		name        string
		customerID  string
		params      map[string]interface{}
		expectError bool
	}{
		{
			name:       "update_name",
			customerID: customerID,
			params: map[string]interface{}{
				"name": "Updated Customer Name",
			},
			expectError: false,
		},
		{
			name:       "update_email",
			customerID: customerID,
			params: map[string]interface{}{
				"email": fmt.Sprintf("updated-%d@theraclosure.com", time.Now().Unix()),
			},
			expectError: false,
		},
		{
			name:       "update_metadata",
			customerID: customerID,
			params: map[string]interface{}{
				"metadata": map[string]string{
					"test_key": "test_value",
				},
			},
			expectError: false,
		},
		{
			name:        "invalid_customer",
			customerID:  "cus_invalid",
			params:      map[string]interface{}{"name": "Test"},
			expectError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.UpdateCustomer(tt.customerID, tt.params)
			
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestStripeClient_CreatePaymentIntent(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	client := setupTestClient()
	
	// Create a customer first
	email := fmt.Sprintf("payment-test-%d@theraclosure.com", time.Now().Unix())
	customerID, err := client.CreateCustomer(email, "Payment Test Customer")
	require.NoError(t, err)
	require.NotEmpty(t, customerID)
	
	tests := []struct {
		name        string
		amount      int64
		currency    string
		customerID  string
		metadata    map[string]string
		expectError bool
	}{
		{
			name:       "valid_payment_intent",
			amount:     2999, // $29.99
			currency:   "usd",
			customerID: customerID,
			metadata: map[string]string{
				"test": "true",
			},
			expectError: false,
		},
		{
			name:        "zero_amount",
			amount:      0,
			currency:    "usd",
			customerID:  customerID,
			metadata:    nil,
			expectError: true,
		},
		{
			name:        "negative_amount",
			amount:      -100,
			currency:    "usd",
			customerID:  customerID,
			metadata:    nil,
			expectError: true,
		},
		{
			name:        "invalid_currency",
			amount:      1000,
			currency:    "invalid",
			customerID:  customerID,
			metadata:    nil,
			expectError: true,
		},
		{
			name:        "invalid_customer",
			amount:      1000,
			currency:    "usd",
			customerID:  "cus_invalid",
			metadata:    nil,
			expectError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paymentIntent, err := client.CreatePaymentIntent(
				tt.amount,
				tt.currency,
				tt.customerID,
				tt.metadata,
			)
			
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, paymentIntent)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, paymentIntent)
				assert.Contains(t, paymentIntent["id"].(string), "pi_")
				assert.Equal(t, tt.amount, int64(paymentIntent["amount"].(float64)))
				assert.Equal(t, tt.currency, paymentIntent["currency"])
			}
		})
	}
}

func TestStripeClient_GetPaymentIntent(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	client := setupTestClient()
	
	// Create a customer and payment intent first
	email := fmt.Sprintf("get-payment-test-%d@theraclosure.com", time.Now().Unix())
	customerID, err := client.CreateCustomer(email, "Get Payment Test Customer")
	require.NoError(t, err)
	
	paymentIntent, err := client.CreatePaymentIntent(1999, "usd", customerID, nil)
	require.NoError(t, err)
	paymentIntentID := paymentIntent["id"].(string)
	
	tests := []struct {
		name              string
		paymentIntentID   string
		expectError       bool
	}{
		{
			name:            "valid_payment_intent",
			paymentIntentID: paymentIntentID,
			expectError:     false,
		},
		{
			name:            "invalid_payment_intent",
			paymentIntentID: "pi_invalid",
			expectError:     true,
		},
		{
			name:            "empty_payment_intent_id",
			paymentIntentID: "",
			expectError:     true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pi, err := client.GetPaymentIntent(tt.paymentIntentID)
			
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, pi)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, pi)
				assert.Equal(t, paymentIntentID, pi["id"])
				assert.Equal(t, int64(1999), int64(pi["amount"].(float64)))
			}
		})
	}
}

func TestStripeClient_RefundPayment(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	// Note: This test requires a successful payment to refund
	// In a real test environment, you would need to simulate a successful payment
	// For now, we'll test with invalid payment intents to check error handling
	
	client := setupTestClient()
	
	tests := []struct {
		name              string
		paymentIntentID   string
		amount            *int64
		expectError       bool
	}{
		{
			name:            "invalid_payment_intent",
			paymentIntentID: "pi_invalid",
			amount:          nil,
			expectError:     true,
		},
		{
			name:            "empty_payment_intent_id",
			paymentIntentID: "",
			amount:          nil,
			expectError:     true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			refund, err := client.RefundPayment(tt.paymentIntentID, tt.amount)
			
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, refund)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, refund)
			}
		})
	}
}

func TestStripeClient_CreateSubscription(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	client := setupTestClient()
	
	// Create a customer first
	email := fmt.Sprintf("sub-test-%d@theraclosure.com", time.Now().Unix())
	customerID, err := client.CreateCustomer(email, "Subscription Test Customer")
	require.NoError(t, err)
	
	// Use Stripe's test price ID for testing
	testPriceID := "price_1HslsBGezWG3dJr0ZK5VcryL" // This might need to be updated with a valid test price
	
	tests := []struct {
		name        string
		customerID  string
		priceID     string
		trialDays   *int
		expectError bool
	}{
		{
			name:        "invalid_price_id", // Testing with invalid price to avoid creating actual subscriptions
			customerID:  customerID,
			priceID:     "price_invalid",
			trialDays:   nil,
			expectError: true,
		},
		{
			name:        "invalid_customer",
			customerID:  "cus_invalid",
			priceID:     testPriceID,
			trialDays:   nil,
			expectError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subscription, err := client.CreateSubscription(tt.customerID, tt.priceID, tt.trialDays)
			
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, subscription)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, subscription)
				assert.Contains(t, subscription["id"].(string), "sub_")
			}
		})
	}
}

func TestStripeClient_ConstructEvent(t *testing.T) {
	client := setupTestClient()
	
	// Test webhook event construction
	testPayload := []byte(`{
		"id": "evt_test_webhook",
		"object": "event",
		"api_version": "2020-08-27",
		"created": 1326853478,
		"data": {
			"object": {
				"id": "cus_test_webhook",
				"object": "customer"
			}
		},
		"livemode": false,
		"pending_webhooks": 1,
		"request": {
			"id": null,
			"idempotency_key": null
		},
		"type": "customer.created"
	}`)
	
	tests := []struct {
		name          string
		payload       []byte
		signature     string
		webhookSecret string
		expectError   bool
	}{
		{
			name:          "invalid_signature",
			payload:       testPayload,
			signature:     "invalid_signature",
			webhookSecret: "whsec_test",
			expectError:   true,
		},
		{
			name:          "empty_payload",
			payload:       []byte{},
			signature:     "t=123,v1=signature",
			webhookSecret: "whsec_test",
			expectError:   true,
		},
		{
			name:          "empty_signature",
			payload:       testPayload,
			signature:     "",
			webhookSecret: "whsec_test",
			expectError:   true,
		},
		{
			name:          "empty_webhook_secret",
			payload:       testPayload,
			signature:     "t=123,v1=signature",
			webhookSecret: "",
			expectError:   true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event, err := client.ConstructEvent(tt.payload, tt.signature, tt.webhookSecret)
			
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, event)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, event)
			}
		})
	}
}

// Benchmark tests for performance
func BenchmarkStripeClient_CreateCustomer(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping benchmark in short mode")
	}
	
	client := setupTestClient()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		email := fmt.Sprintf("bench-%d-%d@theraclosure.com", time.Now().UnixNano(), i)
		_, err := client.CreateCustomer(email, "Benchmark Customer")
		if err != nil {
			b.Fatalf("Failed to create customer: %v", err)
		}
	}
}

func BenchmarkStripeClient_CreatePaymentIntent(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping benchmark in short mode")
	}
	
	client := setupTestClient()
	
	// Create a customer for the benchmark
	email := fmt.Sprintf("bench-payment-%d@theraclosure.com", time.Now().Unix())
	customerID, err := client.CreateCustomer(email, "Benchmark Payment Customer")
	if err != nil {
		b.Fatalf("Failed to create customer: %v", err)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.CreatePaymentIntent(1000, "usd", customerID, nil)
		if err != nil {
			b.Fatalf("Failed to create payment intent: %v", err)
		}
	}
}

// Helper function to skip tests if Stripe test key is not configured
func skipIfNoStripeKey(t *testing.T) {
	if os.Getenv("STRIPE_TEST_SECRET_KEY") == "" {
		t.Skip("STRIPE_TEST_SECRET_KEY environment variable not set, skipping Stripe integration test")
	}
}