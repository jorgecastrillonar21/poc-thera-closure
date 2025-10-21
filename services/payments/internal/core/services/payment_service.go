package services


import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"theraclosure/payments-service/internal/core/domain"
	"theraclosure/payments-service/internal/core/ports"
)

type paymentService struct {
	customerRepo     ports.CustomerRepository
	subscriptionRepo ports.SubscriptionRepository
	paymentRepo      ports.PaymentRepository
	stripeClient     ports.StripeClient
}

// NewPaymentService creates a new payment service instance
func NewPaymentService(
	customerRepo ports.CustomerRepository,
	subscriptionRepo ports.SubscriptionRepository,
	paymentRepo ports.PaymentRepository,
	stripeClient ports.StripeClient,
) ports.PaymentService {
	return &paymentService{
		customerRepo:     customerRepo,
		subscriptionRepo: subscriptionRepo,
		paymentRepo:      paymentRepo,
		stripeClient:     stripeClient,
	}
}

// Customer operations
func (s *paymentService) CreateCustomer(ctx context.Context, req ports.CreateCustomerRequest) (*domain.Customer, error) {
	// Check if customer already exists for this user
	existingCustomer, err := s.customerRepo.GetByUserID(ctx, req.UserID)
	if err == nil && existingCustomer != nil {
		return nil, fmt.Errorf("customer already exists for user ID: %s", req.UserID)
	}

	// Create customer in Stripe
	stripeID, err := s.stripeClient.CreateCustomer(req.Email, req.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to create customer in Stripe: %w", err)
	}

	// Create customer in database
	customer := &domain.Customer{
		UserID:   req.UserID,
		StripeID: stripeID,
		Email:    req.Email,
		Name:     req.Name,
		Active:   true,
	}

	if !customer.IsValid() {
		return nil, fmt.Errorf("invalid customer data")
	}

	if err := s.customerRepo.Create(ctx, customer); err != nil {
		// Try to clean up Stripe customer if database creation fails
		s.stripeClient.DeleteCustomer(stripeID)
		return nil, fmt.Errorf("failed to create customer in database: %w", err)
	}

	return customer, nil
}

func (s *paymentService) GetCustomer(ctx context.Context, id string) (*domain.Customer, error) {
	return s.customerRepo.GetByID(ctx, id)
}

func (s *paymentService) GetCustomerByUserID(ctx context.Context, userID string) (*domain.Customer, error) {
	return s.customerRepo.GetByUserID(ctx, userID)
}

func (s *paymentService) UpdateCustomer(ctx context.Context, id string, req ports.UpdateCustomerRequest) (*domain.Customer, error) {
	customer, err := s.customerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("customer not found: %w", err)
	}

	// Update fields if provided
	if req.Email != "" {
		customer.Email = req.Email
	}
	if req.Name != "" {
		customer.Name = req.Name
	}

	// Update in Stripe if we have a Stripe ID
	if customer.StripeID != "" {
		params := map[string]interface{}{}
		if req.Email != "" {
			params["email"] = req.Email
		}
		if req.Name != "" {
			params["name"] = req.Name
		}

		if len(params) > 0 {
			if err := s.stripeClient.UpdateCustomer(customer.StripeID, params); err != nil {
				return nil, fmt.Errorf("failed to update customer in Stripe: %w", err)
			}
		}
	}

	// Update in database
	if err := s.customerRepo.Update(ctx, customer); err != nil {
		return nil, fmt.Errorf("failed to update customer in database: %w", err)
	}

	return customer, nil
}

func (s *paymentService) DeleteCustomer(ctx context.Context, id string) error {
	customer, err := s.customerRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("customer not found: %w", err)
	}

	// Delete from Stripe if we have a Stripe ID
	if customer.StripeID != "" {
		if err := s.stripeClient.DeleteCustomer(customer.StripeID); err != nil {
			return fmt.Errorf("failed to delete customer from Stripe: %w", err)
		}
	}

	// Delete from database
	if err := s.customerRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete customer from database: %w", err)
	}

	return nil
}

func (s *paymentService) ListCustomers(ctx context.Context, req ports.ListCustomersRequest) (*ports.ListCustomersResponse, error) {
	if req.Limit == 0 {
		req.Limit = 10
	}

	customers, total, err := s.customerRepo.List(ctx, req.Offset, req.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list customers: %w", err)
	}

	return &ports.ListCustomersResponse{
		Customers: customers,
		Total:     total,
		Offset:    req.Offset,
		Limit:     req.Limit,
	}, nil
}

// Subscription operations
func (s *paymentService) CreateSubscription(ctx context.Context, req ports.CreateSubscriptionRequest) (*domain.Subscription, error) {
	// Get customer
	customer, err := s.customerRepo.GetByID(ctx, req.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("customer not found: %w", err)
	}

	if customer.StripeID == "" {
		return nil, fmt.Errorf("customer does not have a Stripe ID")
	}

	// Create subscription in Stripe
	stripeResult, err := s.stripeClient.CreateSubscription(customer.StripeID, req.PriceID, req.TrialDays)
	if err != nil {
		return nil, fmt.Errorf("failed to create subscription in Stripe: %w", err)
	}

	// Parse Stripe response
	stripeID, ok := stripeResult["id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid Stripe subscription response")
	}

	// Extract subscription details from Stripe response
	subscription := &domain.Subscription{
		CustomerID: req.CustomerID,
		StripeID:   stripeID,
		PriceID:    req.PriceID,
		Status:     domain.SubscriptionStatusActive, // Will be updated based on Stripe status
		Active:     true,
		Currency:   "usd",
	}

	// Parse additional fields from Stripe response
	if status, ok := stripeResult["status"].(string); ok {
		subscription.Status = domain.SubscriptionStatus(status)
	}
	if currentPeriodStart, ok := stripeResult["current_period_start"].(float64); ok {
		subscription.CurrentPeriodStart = time.Unix(int64(currentPeriodStart), 0)
	}
	if currentPeriodEnd, ok := stripeResult["current_period_end"].(float64); ok {
		subscription.CurrentPeriodEnd = time.Unix(int64(currentPeriodEnd), 0)
	}

	if !subscription.IsValid() {
		return nil, fmt.Errorf("invalid subscription data")
	}

	// Create subscription in database
	if err := s.subscriptionRepo.Create(ctx, subscription); err != nil {
		// Try to cancel Stripe subscription if database creation fails
		s.stripeClient.CancelSubscription(stripeID, false)
		return nil, fmt.Errorf("failed to create subscription in database: %w", err)
	}

	return subscription, nil
}

func (s *paymentService) GetSubscription(ctx context.Context, id string) (*domain.Subscription, error) {
	return s.subscriptionRepo.GetByID(ctx, id)
}

func (s *paymentService) UpdateSubscription(ctx context.Context, id string, req ports.UpdateSubscriptionRequest) (*domain.Subscription, error) {
	subscription, err := s.subscriptionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("subscription not found: %w", err)
	}

	// Update in Stripe if we have a Stripe ID
	if subscription.StripeID != "" {
		params := map[string]interface{}{}
		if req.PriceID != "" {
			params["price_id"] = req.PriceID
		}
		if req.CancelAt != nil {
			params["cancel_at"] = *req.CancelAt
		}

		if len(params) > 0 {
			stripeResult, err := s.stripeClient.UpdateSubscription(subscription.StripeID, params)
			if err != nil {
				return nil, fmt.Errorf("failed to update subscription in Stripe: %w", err)
			}

			// Update local subscription with Stripe response
			if status, ok := stripeResult["status"].(string); ok {
				subscription.Status = domain.SubscriptionStatus(status)
			}
		}
	}

	// Update fields if provided
	if req.PriceID != "" {
		subscription.PriceID = req.PriceID
	}
	if req.Status != "" {
		subscription.Status = req.Status
	}
	if req.CancelAt != nil {
		cancelTime := time.Unix(*req.CancelAt, 0)
		subscription.CancelAt = &cancelTime
	}

	// Update in database
	if err := s.subscriptionRepo.Update(ctx, subscription); err != nil {
		return nil, fmt.Errorf("failed to update subscription in database: %w", err)
	}

	return subscription, nil
}

func (s *paymentService) CancelSubscription(ctx context.Context, id string) (*domain.Subscription, error) {
	subscription, err := s.subscriptionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("subscription not found: %w", err)
	}

	// Cancel in Stripe if we have a Stripe ID
	if subscription.StripeID != "" {
		_, err := s.stripeClient.CancelSubscription(subscription.StripeID, true)
		if err != nil {
			return nil, fmt.Errorf("failed to cancel subscription in Stripe: %w", err)
		}
	}

	// Update subscription status
	subscription.Status = domain.SubscriptionStatusCanceled
	now := time.Now()
	subscription.CanceledAt = &now

	// Update in database
	if err := s.subscriptionRepo.Update(ctx, subscription); err != nil {
		return nil, fmt.Errorf("failed to update subscription in database: %w", err)
	}

	return subscription, nil
}

func (s *paymentService) ListSubscriptions(ctx context.Context, req ports.ListSubscriptionsRequest) (*ports.ListSubscriptionsResponse, error) {
	if req.Limit == 0 {
		req.Limit = 10
	}

	subscriptions, total, err := s.subscriptionRepo.List(ctx, req.Offset, req.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list subscriptions: %w", err)
	}

	return &ports.ListSubscriptionsResponse{
		Subscriptions: subscriptions,
		Total:         total,
		Offset:        req.Offset,
		Limit:         req.Limit,
	}, nil
}

func (s *paymentService) GetCustomerSubscriptions(ctx context.Context, customerID string) ([]*domain.Subscription, error) {
	return s.subscriptionRepo.GetByCustomerID(ctx, customerID)
}

// Payment operations
func (s *paymentService) CreatePayment(ctx context.Context, req ports.CreatePaymentRequest) (*domain.Payment, error) {
	// Get customer to validate existence
	_, err := s.customerRepo.GetByID(ctx, req.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("customer not found: %w", err)
	}

	// Create payment in database
	payment := &domain.Payment{
		CustomerID:     req.CustomerID,
		SubscriptionID: req.SubscriptionID,
		Amount:         req.Amount,
		Currency:       req.Currency,
		Description:    req.Description,
		Status:         domain.PaymentStatusPending,
	}

	if req.Metadata != nil {
		metadataBytes, _ := json.Marshal(req.Metadata)
		payment.Metadata = string(metadataBytes)
	}

	if !payment.IsValid() {
		return nil, fmt.Errorf("invalid payment data")
	}

	if err := s.paymentRepo.Create(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to create payment in database: %w", err)
	}

	return payment, nil
}

func (s *paymentService) GetPayment(ctx context.Context, id string) (*domain.Payment, error) {
	return s.paymentRepo.GetByID(ctx, id)
}

func (s *paymentService) ListPayments(ctx context.Context, req ports.ListPaymentsRequest) (*ports.ListPaymentsResponse, error) {
	if req.Limit == 0 {
		req.Limit = 10
	}

	payments, total, err := s.paymentRepo.List(ctx, req.Offset, req.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list payments: %w", err)
	}

	return &ports.ListPaymentsResponse{
		Payments: payments,
		Total:    total,
		Offset:   req.Offset,
		Limit:    req.Limit,
	}, nil
}

func (s *paymentService) RefundPayment(ctx context.Context, id string, amount *int64) (*domain.Payment, error) {
	payment, err := s.paymentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("payment not found: %w", err)
	}

	if payment.PaymentIntentID == "" {
		return nil, fmt.Errorf("payment does not have a payment intent ID")
	}

	// Create refund in Stripe
	_, err = s.stripeClient.RefundPayment(payment.PaymentIntentID, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to create refund in Stripe: %w", err)
	}

	// Update payment status
	payment.Status = domain.PaymentStatusRefunded
	now := time.Now()
	payment.RefundedAt = &now

	// Update in database
	if err := s.paymentRepo.Update(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to update payment in database: %w", err)
	}

	return payment, nil
}

// Stripe operations
func (s *paymentService) CreatePaymentIntent(ctx context.Context, req ports.CreatePaymentIntentRequest) (*ports.CreatePaymentIntentResponse, error) {
	// Get customer
	customer, err := s.customerRepo.GetByID(ctx, req.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("customer not found: %w", err)
	}

	if customer.StripeID == "" {
		return nil, fmt.Errorf("customer does not have a Stripe ID")
	}

	// Create payment intent in Stripe
	stripeResult, err := s.stripeClient.CreatePaymentIntent(
		req.Amount,
		req.Currency,
		customer.StripeID,
		req.Metadata,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment intent in Stripe: %w", err)
	}

	// Parse response
	paymentIntentID, ok := stripeResult["id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid payment intent response")
	}

	clientSecret, ok := stripeResult["client_secret"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid payment intent response: missing client secret")
	}

	status, ok := stripeResult["status"].(string)
	if !ok {
		status = "requires_payment_method"
	}

	// Create payment record in database
	payment := &domain.Payment{
		CustomerID:      req.CustomerID,
		PaymentIntentID: paymentIntentID,
		Amount:          req.Amount,
		Currency:        req.Currency,
		Description:     req.Description,
		Status:          domain.PaymentStatusPending,
	}

	if req.Metadata != nil {
		metadataBytes, _ := json.Marshal(req.Metadata)
		payment.Metadata = string(metadataBytes)
	}

	if err := s.paymentRepo.Create(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to create payment record: %w", err)
	}

	return &ports.CreatePaymentIntentResponse{
		PaymentIntentID: paymentIntentID,
		ClientSecret:    clientSecret,
		Status:          status,
	}, nil
}

func (s *paymentService) ConfirmPaymentIntent(ctx context.Context, paymentIntentID string) (*ports.ConfirmPaymentIntentResponse, error) {
	// Confirm payment intent in Stripe
	stripeResult, err := s.stripeClient.ConfirmPaymentIntent(paymentIntentID)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm payment intent in Stripe: %w", err)
	}

	status, ok := stripeResult["status"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid payment intent confirmation response")
	}

	// Find and update payment record
	payment, err := s.paymentRepo.GetByStripeID(ctx, paymentIntentID)
	if err == nil && payment != nil {
		if status == "succeeded" {
			payment.Status = domain.PaymentStatusSucceeded
			now := time.Now()
			payment.ProcessedAt = &now
		} else if status == "failed" {
			payment.Status = domain.PaymentStatusFailed
		}

		s.paymentRepo.Update(ctx, payment)
	}

	response := &ports.ConfirmPaymentIntentResponse{
		PaymentIntentID: paymentIntentID,
		Status:          status,
	}

	if payment != nil {
		response.PaymentID = payment.ID
	}

	return response, nil
}

func (s *paymentService) HandleWebhook(ctx context.Context, payload []byte, signature string) error {
	// TODO: Get webhook secret from configuration
	webhookSecret := "whsec_test_secret" // This should come from config
	
	// Construct and validate the webhook event
	event, err := s.stripeClient.ConstructEvent(payload, signature, webhookSecret)
	if err != nil {
		return fmt.Errorf("failed to verify webhook signature: %w", err)
	}

	// Process different event types
	eventType, ok := event["type"].(string)
	if !ok {
		return fmt.Errorf("invalid event type")
	}

	switch eventType {
	case "payment_intent.succeeded":
		return s.handlePaymentIntentSucceeded(ctx, event)
	case "payment_intent.payment_failed":
		return s.handlePaymentIntentFailed(ctx, event)
	case "invoice.payment_succeeded":
		return s.handleInvoicePaymentSucceeded(ctx, event)
	case "invoice.payment_failed":
		return s.handleInvoicePaymentFailed(ctx, event)
	case "customer.subscription.created":
		return s.handleSubscriptionCreated(ctx, event)
	case "customer.subscription.updated":
		return s.handleSubscriptionUpdated(ctx, event)
	case "customer.subscription.deleted":
		return s.handleSubscriptionDeleted(ctx, event)
	default:
		// Log unhandled event types but don't fail
		return nil
	}
}

func (s *paymentService) handlePaymentIntentSucceeded(ctx context.Context, event map[string]interface{}) error {
	// Extract payment intent data from event
	data, ok := event["data"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid event data")
	}
	
	object, ok := data["object"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid event object")
	}

	paymentIntentID, ok := object["id"].(string)
	if !ok {
		return fmt.Errorf("missing payment intent ID")
	}

	// Find and update payment record
	payment, err := s.paymentRepo.GetByStripeID(ctx, paymentIntentID)
	if err != nil {
		// Payment might not exist in our system, which is okay
		return nil
	}

	// Update payment status
	payment.Status = domain.PaymentStatusSucceeded
	now := time.Now()
	payment.ProcessedAt = &now

	return s.paymentRepo.Update(ctx, payment)
}

func (s *paymentService) handlePaymentIntentFailed(ctx context.Context, event map[string]interface{}) error {
	// Extract payment intent data from event
	data, ok := event["data"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid event data")
	}
	
	object, ok := data["object"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid event object")
	}

	paymentIntentID, ok := object["id"].(string)
	if !ok {
		return fmt.Errorf("missing payment intent ID")
	}

	// Find and update payment record
	payment, err := s.paymentRepo.GetByStripeID(ctx, paymentIntentID)
	if err != nil {
		// Payment might not exist in our system, which is okay
		return nil
	}

	// Update payment status
	payment.Status = domain.PaymentStatusFailed

	return s.paymentRepo.Update(ctx, payment)
}

func (s *paymentService) handleInvoicePaymentSucceeded(ctx context.Context, event map[string]interface{}) error {
	// Handle subscription payment success
	return s.processInvoicePayment(ctx, event, domain.PaymentStatusSucceeded)
}

func (s *paymentService) handleInvoicePaymentFailed(ctx context.Context, event map[string]interface{}) error {
	// Handle subscription payment failure
	return s.processInvoicePayment(ctx, event, domain.PaymentStatusFailed)
}

func (s *paymentService) processInvoicePayment(ctx context.Context, event map[string]interface{}, status domain.PaymentStatus) error {
	// Extract invoice data
	data, ok := event["data"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid event data")
	}
	
	object, ok := data["object"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid event object")
	}

	subscriptionID, ok := object["subscription"].(string)
	if !ok {
		return nil // Not a subscription payment
	}

	// Find subscription
	subscription, err := s.subscriptionRepo.GetByStripeID(ctx, subscriptionID)
	if err != nil {
		return nil // Subscription might not exist in our system
	}

	// Create payment record for the invoice
	amount, _ := object["amount_paid"].(float64)
	currency, _ := object["currency"].(string)
	
	payment := &domain.Payment{
		CustomerID:     subscription.CustomerID,
		SubscriptionID: &subscription.ID,
		Amount:         int64(amount),
		Currency:       currency,
		Description:    "Subscription payment",
		Status:         status,
	}

	if status == domain.PaymentStatusSucceeded {
		now := time.Now()
		payment.ProcessedAt = &now
	}

	return s.paymentRepo.Create(ctx, payment)
}

func (s *paymentService) handleSubscriptionCreated(ctx context.Context, event map[string]interface{}) error {
	return s.syncSubscriptionFromWebhook(ctx, event)
}

func (s *paymentService) handleSubscriptionUpdated(ctx context.Context, event map[string]interface{}) error {
	return s.syncSubscriptionFromWebhook(ctx, event)
}

func (s *paymentService) handleSubscriptionDeleted(ctx context.Context, event map[string]interface{}) error {
	// Extract subscription data
	data, ok := event["data"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid event data")
	}
	
	object, ok := data["object"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid event object")
	}

	subscriptionID, ok := object["id"].(string)
	if !ok {
		return fmt.Errorf("missing subscription ID")
	}

	// Find and update subscription
	subscription, err := s.subscriptionRepo.GetByStripeID(ctx, subscriptionID)
	if err != nil {
		return nil // Subscription might not exist in our system
	}

	// Update subscription status
	subscription.Status = domain.SubscriptionStatusCanceled
	now := time.Now()
	subscription.CanceledAt = &now

	return s.subscriptionRepo.Update(ctx, subscription)
}

func (s *paymentService) syncSubscriptionFromWebhook(ctx context.Context, event map[string]interface{}) error {
	// Extract subscription data
	data, ok := event["data"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid event data")
	}
	
	object, ok := data["object"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid event object")
	}

	subscriptionID, ok := object["id"].(string)
	if !ok {
		return fmt.Errorf("missing subscription ID")
	}

	// Find subscription
	subscription, err := s.subscriptionRepo.GetByStripeID(ctx, subscriptionID)
	if err != nil {
		return nil // Subscription might not exist in our system
	}

	// Update subscription with Stripe data
	if status, ok := object["status"].(string); ok {
		subscription.Status = domain.SubscriptionStatus(status)
	}

	if currentPeriodStart, ok := object["current_period_start"].(float64); ok {
		subscription.CurrentPeriodStart = time.Unix(int64(currentPeriodStart), 0)
	}

	if currentPeriodEnd, ok := object["current_period_end"].(float64); ok {
		subscription.CurrentPeriodEnd = time.Unix(int64(currentPeriodEnd), 0)
	}

	if cancelAt, ok := object["cancel_at"].(float64); ok && cancelAt > 0 {
		cancelTime := time.Unix(int64(cancelAt), 0)
		subscription.CancelAt = &cancelTime
	}

	if canceledAt, ok := object["canceled_at"].(float64); ok && canceledAt > 0 {
		cancelTime := time.Unix(int64(canceledAt), 0)
		subscription.CanceledAt = &cancelTime
	}

	return s.subscriptionRepo.Update(ctx, subscription)
}

func (s *paymentService) HealthCheck(ctx context.Context) error {
	// Simple health check - could be expanded to check database connectivity, Stripe API, etc.
	return nil
}