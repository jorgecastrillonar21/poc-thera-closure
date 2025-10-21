package errors

import (
	"fmt"
	"net/http"
	"time"
)

// ErrorCode represents a unique error identifier
type ErrorCode string

// Predefined error codes for the payments service
const (
	// General errors
	ErrCodeInternal     ErrorCode = "INTERNAL_ERROR"
	ErrCodeValidation   ErrorCode = "VALIDATION_ERROR"
	ErrCodeNotFound     ErrorCode = "NOT_FOUND"
	ErrCodeUnauthorized ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden    ErrorCode = "FORBIDDEN"
	ErrCodeConflict     ErrorCode = "CONFLICT"
	ErrCodeRateLimit    ErrorCode = "RATE_LIMIT_EXCEEDED"

	// Customer errors
	ErrCodeCustomerNotFound ErrorCode = "CUSTOMER_NOT_FOUND"
	ErrCodeCustomerExists   ErrorCode = "CUSTOMER_EXISTS"
	ErrCodeInvalidCustomer  ErrorCode = "INVALID_CUSTOMER"
	ErrCodeCustomerInactive ErrorCode = "CUSTOMER_INACTIVE"

	// Subscription errors
	ErrCodeSubscriptionNotFound  ErrorCode = "SUBSCRIPTION_NOT_FOUND"
	ErrCodeSubscriptionExists    ErrorCode = "SUBSCRIPTION_EXISTS"
	ErrCodeInvalidSubscription   ErrorCode = "INVALID_SUBSCRIPTION"
	ErrCodeSubscriptionCancelled ErrorCode = "SUBSCRIPTION_CANCELLED"
	ErrCodeSubscriptionExpired   ErrorCode = "SUBSCRIPTION_EXPIRED"

	// Payment errors
	ErrCodePaymentNotFound         ErrorCode = "PAYMENT_NOT_FOUND"
	ErrCodePaymentFailed           ErrorCode = "PAYMENT_FAILED"
	ErrCodeInvalidAmount           ErrorCode = "INVALID_AMOUNT"
	ErrCodeInsufficientFunds       ErrorCode = "INSUFFICIENT_FUNDS"
	ErrCodePaymentMethodRequired   ErrorCode = "PAYMENT_METHOD_REQUIRED"
	ErrCodePaymentAlreadyProcessed ErrorCode = "PAYMENT_ALREADY_PROCESSED"

	// Stripe errors
	ErrCodeStripeAPI          ErrorCode = "STRIPE_API_ERROR"
	ErrCodeStripeCardDeclined ErrorCode = "STRIPE_CARD_DECLINED"
	ErrCodeStripeInvalidCard  ErrorCode = "STRIPE_INVALID_CARD"
	ErrCodeStripeWebhook      ErrorCode = "STRIPE_WEBHOOK_ERROR"

	// Database errors
	ErrCodeDatabaseConnection ErrorCode = "DATABASE_CONNECTION_ERROR"
	ErrCodeDatabaseQuery      ErrorCode = "DATABASE_QUERY_ERROR"
	ErrCodeDatabaseConstraint ErrorCode = "DATABASE_CONSTRAINT_ERROR"
)

// ServiceError represents a structured error with code, message, and context
type ServiceError struct {
	Code       ErrorCode              `json:"code"`
	Message    string                 `json:"message"`
	Details    string                 `json:"details,omitempty"`
	Context    map[string]interface{} `json:"context,omitempty"`
	Timestamp  time.Time              `json:"timestamp"`
	RequestID  string                 `json:"request_id,omitempty"`
	HTTPStatus int                    `json:"-"`
	Cause      error                  `json:"-"`
}

// Error implements the error interface
func (e *ServiceError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s - %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap allows error unwrapping for Go 1.13+ error handling
func (e *ServiceError) Unwrap() error {
	return e.Cause
}

// Is allows error comparison for Go 1.13+ error handling
func (e *ServiceError) Is(target error) bool {
	if target == nil {
		return false
	}

	if serviceErr, ok := target.(*ServiceError); ok {
		return e.Code == serviceErr.Code
	}

	return false
}

// NewServiceError creates a new service error
func NewServiceError(code ErrorCode, message string) *ServiceError {
	return &ServiceError{
		Code:       code,
		Message:    message,
		Timestamp:  time.Now(),
		HTTPStatus: getHTTPStatusForErrorCode(code),
	}
}

// WithDetails adds additional details to the error
func (e *ServiceError) WithDetails(details string) *ServiceError {
	e.Details = details
	return e
}

// WithContext adds context information to the error
func (e *ServiceError) WithContext(key string, value interface{}) *ServiceError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// WithCause adds the underlying cause of the error
func (e *ServiceError) WithCause(err error) *ServiceError {
	e.Cause = err
	return e
}

// WithRequestID adds request ID for tracing
func (e *ServiceError) WithRequestID(requestID string) *ServiceError {
	e.RequestID = requestID
	return e
}

// WithHTTPStatus overrides the default HTTP status
func (e *ServiceError) WithHTTPStatus(status int) *ServiceError {
	e.HTTPStatus = status
	return e
}

// Helper functions to create common errors

// NewValidationError creates a validation error
func NewValidationError(field string, reason string) *ServiceError {
	return NewServiceError(ErrCodeValidation, "Validation failed").
		WithDetails(fmt.Sprintf("Field '%s': %s", field, reason)).
		WithContext("field", field)
}

// NewNotFoundError creates a not found error
func NewNotFoundError(resource string, id string) *ServiceError {
	return NewServiceError(ErrCodeNotFound, fmt.Sprintf("%s not found", resource)).
		WithContext("resource", resource).
		WithContext("id", id)
}

// NewCustomerError creates customer-related errors
func NewCustomerError(code ErrorCode, customerID string) *ServiceError {
	messages := map[ErrorCode]string{
		ErrCodeCustomerNotFound: "Customer not found",
		ErrCodeCustomerExists:   "Customer already exists",
		ErrCodeInvalidCustomer:  "Invalid customer data",
		ErrCodeCustomerInactive: "Customer account is inactive",
	}

	message, exists := messages[code]
	if !exists {
		message = "Customer error"
	}

	return NewServiceError(code, message).
		WithContext("customer_id", customerID)
}

// NewPaymentError creates payment-related errors
func NewPaymentError(code ErrorCode, paymentID string) *ServiceError {
	messages := map[ErrorCode]string{
		ErrCodePaymentNotFound:         "Payment not found",
		ErrCodePaymentFailed:           "Payment processing failed",
		ErrCodeInvalidAmount:           "Invalid payment amount",
		ErrCodeInsufficientFunds:       "Insufficient funds",
		ErrCodePaymentMethodRequired:   "Payment method required",
		ErrCodePaymentAlreadyProcessed: "Payment already processed",
	}

	message, exists := messages[code]
	if !exists {
		message = "Payment error"
	}

	return NewServiceError(code, message).
		WithContext("payment_id", paymentID)
}

// NewStripeError creates Stripe-related errors
func NewStripeError(code ErrorCode, stripeErr error) *ServiceError {
	messages := map[ErrorCode]string{
		ErrCodeStripeAPI:          "Stripe API error",
		ErrCodeStripeCardDeclined: "Card was declined",
		ErrCodeStripeInvalidCard:  "Invalid card information",
		ErrCodeStripeWebhook:      "Stripe webhook processing error",
	}

	message, exists := messages[code]
	if !exists {
		message = "Stripe error"
	}

	return NewServiceError(code, message).
		WithCause(stripeErr).
		WithContext("stripe_error", stripeErr.Error())
}

// NewDatabaseError creates database-related errors
func NewDatabaseError(code ErrorCode, operation string, err error) *ServiceError {
	messages := map[ErrorCode]string{
		ErrCodeDatabaseConnection: "Database connection error",
		ErrCodeDatabaseQuery:      "Database query error",
		ErrCodeDatabaseConstraint: "Database constraint violation",
	}

	message, exists := messages[code]
	if !exists {
		message = "Database error"
	}

	return NewServiceError(code, message).
		WithCause(err).
		WithContext("operation", operation).
		WithContext("database_error", err.Error())
}

// getHTTPStatusForErrorCode maps error codes to HTTP status codes
func getHTTPStatusForErrorCode(code ErrorCode) int {
	switch code {
	case ErrCodeValidation, ErrCodeInvalidCustomer, ErrCodeInvalidSubscription, ErrCodeInvalidAmount:
		return http.StatusBadRequest
	case ErrCodeUnauthorized:
		return http.StatusUnauthorized
	case ErrCodeForbidden:
		return http.StatusForbidden
	case ErrCodeNotFound, ErrCodeCustomerNotFound, ErrCodeSubscriptionNotFound, ErrCodePaymentNotFound:
		return http.StatusNotFound
	case ErrCodeConflict, ErrCodeCustomerExists, ErrCodeSubscriptionExists, ErrCodePaymentAlreadyProcessed:
		return http.StatusConflict
	case ErrCodeRateLimit:
		return http.StatusTooManyRequests
	case ErrCodePaymentFailed, ErrCodeStripeCardDeclined, ErrCodeInsufficientFunds:
		return http.StatusPaymentRequired
	default:
		return http.StatusInternalServerError
	}
}

// ErrorResponse represents the JSON structure returned to clients
type ErrorResponse struct {
	Error     ErrorCode              `json:"error"`
	Message   string                 `json:"message"`
	Details   string                 `json:"details,omitempty"`
	Context   map[string]interface{} `json:"context,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	RequestID string                 `json:"request_id,omitempty"`
}

// ToErrorResponse converts a ServiceError to an ErrorResponse
func (e *ServiceError) ToErrorResponse() *ErrorResponse {
	return &ErrorResponse{
		Error:     e.Code,
		Message:   e.Message,
		Details:   e.Details,
		Context:   e.Context,
		Timestamp: e.Timestamp,
		RequestID: e.RequestID,
	}
}
