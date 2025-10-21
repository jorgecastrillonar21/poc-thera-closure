package stripe

import (
	"encoding/json"
	"fmt"

	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/paymentintent"
	"github.com/stripe/stripe-go/v72/refund"
	"github.com/stripe/stripe-go/v72/sub"
	"github.com/stripe/stripe-go/v72/webhook"
)

type StripeClient struct {
	secretKey string
}

func NewStripeClient(secretKey string) *StripeClient {
	stripe.Key = secretKey
	return &StripeClient{
		secretKey: secretKey,
	}
}

func (s *StripeClient) CreateCustomer(email, name string) (string, error) {
	params := &stripe.CustomerParams{
		Email: stripe.String(email),
		Name:  stripe.String(name),
	}

	c, err := customer.New(params)
	if err != nil {
		return "", fmt.Errorf("failed to create customer: %w", err)
	}

	return c.ID, nil
}

func (s *StripeClient) GetCustomer(stripeID string) (map[string]interface{}, error) {
	c, err := customer.Get(stripeID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	// Convert to map
	data, err := json.Marshal(c)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal customer: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal customer: %w", err)
	}

	return result, nil
}

func (s *StripeClient) UpdateCustomer(stripeID string, params map[string]interface{}) error {
	updateParams := &stripe.CustomerParams{}

	if email, ok := params["email"].(string); ok {
		updateParams.Email = stripe.String(email)
	}
	if name, ok := params["name"].(string); ok {
		updateParams.Name = stripe.String(name)
	}

	_, err := customer.Update(stripeID, updateParams)
	if err != nil {
		return fmt.Errorf("failed to update customer: %w", err)
	}

	return nil
}

func (s *StripeClient) DeleteCustomer(stripeID string) error {
	_, err := customer.Del(stripeID, nil)
	if err != nil {
		return fmt.Errorf("failed to delete customer: %w", err)
	}
	return nil
}

func (s *StripeClient) CreateSubscription(customerID, priceID string, trialDays *int) (map[string]interface{}, error) {
	params := &stripe.SubscriptionParams{
		Customer: stripe.String(customerID),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Price: stripe.String(priceID),
			},
		},
	}

	if trialDays != nil && *trialDays > 0 {
		params.TrialPeriodDays = stripe.Int64(int64(*trialDays))
	}

	subscription, err := sub.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	// Convert to map
	data, err := json.Marshal(subscription)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal subscription: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal subscription: %w", err)
	}

	return result, nil
}

func (s *StripeClient) GetSubscription(subscriptionID string) (map[string]interface{}, error) {
	subscription, err := sub.Get(subscriptionID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	// Convert to map
	data, err := json.Marshal(subscription)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal subscription: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal subscription: %w", err)
	}

	return result, nil
}

func (s *StripeClient) UpdateSubscription(subscriptionID string, params map[string]interface{}) (map[string]interface{}, error) {
	updateParams := &stripe.SubscriptionParams{}

	if priceID, ok := params["price_id"].(string); ok {
		updateParams.Items = []*stripe.SubscriptionItemsParams{
			{
				Price: stripe.String(priceID),
			},
		}
	}

	if cancelAt, ok := params["cancel_at"].(int64); ok {
		updateParams.CancelAt = stripe.Int64(cancelAt)
	}

	subscription, err := sub.Update(subscriptionID, updateParams)
	if err != nil {
		return nil, fmt.Errorf("failed to update subscription: %w", err)
	}

	// Convert to map
	data, err := json.Marshal(subscription)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal subscription: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal subscription: %w", err)
	}

	return result, nil
}

func (s *StripeClient) CancelSubscription(subscriptionID string, cancelAtPeriodEnd bool) (map[string]interface{}, error) {
	// Use nil params for immediate cancellation, or implement with proper Stripe API
	subscription, err := sub.Cancel(subscriptionID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to cancel subscription: %w", err)
	}

	// Convert to map
	data, err := json.Marshal(subscription)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal subscription: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal subscription: %w", err)
	}

	return result, nil
}

func (s *StripeClient) CreatePaymentIntent(amount int64, currency, customerID string, metadata map[string]string) (map[string]interface{}, error) {
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amount),
		Currency: stripe.String(currency),
		Customer: stripe.String(customerID),
	}

	if metadata != nil {
		params.Metadata = metadata
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment intent: %w", err)
	}

	// Convert to map
	data, err := json.Marshal(pi)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payment intent: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payment intent: %w", err)
	}

	return result, nil
}

func (s *StripeClient) GetPaymentIntent(paymentIntentID string) (map[string]interface{}, error) {
	pi, err := paymentintent.Get(paymentIntentID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment intent: %w", err)
	}

	// Convert to map
	data, err := json.Marshal(pi)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payment intent: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payment intent: %w", err)
	}

	return result, nil
}

func (s *StripeClient) ConfirmPaymentIntent(paymentIntentID string) (map[string]interface{}, error) {
	pi, err := paymentintent.Confirm(paymentIntentID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm payment intent: %w", err)
	}

	// Convert to map
	data, err := json.Marshal(pi)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payment intent: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payment intent: %w", err)
	}

	return result, nil
}

func (s *StripeClient) RefundPayment(paymentIntentID string, amount *int64) (map[string]interface{}, error) {
	params := &stripe.RefundParams{
		PaymentIntent: stripe.String(paymentIntentID),
	}

	if amount != nil {
		params.Amount = stripe.Int64(*amount)
	}

	r, err := refund.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create refund: %w", err)
	}

	// Convert to map
	data, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal refund: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal refund: %w", err)
	}

	return result, nil
}

func (s *StripeClient) ConstructEvent(payload []byte, signature, webhookSecret string) (map[string]interface{}, error) {
	event, err := webhook.ConstructEvent(payload, signature, webhookSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to construct webhook event: %w", err)
	}

	// Convert to map
	data, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal event: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal event: %w", err)
	}

	return result, nil
}