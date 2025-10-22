package http

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	_ "theraclosure/payments-service/docs"
	"theraclosure/payments-service/internal/adapters/config"
	"theraclosure/payments-service/internal/adapters/http/middleware"
	"theraclosure/payments-service/internal/adapters/logging"
	"theraclosure/payments-service/internal/adapters/monitoring"
	_ "theraclosure/payments-service/internal/core/domain"
	"theraclosure/payments-service/internal/core/ports"
)

type Server struct {
	Router     *gin.Engine
	Config     *config.Config
	service    ports.PaymentService
	db         *gorm.DB
	redis      *redis.Client
	stripeKey  string
	logger     *logging.Logger
	monitoring *monitoring.Metrics
}

func NewServer(service ports.PaymentService, cfg *config.Config, db *gorm.DB, redis *redis.Client, logger *logging.Logger) *Server {
	// Initialize monitoring
	metrics := monitoring.NewMetrics("payments_service", "1.0.0")

	server := &Server{
		service:    service,
		Config:     cfg,
		db:         db,
		redis:      redis,
		stripeKey:  cfg.Stripe.SecretKey,
		logger:     logger,
		monitoring: metrics,
	}

	server.setupRouter()
	return server
}

func (s *Server) setupRouter() {
	// Set Gin mode based on config
	if s.Config.App.LogLevel == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	s.Router = gin.New()

	// Add performance monitoring middleware
	perfMonitor := middleware.NewPerformanceMonitor(s.monitoring, s.logger)
	s.Router.Use(perfMonitor.Middleware())

	// Add custom middleware for logging and error handling
	s.Router.Use(middleware.RequestIDMiddleware())
	s.Router.Use(middleware.RecoveryMiddleware(s.logger))
	s.Router.Use(middleware.LoggingMiddleware(s.logger))
	s.Router.Use(middleware.ErrorHandlingMiddleware(s.logger))
	s.Router.Use(middleware.SecurityLoggingMiddleware(s.logger))

	// Security middleware
	s.setupSecurityMiddleware()

	// CORS middleware
	corsConfig := cors.Config{
		AllowOrigins:     s.Config.CORS.AllowedOrigins,
		AllowMethods:     s.Config.CORS.AllowedMethods,
		AllowHeaders:     s.Config.CORS.AllowedHeaders,
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	s.Router.Use(cors.New(corsConfig))

	// Add Swagger endpoint
	s.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Add metrics endpoint
	s.Router.GET("/metrics", gin.WrapH(s.monitoring.GetHandler()))

	// Health check endpoints
	s.Router.GET("/health", s.healthCheck)                  // Basic health check
	s.Router.GET("/health/detailed", s.detailedHealthCheck) // Detailed health with components
	s.Router.GET("/health/ready", s.readinessProbe)         // Kubernetes readiness probe
	s.Router.GET("/health/live", s.livenessProbe)           // Kubernetes liveness probe

	// API v1 routes
	v1 := s.Router.Group("/api/v1")
	{
		// Health check endpoint
		v1.GET("/health", s.healthCheck)

		// Customer routes
		customers := v1.Group("/customers")
		{
			customers.POST("", s.createCustomer)
			customers.GET("", s.listCustomers)
			customers.GET("/:id", s.getCustomer)
			customers.PUT("/:id", s.updateCustomer)
			customers.DELETE("/:id", s.deleteCustomer)
			customers.GET("/user/:userID", s.getCustomerByUserID)
		}

		// Subscription routes
		subscriptions := v1.Group("/subscriptions")
		{
			subscriptions.POST("", s.createSubscription)
			subscriptions.GET("", s.listSubscriptions)
			subscriptions.GET("/:id", s.getSubscription)
			subscriptions.PUT("/:id", s.updateSubscription)
			subscriptions.DELETE("/:id", s.cancelSubscription)
		}

		// Payment routes
		payments := v1.Group("/payments")
		{
			payments.POST("", s.createPayment)
			payments.GET("", s.listPayments)
			payments.GET("/:id", s.getPayment)
			payments.POST("/:id/refund", s.refundPayment)
		}

		// Payment Intent routes
		paymentIntents := v1.Group("/payment-intents")
		{
			paymentIntents.POST("", s.createPaymentIntent)
			paymentIntents.POST("/:id/confirm", s.confirmPaymentIntent)
		}

		// Webhook endpoint
		v1.POST("/webhooks/stripe", s.handleStripeWebhook)
	}
}

func (s *Server) setupSecurityMiddleware() {
	// Request size middleware
	s.Router.Use(middleware.RequestSizeMiddleware(s.Config.Security.MaxRequestSize))

	// Rate limiting middleware
	rateLimiterConfig := &middleware.RateLimiterConfig{
		RequestsPerWindow: s.Config.Security.RateLimitRPS,
		WindowDuration:    parseRateLimitWindow(s.Config.Security.RateLimitWindow),
		SkipPaths: []string{
			"/health",
			"/swagger",
			"/metrics",
		},
		KeyGenerator: middleware.EndpointBasedKeyGenerator,
	}
	rateLimiter := middleware.NewRateLimiter(rateLimiterConfig, s.redis, s.logger)
	s.Router.Use(rateLimiter.Middleware())

	// Request validation middleware
	validationConfig := middleware.DefaultValidationConfig()
	validationConfig.MaxBodySize = s.Config.Security.MaxRequestSize
	validator := middleware.NewRequestValidator(validationConfig, s.logger)
	s.Router.Use(validator.Middleware())

	// Webhook signature verification (for Stripe webhooks)
	webhookVerifier := middleware.NewWebhookVerifier(s.Config.Stripe.WebhookSecret, s.logger)
	s.Router.Use(webhookVerifier.VerifyStripeSignature())

	// Optional authentication middleware (doesn't block unauthenticated requests by default)
	authConfig := &middleware.AuthConfig{
		JWTSecret:      s.Config.Security.JWTSecret,
		APIKeys:        s.Config.Security.APIKeys,
		BasicAuthUsers: s.Config.Security.BasicAuthUsers,
		SkipPaths: []string{
			"/health",
			"/swagger",
			"/metrics",
			"/api/v1/webhooks", // Webhooks use signature verification
		},
	}
	authenticator := middleware.NewAuthenticator(authConfig, s.logger)

	if s.Config.Security.EnableAuthentication {
		// Strict authentication (blocks unauthenticated requests)
		s.Router.Use(authenticator.JWTMiddleware())
	} else {
		// Optional authentication (allows unauthenticated requests)
		s.Router.Use(authenticator.OptionalAuthMiddleware())
	}
}

func parseRateLimitWindow(window string) time.Duration {
	if duration, err := time.ParseDuration(window); err == nil {
		return duration
	}
	// Default to 1 minute if parsing fails
	return time.Minute
}

func (s *Server) Start() error {
	return s.Router.Run(s.Config.GetServerAddress())
}

// AddSwaggerEndpoint adds Swagger documentation endpoint
func (s *Server) AddSwaggerEndpoint(handler gin.HandlerFunc) {
	s.Router.GET("/swagger/*any", handler)
}

// Customer handlers

// createCustomer godoc
// @Summary Create a new customer
// @Description Creates a new customer in the payments system with Stripe integration
// @Tags customers
// @Accept json
// @Produce json
// @Param customer body ports.CreateCustomerRequest true "Customer creation request"
// @Success 201 {object} domain.Customer "Customer created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body or validation error"
// @Failure 409 {object} map[string]interface{} "Customer already exists"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/customers [post]
func (s *Server) createCustomer(c *gin.Context) {
	var req ports.CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	customer, err := s.service.CreateCustomer(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create customer", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": customer})
}

// listCustomers godoc
// @Summary List customers
// @Description Retrieves a paginated list of customers with optional filtering
// @Tags customers
// @Accept json
// @Produce json
// @Param offset query int false "Number of records to skip" default(0)
// @Param limit query int false "Maximum number of records to return" default(10)
// @Param active query bool false "Filter by active status"
// @Param email query string false "Filter by email"
// @Param user_id query string false "Filter by user ID"
// @Success 200 {object} ports.ListCustomersResponse "List of customers"
// @Failure 400 {object} map[string]interface{} "Invalid query parameters"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/customers [get]
func (s *Server) listCustomers(c *gin.Context) {
	var req ports.ListCustomersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters", "details": err.Error()})
		return
	}

	if req.Limit == 0 {
		req.Limit = 10
	}

	response, err := s.service.ListCustomers(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list customers", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// getCustomer godoc
// @Summary Get customer by ID
// @Description Retrieves a specific customer by their ID
// @Tags customers
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Success 200 {object} domain.Customer "Customer details"
// @Failure 400 {object} map[string]interface{} "Invalid customer ID"
// @Failure 404 {object} map[string]interface{} "Customer not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/customers/{id} [get]
func (s *Server) getCustomer(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Customer ID is required"})
		return
	}

	customer, err := s.service.GetCustomer(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": customer})
}

// updateCustomer godoc
// @Summary Update customer
// @Description Updates an existing customer's information
// @Tags customers
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Param customer body ports.UpdateCustomerRequest true "Customer update request"
// @Success 200 {object} domain.Customer "Updated customer details"
// @Failure 400 {object} map[string]interface{} "Invalid request body or customer ID"
// @Failure 404 {object} map[string]interface{} "Customer not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/customers/{id} [put]
func (s *Server) updateCustomer(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Customer ID is required"})
		return
	}

	var req ports.UpdateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	customer, err := s.service.UpdateCustomer(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update customer", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": customer})
}

// deleteCustomer godoc
// @Summary Delete customer
// @Description Deletes a customer and cancels all their subscriptions
// @Tags customers
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Success 200 {object} map[string]string "Customer deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid customer ID"
// @Failure 404 {object} map[string]interface{} "Customer not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/customers/{id} [delete]
func (s *Server) deleteCustomer(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Customer ID is required"})
		return
	}

	err := s.service.DeleteCustomer(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete customer", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Customer deleted successfully"})
}

// getCustomerByUserID godoc
// @Summary Get customer by user ID
// @Description Retrieves a customer by their associated user ID
// @Tags customers
// @Accept json
// @Produce json
// @Param userID path string true "User ID"
// @Success 200 {object} domain.Customer "Customer details"
// @Failure 400 {object} map[string]interface{} "Invalid user ID"
// @Failure 404 {object} map[string]interface{} "Customer not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/customers/user/{userID} [get]
func (s *Server) getCustomerByUserID(c *gin.Context) {
	userID := c.Param("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	customer, err := s.service.GetCustomerByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": customer})
}

// Subscription handlers

// createSubscription godoc
// @Summary Create a new subscription
// @Description Creates a new subscription for a customer with Stripe integration
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body ports.CreateSubscriptionRequest true "Subscription creation request"
// @Success 201 {object} domain.Subscription "Subscription created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body or validation error"
// @Failure 404 {object} map[string]interface{} "Customer not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/subscriptions [post]
func (s *Server) createSubscription(c *gin.Context) {
	var req ports.CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	subscription, err := s.service.CreateSubscription(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subscription", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": subscription})
}

// listSubscriptions godoc
// @Summary List subscriptions
// @Description Retrieves a paginated list of subscriptions with optional filtering
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param offset query int false "Number of records to skip" default(0)
// @Param limit query int false "Maximum number of records to return" default(10)
// @Param customer_id query string false "Filter by customer ID"
// @Param status query string false "Filter by status (active, canceled, etc.)"
// @Param active query bool false "Filter by active status"
// @Success 200 {object} ports.ListSubscriptionsResponse "List of subscriptions"
// @Failure 400 {object} map[string]interface{} "Invalid query parameters"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/subscriptions [get]
func (s *Server) listSubscriptions(c *gin.Context) {
	var req ports.ListSubscriptionsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters", "details": err.Error()})
		return
	}

	if req.Limit == 0 {
		req.Limit = 10
	}

	response, err := s.service.ListSubscriptions(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list subscriptions", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// getSubscription godoc
// @Summary Get subscription by ID
// @Description Retrieves a specific subscription by its ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 200 {object} domain.Subscription "Subscription details"
// @Failure 400 {object} map[string]interface{} "Invalid subscription ID"
// @Failure 404 {object} map[string]interface{} "Subscription not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/subscriptions/{id} [get]
func (s *Server) getSubscription(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Subscription ID is required"})
		return
	}

	subscription, err := s.service.GetSubscription(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": subscription})
}

// updateSubscription godoc
// @Summary Update subscription
// @Description Updates an existing subscription (e.g., change price)
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID"
// @Param subscription body ports.UpdateSubscriptionRequest true "Subscription update request"
// @Success 200 {object} domain.Subscription "Updated subscription details"
// @Failure 400 {object} map[string]interface{} "Invalid request body or subscription ID"
// @Failure 404 {object} map[string]interface{} "Subscription not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/subscriptions/{id} [put]
func (s *Server) updateSubscription(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Subscription ID is required"})
		return
	}

	var req ports.UpdateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	subscription, err := s.service.UpdateSubscription(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update subscription", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": subscription})
}

// cancelSubscription godoc
// @Summary Cancel subscription
// @Description Cancels an active subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 200 {object} domain.Subscription "Canceled subscription details"
// @Failure 400 {object} map[string]interface{} "Invalid subscription ID"
// @Failure 404 {object} map[string]interface{} "Subscription not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/subscriptions/{id} [delete]
func (s *Server) cancelSubscription(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Subscription ID is required"})
		return
	}

	subscription, err := s.service.CancelSubscription(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel subscription", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": subscription})
}

// Payment handlers

// createPayment godoc
// @Summary Create a new payment
// @Description Creates a new one-time payment for a customer
// @Tags payments
// @Accept json
// @Produce json
// @Param payment body ports.CreatePaymentRequest true "Payment creation request"
// @Success 201 {object} domain.Payment "Payment created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body or validation error"
// @Failure 404 {object} map[string]interface{} "Customer not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/payments [post]
func (s *Server) createPayment(c *gin.Context) {
	var req ports.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	payment, err := s.service.CreatePayment(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": payment})
}

// listPayments godoc
// @Summary List payments
// @Description Retrieves a paginated list of payments with optional filtering
// @Tags payments
// @Accept json
// @Produce json
// @Param offset query int false "Number of records to skip" default(0)
// @Param limit query int false "Maximum number of records to return" default(10)
// @Param customer_id query string false "Filter by customer ID"
// @Param subscription_id query string false "Filter by subscription ID"
// @Param status query string false "Filter by status (succeeded, failed, etc.)"
// @Success 200 {object} ports.ListPaymentsResponse "List of payments"
// @Failure 400 {object} map[string]interface{} "Invalid query parameters"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/payments [get]
func (s *Server) listPayments(c *gin.Context) {
	var req ports.ListPaymentsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters", "details": err.Error()})
		return
	}

	if req.Limit == 0 {
		req.Limit = 10
	}

	response, err := s.service.ListPayments(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list payments", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// getPayment godoc
// @Summary Get payment by ID
// @Description Retrieves a specific payment by its ID
// @Tags payments
// @Accept json
// @Produce json
// @Param id path string true "Payment ID"
// @Success 200 {object} domain.Payment "Payment details"
// @Failure 400 {object} map[string]interface{} "Invalid payment ID"
// @Failure 404 {object} map[string]interface{} "Payment not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/payments/{id} [get]
func (s *Server) getPayment(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payment ID is required"})
		return
	}

	payment, err := s.service.GetPayment(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": payment})
}

// refundPayment godoc
// @Summary Refund payment
// @Description Issues a refund for an existing payment
// @Tags payments
// @Accept json
// @Produce json
// @Param id path string true "Payment ID"
// @Param refund body object false "Refund request (optional amount and reason)"
// @Success 200 {object} domain.Payment "Refunded payment details"
// @Failure 400 {object} map[string]interface{} "Invalid payment ID or refund request"
// @Failure 404 {object} map[string]interface{} "Payment not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/payments/{id}/refund [post]
func (s *Server) refundPayment(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payment ID is required"})
		return
	}

	var reqBody struct {
		Amount *int64 `json:"amount,omitempty"`
	}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	payment, err := s.service.RefundPayment(c.Request.Context(), id, reqBody.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to refund payment", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": payment})
}

// Payment Intent handlers

// createPaymentIntent godoc
// @Summary Create payment intent
// @Description Creates a Stripe Payment Intent for collecting payment from customers
// @Tags payment-intents
// @Accept json
// @Produce json
// @Param paymentIntent body ports.CreatePaymentIntentRequest true "Payment intent creation request"
// @Success 201 {object} ports.CreatePaymentIntentResponse "Payment intent created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body or validation error"
// @Failure 404 {object} map[string]interface{} "Customer not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/payment-intents [post]
func (s *Server) createPaymentIntent(c *gin.Context) {
	var req ports.CreatePaymentIntentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	response, err := s.service.CreatePaymentIntent(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment intent", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": response})
}

// confirmPaymentIntent godoc
// @Summary Confirm payment intent
// @Description Confirms a Stripe Payment Intent to complete the payment process
// @Tags payment-intents
// @Accept json
// @Produce json
// @Param id path string true "Payment Intent ID"
// @Success 200 {object} ports.ConfirmPaymentIntentResponse "Payment intent confirmed successfully"
// @Failure 400 {object} map[string]interface{} "Invalid payment intent ID or confirmation request"
// @Failure 404 {object} map[string]interface{} "Payment intent not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/payment-intents/{id}/confirm [post]
func (s *Server) confirmPaymentIntent(c *gin.Context) {
	paymentIntentID := c.Param("id")
	if paymentIntentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payment intent ID is required"})
		return
	}

	response, err := s.service.ConfirmPaymentIntent(c.Request.Context(), paymentIntentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm payment intent", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// Webhook handler

// handleStripeWebhook godoc
// @Summary Handle Stripe webhook
// @Description Processes incoming Stripe webhook events for payment updates
// @Tags webhooks
// @Accept json
// @Produce json
// @Param Stripe-Signature header string true "Stripe webhook signature for verification"
// @Param webhook body object true "Stripe webhook payload"
// @Success 200 {object} map[string]string "Webhook processed successfully"
// @Failure 400 {object} map[string]interface{} "Invalid webhook signature or payload"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/webhooks/stripe [post]
func (s *Server) handleStripeWebhook(c *gin.Context) {
	payload, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	signature := c.GetHeader("Stripe-Signature")
	if signature == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing Stripe signature header"})
		return
	}

	err = s.service.HandleWebhook(c.Request.Context(), payload, signature)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to process webhook", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webhook processed successfully"})
}
