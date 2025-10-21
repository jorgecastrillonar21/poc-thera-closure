package http

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"theraclosure/payments-service/internal/adapters/config"
	"theraclosure/payments-service/internal/core/ports"
)

type Server struct {
	service ports.PaymentService
	config  *config.Config
	router  *gin.Engine
}

func NewServer(service ports.PaymentService, cfg *config.Config) *Server {
	server := &Server{
		service: service,
		config:  cfg,
	}

	server.setupRouter()
	return server
}

func (s *Server) setupRouter() {
	// Set Gin mode based on config
	if s.config.App.LogLevel == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	s.router = gin.New()
	s.router.Use(gin.Logger())
	s.router.Use(gin.Recovery())

	// CORS middleware
	corsConfig := cors.Config{
		AllowOrigins:     s.config.CORS.AllowedOrigins,
		AllowMethods:     s.config.CORS.AllowedMethods,
		AllowHeaders:     s.config.CORS.AllowedHeaders,
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	s.router.Use(cors.New(corsConfig))

	// Health check endpoint
	s.router.GET("/health", s.healthCheck)

	// API v1 routes
	v1 := s.router.Group("/api/v1")
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

func (s *Server) Start() error {
	return s.router.Run(s.config.GetServerAddress())
}

// Health check handler
func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"service":   s.config.App.Name,
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
	})
}

// Customer handlers
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