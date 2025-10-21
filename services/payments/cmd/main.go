package main

import (
	"log"

	"theraclosure/payments-service/internal/adapters/config"
	"theraclosure/payments-service/internal/adapters/http"
	"theraclosure/payments-service/internal/adapters/persistence"
	"theraclosure/payments-service/internal/adapters/stripe"
	"theraclosure/payments-service/internal/core/services"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := persistence.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize repositories
	customerRepo := persistence.NewCustomerRepository(db)
	subscriptionRepo := persistence.NewSubscriptionRepository(db)
	paymentRepo := persistence.NewPaymentRepository(db)

	// Initialize Stripe client
	stripeClient := stripe.NewStripeClient(cfg.Stripe.SecretKey)

	// Initialize services
	paymentService := services.NewPaymentService(
		customerRepo,
		subscriptionRepo,
		paymentRepo,
		stripeClient,
	)

	// Initialize HTTP server
	server := http.NewServer(paymentService, cfg)

	// Start server
	log.Printf("Payments service starting on %s", cfg.GetServerAddress())
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}