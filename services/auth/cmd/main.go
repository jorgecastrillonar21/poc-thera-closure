package main

import (
	"log"
	"theraclosure/auth-service/internal/adapters/config"
	"theraclosure/auth-service/internal/adapters/http"
	"theraclosure/auth-service/internal/adapters/persistence"
	"theraclosure/auth-service/internal/core/services"
)

// @title TheraClosure Auth Service API
// @version 1.0
// @description Authentication microservice for TheraClosure application
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.theraclosure.com/support
// @contact.email support@theraclosure.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3001
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize database
	db, err := persistence.NewPostgresDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize repositories (adapters)
	userRepo := persistence.NewUserRepository(db)
	sessionRepo := persistence.NewSessionRepository(db)

	// Initialize domain services (core)
	authService := services.NewAuthService(userRepo, sessionRepo, cfg)
	userService := services.NewUserService(userRepo)

	// Initialize HTTP server (adapter)
	server := http.NewServer(authService, userService, cfg)

	// Start server
	log.Printf("Starting %s on port %s", cfg.App.Name, cfg.Server.Port)
	if err := server.Start(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
