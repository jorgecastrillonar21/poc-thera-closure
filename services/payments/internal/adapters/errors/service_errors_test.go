package errors

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServiceError_Creation(t *testing.T) {
	t.Run("Basic ServiceError creation", func(t *testing.T) {
		err := NewServiceError(ErrCodeValidation, "Invalid input")

		assert.Equal(t, ErrCodeValidation, err.Code)
		assert.Equal(t, "Invalid input", err.Message)
		assert.Equal(t, 400, err.HTTPStatus)
		assert.WithinDuration(t, time.Now(), err.Timestamp, time.Second)
	})

	t.Run("ServiceError with details", func(t *testing.T) {
		err := NewServiceError(ErrCodeCustomerNotFound, "Customer not found").
			WithDetails("Customer ID: 123 does not exist")

		assert.Equal(t, "Customer ID: 123 does not exist", err.Details)
	})

	t.Run("ServiceError with context", func(t *testing.T) {
		err := NewServiceError(ErrCodePaymentFailed, "Payment failed").
			WithContext("payment_id", "pay_123").
			WithContext("amount", 1000)

		assert.Equal(t, "pay_123", err.Context["payment_id"])
		assert.Equal(t, 1000, err.Context["amount"])
	})

	t.Run("ServiceError with cause", func(t *testing.T) {
		originalErr := errors.New("database connection failed")
		err := NewServiceError(ErrCodeDatabaseConnection, "Database error").
			WithCause(originalErr)

		assert.Equal(t, originalErr, err.Cause)
		assert.Equal(t, originalErr, errors.Unwrap(err))
	})
}

func TestServiceError_Methods(t *testing.T) {
	t.Run("Error method", func(t *testing.T) {
		err := NewServiceError(ErrCodeValidation, "Validation failed")
		expected := "VALIDATION_ERROR: Validation failed"
		assert.Equal(t, expected, err.Error())

		errWithDetails := err.WithDetails("Field 'email' is required")
		expectedWithDetails := "VALIDATION_ERROR: Validation failed - Field 'email' is required"
		assert.Equal(t, expectedWithDetails, errWithDetails.Error())
	})

	t.Run("Is method", func(t *testing.T) {
		err1 := NewServiceError(ErrCodeValidation, "Validation failed")
		err2 := NewServiceError(ErrCodeValidation, "Another validation error")
		err3 := NewServiceError(ErrCodeNotFound, "Not found")

		assert.True(t, errors.Is(err1, err2))
		assert.False(t, errors.Is(err1, err3))
	})

	t.Run("ToErrorResponse method", func(t *testing.T) {
		err := NewServiceError(ErrCodePaymentFailed, "Payment processing failed").
			WithDetails("Card declined").
			WithContext("payment_id", "pay_123").
			WithRequestID("req_456")

		response := err.ToErrorResponse()

		assert.Equal(t, ErrCodePaymentFailed, response.Error)
		assert.Equal(t, "Payment processing failed", response.Message)
		assert.Equal(t, "Card declined", response.Details)
		assert.Equal(t, "pay_123", response.Context["payment_id"])
		assert.Equal(t, "req_456", response.RequestID)
	})
}

func TestErrorCreationHelpers(t *testing.T) {
	t.Run("NewValidationError", func(t *testing.T) {
		err := NewValidationError("email", "is required")

		assert.Equal(t, ErrCodeValidation, err.Code)
		assert.Equal(t, "Validation failed", err.Message)
		assert.Contains(t, err.Details, "email")
		assert.Contains(t, err.Details, "is required")
		assert.Equal(t, "email", err.Context["field"])
	})

	t.Run("NewNotFoundError", func(t *testing.T) {
		err := NewNotFoundError("Customer", "123")

		assert.Equal(t, ErrCodeNotFound, err.Code)
		assert.Contains(t, err.Message, "Customer not found")
		assert.Equal(t, "Customer", err.Context["resource"])
		assert.Equal(t, "123", err.Context["id"])
	})

	t.Run("NewCustomerError", func(t *testing.T) {
		err := NewCustomerError(ErrCodeCustomerInactive, "cust_123")

		assert.Equal(t, ErrCodeCustomerInactive, err.Code)
		assert.Equal(t, "Customer account is inactive", err.Message)
		assert.Equal(t, "cust_123", err.Context["customer_id"])
	})

	t.Run("NewPaymentError", func(t *testing.T) {
		err := NewPaymentError(ErrCodeInvalidAmount, "pay_123")

		assert.Equal(t, ErrCodeInvalidAmount, err.Code)
		assert.Equal(t, "Invalid payment amount", err.Message)
		assert.Equal(t, "pay_123", err.Context["payment_id"])
	})

	t.Run("NewStripeError", func(t *testing.T) {
		originalErr := errors.New("card_declined")
		err := NewStripeError(ErrCodeStripeCardDeclined, originalErr)

		assert.Equal(t, ErrCodeStripeCardDeclined, err.Code)
		assert.Equal(t, "Card was declined", err.Message)
		assert.Equal(t, originalErr, err.Cause)
		assert.Equal(t, "card_declined", err.Context["stripe_error"])
	})

	t.Run("NewDatabaseError", func(t *testing.T) {
		originalErr := errors.New("connection timeout")
		err := NewDatabaseError(ErrCodeDatabaseConnection, "SELECT customers", originalErr)

		assert.Equal(t, ErrCodeDatabaseConnection, err.Code)
		assert.Equal(t, "Database connection error", err.Message)
		assert.Equal(t, originalErr, err.Cause)
		assert.Equal(t, "SELECT customers", err.Context["operation"])
		assert.Equal(t, "connection timeout", err.Context["database_error"])
	})
}

func TestHTTPStatusMapping(t *testing.T) {
	testCases := []struct {
		errorCode      ErrorCode
		expectedStatus int
	}{
		{ErrCodeValidation, 400},
		{ErrCodeUnauthorized, 401},
		{ErrCodeForbidden, 403},
		{ErrCodeNotFound, 404},
		{ErrCodeConflict, 409},
		{ErrCodeRateLimit, 429},
		{ErrCodePaymentFailed, 402},
		{ErrCodeInternal, 500},
	}

	for _, tc := range testCases {
		t.Run(string(tc.errorCode), func(t *testing.T) {
			err := NewServiceError(tc.errorCode, "Test message")
			assert.Equal(t, tc.expectedStatus, err.HTTPStatus)
		})
	}
}

func TestErrorConstants(t *testing.T) {
	// Test that all error codes are properly defined
	errorCodes := []ErrorCode{
		ErrCodeInternal,
		ErrCodeValidation,
		ErrCodeNotFound,
		ErrCodeCustomerNotFound,
		ErrCodePaymentFailed,
		ErrCodeStripeAPI,
		ErrCodeDatabaseConnection,
	}

	for _, code := range errorCodes {
		t.Run(string(code), func(t *testing.T) {
			assert.NotEmpty(t, string(code))

			// Test that each error code produces a valid HTTP status
			err := NewServiceError(code, "Test")
			assert.Greater(t, err.HTTPStatus, 0)
			assert.Less(t, err.HTTPStatus, 600)
		})
	}
}
