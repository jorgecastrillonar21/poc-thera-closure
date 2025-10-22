package main

import (
	"log"

	"github.com/redis/go-redis/v9"

	"theraclosure/payments-service/internal/adapters/config"
	"theraclosure/payments-service/internal/adapters/http"
	"theraclosure/payments-service/internal/adapters/logging"
	"theraclosure/payments-service/internal/adapters/persistence"
	"theraclosure/payments-service/internal/adapters/stripe"
	"theraclosure/payments-service/internal/core/services"
)

func main_backup() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize structured logger
	logger := logging.NewLogger(cfg.App.LogLevel)
	logger.WithFields(logging.Fields{
		{Key: "service", Value: cfg.App.Name},
		{Key: "version", Value: cfg.App.Version},
	}).Info("Starting payments service")

	// Initialize database
	db, err := persistence.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize repositories
	customerRepo := persistence.NewCustomerRepository(db)
	subscriptionRepo := persistence.NewSubscriptionRepository(db)
	paymentRepo := persistence.NewPaymentRepository(db)

	// Initialize Redis client (optional)
	var redisClient *redis.Client
	if cfg.Redis.Address != "" {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     cfg.Redis.Address,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		})
		log.Println("Redis client initialized")
	} else {
		log.Println("Redis not configured, health checks will show degraded status")
	}

	// Initialize Stripe client
	stripeClient := stripe.NewStripeClient(cfg.Stripe.SecretKey)

	// Initialize services
	paymentService := services.NewPaymentService(
		customerRepo,
		subscriptionRepo,
		paymentRepo,
		stripeClient,
	)

	// Initialize HTTP server with database, Redis, and logger
	server := http.NewServer(paymentService, cfg, db.GetDB(), redisClient, logger)

	// Start server
	log.Printf("Payments service starting on %s", cfg.GetServerAddress())
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
