package ports

import (
	"context"

	"theraclosure/payments-service/internal/core/domain"
)

// PaymentService defines the interface for payment operations
type PaymentService interface {
	// Customer operations
	CreateCustomer(ctx context.Context, req CreateCustomerRequest) (*domain.Customer, error)
	GetCustomer(ctx context.Context, id string) (*domain.Customer, error)
	GetCustomerByUserID(ctx context.Context, userID string) (*domain.Customer, error)
	UpdateCustomer(ctx context.Context, id string, req UpdateCustomerRequest) (*domain.Customer, error)
	DeleteCustomer(ctx context.Context, id string) error
	ListCustomers(ctx context.Context, req ListCustomersRequest) (*ListCustomersResponse, error)

	// Subscription operations
	CreateSubscription(ctx context.Context, req CreateSubscriptionRequest) (*domain.Subscription, error)
	GetSubscription(ctx context.Context, id string) (*domain.Subscription, error)
	UpdateSubscription(ctx context.Context, id string, req UpdateSubscriptionRequest) (*domain.Subscription, error)
	CancelSubscription(ctx context.Context, id string) (*domain.Subscription, error)
	ListSubscriptions(ctx context.Context, req ListSubscriptionsRequest) (*ListSubscriptionsResponse, error)
	GetCustomerSubscriptions(ctx context.Context, customerID string) ([]*domain.Subscription, error)

	// Payment operations
	CreatePayment(ctx context.Context, req CreatePaymentRequest) (*domain.Payment, error)
	GetPayment(ctx context.Context, id string) (*domain.Payment, error)
	ListPayments(ctx context.Context, req ListPaymentsRequest) (*ListPaymentsResponse, error)
	RefundPayment(ctx context.Context, id string, amount *int64) (*domain.Payment, error)

	// Stripe operations
	CreatePaymentIntent(ctx context.Context, req CreatePaymentIntentRequest) (*CreatePaymentIntentResponse, error)
	ConfirmPaymentIntent(ctx context.Context, paymentIntentID string) (*ConfirmPaymentIntentResponse, error)
	HandleWebhook(ctx context.Context, payload []byte, signature string) error

	// Health check
	HealthCheck(ctx context.Context) error
}

// CustomerRepository defines the interface for customer data operations
type CustomerRepository interface {
	Create(ctx context.Context, customer *domain.Customer) error
	GetByID(ctx context.Context, id string) (*domain.Customer, error)
	GetByUserID(ctx context.Context, userID string) (*domain.Customer, error)
	GetByStripeID(ctx context.Context, stripeID string) (*domain.Customer, error)
	Update(ctx context.Context, customer *domain.Customer) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, offset, limit int) ([]*domain.Customer, int64, error)
}

// SubscriptionRepository defines the interface for subscription data operations
type SubscriptionRepository interface {
	Create(ctx context.Context, subscription *domain.Subscription) error
	GetByID(ctx context.Context, id string) (*domain.Subscription, error)
	GetByStripeID(ctx context.Context, stripeID string) (*domain.Subscription, error)
	GetByCustomerID(ctx context.Context, customerID string) ([]*domain.Subscription, error)
	Update(ctx context.Context, subscription *domain.Subscription) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, offset, limit int) ([]*domain.Subscription, int64, error)
}

// PaymentRepository defines the interface for payment data operations
type PaymentRepository interface {
	Create(ctx context.Context, payment *domain.Payment) error
	GetByID(ctx context.Context, id string) (*domain.Payment, error)
	GetByStripeID(ctx context.Context, stripeID string) (*domain.Payment, error)
	GetByCustomerID(ctx context.Context, customerID string) ([]*domain.Payment, error)
	GetBySubscriptionID(ctx context.Context, subscriptionID string) ([]*domain.Payment, error)
	Update(ctx context.Context, payment *domain.Payment) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, offset, limit int) ([]*domain.Payment, int64, error)
}

// StripeClient defines the interface for Stripe operations
type StripeClient interface {
	CreateCustomer(email, name string) (string, error)
	GetCustomer(stripeID string) (map[string]interface{}, error)
	UpdateCustomer(stripeID string, params map[string]interface{}) error
	DeleteCustomer(stripeID string) error

	CreateSubscription(customerID, priceID string, trialDays *int) (map[string]interface{}, error)
	GetSubscription(subscriptionID string) (map[string]interface{}, error)
	UpdateSubscription(subscriptionID string, params map[string]interface{}) (map[string]interface{}, error)
	CancelSubscription(subscriptionID string, cancelAtPeriodEnd bool) (map[string]interface{}, error)

	CreatePaymentIntent(amount int64, currency, customerID string, metadata map[string]string) (map[string]interface{}, error)
	GetPaymentIntent(paymentIntentID string) (map[string]interface{}, error)
	ConfirmPaymentIntent(paymentIntentID string) (map[string]interface{}, error)
	RefundPayment(paymentIntentID string, amount *int64) (map[string]interface{}, error)

	ConstructEvent(payload []byte, signature, webhookSecret string) (map[string]interface{}, error)
}

// DTOs and Request/Response structs

type CreateCustomerRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Email  string `json:"email" binding:"required,email"`
	Name   string `json:"name" binding:"required"`
}

type UpdateCustomerRequest struct {
	Email string `json:"email,omitempty" binding:"omitempty,email"`
	Name  string `json:"name,omitempty"`
}

type ListCustomersRequest struct {
	Offset int `json:"offset" form:"offset"`
	Limit  int `json:"limit" form:"limit" binding:"min=1,max=100"`
}

type ListCustomersResponse struct {
	Customers []*domain.Customer `json:"customers"`
	Total     int64              `json:"total"`
	Offset    int                `json:"offset"`
	Limit     int                `json:"limit"`
}

type CreateSubscriptionRequest struct {
	CustomerID string `json:"customer_id" binding:"required"`
	PriceID    string `json:"price_id" binding:"required"`
	TrialDays  *int   `json:"trial_days,omitempty"`
}

type UpdateSubscriptionRequest struct {
	PriceID  string                     `json:"price_id,omitempty"`
	Status   domain.SubscriptionStatus  `json:"status,omitempty"`
	CancelAt *int64                     `json:"cancel_at,omitempty"` // Unix timestamp
}

type ListSubscriptionsRequest struct {
	CustomerID string `json:"customer_id" form:"customer_id"`
	Status     string `json:"status" form:"status"`
	Offset     int    `json:"offset" form:"offset"`
	Limit      int    `json:"limit" form:"limit" binding:"min=1,max=100"`
}

type ListSubscriptionsResponse struct {
	Subscriptions []*domain.Subscription `json:"subscriptions"`
	Total         int64                  `json:"total"`
	Offset        int                    `json:"offset"`
	Limit         int                    `json:"limit"`
}

type CreatePaymentRequest struct {
	CustomerID     string            `json:"customer_id" binding:"required"`
	SubscriptionID *string           `json:"subscription_id,omitempty"`
	Amount         int64             `json:"amount" binding:"required,min=1"`
	Currency       string            `json:"currency" binding:"required"`
	Description    string            `json:"description,omitempty"`
	Metadata       map[string]string `json:"metadata,omitempty"`
}

type CreatePaymentIntentRequest struct {
	CustomerID  string            `json:"customer_id" binding:"required"`
	Amount      int64             `json:"amount" binding:"required,min=1"`
	Currency    string            `json:"currency" binding:"required"`
	Description string            `json:"description,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

type CreatePaymentIntentResponse struct {
	PaymentIntentID string `json:"payment_intent_id"`
	ClientSecret    string `json:"client_secret"`
	Status          string `json:"status"`
}

type ConfirmPaymentIntentResponse struct {
	PaymentIntentID string `json:"payment_intent_id"`
	Status          string `json:"status"`
	PaymentID       string `json:"payment_id,omitempty"`
}

type ListPaymentsRequest struct {
	CustomerID     string `json:"customer_id" form:"customer_id"`
	SubscriptionID string `json:"subscription_id" form:"subscription_id"`
	Status         string `json:"status" form:"status"`
	Offset         int    `json:"offset" form:"offset"`
	Limit          int    `json:"limit" form:"limit" binding:"min=1,max=100"`
}

type ListPaymentsResponse struct {
	Payments []*domain.Payment `json:"payments"`
	Total    int64             `json:"total"`
	Offset   int               `json:"offset"`
	Limit    int               `json:"limit"`
}