package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"theraclosure/payments-service/internal/adapters/logging"
)

const (
	// Maximum age for webhook signatures (5 minutes)
	MaxWebhookAge = 5 * time.Minute

	// Stripe signature header format: t=timestamp,v1=signature1,v1=signature2
	StripeSignatureHeader = "Stripe-Signature"
)

// WebhookVerifier provides webhook signature verification
type WebhookVerifier struct {
	webhookSecret string
	logger        *logging.Logger
}

// NewWebhookVerifier creates a new webhook signature verifier
func NewWebhookVerifier(webhookSecret string, logger *logging.Logger) *WebhookVerifier {
	return &WebhookVerifier{
		webhookSecret: webhookSecret,
		logger:        logger,
	}
}

// VerifyStripeSignature returns middleware that verifies Stripe webhook signatures
func (wv *WebhookVerifier) VerifyStripeSignature() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only verify for Stripe webhook endpoints
		if !strings.HasPrefix(c.Request.URL.Path, "/api/v1/webhooks/stripe") {
			c.Next()
			return
		}

		signature := c.GetHeader(StripeSignatureHeader)
		if signature == "" {
			wv.logger.WithFields(logging.Fields{
				{Key: "path", Value: c.Request.URL.Path},
				{Key: "ip", Value: c.ClientIP()},
			}).Warn("Missing Stripe signature header")

			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "MISSING_SIGNATURE",
				"message": "Stripe signature header is required",
			})
			c.Abort()
			return
		}

		// Read the request body
		body, err := c.GetRawData()
		if err != nil {
			wv.logger.WithFields(logging.Fields{
				{Key: "error", Value: err.Error()},
				{Key: "ip", Value: c.ClientIP()},
			}).Error("Failed to read webhook body")

			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "INVALID_REQUEST",
				"message": "Failed to read request body",
			})
			c.Abort()
			return
		}

		// Verify the signature
		if !wv.verifyStripeSignature(body, signature) {
			wv.logger.WithFields(logging.Fields{
				{Key: "signature", Value: signature},
				{Key: "ip", Value: c.ClientIP()},
			}).Error("Invalid Stripe webhook signature")

			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "INVALID_SIGNATURE",
				"message": "Webhook signature verification failed",
			})
			c.Abort()
			return
		}

		// Store the body for the handler to use
		c.Set("webhook_body", body)

		wv.logger.WithFields(logging.Fields{
			{Key: "path", Value: c.Request.URL.Path},
			{Key: "body_size", Value: len(body)},
		}).Info("Stripe webhook signature verified successfully")

		c.Next()
	}
}

// verifyStripeSignature verifies the Stripe webhook signature
func (wv *WebhookVerifier) verifyStripeSignature(payload []byte, signature string) bool {
	// Parse the signature header
	timestamp, signatures, err := wv.parseStripeSignature(signature)
	if err != nil {
		wv.logger.WithFields(logging.Fields{
			{Key: "error", Value: err.Error()},
			{Key: "signature", Value: signature},
		}).Error("Failed to parse Stripe signature")
		return false
	}

	// Check timestamp freshness (replay attack protection)
	if time.Since(time.Unix(timestamp, 0)) > MaxWebhookAge {
		wv.logger.WithFields(logging.Fields{
			{Key: "timestamp", Value: timestamp},
			{Key: "age", Value: time.Since(time.Unix(timestamp, 0)).String()},
		}).Warn("Webhook timestamp too old")
		return false
	}

	// Create the expected signature
	expectedSig := wv.computeSignature(timestamp, payload)

	// Compare with provided signatures
	for _, sig := range signatures {
		if hmac.Equal([]byte(expectedSig), []byte(sig)) {
			return true
		}
	}

	wv.logger.WithFields(logging.Fields{
		{Key: "expected", Value: expectedSig},
		{Key: "provided", Value: signatures},
	}).Error("Signature mismatch")

	return false
}

// parseStripeSignature parses the Stripe-Signature header
func (wv *WebhookVerifier) parseStripeSignature(signature string) (int64, []string, error) {
	var timestamp int64
	var signatures []string

	pairs := strings.Split(signature, ",")
	for _, pair := range pairs {
		parts := strings.Split(pair, "=")
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "t":
			ts, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return 0, nil, fmt.Errorf("invalid timestamp: %v", err)
			}
			timestamp = ts
		case "v1":
			signatures = append(signatures, value)
		}
	}

	if timestamp == 0 {
		return 0, nil, fmt.Errorf("no timestamp found in signature")
	}

	if len(signatures) == 0 {
		return 0, nil, fmt.Errorf("no signatures found")
	}

	return timestamp, signatures, nil
}

// computeSignature computes the expected signature for the webhook
func (wv *WebhookVerifier) computeSignature(timestamp int64, payload []byte) string {
	// Create the signed payload string
	signedPayload := fmt.Sprintf("%d.%s", timestamp, string(payload))

	// Compute HMAC-SHA256
	h := hmac.New(sha256.New, []byte(wv.webhookSecret))
	h.Write([]byte(signedPayload))

	return hex.EncodeToString(h.Sum(nil))
}

// VerifyGenericWebhook provides generic webhook verification for other services
func (wv *WebhookVerifier) VerifyGenericWebhook(secretHeader, signatureHeader string) gin.HandlerFunc {
	return func(c *gin.Context) {
		expectedSecret := c.GetHeader(secretHeader)
		providedSignature := c.GetHeader(signatureHeader)

		if expectedSecret == "" || providedSignature == "" {
			wv.logger.WithFields(logging.Fields{
				{Key: "path", Value: c.Request.URL.Path},
				{Key: "ip", Value: c.ClientIP()},
			}).Warn("Missing webhook authentication headers")

			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "MISSING_AUTHENTICATION",
				"message": "Webhook authentication headers required",
			})
			c.Abort()
			return
		}

		// Simple secret comparison for generic webhooks
		if !hmac.Equal([]byte(expectedSecret), []byte(wv.webhookSecret)) {
			wv.logger.WithFields(logging.Fields{
				{Key: "path", Value: c.Request.URL.Path},
				{Key: "ip", Value: c.ClientIP()},
			}).Error("Invalid webhook secret")

			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "INVALID_SECRET",
				"message": "Webhook secret verification failed",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
