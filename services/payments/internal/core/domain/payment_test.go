package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCustomer_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		customer Customer
		expected bool
	}{
		{
			name: "valid customer",
			customer: Customer{
				UserID: uuid.New().String(),
				Email:  "test@example.com",
				Name:   "Test User",
			},
			expected: true,
		},
		{
			name: "missing user ID",
			customer: Customer{
				Email: "test@example.com",
				Name:  "Test User",
			},
			expected: false,
		},
		{
			name: "missing email",
			customer: Customer{
				UserID: uuid.New().String(),
				Name:   "Test User",
			},
			expected: false,
		},
		{
			name: "missing name",
			customer: Customer{
				UserID: uuid.New().String(),
				Email:  "test@example.com",
			},
			expected: false,
		},
		{
			name:     "empty customer",
			customer: Customer{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.customer.IsValid())
		})
	}
}

func TestCustomer_BeforeCreate(t *testing.T) {
	customer := &Customer{
		UserID: uuid.New().String(),
		Email:  "test@example.com",
		Name:   "Test User",
	}

	// ID should be empty initially
	assert.Empty(t, customer.ID)

	// Call BeforeCreate
	err := customer.BeforeCreate(nil)
	require.NoError(t, err)

	// ID should now be set
	assert.NotEmpty(t, customer.ID)
	
	// Verify it's a valid UUID
	_, err = uuid.Parse(customer.ID)
	assert.NoError(t, err)
}

func TestCustomer_TableName(t *testing.T) {
	customer := Customer{}
	assert.Equal(t, "customers", customer.TableName())
}

func TestSubscription_IsValid(t *testing.T) {
	customerID := uuid.New().String()

	tests := []struct {
		name         string
		subscription Subscription
		expected     bool
	}{
		{
			name: "valid subscription",
			subscription: Subscription{
				CustomerID: customerID,
				PriceID:    "price_123",
				Status:     SubscriptionStatusActive,
			},
			expected: true,
		},
		{
			name: "missing customer ID",
			subscription: Subscription{
				PriceID: "price_123",
				Status:  SubscriptionStatusActive,
			},
			expected: false,
		},
		{
			name: "missing price ID",
			subscription: Subscription{
				CustomerID: customerID,
				Status:     SubscriptionStatusActive,
			},
			expected: false,
		},
		{
			name: "missing status",
			subscription: Subscription{
				CustomerID: customerID,
				PriceID:    "price_123",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.subscription.IsValid())
		})
	}
}

func TestSubscription_BeforeCreate(t *testing.T) {
	subscription := &Subscription{
		CustomerID: uuid.New().String(),
		PriceID:    "price_123",
		Status:     SubscriptionStatusActive,
	}

	// ID should be empty initially
	assert.Empty(t, subscription.ID)

	// Call BeforeCreate
	err := subscription.BeforeCreate(nil)
	require.NoError(t, err)

	// ID should now be set
	assert.NotEmpty(t, subscription.ID)
	
	// Verify it's a valid UUID
	_, err = uuid.Parse(subscription.ID)
	assert.NoError(t, err)
}

func TestSubscription_TableName(t *testing.T) {
	subscription := Subscription{}
	assert.Equal(t, "subscriptions", subscription.TableName())
}

func TestPayment_IsValid(t *testing.T) {
	customerID := uuid.New().String()

	tests := []struct {
		name     string
		payment  Payment
		expected bool
	}{
		{
			name: "valid payment",
			payment: Payment{
				CustomerID: customerID,
				Amount:     1000,
				Currency:   "usd",
			},
			expected: true,
		},
		{
			name: "missing customer ID",
			payment: Payment{
				Amount:   1000,
				Currency: "usd",
			},
			expected: false,
		},
		{
			name: "zero amount",
			payment: Payment{
				CustomerID: customerID,
				Amount:     0,
				Currency:   "usd",
			},
			expected: false,
		},
		{
			name: "negative amount",
			payment: Payment{
				CustomerID: customerID,
				Amount:     -100,
				Currency:   "usd",
			},
			expected: false,
		},
		{
			name: "missing currency",
			payment: Payment{
				CustomerID: customerID,
				Amount:     1000,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.payment.IsValid())
		})
	}
}

func TestPayment_BeforeCreate(t *testing.T) {
	payment := &Payment{
		CustomerID: uuid.New().String(),
		Amount:     1000,
		Currency:   "usd",
	}

	// ID should be empty initially
	assert.Empty(t, payment.ID)

	// Call BeforeCreate
	err := payment.BeforeCreate(nil)
	require.NoError(t, err)

	// ID should now be set
	assert.NotEmpty(t, payment.ID)
	
	// Verify it's a valid UUID
	_, err = uuid.Parse(payment.ID)
	assert.NoError(t, err)
}

func TestPayment_TableName(t *testing.T) {
	payment := Payment{}
	assert.Equal(t, "payments", payment.TableName())
}

func TestSubscriptionStatus_Constants(t *testing.T) {
	// Test that all status constants are defined correctly
	assert.Equal(t, SubscriptionStatus("active"), SubscriptionStatusActive)
	assert.Equal(t, SubscriptionStatus("canceled"), SubscriptionStatusCanceled)
	assert.Equal(t, SubscriptionStatus("trialing"), SubscriptionStatusTrialing)
	assert.Equal(t, SubscriptionStatus("past_due"), SubscriptionStatusPastDue)
	assert.Equal(t, SubscriptionStatus("unpaid"), SubscriptionStatusUnpaid)
	assert.Equal(t, SubscriptionStatus("paused"), SubscriptionStatusPaused)
}

func TestPaymentStatus_Constants(t *testing.T) {
	// Test that all status constants are defined correctly
	assert.Equal(t, PaymentStatus("pending"), PaymentStatusPending)
	assert.Equal(t, PaymentStatus("succeeded"), PaymentStatusSucceeded)
	assert.Equal(t, PaymentStatus("failed"), PaymentStatusFailed)
	assert.Equal(t, PaymentStatus("refunded"), PaymentStatusRefunded)
	assert.Equal(t, PaymentStatus("canceled"), PaymentStatusCanceled)
}

func TestCustomer_Relationships(t *testing.T) {
	now := time.Now()
	customer := Customer{
		ID:        uuid.New().String(),
		UserID:    uuid.New().String(),
		StripeID:  "cus_test123",
		Email:     "test@example.com",
		Name:      "Test User",
		Active:    true,
		CreatedAt: now,
		UpdatedAt: now,
		Subscriptions: []Subscription{
			{
				ID:         uuid.New().String(),
				PriceID:    "price_123",
				Status:     SubscriptionStatusActive,
				CreatedAt:  now,
				UpdatedAt:  now,
			},
		},
		Payments: []Payment{
			{
				ID:        uuid.New().String(),
				Amount:    1000,
				Currency:  "usd",
				Status:    PaymentStatusSucceeded,
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
	}

	// Test relationships are properly set
	assert.Len(t, customer.Subscriptions, 1)
	assert.Len(t, customer.Payments, 1)
	assert.Equal(t, SubscriptionStatusActive, customer.Subscriptions[0].Status)
	assert.Equal(t, PaymentStatusSucceeded, customer.Payments[0].Status)
}

func TestSubscription_Relationships(t *testing.T) {
	now := time.Now()
	customerID := uuid.New().String()
	subscriptionID := uuid.New().String()

	subscription := Subscription{
		ID:         subscriptionID,
		CustomerID: customerID,
		PriceID:    "price_123",
		Status:     SubscriptionStatusActive,
		Amount:     2000,
		Currency:   "usd",
		CreatedAt:  now,
		UpdatedAt:  now,
		Customer: Customer{
			ID:     customerID,
			UserID: uuid.New().String(),
			Email:  "test@example.com",
			Name:   "Test User",
		},
		Payments: []Payment{
			{
				ID:             uuid.New().String(),
				CustomerID:     customerID,
				SubscriptionID: &subscriptionID,
				Amount:         2000,
				Currency:       "usd",
				Status:         PaymentStatusSucceeded,
				CreatedAt:      now,
				UpdatedAt:      now,
			},
		},
	}

	// Test relationships
	assert.Equal(t, customerID, subscription.Customer.ID)
	assert.Len(t, subscription.Payments, 1)
	assert.Equal(t, &subscriptionID, subscription.Payments[0].SubscriptionID)
}