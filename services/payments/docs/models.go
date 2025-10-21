package docs


// This file contains all the API models used in Swagger documentation
// These structs are used for generating proper OpenAPI specifications

// Customer Models

// CreateCustomerRequest represents the request body for creating a customer
type CreateCustomerRequest struct {
	UserID string `json:"user_id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Email  string `json:"email" binding:"required,email" example:"john.doe@example.com"`
	Name   string `json:"name" binding:"required" example:"John Doe"`
} // @name CreateCustomerRequest

// UpdateCustomerRequest represents the request body for updating a customer
type UpdateCustomerRequest struct {
	Email string `json:"email,omitempty" example:"john.doe@example.com"`
	Name  string `json:"name,omitempty" example:"John Doe"`
} // @name UpdateCustomerRequest

// CustomerResponse represents a customer in API responses
type CustomerResponse struct {
	ID                      string `json:"id" example:"cus_1234567890"`
	UserID                  string `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	StripeID                string `json:"stripe_id" example:"cus_stripe_1234567890"`
	Email                   string `json:"email" example:"john.doe@example.com"`
	Name                    string `json:"name" example:"John Doe"`
	DefaultPaymentMethodID  string `json:"default_payment_method_id,omitempty" example:"pm_1234567890"`
	Active                  bool   `json:"active" example:"true"`
	CreatedAt               string `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt               string `json:"updated_at" example:"2023-01-01T00:00:00Z"`
} // @name CustomerResponse

// ListCustomersResponse represents the response for listing customers
type ListCustomersResponse struct {
	Customers []CustomerResponse `json:"customers"`
	Total     int                `json:"total" example:"100"`
	Offset    int                `json:"offset" example:"0"`
	Limit     int                `json:"limit" example:"10"`
} // @name ListCustomersResponse

// Subscription Models

// CreateSubscriptionRequest represents the request body for creating a subscription
type CreateSubscriptionRequest struct {
	CustomerID string `json:"customer_id" binding:"required" example:"cus_1234567890"`
	PriceID    string `json:"price_id" binding:"required" example:"price_1234567890"`
	TrialDays  *int   `json:"trial_days,omitempty" example:"7"`
} // @name CreateSubscriptionRequest

// UpdateSubscriptionRequest represents the request body for updating a subscription
type UpdateSubscriptionRequest struct {
	PriceID string `json:"price_id,omitempty" example:"price_1234567890"`
} // @name UpdateSubscriptionRequest

// SubscriptionResponse represents a subscription in API responses
type SubscriptionResponse struct {
	ID                 string  `json:"id" example:"sub_1234567890"`
	CustomerID         string  `json:"customer_id" example:"cus_1234567890"`
	StripeID           string  `json:"stripe_id" example:"sub_stripe_1234567890"`
	PriceID            string  `json:"price_id" example:"price_1234567890"`
	Status             string  `json:"status" example:"active"`
	CurrentPeriodStart string  `json:"current_period_start" example:"2023-01-01T00:00:00Z"`
	CurrentPeriodEnd   string  `json:"current_period_end" example:"2023-02-01T00:00:00Z"`
	Amount             int64   `json:"amount" example:"2000"`
	Currency           string  `json:"currency" example:"usd"`
	TrialStart         *string `json:"trial_start,omitempty" example:"2023-01-01T00:00:00Z"`
	TrialEnd           *string `json:"trial_end,omitempty" example:"2023-01-08T00:00:00Z"`
	Active             bool    `json:"active" example:"true"`
	CreatedAt          string  `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt          string  `json:"updated_at" example:"2023-01-01T00:00:00Z"`
} // @name SubscriptionResponse

// ListSubscriptionsResponse represents the response for listing subscriptions
type ListSubscriptionsResponse struct {
	Subscriptions []SubscriptionResponse `json:"subscriptions"`
	Total         int                    `json:"total" example:"50"`
	Offset        int                    `json:"offset" example:"0"`
	Limit         int                    `json:"limit" example:"10"`
} // @name ListSubscriptionsResponse

// Payment Models

// CreatePaymentRequest represents the request body for creating a payment
type CreatePaymentRequest struct {
	CustomerID     string  `json:"customer_id" binding:"required" example:"cus_1234567890"`
	Amount         int64   `json:"amount" binding:"required,min=1" example:"1000"`
	Currency       string  `json:"currency" binding:"required" example:"usd"`
	Description    string  `json:"description,omitempty" example:"Payment for services"`
	SubscriptionID *string `json:"subscription_id,omitempty" example:"sub_1234567890"`
} // @name CreatePaymentRequest

// RefundPaymentRequest represents the request body for refunding a payment
type RefundPaymentRequest struct {
	Amount string `json:"amount,omitempty" example:"500"`
	Reason string `json:"reason,omitempty" example:"requested_by_customer"`
} // @name RefundPaymentRequest

// PaymentResponse represents a payment in API responses
type PaymentResponse struct {
	ID             string  `json:"id" example:"pay_1234567890"`
	CustomerID     string  `json:"customer_id" example:"cus_1234567890"`
	SubscriptionID *string `json:"subscription_id,omitempty" example:"sub_1234567890"`
	StripeID       string  `json:"stripe_id" example:"pi_stripe_1234567890"`
	Amount         int64   `json:"amount" example:"1000"`
	Currency       string  `json:"currency" example:"usd"`
	Status         string  `json:"status" example:"succeeded"`
	Description    string  `json:"description,omitempty" example:"Payment for services"`
	CreatedAt      string  `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt      string  `json:"updated_at" example:"2023-01-01T00:00:00Z"`
} // @name PaymentResponse

// ListPaymentsResponse represents the response for listing payments
type ListPaymentsResponse struct {
	Payments []PaymentResponse `json:"payments"`
	Total    int               `json:"total" example:"200"`
	Offset   int               `json:"offset" example:"0"`
	Limit    int               `json:"limit" example:"10"`
} // @name ListPaymentsResponse

// Payment Intent Models

// CreatePaymentIntentRequest represents the request body for creating a payment intent
type CreatePaymentIntentRequest struct {
	CustomerID      string            `json:"customer_id" binding:"required" example:"cus_1234567890"`
	Amount          int64             `json:"amount" binding:"required,min=1" example:"1000"`
	Currency        string            `json:"currency" binding:"required" example:"usd"`
	PaymentMethodID string            `json:"payment_method_id,omitempty" example:"pm_1234567890"`
	Metadata        map[string]string `json:"metadata,omitempty"`
} // @name CreatePaymentIntentRequest

// PaymentIntentResponse represents a payment intent in API responses
type PaymentIntentResponse struct {
	ID           string            `json:"id" example:"pi_1234567890"`
	CustomerID   string            `json:"customer_id" example:"cus_1234567890"`
	Amount       int64             `json:"amount" example:"1000"`
	Currency     string            `json:"currency" example:"usd"`
	Status       string            `json:"status" example:"requires_confirmation"`
	ClientSecret string            `json:"client_secret" example:"pi_1234567890_secret_xyz"`
	Metadata     map[string]string `json:"metadata,omitempty"`
	CreatedAt    string            `json:"created_at" example:"2023-01-01T00:00:00Z"`
} // @name PaymentIntentResponse

// ConfirmPaymentIntentRequest represents the request body for confirming a payment intent
type ConfirmPaymentIntentRequest struct {
	PaymentMethodID string `json:"payment_method_id,omitempty" example:"pm_1234567890"`
	ReturnURL       string `json:"return_url,omitempty" example:"https://example.com/return"`
} // @name ConfirmPaymentIntentRequest

// Health Check Models

// HealthResponse represents the basic health check response
type HealthResponse struct {
	Status    string `json:"status" example:"healthy"`
	Timestamp string `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	Service   string `json:"service" example:"payments-service"`
} // @name HealthResponse

// DetailedHealthResponse represents the detailed health check response
type DetailedHealthResponse struct {
	Status     string                     `json:"status" example:"healthy"`
	Version    string                     `json:"version" example:"1.0.0"`
	Timestamp  string                     `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	Components map[string]ComponentHealth `json:"components"`
} // @name DetailedHealthResponse

// ComponentHealth represents the health status of a component
type ComponentHealth struct {
	Status       string `json:"status" example:"healthy"`
	Message      string `json:"message,omitempty" example:"Database connection successful"`
	LastChecked  string `json:"last_checked" example:"2023-01-01T00:00:00Z"`
	ResponseTime string `json:"response_time,omitempty" example:"5ms"`
} // @name ComponentHealth

// Error Models

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error     string            `json:"error" example:"VALIDATION_ERROR"`
	Message   string            `json:"message" example:"Validation failed"`
	Details   string            `json:"details,omitempty" example:"Field 'email' is required"`
	Context   map[string]interface{} `json:"context,omitempty"`
	Timestamp string            `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	RequestID string            `json:"request_id,omitempty" example:"req_1234567890"`
} // @name ErrorResponse

// Common Query Parameters

// PaginationParams represents pagination query parameters
type PaginationParams struct {
	Offset int `form:"offset" example:"0"`
	Limit  int `form:"limit" example:"10"`
} // @name PaginationParams

// CustomerFilters represents customer filtering parameters
type CustomerFilters struct {
	Active *bool   `form:"active" example:"true"`
	Email  string  `form:"email" example:"john@example.com"`
	UserID string  `form:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
} // @name CustomerFilters

// SubscriptionFilters represents subscription filtering parameters
type SubscriptionFilters struct {
	CustomerID string `form:"customer_id" example:"cus_1234567890"`
	Status     string `form:"status" example:"active"`
	Active     *bool  `form:"active" example:"true"`
} // @name SubscriptionFilters

// PaymentFilters represents payment filtering parameters
type PaymentFilters struct {
	CustomerID     string `form:"customer_id" example:"cus_1234567890"`
	SubscriptionID string `form:"subscription_id" example:"sub_1234567890"`
	Status         string `form:"status" example:"succeeded"`
} // @name PaymentFilters